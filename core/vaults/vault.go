package vaults

import (
	"os"
	"path"
	"regexp"
	"strings"

	"savesecrets.org/slv/core/commons"
	"savesecrets.org/slv/core/crypto"
)

type vaultConfig struct {
	PublicKey   string   `json:"publicKey" yaml:"publicKey"`
	HashLength  uint8    `json:"hashLength,omitempty" yaml:"hashLength,omitempty"`
	WrappedKeys []string `json:"wrappedKeys" yaml:"wrappedKeys"`
}

type Vault struct {
	Secrets             map[string]string `json:"slvSecrets" yaml:"slvSecrets"`
	Config              vaultConfig       `json:"slvConfig" yaml:"slvConfig"`
	path                string            `json:"-"`
	publicKey           *crypto.PublicKey `json:"-"`
	secretKey           *crypto.SecretKey `json:"-"`
	unlockedBy          *string           `json:"-"`
	decryptedSecrets    map[string][]byte `json:"-"`
	vaultSecretRefRegex *regexp.Regexp    `json:"-"`
	objectField         string            `json:"-"`
}

func (vlt *Vault) Id() string {
	return vlt.Config.PublicKey
}

func (vlt *Vault) getPublicKey() (publicKey *crypto.PublicKey, err error) {
	if vlt.publicKey == nil {
		if vlt.Config.PublicKey == "" {
			return nil, errVaultPublicKeyNotFound
		}
		publicKey, err = crypto.PublicKeyFromString(vlt.Config.PublicKey)
		if err == nil {
			vlt.publicKey = publicKey
		}
	}
	return vlt.publicKey, err
}

func isValidVaultFileName(fileName string) bool {
	return strings.HasSuffix(fileName, "."+vaultFileNameEnding) ||
		strings.HasSuffix(fileName, vaultFileNameEnding+".yaml") ||
		strings.HasSuffix(fileName, vaultFileNameEnding+".yml")
}

// Returns new vault instance and the vault contents set into the specified field. The vault file name must end with .slv.yml or .slv.yaml.
func New(filePath, objectField string, hashLength uint8, rootPublicKey *crypto.PublicKey, publicKeys ...*crypto.PublicKey) (vlt *Vault, err error) {
	if !isValidVaultFileName(filePath) {
		return nil, errInvalidVaultFileName
	}
	if commons.FileExists(filePath) {
		return nil, errVaultExists
	}
	if os.MkdirAll(path.Dir(filePath), os.FileMode(0755)) != nil {
		return nil, errVaultDirPathCreation
	}
	vaultSecretKey, err := crypto.NewSecretKey(VaultKey)
	if err != nil {
		return nil, err
	}
	vaultPublicKey, err := vaultSecretKey.PublicKey()
	if err != nil {
		return nil, err
	}
	vlt = &Vault{
		publicKey: vaultPublicKey,
		Config: vaultConfig{
			PublicKey:  vaultPublicKey.String(),
			HashLength: hashLength,
		},
		path:        filePath,
		secretKey:   vaultSecretKey,
		objectField: objectField,
	}
	if rootPublicKey != nil {
		if _, err := vlt.Share(rootPublicKey); err != nil {
			return nil, err
		}
	}
	for _, pubKey := range publicKeys {
		if _, err := vlt.Share(pubKey); err != nil {
			return nil, err
		}
	}
	return vlt, vlt.commit()
}

// Returns the vault instance from a given yaml. The vault file name must end with .slv.yml or .slv.yaml.
func Get(filePath string) (vlt *Vault, err error) {
	return GetFromField(filePath, "")
}

// Returns the vault instance from a given yaml file considering a field as vault. The vault file name must end with .slv.yml or .slv.yaml.
func GetFromField(filePath, fieldName string) (vlt *Vault, err error) {
	if !isValidVaultFileName(filePath) {
		return nil, errInvalidVaultFileName
	}
	if !commons.FileExists(filePath) {
		return nil, errVaultNotFound
	}
	vlt = &Vault{}
	if fieldName == "" {
		if err = commons.ReadFromYAML(filePath, &vlt); err != nil {
			return nil, err
		}
	} else {
		if err = commons.ReadChildFromYAML(filePath, fieldName, &vlt); err != nil {
			return nil, err
		}
		vlt.objectField = fieldName
	}
	vlt.path = filePath
	return vlt, nil
}

func (vlt *Vault) IsLocked() bool {
	return vlt.secretKey == nil
}

func (vlt *Vault) Lock() {
	vlt.clearSecretCache()
	vlt.secretKey = nil
}

func (vlt *Vault) UnlockedBy() (id *string) {
	return vlt.unlockedBy
}

func (vlt *Vault) ListAccessors() ([]crypto.PublicKey, error) {
	var accessors []crypto.PublicKey
	for _, wrappedKeyStr := range vlt.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		err := wrappedKey.FromString(wrappedKeyStr)
		if err != nil {
			return nil, err
		}
		accessors = append(accessors, wrappedKey.EncryptedBy())
	}
	return accessors, nil
}

func (vlt *Vault) Unlock(secretKey crypto.SecretKey) error {
	publicKey, err := secretKey.PublicKey()
	if err != nil || (!vlt.IsLocked() && *vlt.unlockedBy == publicKey.String()) {
		return err
	}
	for _, wrappedKeyStr := range vlt.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		if err = wrappedKey.FromString(wrappedKeyStr); err != nil {
			return err
		}
		decryptedKey, err := secretKey.DecryptKey(*wrappedKey)
		if err == nil {
			vlt.secretKey = decryptedKey
			vlt.unlockedBy = new(string)
			*vlt.unlockedBy = publicKey.String()
			return nil
		}
	}
	return errVaultNotAccessible
}

func (vlt *Vault) commit() error {
	if vlt.objectField == "" {
		return commons.WriteToYAML(vlt.path,
			"# Use the pattern "+vlt.getSecretRef("YOUR_SECRET_NAME")+" as placeholder to reference secrets from this vault into files\n", vlt)
	} else {
		var obj map[string]interface{}
		if commons.FileExists(vlt.path) {
			if err := commons.ReadFromYAML(vlt.path, &obj); err != nil {
				return err
			}
		} else {
			obj = make(map[string]interface{})
		}
		obj[vlt.objectField] = vlt
		return commons.WriteToYAML(vlt.path,
			"# Use the pattern "+vlt.getSecretRef("YOUR_SECRET_NAME")+" as placeholder to reference secrets from this vault into files\n", obj)
	}
}

func (vlt *Vault) reset() error {
	vlt.clearSecretCache()
	return commons.ReadFromYAML(vlt.path, &vlt)
}

func (vlt *Vault) Share(publicKey *crypto.PublicKey) (bool, error) {
	if vlt.IsLocked() {
		return false, errVaultLocked
	}
	if publicKey.Type() == VaultKey {
		return false, errVaultCannotBeSharedWithVault
	}
	for _, wrappedKeyStr := range vlt.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		if err := wrappedKey.FromString(wrappedKeyStr); err != nil {
			return false, err
		}
		if wrappedKey.IsEncryptedBy(publicKey) {
			return false, nil
		}
	}
	wrappedKey, err := publicKey.EncryptKey(*vlt.secretKey)
	if err == nil {
		vlt.Config.WrappedKeys = append(vlt.Config.WrappedKeys, wrappedKey.String())
		err = vlt.commit()
	}
	return err == nil, err
}

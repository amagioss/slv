package vaults

import (
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"gopkg.in/yaml.v3"
)

type config struct {
	PublicKey   *string   `yaml:"publicKey"`
	HashLength  *uint32   `yaml:"hashLength,omitempty"`
	WrappedKeys []*string `yaml:"wrappedKeys,omitempty"`
	publicKey   *crypto.PublicKey
}

type vault struct {
	Secrets map[string]*string `yaml:"slvSecrets,omitempty"`
	Config  config             `yaml:"slvConfig,omitempty"`
}

type Vault struct {
	*vault
	path                 string
	secretKey            *crypto.SecretKey
	unlockedBy           *string
	decryptedSecretCache map[string][]byte
	vaultSecretRefRegex  *regexp.Regexp
}

func (vlt *Vault) Id() string {
	return *vlt.Config.PublicKey
}

func (vlt *Vault) getPublicKey() (publicKey *crypto.PublicKey, err error) {
	if vlt.Config.publicKey == nil {
		if vlt.Config.PublicKey == nil {
			return nil, ErrVaultPublicKeyNotFound
		}
		publicKey, err = crypto.PublicKeyFromString(*vlt.Config.PublicKey)
		if err == nil {
			vlt.Config.publicKey = publicKey
		}
	}
	return vlt.Config.publicKey, err
}

func (vlt Vault) MarshalYAML() (interface{}, error) {
	return vlt.vault, nil
}

func (vlt *Vault) UnmarshalYAML(value *yaml.Node) (err error) {
	return value.Decode(&vlt.vault)
}

// Returns new vault instance. The vault name should end with .vlt.slv
func New(vaultFile string, hashLength uint32, rootPublicKey *crypto.PublicKey, publicKeys ...crypto.PublicKey) (vlt *Vault, err error) {
	if !strings.HasSuffix(vaultFile, vaultFileExtension) {
		return nil, ErrInvalidVaultFileName
	}
	if commons.FileExists(vaultFile) {
		return nil, ErrVaultExists
	}
	if os.MkdirAll(path.Dir(vaultFile), os.FileMode(0755)) != nil {
		return nil, ErrVaultDirPathCreation
	}
	vaultSecretKey, err := crypto.NewSecretKey(VaultKey)
	if err != nil {
		return nil, err
	}
	vaultPublicKey, err := vaultSecretKey.PublicKey()
	if err != nil {
		return nil, err
	}
	var hashLen *uint32
	if hashLength > 0 {
		hashLen = &hashLength
	}
	vlt = &Vault{
		vault: &vault{
			Config: config{
				publicKey:  vaultPublicKey,
				PublicKey:  commons.StringPtr(vaultPublicKey.String()),
				HashLength: hashLen,
			},
		},
		path:      vaultFile,
		secretKey: vaultSecretKey,
	}
	if rootPublicKey != nil {
		if _, err := vlt.Share(*rootPublicKey); err != nil {
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

// Returns the vault instance for a given vault file. The vault name should end with .slv
func Get(vaultFile string) (vlt *Vault, err error) {
	if !strings.HasSuffix(vaultFile, vaultFileExtension) {
		return nil, ErrInvalidVaultFileName
	}
	if !commons.FileExists(vaultFile) {
		return nil, ErrVaultNotFound
	}
	vlt = &Vault{
		path: vaultFile,
	}
	if err = commons.ReadFromYAML(vlt.path, &vlt.vault); err != nil {
		return nil, err
	}
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

func (vlt *Vault) Unlock(secretKey crypto.SecretKey) error {
	publicKey, err := secretKey.PublicKey()
	if err != nil || (!vlt.IsLocked() && *vlt.unlockedBy == publicKey.String()) {
		return err
	}
	for _, wrappedKeyStr := range vlt.vault.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		if err = wrappedKey.FromString(*wrappedKeyStr); err != nil {
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
	return ErrVaultNotAccessible
}

func (vlt *Vault) commit() error {
	return commons.WriteToYAML(vlt.path,
		"# Use the pattern "+vlt.getSecretRef("YOUR_SECRET_NAME")+" to reference secrets from this vault into files\n", *vlt.vault)
}

func (vlt *Vault) reset() error {
	vlt.clearSecretCache()
	return commons.ReadFromYAML(vlt.path, &vlt.vault)
}

func (vlt *Vault) Share(publicKey crypto.PublicKey) (bool, error) {
	if vlt.IsLocked() {
		return false, ErrVaultLocked
	}
	if publicKey.Type() == VaultKey {
		return false, ErrVaultCannotBeSharedWithVault
	}
	for _, wrappedKeyStr := range vlt.vault.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		if err := wrappedKey.FromString(*wrappedKeyStr); err != nil {
			return false, err
		}
		if wrappedKey.IsEncryptedBy(&publicKey) {
			return false, nil
		}
	}
	wrappedKey, err := publicKey.EncryptKey(*vlt.secretKey)
	if err == nil {
		vlt.vault.Config.WrappedKeys = append(vlt.vault.Config.WrappedKeys, commons.StringPtr(wrappedKey.String()))
		err = vlt.commit()
	}
	return err == nil, err
}

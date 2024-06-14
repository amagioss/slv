package vaults

import (
	"crypto/rand"
	"os"
	"path"
	"regexp"
	"strings"

	"golang.org/x/mod/semver"
	"oss.amagi.com/slv/internal/core/commons"
	"oss.amagi.com/slv/internal/core/config"
	"oss.amagi.com/slv/internal/core/crypto"
)

type vaultConfig struct {
	Id          string   `json:"id" yaml:"id"`
	PublicKey   string   `json:"publicKey" yaml:"publicKey"`
	HashLength  uint8    `json:"hashLength,omitempty" yaml:"hashLength,omitempty"`
	WrappedKeys []string `json:"wrappedKeys" yaml:"wrappedKeys"`
}

type Vault struct {
	Version             string            `json:"version,omitempty" yaml:"version,omitempty"`
	Secrets             map[string]string `json:"slvSecrets" yaml:"slvSecrets"`
	Config              vaultConfig       `json:"slvConfig" yaml:"slvConfig"`
	path                string            `json:"-"`
	publicKey           *crypto.PublicKey `json:"-"`
	secretKey           *crypto.SecretKey `json:"-"`
	decryptedSecrets    map[string][]byte `json:"-"`
	vaultSecretRefRegex *regexp.Regexp    `json:"-"`
	k8s                 *k8slv            `json:"-"`
}

func (vlt *Vault) Id() string {
	return vlt.Config.Id
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

func newVaultId() (string, error) {
	idBytes := make([]byte, vaultIdLength)
	if _, err := rand.Read(idBytes); err != nil {
		return "", errGeneratingVaultId
	}
	return config.AppNameUpperCase + "_" + vaultIdAbbrev + "_" + commons.Encode(idBytes), nil
}

// Returns new vault instance and the vault contents set into the specified field. The vault file name must end with .slv.yml or .slv.yaml.
func New(filePath, k8sName, k8SecretFile string, hashLength uint8, quantumSafe bool, rootPublicKey *crypto.PublicKey, publicKeys ...*crypto.PublicKey) (vlt *Vault, err error) {
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
	vaultPublicKey, err := vaultSecretKey.PublicKey(quantumSafe)
	if err != nil {
		return nil, err
	}
	vauldId, err := newVaultId()
	if err != nil {
		return nil, err
	}
	vaultPubKeyStr, err := vaultPublicKey.String()
	if err != nil {
		return nil, err
	}
	vlt = &Vault{
		Version:   config.Version,
		publicKey: vaultPublicKey,
		Config: vaultConfig{
			Id:         vauldId,
			PublicKey:  vaultPubKeyStr,
			HashLength: hashLength,
		},
		path:      filePath,
		secretKey: vaultSecretKey,
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
	if k8sName == "" && k8SecretFile == "" {
		return vlt, vlt.commit()
	} else {
		return vlt, vlt.ToK8s(k8sName, k8SecretFile)
	}
}

// Returns the vault instance from a given yaml. The vault file name must end with .slv.yml or .slv.yaml.
func Get(filePath string) (vlt *Vault, err error) {
	obj := make(map[string]interface{})
	if err := commons.ReadFromYAML(filePath, &obj); err != nil {
		return nil, err
	}
	if obj[k8sVaultField] != nil {
		return getFromField(filePath, true)
	}
	return getFromField(filePath, false)
}

func getFromField(filePath string, k8s bool) (vlt *Vault, err error) {
	if !isValidVaultFileName(filePath) {
		return nil, errInvalidVaultFileName
	}
	if !commons.FileExists(filePath) {
		return nil, errVaultNotFound
	}
	if k8s {
		k8sVault := &k8slv{}
		if err = commons.ReadFromYAML(filePath, k8sVault); err != nil {
			return nil, err
		}
		vlt = &k8sVault.Spec
		vlt.k8s = k8sVault
	} else {
		vlt = &Vault{}
		if err = commons.ReadFromYAML(filePath, vlt); err != nil {
			return nil, err
		}
	}
	vlt.path = filePath
	vaultVersion := vlt.Version
	if vaultVersion != "" && (!semver.IsValid(vaultVersion) || semver.Compare(config.Version, vaultVersion) < 0) {
		return nil, errVaultVersionNotRecognized
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

func (vlt *Vault) Delete() error {
	vlt.clearSecretCache()
	return os.Remove(vlt.path)
}

func (vlt *Vault) commit() error {
	var data interface{}
	if vlt.k8s != nil {
		data = vlt.k8s
	} else {
		data = vlt
	}
	return commons.WriteToYAML(vlt.path,
		"# Use the pattern "+vlt.getSecretRef("YOUR_SECRET_NAME")+" as placeholder to reference secrets from this vault into files\n", data)
}

func (vlt *Vault) reset() error {
	vlt.clearSecretCache()
	return commons.ReadFromYAML(vlt.path, &vlt)
}

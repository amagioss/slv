package vaults

import (
	"crypto/rand"
	"encoding/json"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"golang.org/x/mod/semver"
	"oss.amagi.com/slv/internal/core/commons"
	"oss.amagi.com/slv/internal/core/config"
	"oss.amagi.com/slv/internal/core/crypto"
)

type vaultConfig struct {
	Version     string   `json:"version,omitempty" yaml:"version,omitempty"`
	Id          string   `json:"id" yaml:"id"`
	PublicKey   string   `json:"publicKey" yaml:"publicKey"`
	Hash        bool     `json:"hash,omitempty" yaml:"hash,omitempty"`
	WrappedKeys []string `json:"wrappedKeys" yaml:"wrappedKeys"`
}

type Vault struct {
	Secrets             map[string]string     `json:"slvSecrets,omitempty" yaml:"slvSecrets,omitempty"`
	Data                map[string]string     `json:"slvData,omitempty" yaml:"slvData,omitempty"`
	Config              vaultConfig           `json:"slvConfig" yaml:"slvConfig"`
	path                string                `json:"-"`
	publicKey           *crypto.PublicKey     `json:"-"`
	secretKey           *crypto.SecretKey     `json:"-"`
	cache               map[string]*VaultData `json:"-"`
	vaultSecretRefRegex *regexp.Regexp        `json:"-"`
	k8s                 *k8slv                `json:"-"`
}

type VaultData struct {
	value     []byte     `json:"-"`
	isSecret  bool       `json:"-"`
	updatedAt *time.Time `json:"-"`
	hash      string     `json:"-"`
}

func (vd *VaultData) Value() []byte {
	return vd.value
}

func (vd *VaultData) IsSecret() bool {
	return vd.isSecret
}

func (vd *VaultData) UpdatedAt() *time.Time {
	return vd.updatedAt
}

func (vd *VaultData) Hash() string {
	return vd.hash
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
func New(filePath, k8sName, k8sNamespace string, k8SecretContent []byte, hash, quantumSafe bool, rootPublicKey *crypto.PublicKey, publicKeys ...*crypto.PublicKey) (vlt *Vault, err error) {
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
		publicKey: vaultPublicKey,
		Config: vaultConfig{
			Version:   semver.MajorMinor(config.Version),
			Id:        vauldId,
			PublicKey: vaultPubKeyStr,
			Hash:      hash,
		},
		path:      filePath,
		secretKey: vaultSecretKey,
	}
	if rootPublicKey != nil {
		if _, err := vlt.share(rootPublicKey, false); err != nil {
			return nil, err
		}
	}
	for _, pubKey := range publicKeys {
		if _, err := vlt.share(pubKey, false); err != nil {
			return nil, err
		}
	}
	if k8sName == "" && k8SecretContent == nil {
		return vlt, vlt.commit()
	} else {
		return vlt, vlt.ToK8s(k8sName, k8sNamespace, k8SecretContent)
	}
}

// Returns the vault instance from a given yaml. The vault file name must end with .slv.yml or .slv.yaml.
func Get(filePath string) (vlt *Vault, err error) {
	if !isValidVaultFileName(filePath) {
		return nil, errInvalidVaultFileName
	}
	if !commons.FileExists(filePath) {
		return nil, errVaultNotFound
	}
	obj := make(map[string]interface{})
	if err := commons.ReadFromYAML(filePath, &obj); err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return getFromField(jsonData, filePath, obj[k8sVaultField] != nil)
}

func getFromField(jsonData []byte, filePath string, k8s bool) (vlt *Vault, err error) {
	if k8s {
		k8sVault := &k8slv{}
		if err = json.Unmarshal(jsonData, k8sVault); err != nil {
			return nil, err
		}
		vlt = k8sVault.Spec
		vlt.k8s = k8sVault
	} else {
		vlt = &Vault{}
		if err = json.Unmarshal(jsonData, vlt); err != nil {
			return nil, err
		}
	}
	vlt.path = filePath
	if vlt.Config.Version != "" && config.Version != "" &&
		(!semver.IsValid(vlt.Config.Version) || !semver.IsValid(config.Version) ||
			semver.Compare(semver.MajorMinor(config.Version), semver.MajorMinor(vlt.Config.Version)) < 0) {
		return nil, errVaultVersionNotRecognized
	}
	if vlt.Secrets != nil {
		if vlt.Data != nil {
			for key, value := range vlt.Data {
				vlt.Secrets[key] = value
			}
		}
		vlt.Data = vlt.Secrets
		vlt.Secrets = nil
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
	if err := vlt.validate(); err != nil {
		return err
	}
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

func (vlt *Vault) validate() error {
	if vlt.Config.PublicKey == "" {
		return errVaultPublicKeyNotFound
	}
	if len(vlt.Config.WrappedKeys) == 0 {
		return errVaultWrappedKeysNotFound
	}
	return nil
}

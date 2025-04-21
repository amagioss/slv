package vaults

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"slv.sh/slv/internal/core/commons"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
)

type vaultConfig struct {
	Id          string   `json:"id" yaml:"id"`
	PublicKey   string   `json:"publicKey" yaml:"publicKey"`
	Hash        bool     `json:"hash,omitempty" yaml:"hash,omitempty"`
	WrappedKeys []string `json:"wrappedKeys" yaml:"wrappedKeys"`
}

type Vault struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Type string     `json:"type,omitempty" yaml:"type,omitempty"`
	Spec *VaultSpec `json:"spec" yaml:"spec"`
}

type VaultSpec struct {
	Secrets             map[string]string     `json:"slvSecrets,omitempty" yaml:"slvSecrets,omitempty"`
	Data                map[string]string     `json:"slvData,omitempty" yaml:"slvData,omitempty"`
	Config              vaultConfig           `json:"slvConfig" yaml:"slvConfig"`
	path                string                `json:"-"`
	publicKey           *crypto.PublicKey     `json:"-"`
	secretKey           *crypto.SecretKey     `json:"-"`
	cache               map[string]*VaultItem `json:"-"`
	vaultSecretRefRegex *regexp.Regexp        `json:"-"`
}

type VaultItem struct {
	value     []byte     `json:"-"`
	isSecret  bool       `json:"-"`
	updatedAt *time.Time `json:"-"`
	hash      string     `json:"-"`
}

func (vi *VaultItem) Value() []byte {
	return vi.value
}

func (vi *VaultItem) IsSecret() bool {
	return vi.isSecret
}

func (vi *VaultItem) UpdatedAt() *time.Time {
	return vi.updatedAt
}

func (vi *VaultItem) Hash() string {
	return vi.hash
}

func (vlt *Vault) Id() string {
	return vlt.Spec.Config.Id
}

func (vlt *Vault) getPublicKey() (publicKey *crypto.PublicKey, err error) {
	if vlt.Spec.publicKey == nil {
		if vlt.Spec.Config.PublicKey == "" {
			return nil, errVaultPublicKeyNotFound
		}
		publicKey, err = crypto.PublicKeyFromString(vlt.Spec.Config.PublicKey)
		if err == nil {
			vlt.Spec.publicKey = publicKey
		}
	}
	return vlt.Spec.publicKey, err
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
func New(filePath, vaultName, k8sNamespace string, k8SecretContent []byte, hash, quantumSafe bool, publicKeys ...*crypto.PublicKey) (vlt *Vault, err error) {
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
	if vaultName == "" {
		vaultName = getNameFromFilePath(filePath)
	}
	vlt = &Vault{
		TypeMeta: metav1.TypeMeta{
			APIVersion: k8sApiVersion,
			Kind:       k8sKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: vaultName,
		},
		Spec: &VaultSpec{
			publicKey: vaultPublicKey,
			Config: vaultConfig{
				Id:        vauldId,
				PublicKey: vaultPubKeyStr,
				Hash:      hash,
			},
			path:      filePath,
			secretKey: vaultSecretKey,
		},
	}
	for _, pubKey := range publicKeys {
		if _, err := vlt.share(pubKey, false); err != nil {
			return nil, err
		}
	}
	return vlt, vlt.commit()
}

// Returns the vault instance from a given yaml. The vault file name must end with .slv.yml or .slv.yaml.
func Get(filePath string) (vlt *Vault, err error) {
	if !isValidVaultFileName(filePath) {
		return nil, errInvalidVaultFileName
	}
	if !commons.FileExists(filePath) {
		return nil, errVaultNotFound
	}
	obj := make(map[string]any)
	if err := commons.ReadFromYAML(filePath, &obj); err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return get(jsonData, filePath, obj[k8sVaultSpecField] != nil)
}

func get(jsonData []byte, filePath string, fullVault bool) (vlt *Vault, err error) {
	if fullVault {
		vlt = &Vault{}
		if err = json.Unmarshal(jsonData, vlt); err != nil {
			return nil, err
		}
	} else {
		vs := &VaultSpec{}
		if err = json.Unmarshal(jsonData, vs); err != nil {
			return nil, err
		}
		vlt = &Vault{
			Spec: vs,
		}
	}
	vlt.validate()
	vlt.Spec.path = filePath
	vlt.init()
	return vlt, nil
}

func (vlt *Vault) init() {
	if vlt.Spec.Secrets != nil {
		if vlt.Spec.Data != nil {
			for key, value := range vlt.Spec.Data {
				vlt.Spec.Secrets[key] = value
			}
		}
		vlt.Spec.Data = vlt.Spec.Secrets
		vlt.Spec.Secrets = nil
	}
}

func (vlt *Vault) IsLocked() bool {
	return vlt.Spec.secretKey == nil
}

func (vlt *Vault) Lock() {
	vlt.clearSecretCache()
	vlt.Spec.secretKey = nil
}

func (vlt *Vault) Delete() error {
	vlt.clearSecretCache()
	return os.Remove(vlt.Spec.path)
}

func (vlt *Vault) commit() error {
	if err := vlt.validate(); err != nil {
		return err
	}
	jsonData, err := json.Marshal(vlt)
	if err != nil {
		return err
	}
	var data any
	if err = json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}
	return commons.WriteToYAML(vlt.Spec.path,
		"# Use the pattern "+vlt.getDataRef("YOUR_SECRET_NAME")+
			" as placeholder to reference data from this vault into files\n", data)
}

func (vlt *Vault) reset() error {
	vlt.clearSecretCache()
	return commons.ReadFromYAML(vlt.Spec.path, &vlt)
}

func getNameFromFilePath(path string) string {
	name := filepath.Base(path)
	return strings.TrimSuffix(name, "."+vaultFileNameEnding+filepath.Ext(name))
}

func (vlt *Vault) validate() error {
	vlt.APIVersion = k8sApiVersion
	vlt.Kind = k8sKind
	if vlt.ObjectMeta.Name == "" {
		vlt.ObjectMeta.Name = getNameFromFilePath(vlt.Spec.path)
	}
	if vlt.Annotations == nil {
		vlt.Annotations = make(map[string]string)
	}
	if vlt.Annotations[k8sVersionAnnotationKey] == "" {
		vlt.Annotations[k8sVersionAnnotationKey] = config.Version
	}
	if vlt.ObjectMeta.CreationTimestamp.IsZero() {
		vlt.ObjectMeta.CreationTimestamp = metav1.Now()
	}
	if vlt.Spec.Config.PublicKey == "" {
		return errVaultPublicKeyNotFound
	}
	if len(vlt.Spec.Config.WrappedKeys) == 0 {
		return errVaultWrappedKeysNotFound
	}
	return nil
}

package vaults

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"slv.sh/slv/internal/core/commons"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
)

type vaultConfig struct {
	PublicKey   string   `json:"publicKey" yaml:"publicKey"`
	Hash        bool     `json:"hash,omitempty" yaml:"hash,omitempty"`
	WrappedKeys []string `json:"wrappedKeys" yaml:"wrappedKeys"`
}

type Vault struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`

	Type string     `json:"type,omitempty" yaml:"type,omitempty"`
	Spec *VaultSpec `json:"spec" yaml:"spec"`
}

type VaultSpec struct {
	Data                map[string]string     `json:"slvData,omitempty" yaml:"slvData,omitempty"`
	Config              vaultConfig           `json:"slvConfig" yaml:"slvConfig"`
	writable            bool                  `json:"-" yaml:"-"`
	path                string                `json:"-" yaml:"-"`
	publicKey           *crypto.PublicKey     `json:"-" yaml:"-"`
	secretKey           *crypto.SecretKey     `json:"-" yaml:"-"`
	cache               map[string]*VaultItem `json:"-" yaml:"-"`
	vaultSecretRefRegex *regexp.Regexp        `json:"-" yaml:"-"`
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
	return strings.HasSuffix(fileName, "."+vaultFileNameRawExt) ||
		strings.HasSuffix(fileName, "."+vaultFileNameRawExt+".yaml") ||
		strings.HasSuffix(fileName, "."+vaultFileNameRawExt+".yml")
}

// Returns new vault instance and the vault contents set into the specified field. The vault file name must end with .slv.yaml or .slv.yml.
func New(vaultFile, name, k8sNamespace string, hash, quantumSafe bool, publicKeys ...*crypto.PublicKey) (vlt *Vault, err error) {
	if !isValidVaultFileName(vaultFile) {
		vaultFile = vaultFile + vaultFileNameDesiredExt
	}
	if commons.FileExists(vaultFile) {
		return nil, errVaultExists
	}
	if os.MkdirAll(path.Dir(vaultFile), os.FileMode(0755)) != nil {
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
	vaultPubKeyStr, err := vaultPublicKey.String()
	if err != nil {
		return nil, err
	}
	if name == "" {
		name = getNameFromFilePath(vaultFile)
	}
	vlt = &Vault{
		TypeMeta: metav1.TypeMeta{
			APIVersion: k8sApiVersion,
			Kind:       k8sKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: k8sNamespace,
		},
		Spec: &VaultSpec{
			writable:  true,
			publicKey: vaultPublicKey,
			Config: vaultConfig{
				PublicKey: vaultPubKeyStr,
				Hash:      hash,
			},
			path:      vaultFile,
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

// Returns the vault instance for a given file or URL. The vault file path must end with .slv.yaml or .slv.yml.
func Get(vaultFileOrURL string) (vlt *Vault, err error) {
	var contents []byte
	var writable bool
	if strings.HasPrefix(vaultFileOrURL, "http://") || strings.HasPrefix(vaultFileOrURL, "https://") {
		headers := make(map[string]string)
		headers["User-Agent"] = config.AppNameUpperCase + "-" + config.Version + " (" + runtime.GOOS + "/" + runtime.GOARCH + ")"
		contents, err = commons.GetURLContents(vaultFileOrURL, headers)
	} else if !commons.FileExists(vaultFileOrURL) {
		return nil, errVaultNotFound
	} else {
		writable = true
		contents, err = os.ReadFile(vaultFileOrURL)
	}
	if err != nil {
		return nil, err
	}
	obj := make(map[string]any)
	if err = yaml.Unmarshal(contents, &obj); err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return get(jsonData, vaultFileOrURL, obj[k8sVaultSpecField] != nil, writable)
}

func get(jsonData []byte, filePath string, fullVault, writable bool) (vlt *Vault, err error) {
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
	vlt.Spec.path = filePath
	vlt.Spec.writable = writable
	err = vlt.validateAndUpdate()
	return
}

func (vlt *Vault) IsLocked() bool {
	return vlt.Spec.secretKey == nil
}

func (vlt *Vault) Lock() {
	vlt.clearCache()
	vlt.Spec.secretKey = nil
}

func (vlt *Vault) Delete() error {
	vlt.clearCache()
	return os.Remove(vlt.Spec.path)
}

func (vlt *Vault) commit() error {
	if !vlt.Spec.writable {
		return errVaultNotWritable
	}
	if err := vlt.validateAndUpdate(); err != nil {
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
	return commons.WriteToYAML(vlt.Spec.path, data)
}

func (vlt *Vault) reload() error {
	vlt.clearCache()
	return commons.ReadFromYAML(vlt.Spec.path, &vlt)
}

func getNameFromFilePath(path string) string {
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, "."+vaultFileNameRawExt+filepath.Ext(name))
	return regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(name, "_")
}

func (vlt *Vault) validateAndUpdate() error {
	vlt.APIVersion = k8sApiVersion
	vlt.Kind = k8sKind
	if vlt.ObjectMeta.Name == "" {
		vlt.ObjectMeta.Name = getNameFromFilePath(vlt.Spec.path)
	}
	if vlt.Annotations == nil {
		vlt.Annotations = make(map[string]string)
	}
	vlt.Annotations[k8sVersionAnnotationKey] = config.Version
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

func (vlt *Vault) Unlock(secretKey *crypto.SecretKey) error {
	if !vlt.IsLocked() {
		return nil
	}
	for _, wrappedKeyStr := range vlt.Spec.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		if err := wrappedKey.FromString(wrappedKeyStr); err != nil {
			return err
		}
		decryptedKey, err := secretKey.DecryptKey(*wrappedKey)
		if err == nil {
			vlt.Spec.secretKey = decryptedKey
			return nil
		}
	}
	return errVaultNotAccessible
}

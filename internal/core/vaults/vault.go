package vaults

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

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
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Type string     `json:"type,omitempty" yaml:"type,omitempty"`
	Spec *VaultSpec `json:"spec" yaml:"spec"`
}

type VaultSpec struct {
	Data                map[string]string     `json:"slvData,omitempty" yaml:"slvData,omitempty"`
	Config              vaultConfig           `json:"slvConfig" yaml:"slvConfig"`
	path                string                `json:"-"`
	publicKey           *crypto.PublicKey     `json:"-"`
	secretKey           *crypto.SecretKey     `json:"-"`
	cache               map[string]*VaultItem `json:"-"`
	vaultSecretRefRegex *regexp.Regexp        `json:"-"`
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
func New(vaultFile, name, k8sNamespace string, k8SecretContent []byte, hash, quantumSafe bool, publicKeys ...*crypto.PublicKey) (vlt *Vault, err error) {
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
	if k8SecretContent != nil {
		err = vlt.Update(name, k8sNamespace, k8SecretContent)
	} else {
		err = vlt.commit()
	}
	return
}

// Returns the vault instance from a given yaml. The vault file name must end with .slv.yaml or .slv.yml.
func Get(vaultFile string) (vlt *Vault, err error) {
	if !isValidVaultFileName(vaultFile) {
		return nil, errInvalidVaultFileName
	}
	if !commons.FileExists(vaultFile) {
		return nil, errVaultNotFound
	}
	obj := make(map[string]any)
	if err := commons.ReadFromYAML(vaultFile, &obj); err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return get(jsonData, vaultFile, obj[k8sVaultSpecField] != nil)
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
	vlt.Spec.path = filePath
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

func ListVaultFiles() ([]string, error) {
	var vaultFiles []string
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	err = filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), "."+vaultFileNameRawExt) ||
			strings.HasSuffix(d.Name(), vaultFileNameRawExt+".yaml") ||
			strings.HasSuffix(d.Name(), vaultFileNameRawExt+".yml") {
			if relPath, err := filepath.Rel(wd, path); err == nil {
				vaultFiles = append(vaultFiles, relPath)
			} else {
				vaultFiles = append(vaultFiles, path)
			}
		}
		return nil
	})
	return vaultFiles, err
}

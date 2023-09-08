package vaults

import (
	"bytes"
	"os"
	"path"
	"strings"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"gopkg.in/yaml.v3"
)

type secrets struct {
	Direct     map[string]*crypto.SealedSecret `yaml:"direct,omitempty"`
	Referenced map[string]*crypto.SealedSecret `yaml:"referenced,omitempty"`
}

type config struct {
	Version    string               `yaml:"version,omitempty"`
	PublicKey  *crypto.PublicKey    `yaml:"publicKey"`
	HashLength *uint32              `yaml:"hashLength,omitempty"`
	KeyWraps   []*crypto.WrappedKey `yaml:"wrappedKeys"`
}

type vault struct {
	Secrets secrets `yaml:"secrets,omitempty"`
	Config  config  `yaml:"config,omitempty"`
}

type Vault struct {
	*vault
	path       string
	secretKey  *crypto.SecretKey
	unlockedBy *string
}

func (vlt Vault) MarshalYAML() (interface{}, error) {
	return vlt.vault, nil
}

func (vlt *Vault) UnmarshalYAML(value *yaml.Node) (err error) {
	return value.Decode(&vlt.vault)
}

// Returns new vault instance. The vault name should end with .vlt.slv
func New(vaultFile string, hashLength uint32, publicKeys ...crypto.PublicKey) (vlt *Vault, err error) {
	if !strings.HasSuffix(vaultFile, vaultFileExtension) {
		return nil, ErrInvalidVaultFileName
	}
	if commons.FileExists(vaultFile) {
		return nil, ErrVaultExists
	}
	if os.MkdirAll(path.Dir(vaultFile), os.FileMode(0755)) != nil {
		return nil, ErrVaultDirPathCreation
	}
	vaultPublicKey, vaultSecretKey, err := crypto.NewKeyPair(VaultKey)
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
				Version:    commons.Version,
				PublicKey:  vaultPublicKey,
				HashLength: hashLen,
			},
		},
		path:      vaultFile,
		secretKey: vaultSecretKey,
	}
	for _, pubKey := range publicKeys {
		vlt.ShareAccessToKey(pubKey)
	}
	vlt.commit()
	return vlt, nil
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
		return nil, ErrReadingVault
	}
	return vlt, nil
}

func (vlt *Vault) IsLocked() bool {
	return vlt.secretKey == nil
}

func (vlt *Vault) Lock() {
	vlt.secretKey = nil
}

func (vlt *Vault) UnlockedBy() (id *string) {
	return vlt.unlockedBy
}

func (vlt *Vault) Unlock(secretKey crypto.SecretKey) (err error) {
	if err != nil || (!vlt.IsLocked() && *vlt.unlockedBy == secretKey.PublicKey.String()) {
		return
	}
	for _, secretKeyWrapping := range vlt.vault.Config.KeyWraps {
		decryptedKey, err := secretKey.DecryptKey(*secretKeyWrapping)
		if err == nil {
			vlt.secretKey = decryptedKey
			vlt.unlockedBy = new(string)
			*vlt.unlockedBy = secretKey.PublicKey.String()
			return nil
		}
	}
	return ErrVaultNotAccessible
}

func (vlt *Vault) commit() error {
	return commons.WriteToYAML(vlt.path, *vlt.vault)
}

func (vlt *Vault) GetVersion() string {
	return vlt.vault.Config.Version
}

func (vlt *Vault) share(targetPublicKey crypto.PublicKey, checkForAccess bool) (err error) {
	if vlt.IsLocked() {
		err = ErrVaultLocked
		return
	}
	if targetPublicKey.Type() == VaultKey {
		return ErrVaultCannotBeSharedWithVault
	}
	if checkForAccess {
		for _, keyWrappings := range vlt.vault.Config.KeyWraps {
			if bytes.Equal(keyWrappings.GetKeyId()[:], targetPublicKey.Id()[:]) {
				return ErrVaultAlreadySharedWithKey
			}
		}
	}
	encryptedKey, err := targetPublicKey.EncryptKey(*vlt.secretKey)
	if err == nil {
		vlt.vault.Config.KeyWraps = append(vlt.vault.Config.KeyWraps, encryptedKey)
		err = vlt.commit()
	}
	return
}

func (vlt *Vault) ShareAccessToKey(envPublicKey crypto.PublicKey) (err error) {
	return vlt.share(envPublicKey, true)
}

func (vlt *Vault) ForceShareAccessToKey(envPublicKey crypto.PublicKey) (err error) {
	return vlt.share(envPublicKey, false)
}

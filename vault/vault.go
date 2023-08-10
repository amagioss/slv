package vault

import (
	"os"
	"path"
	"strings"

	"github.com/shibme/slv/commons"
	"github.com/shibme/slv/crypto"
)

type Secrets struct {
	Direct map[string]crypto.SealedData `yaml:"direct"`
}

type Meta struct {
	Version   string             `yaml:"svl_version"`
	PublicKey crypto.PublicKey   `yaml:"public_key"`
	KeyWraps  []crypto.SealedKey `yaml:"sealed_keys"`
}

type VaultData struct {
	Secrets Secrets `yaml:"secrets"`
	Meta    Meta    `yaml:"meta"`
}

type Vault struct {
	path           string
	pendingChanges bool
	data           *VaultData
	encrypter      *crypto.Encrypter
	privateKey     *crypto.PrivateKey
	unlockedBy     *string
}

// Returns new vault instance. The vault name should end with .slv
func New(vaultFile string, shareWithPublicKeys ...crypto.PublicKey) (vault *Vault, err error) {
	if !strings.HasSuffix(vaultFile, vaultFileExtension) {
		return nil, ErrInvalidVaultFileName
	}
	if commons.FileExists(vaultFile) {
		return nil, ErrVaultExists
	}
	vault = &Vault{
		path:           vaultFile,
		pendingChanges: true,
	}
	if os.MkdirAll(path.Dir(vaultFile), os.FileMode(0755)) != nil {
		return nil, ErrVaultDirPathCreation
	}
	var vaultKeyPair *crypto.KeyPair
	vaultKeyPair, err = crypto.NewKeyPair(VaultKey)
	if err != nil {
		return nil, err
	}
	vault.privateKey = new(crypto.PrivateKey)
	*vault.privateKey = vaultKeyPair.PrivateKey()
	vault.data = &VaultData{
		Secrets: Secrets{
			Direct: make(map[string]crypto.SealedData),
		},
		Meta: Meta{
			Version:   commons.Version,
			PublicKey: vaultKeyPair.PublicKey(),
		},
	}
	for _, pubKey := range shareWithPublicKeys {
		vault.ShareAccessToKey(pubKey)
	}
	vault.Lock()
	vault.Commit()
	return vault, nil
}

// Returns the vault instance for a given vault file. The vault name should end with .slv
func Read(vaultFile string) (vault *Vault, err error) {
	if !strings.HasSuffix(vaultFile, vaultFileExtension) {
		return nil, ErrInvalidVaultFileName
	}
	if !commons.FileExists(vaultFile) {
		return nil, ErrVaultNotFound
	}
	vault = &Vault{
		path:           vaultFile,
		pendingChanges: false,
	}
	if err = commons.ReadFromYAML(vault.path, &vault.data); err != nil {
		return nil, ErrReadingVault
	}
	if err = vault.initEncrypter(); err != nil {
		return nil, err
	}
	return vault, nil
}

func (vault *Vault) initEncrypter() (err error) {
	if vault.encrypter == nil {
		vault.encrypter, err = vault.data.Meta.PublicKey.GetEncrypter()
	}
	return
}

func (vault *Vault) IsLocked() bool {
	return vault.privateKey == nil
}

func (vault *Vault) Lock() {
	vault.privateKey = nil
}

func (vault *Vault) UnlockedBy() (id *string) {
	return vault.unlockedBy
}

func (vault *Vault) Unlock(privateKey crypto.PrivateKey) (err error) {
	if err != nil || (!vault.IsLocked() && *vault.unlockedBy == privateKey.Id()) {
		return
	}
	envDecrypter := privateKey.GetDecrypter()
	for _, privateKeyWrapping := range vault.data.Meta.KeyWraps {
		decryptedKey, err := envDecrypter.DecryptKey(privateKeyWrapping)
		if err == nil {
			vault.privateKey = &decryptedKey
			vault.unlockedBy = new(string)
			*vault.unlockedBy = privateKey.Id()
			return nil
		}
	}
	return ErrVaultNotAccessible
}

func (vault *Vault) Commit() {
	commons.WriteToYAML(vault.path, *vault.data)
	vault.pendingChanges = false
}

func (vault *Vault) GetVersion() string {
	return vault.data.Meta.Version
}

func (vault *Vault) share(targetPublicKey crypto.PublicKey, checkForAccess bool) (err error) {
	if vault.IsLocked() {
		err = ErrVaultLocked
		return
	}
	if targetPublicKey.Type() == VaultKey {
		return ErrVaultCannotBeSharedWithVault
	}
	if checkForAccess {
		for _, privateKeyWrapping := range vault.data.Meta.KeyWraps {
			if privateKeyWrapping.GetAccessKeyId() == targetPublicKey.Id() {
				return ErrVaultAlreadySharedWithKey
			}
		}
	}
	encrypter, err := targetPublicKey.GetEncrypter()
	if err != nil {
		return
	}
	encryptedKey, err := encrypter.EncryptKey(*vault.privateKey)
	if err == nil {
		vault.data.Meta.KeyWraps = append(vault.data.Meta.KeyWraps, encryptedKey)
		vault.pendingChanges = true
	}
	return
}

func (vault *Vault) ShareAccessToKey(envPublicKey crypto.PublicKey) (err error) {
	return vault.share(envPublicKey, true)
}

func (vault *Vault) ShareAccessToKeyForced(envPublicKey crypto.PublicKey) (err error) {
	return vault.share(envPublicKey, false)
}

func (vault *Vault) AddDirectSecret(secretName string, secretValue string) (err error) {
	err = vault.initEncrypter()
	if err == nil {
		var cipherData crypto.SealedData
		cipherData, err = vault.encrypter.EncryptString(secretValue)
		if err == nil {
			vault.data.Secrets.Direct[secretName] = cipherData
			vault.pendingChanges = true
		}
	}
	return
}

func (vault *Vault) ListDirectSecretNames() []string {
	names := make([]string, 0, len(vault.data.Secrets.Direct))
	for name := range vault.data.Secrets.Direct {
		names = append(names, name)
	}
	return names
}

func (vault *Vault) GetDirectSecret(secretName string) (secretValue string, err error) {
	if vault.IsLocked() {
		return secretValue, ErrVaultLocked
	}
	encryptedData, ok := vault.data.Secrets.Direct[secretName]
	if !ok {
		return secretValue, ErrVaultSecretNotFound
	}
	decrypter := vault.privateKey.GetDecrypter()
	return decrypter.DecrypToString(encryptedData)
}

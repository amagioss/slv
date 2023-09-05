package vaults

import (
	"crypto/rand"
	"io"
	"os"
	"path"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"gopkg.in/yaml.v3"
)

type secrets struct {
	Direct     map[string]*crypto.SealedData `yaml:"direct,omitempty"`
	Referenced map[string]*crypto.SealedData `yaml:"referenced,omitempty"`
}

type meta struct {
	Version   string              `yaml:"version,omitempty"`
	PublicKey *crypto.PublicKey   `yaml:"public_key,omitempty"`
	KeyWraps  []*crypto.SealedKey `yaml:"sealed_keys,omitempty"`
}

type vault struct {
	Secrets secrets `yaml:"secrets,omitempty"`
	Meta    meta    `yaml:"meta,omitempty"`
}

type Vault struct {
	*vault
	path       string
	encrypter  *crypto.Encrypter
	privateKey *crypto.PrivateKey
	unlockedBy *string
}

func (vlt Vault) MarshalYAML() (interface{}, error) {
	return vlt.vault, nil
}

func (vlt *Vault) UnmarshalYAML(value *yaml.Node) (err error) {
	return value.Decode(&vlt.vault)
}

// Returns new vault instance. The vault name should end with .vlt.slv
func New(vaultFile string, publicKeys ...crypto.PublicKey) (vlt *Vault, err error) {
	if !strings.HasSuffix(vaultFile, vaultFileExtension) {
		return nil, ErrInvalidVaultFileName
	}
	if commons.FileExists(vaultFile) {
		return nil, ErrVaultExists
	}
	vlt = &Vault{
		path: vaultFile,
	}
	if os.MkdirAll(path.Dir(vaultFile), os.FileMode(0755)) != nil {
		return nil, ErrVaultDirPathCreation
	}
	var vaultKeyPair *crypto.KeyPair
	vaultKeyPair, err = crypto.NewKeyPair(VaultKey)
	if err != nil {
		return nil, err
	}
	vlt.privateKey = new(crypto.PrivateKey)
	*vlt.privateKey = vaultKeyPair.PrivateKey()
	vaultPublicKey := vaultKeyPair.PublicKey()
	vlt.vault = &vault{
		Secrets: secrets{
			Direct: make(map[string]*crypto.SealedData),
		},
		Meta: meta{
			Version:   commons.Version,
			PublicKey: &vaultPublicKey,
		},
	}
	for _, pubKey := range publicKeys {
		vlt.ShareAccessToKey(pubKey)
	}
	vlt.Lock()
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
	if err = vlt.initEncrypter(); err != nil {
		return nil, err
	}
	return vlt, nil
}

func (vlt *Vault) initEncrypter() (err error) {
	if vlt.encrypter == nil {
		if vlt.Meta.PublicKey == nil {
			return ErrMissingVaultPublicKey
		}
		vlt.encrypter, err = vlt.vault.Meta.PublicKey.GetEncrypter()
	}
	return
}

func (vlt *Vault) IsLocked() bool {
	return vlt.privateKey == nil
}

func (vlt *Vault) Lock() {
	vlt.privateKey = nil
}

func (vlt *Vault) UnlockedBy() (id *string) {
	return vlt.unlockedBy
}

func (vlt *Vault) Unlock(privateKey crypto.PrivateKey) (err error) {
	if err != nil || (!vlt.IsLocked() && *vlt.unlockedBy == privateKey.Id()) {
		return
	}
	envDecrypter := privateKey.GetDecrypter()
	for _, privateKeyWrapping := range vlt.vault.Meta.KeyWraps {
		decryptedKey, err := envDecrypter.DecryptKey(*privateKeyWrapping)
		if err == nil {
			vlt.privateKey = &decryptedKey
			vlt.unlockedBy = new(string)
			*vlt.unlockedBy = privateKey.Id()
			return nil
		}
	}
	return ErrVaultNotAccessible
}

func (vlt *Vault) commit() error {
	return commons.WriteToYAML(vlt.path, *vlt.vault)
}

func (vlt *Vault) GetVersion() string {
	return vlt.vault.Meta.Version
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
		for _, privateKeyWrapping := range vlt.vault.Meta.KeyWraps {
			if privateKeyWrapping.GetAccessKeyId() == targetPublicKey.Id() {
				return ErrVaultAlreadySharedWithKey
			}
		}
	}
	encrypter, err := targetPublicKey.GetEncrypter()
	if err != nil {
		return
	}
	encryptedKey, err := encrypter.EncryptKey(*vlt.privateKey)
	if err == nil {
		vlt.vault.Meta.KeyWraps = append(vlt.vault.Meta.KeyWraps, &encryptedKey)
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

func (vlt *Vault) AddDirectSecret(secretName string, secretValue string) (err error) {
	err = vlt.initEncrypter()
	if err == nil {
		var cipherData crypto.SealedData
		cipherData, err = vlt.encrypter.EncryptString(secretValue)
		if err == nil {
			if vlt.vault.Secrets.Direct == nil {
				vlt.vault.Secrets.Direct = make(map[string]*crypto.SealedData)
			}
			vlt.vault.Secrets.Direct[secretName] = &cipherData
			err = vlt.commit()
		}
	}
	return
}

func (vlt *Vault) ListDirectSecretNames() []string {
	names := make([]string, 0, len(vlt.vault.Secrets.Direct))
	for name := range vlt.vault.Secrets.Direct {
		names = append(names, name)
	}
	return names
}

func (vlt *Vault) GetDirectSecret(secretName string) (secretValue string, err error) {
	if vlt.IsLocked() {
		return secretValue, ErrVaultLocked
	}
	encryptedData, ok := vlt.vault.Secrets.Direct[secretName]
	if !ok {
		return "", ErrVaultSecretNotFound
	}
	decrypter := vlt.privateKey.GetDecrypter()
	return decrypter.DecrypToString(*encryptedData)
}

func (vlt *Vault) DeleteDirecetSecret(secretName string) error {
	delete(vlt.vault.Secrets.Direct, secretName)
	return vlt.commit()
}

func randomStr(bytecount uint8) (string, error) {
	randBytes := make([]byte, bytecount)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}
	randomString := base58.Encode(randBytes)
	return secretRefPrefix + randomString, nil
}

func (vlt *Vault) addReferencedSecret(secretValue string) (secretReference string, err error) {
	err = vlt.initEncrypter()
	if err == nil {
		var cipherData crypto.SealedData
		cipherData, err = vlt.encrypter.EncryptString(secretValue)
		if err == nil {
			if vlt.vault.Secrets.Referenced == nil {
				vlt.vault.Secrets.Referenced = make(map[string]*crypto.SealedData)
			}
			secretReference, err = randomStr(16)
			attempts := 0
			for err == nil && vlt.Secrets.Referenced[secretReference] != nil && attempts < maxRefNameAttempts {
				secretReference, err = randomStr(16)
				attempts++
			}
			if err == nil && attempts >= maxRefNameAttempts {
				err = ErrMaximumReferenceAttemptsReached
			}
			if err == nil {
				vlt.vault.Secrets.Referenced[secretReference] = &cipherData
				err = vlt.commit()
			}
		}
	}
	return
}

func (vlt *Vault) getReferencedSecret(secretReference string) (secretValue string, err error) {
	if vlt.IsLocked() {
		return secretValue, ErrVaultLocked
	}
	encryptedData, ok := vlt.vault.Secrets.Referenced[secretReference]
	if !ok {
		return "", ErrVaultSecretNotFound
	}
	decrypter := vlt.privateKey.GetDecrypter()
	return decrypter.DecrypToString(*encryptedData)
}

func (vlt *Vault) deleteReferencedSecret(secretReference string) {
	delete(vlt.vault.Secrets.Referenced, secretReference)
}

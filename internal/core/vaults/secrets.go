package vaults

import (
	"fmt"

	"gopkg.in/yaml.v3"
	"oss.amagi.com/slv/internal/core/crypto"
)

type vaultDataValue struct {
	SealedSecret *crypto.SealedSecret
	PlainText    []byte
}

func (vlt *Vault) putWithoutCommit(secretName string, secretValue []byte) (err error) {
	if !secretNameRegex.MatchString(secretName) {
		return errInvalidSecretName
	}
	var sealedSecret *crypto.SealedSecret
	vaultPublicKey, err := vlt.getPublicKey()
	if err != nil {
		return err
	}
	sealedSecret, err = vaultPublicKey.EncryptSecret(secretValue, vlt.Config.Hash)
	if err == nil {
		if vlt.Data == nil {
			vlt.Data = make(map[string]string)
		}
		vlt.Data[secretName] = sealedSecret.String()
	}
	return
}

func (vlt *Vault) Put(secretName string, secretValue []byte) (err error) {
	if err = vlt.putWithoutCommit(secretName, secretValue); err == nil {
		err = vlt.commit()
	}
	return
}

func (vlt *Vault) Import(importData []byte, force bool) (err error) {
	secretsMap := make(map[string]string)
	if err = yaml.Unmarshal(importData, &secretsMap); err != nil {
		return errInvalidImportDataFormat
	}
	if !force {
		for secretName := range secretsMap {
			if vlt.Exists(secretName) {
				return fmt.Errorf("secret %s already exists", secretName)
			}
		}
	}
	for secretName, secretValue := range secretsMap {
		if err = vlt.putWithoutCommit(secretName, []byte(secretValue)); err != nil {
			return err
		}
	}
	return vlt.commit()
}

func (vlt *Vault) Exists(secretName string) (exists bool) {
	if vlt.Data != nil {
		_, exists = vlt.Data[secretName]
	}
	return exists
}

func (vlt *Vault) List() (map[string]vaultDataValue, error) {
	sealedSecretsMap := make(map[string]vaultDataValue)
	for name, value := range vlt.Data {
		sealedSecret := crypto.SealedSecret{}
		if err := sealedSecret.FromString(value); err != nil {
			sealedSecretsMap[name] = vaultDataValue{
				PlainText: []byte(value),
			}
		} else {
			sealedSecretsMap[name] = vaultDataValue{
				SealedSecret: &sealedSecret,
			}
		}
	}
	return sealedSecretsMap, nil
}

func (vlt *Vault) GetAll() (secretsMap map[string][]byte, err error) {
	if vlt.IsLocked() {
		return secretsMap, errVaultLocked
	}
	secretsMap = make(map[string][]byte)
	for secretName := range vlt.Data {
		var decryptedSecret []byte
		if decryptedSecret, err = vlt.Get(secretName); err == nil {
			secretsMap[secretName] = decryptedSecret
		} else {
			return nil, fmt.Errorf("error decrypting secret %s: %w", secretName, err)
		}
	}
	return
}

func (vlt *Vault) Get(name string) (value []byte, err error) {
	if vlt.IsLocked() {
		return nil, errVaultLocked
	}
	rawValue := vlt.Data[name]
	if rawValue == "" {
		return nil, errVaultSecretNotFound
	}
	if value = vlt.getFromCache(name); value == nil {
		sealedSecret := &crypto.SealedSecret{}
		if err = sealedSecret.FromString(rawValue); err == nil {
			if value, err = vlt.secretKey.DecryptSecret(*sealedSecret); err != nil {
				return nil, err
			}
		}
		if value == nil {
			value = []byte(rawValue)
			err = nil
		}
		vlt.putToCache(name, value)
	}
	return
}

func (vlt *Vault) DeleteSecret(secretName string) error {
	return vlt.DeleteSecrets([]string{secretName})
}

func (vlt *Vault) DeleteSecrets(secretNames []string) error {
	for _, secretName := range secretNames {
		delete(vlt.Data, secretName)
		vlt.deleteFromCache(secretName)
	}
	return vlt.commit()
}

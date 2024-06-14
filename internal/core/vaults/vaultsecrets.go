package vaults

import (
	"fmt"

	"gopkg.in/yaml.v3"
	"oss.amagi.com/slv/internal/core/crypto"
)

func (vlt *Vault) putSecretWithoutCommit(secretName string, secretValue []byte) (err error) {
	if !secretNameRegex.MatchString(secretName) {
		return errInvalidSecretName
	}
	var sealedSecret *crypto.SealedSecret
	vaultPublicKey, err := vlt.getPublicKey()
	if err != nil {
		return err
	}
	sealedSecret, err = vaultPublicKey.EncryptSecret(secretValue, vlt.Config.HashLength)
	if err == nil {
		if vlt.Secrets == nil {
			vlt.Secrets = make(map[string]string)
		}
		vlt.Secrets[secretName] = sealedSecret.String()
	}
	return
}

func (vlt *Vault) PutSecret(secretName string, secretValue []byte) (err error) {
	if err = vlt.putSecretWithoutCommit(secretName, secretValue); err == nil {
		err = vlt.commit()
	}
	return
}

func (vlt *Vault) ImportSecrets(importData []byte, force bool) (err error) {
	secretsMap := make(map[string]string)
	if err = yaml.Unmarshal(importData, &secretsMap); err != nil {
		return errInvalidImportDataFormat
	}
	if !force {
		for secretName := range secretsMap {
			if vlt.SecretExists(secretName) {
				return fmt.Errorf("secret %s already exists", secretName)
			}
		}
	}
	for secretName, secretValue := range secretsMap {
		if err = vlt.putSecretWithoutCommit(secretName, []byte(secretValue)); err != nil {
			return err
		}
	}
	return vlt.commit()
}

func (vlt *Vault) SecretExists(secretName string) (exists bool) {
	if vlt.Secrets != nil {
		_, exists = vlt.Secrets[secretName]
	}
	return exists
}

func (vlt *Vault) ListSealedSecrets() (map[string]crypto.SealedSecret, error) {
	sealedSecretsMap := make(map[string]crypto.SealedSecret)
	for name, value := range vlt.Secrets {
		sealedSecret := crypto.SealedSecret{}
		if err := sealedSecret.FromString(value); err != nil {
			return nil, err
		}
		sealedSecretsMap[name] = sealedSecret
	}
	return sealedSecretsMap, nil
}

func (vlt *Vault) GetAllSecrets() (secretsMap map[string][]byte, err error) {
	if vlt.IsLocked() {
		return secretsMap, errVaultLocked
	}
	secretsMap = make(map[string][]byte)
	for secretName := range vlt.Secrets {
		var decryptedSecret []byte
		if decryptedSecret, err = vlt.GetSecret(secretName); err == nil {
			secretsMap[secretName] = decryptedSecret
		} else {
			return nil, err
		}
	}
	return
}

func (vlt *Vault) GetSecret(secretName string) (decryptedSecret []byte, err error) {
	if vlt.IsLocked() {
		return decryptedSecret, errVaultLocked
	}
	sealedSecretData := vlt.Secrets[secretName]
	if sealedSecretData == "" {
		return nil, errVaultSecretNotFound
	}
	if decryptedSecret = vlt.getSecretFromCache(secretName); decryptedSecret == nil {
		sealedSecret := &crypto.SealedSecret{}
		if err = sealedSecret.FromString(sealedSecretData); err == nil {
			if decryptedSecret, err = vlt.secretKey.DecryptSecret(*sealedSecret); err == nil {
				vlt.putSecretToCache(secretName, decryptedSecret)
			}
		}
	}
	return
}

func (vlt *Vault) DeleteSecret(secretName string) error {
	return vlt.DeleteSecrets([]string{secretName})
}

func (vlt *Vault) DeleteSecrets(secretNames []string) error {
	for _, secretName := range secretNames {
		delete(vlt.Secrets, secretName)
		vlt.deleteSecretFromCache(secretName)
	}
	return vlt.commit()
}

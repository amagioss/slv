package vaults

import "github.com/shibme/slv/core/crypto"

func (vlt *Vault) AddDirectSecret(secretName string, secretValue string) (err error) {
	var sealedSecret *crypto.SealedSecret
	sealedSecret, err = vlt.Config.PublicKey.EncryptSecretString(secretValue, vlt.Config.HashLength)
	if err == nil {
		if vlt.vault.Secrets.Direct == nil {
			vlt.vault.Secrets.Direct = make(map[string]*crypto.SealedSecret)
		}
		vlt.vault.Secrets.Direct[secretName] = sealedSecret
		err = vlt.commit()
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
	sealedSecret, ok := vlt.vault.Secrets.Direct[secretName]
	if !ok {
		return "", ErrVaultSecretNotFound
	}
	return vlt.secretKey.DecryptSecretString(*sealedSecret)
}

func (vlt *Vault) DeleteDirecetSecret(secretName string) error {
	delete(vlt.vault.Secrets.Direct, secretName)
	return vlt.commit()
}

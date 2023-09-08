package vaults

import (
	"crypto/rand"
	"io"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
)

func randomStr(bytecount uint8) (string, error) {
	randBytes := make([]byte, bytecount)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}
	return commons.Encode(randBytes), nil
}

func (vlt *Vault) getReferenceName(nodePath string) (string, error) {
	// return autoReferencedPrefix + vlt.Config.PublicKey.String() + "_" + commons.Encode([]byte(nodePath))
	randomStr, err := randomStr(autoReferenceLength)
	if err != nil {
		return "", err
	}
	return autoReferencedPrefix + vlt.Config.PublicKey.Id() + "_" + randomStr, nil
}

func (vlt *Vault) addReferencedSecret(nodePath, secret string) (secretRef string, err error) {
	var sealedSecret *crypto.SealedSecret
	sealedSecret, err = vlt.Config.PublicKey.EncryptSecretString(secret, vlt.Config.HashLength)
	if err == nil {
		if vlt.vault.Secrets.Referenced == nil {
			vlt.vault.Secrets.Referenced = make(map[string]*crypto.SealedSecret)
		}
		secretRef, err = vlt.getReferenceName(nodePath)
		attempts := 0
		for err == nil && vlt.Secrets.Referenced[secretRef] != nil && attempts < maxRefNameAttempts {
			secretRef, err = randomStr(autoReferenceLength)
			attempts++
		}
		if err == nil && attempts >= maxRefNameAttempts {
			err = ErrMaximumReferenceAttemptsReached
		}
		if err == nil {
			vlt.vault.Secrets.Referenced[secretRef] = sealedSecret
			err = vlt.commit()
		}
	}
	return
}

func (vlt *Vault) getReferencedSecret(secretRef string) (secret string, err error) {
	if vlt.IsLocked() {
		return secret, ErrVaultLocked
	}
	encryptedData, ok := vlt.vault.Secrets.Referenced[secretRef]
	if !ok {
		return "", ErrVaultSecretNotFound
	}
	return vlt.secretKey.DecryptSecretString(*encryptedData)
}

func (vlt *Vault) deleteReferencedSecret(secretReference string) {
	delete(vlt.vault.Secrets.Referenced, secretReference)
}

package crypto

import (
	"bytes"

	"github.com/shibme/slv/core/commons"
	"golang.org/x/crypto/nacl/box"
)

func (secretKey *SecretKey) decrypt(ciphertext []byte) (data []byte, err error) {
	ephemeralPublicKey := [keyLength]byte(ciphertext[:keyLength])
	ciphertext = ciphertext[keyLength:]
	nonce := [nonceLength]byte(ciphertext[:nonceLength])
	encryptedBytes := ciphertext[nonceLength:]
	decryptedData, success := box.Open(nil, encryptedBytes, &nonce, &ephemeralPublicKey, secretKey.key)
	if !success {
		return nil, ErrDecryptionFailed
	}
	return commons.Decompress(decryptedData)
}

func (secretKey *SecretKey) DecryptSecret(sealedSecret SealedSecret) (secret []byte, err error) {
	if !bytes.Equal(sealedSecret.keyId[:], secretKey.ShortId()[:]) {
		return nil, ErrSecretKeyMismatch
	}
	if *sealedSecret.keyType != *secretKey.keyType {
		return nil, ErrSecretKeyTypeMismatch
	}
	return secretKey.decrypt(*sealedSecret.ciphertext)
}

func (secretKey *SecretKey) DecryptSecretString(sealedSecret SealedSecret) (string, error) {
	data, err := secretKey.DecryptSecret(sealedSecret)
	return string(data), err
}

func (secretKey *SecretKey) DecryptKey(wrappedKey WrappedKey) (*SecretKey, error) {
	if !bytes.Equal(wrappedKey.keyId[:], secretKey.ShortId()[:]) {
		return nil, ErrSecretKeyMismatch
	}
	secretKeyBytes, err := secretKey.decrypt(*wrappedKey.ciphertext)
	if err == nil {
		var sKey SecretKey
		if err = sKey.fromBytes(secretKeyBytes); err == nil {
			return &sKey, nil
		}
	}
	return nil, err
}

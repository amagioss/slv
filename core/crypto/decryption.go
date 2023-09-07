package crypto

import (
	"bytes"

	"github.com/shibme/slv/core/commons"
	"golang.org/x/crypto/nacl/box"
)

func (secretKey *SecretKey) decrypt(ciphertext []byte) (data []byte, err error) {
	ephemeralPublicKey := [32]byte(ciphertext[0:32])
	nonce := [24]byte(ciphertext[32:56])
	encryptedBytes := ciphertext[56:]
	decryptedData, success := box.Open(nil, encryptedBytes, &nonce, &ephemeralPublicKey, secretKey.key)
	if !success {
		return nil, ErrDecryptionFailed
	}
	return commons.Decompress(decryptedData)
}

func (secretKey *SecretKey) DecryptSecret(sealedSecret SealedSecret) (secret []byte, err error) {
	if !bytes.Equal(sealedSecret.keyId[:], secretKey.Id()[:]) {
		return nil, ErrSecretKeyMismatch
	}
	if sealedSecret.keyType != secretKey.keyType {
		return nil, ErrSecretKeyMismatch
	}
	return secretKey.decrypt(*sealedSecret.ciphertext)
}

func (secretKey *SecretKey) DecryptSecretString(sealedSecret SealedSecret) (string, error) {
	data, err := secretKey.DecryptSecret(sealedSecret)
	return string(data), err
}

func (secretKey *SecretKey) DecryptKey(wrappedKey WrappedKey) (*SecretKey, error) {
	if !bytes.Equal(wrappedKey.keyId[:], secretKey.Id()[:]) {
		return nil, ErrSecretKeyMismatch
	}
	secretKeyBytes, err := secretKey.decrypt(*wrappedKey.ciphertext)
	if err != nil {
		return nil, err
	}
	return secretKeyFromBytes(secretKeyBytes)
}

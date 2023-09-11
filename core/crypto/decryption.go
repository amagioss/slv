package crypto

import (
	"bytes"

	"github.com/shibme/slv/core/commons"
	"golang.org/x/crypto/nacl/box"
)

func (secretKey *SecretKey) decrypt(ciphered *ciphered) (data []byte, err error) {
	ciphertext := *ciphered.ciphertext
	ephemeralPublicKey := [keyLength]byte(ciphertext[:keyLength])
	ciphertext = ciphertext[keyLength:]
	nonce := [nonceLength]byte(ciphertext[:nonceLength])
	encryptedBytes := ciphertext[nonceLength:]
	if *ciphered.keyType != *secretKey.keyType {
		return nil, ErrKeyTypeMismatch
	}
	if !bytes.Equal(ciphered.shortKeyId[:], secretKey.ShortId()[:]) {
		return nil, ErrSecretKeyMismatch
	}
	if ciphered.version == nil || *ciphered.version != *secretKey.version {
		return nil, ErrKeyTypeMismatch
	}
	decryptedData, success := box.Open(nil, encryptedBytes, &nonce, &ephemeralPublicKey, secretKey.key)
	if !success {
		return nil, ErrDecryptionFailed
	}
	return commons.Decompress(decryptedData)
}

func (secretKey *SecretKey) DecryptSecret(sealedSecret SealedSecret) (secret []byte, err error) {
	return secretKey.decrypt(sealedSecret.ciphered)
}

func (secretKey *SecretKey) DecryptKey(wrappedKey WrappedKey) (*SecretKey, error) {
	secretKeyBytes, err := secretKey.decrypt(wrappedKey.ciphered)
	if err == nil {
		var sKey SecretKey
		if err = sKey.fromBytes(secretKeyBytes); err == nil {
			return &sKey, nil
		}
	}
	return nil, err
}

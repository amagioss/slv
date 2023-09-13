package crypto

import (
	"bytes"

	"github.com/shibme/slv/core/commons"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

func (secretKey *SecretKey) decrypt(ciphered *ciphered) (data []byte, err error) {
	encryptedBytes := *ciphered.ciphertext
	ephemeralPubKey := encryptedBytes[:curve25519.ScalarSize]
	encryptedBytes = encryptedBytes[curve25519.ScalarSize:]
	nonce := encryptedBytes[:chacha20poly1305.NonceSize]
	ciphertext := encryptedBytes[chacha20poly1305.NonceSize:]
	if *ciphered.keyType != *secretKey.keyType {
		return nil, ErrKeyTypeMismatch
	}
	if !bytes.Equal(*ciphered.pubKeyBytes, secretKey.Id()) {
		return nil, ErrSecretKeyMismatch
	}
	if ciphered.version == nil || *ciphered.version != *secretKey.version {
		return nil, ErrUnsupportedCryptoVersion
	}
	sharedKey, err := curve25519.X25519(*secretKey.key, ephemeralPubKey)
	if err != nil {
		return nil, ErrGeneratingKey
	}
	aead, err := chacha20poly1305.New(sharedKey)
	if err != nil {
		return nil, ErrGeneratingKey
	}
	decryptedData, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
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
		var keyBase *keyBase
		if keyBase, err = keyBaseFromBytes(secretKeyBytes); err == nil {
			return &SecretKey{
				keyBase: keyBase,
			}, nil
		}
	}
	return nil, err
}

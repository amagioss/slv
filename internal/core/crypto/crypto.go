package crypto

import (
	"time"

	"golang.org/x/crypto/argon2"
)

func (publicKey *PublicKey) encrypt(data []byte) (*ciphered, error) {
	ciphertext, err := publicKey.pubKey.Encrypt(data, true)
	if err != nil {
		return nil, errEncryptionFailed
	}
	pubKeyBytes, err := publicKey.toBytes()
	if err != nil {
		return nil, err
	}
	currentTime := time.Now()
	return &ciphered{
		version:     publicKey.version,
		keyType:     publicKey.keyType,
		encryptedAt: &currentTime,
		encryptedBy: pubKeyBytes,
		ciphertext:  ciphertext,
	}, nil
}

func (publicKey *PublicKey) EncryptKey(secretKey SecretKey) (wrappedKey *WrappedKey, err error) {
	ciphered, err := publicKey.encrypt(secretKey.Bytes())
	if err == nil {
		wrappedKey = &WrappedKey{
			ciphered: ciphered,
		}
	}
	return
}

func hash(data []byte, length uint8) []byte {
	return argon2.IDKey(data, nil, 16, 64, 1, uint32(length))
}

func (publicKey *PublicKey) EncryptSecret(secret []byte, hashLength uint8) (sealedSecret *SealedSecret, err error) {
	ciphered, err := publicKey.encrypt(secret)
	if err == nil {
		sealedSecret = &SealedSecret{
			ciphered: ciphered,
		}
		if hashLength > 0 {
			if hashLength > hashMaxLength {
				hashLength = hashMaxLength
			}
			hash := hash(secret, hashLength)
			sealedSecret.hash = &hash
		}
	}
	return
}

func (secretKey *SecretKey) decrypt(ciphered *ciphered) (data []byte, err error) {
	var pqPubKey, eccPubKey *PublicKey
	pqPubKey, err = secretKey.PublicKey(true)
	if err == nil {
		eccPubKey, err = secretKey.PublicKey(false)
	}
	if err != nil || (!ciphered.IsEncryptedBy(pqPubKey) && !ciphered.IsEncryptedBy(eccPubKey)) {
		return nil, errSecretKeyMismatch
	}
	data, err = secretKey.privKey.Decrypt(ciphered.ciphertext)
	if err != nil {
		return nil, errDecryptionFailed
	}
	return
}

func (secretKey *SecretKey) DecryptSecret(sealedSecret SealedSecret) (secret []byte, err error) {
	return secretKey.decrypt(sealedSecret.ciphered)
}

func (secretKey *SecretKey) DecryptKey(wrappedKey WrappedKey) (*SecretKey, error) {
	decryptedSecretKeyBytes, err := secretKey.decrypt(wrappedKey.ciphered)
	if err != nil {
		return nil, err
	}
	return SecretKeyFromBytes(decryptedSecretKeyBytes)
}

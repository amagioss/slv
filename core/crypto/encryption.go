package crypto

import (
	"crypto/rand"

	"github.com/shibme/slv/core/commons"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

func (publicKey *PublicKey) getEncrypter() (*encrypter, error) {
	if publicKey.encrypter == nil {
		ephPrivKey := make([]byte, curve25519.ScalarSize)
		if _, err := rand.Read(ephPrivKey); err != nil {
			return nil, ErrGeneratingKey
		}
		ephPubKey, err := curve25519.X25519(ephPrivKey, curve25519.Basepoint)
		if err != nil {
			return nil, ErrGeneratingKey
		}
		sharedKey, err := curve25519.X25519(ephPrivKey, *publicKey.key)
		if err != nil {
			return nil, ErrGeneratingKey
		}
		aead, err := chacha20poly1305.New(sharedKey)
		if err != nil {
			return nil, ErrGeneratingKey
		}
		publicKey.encrypter = &encrypter{
			ephpublicKey: &ephPubKey,
			aead:         &aead,
		}
	}
	return publicKey.encrypter, nil
}

func (publicKey *PublicKey) encrypt(data []byte) (ciphd *ciphered, err error) {
	if publicKey.encrypter == nil {
		if publicKey.encrypter, err = publicKey.getEncrypter(); err != nil {
			return nil, err
		}
	}
	compressedData, err := commons.Compress(data)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	aead := *publicKey.encrypter.aead
	encrypted := aead.Seal(nil, nonce, compressedData, nil)
	if err == nil {
		ciphertext := append(*publicKey.encrypter.ephpublicKey, nonce...)
		ciphertext = append(ciphertext, encrypted...)
		ciphd = &ciphered{
			version:    publicKey.version,
			keyType:    publicKey.keyType,
			ciphertext: &ciphertext,
			shortKeyId: publicKey.ShortId(),
		}
	}
	return
}

func (publicKey *PublicKey) getHashForSecret(secret []byte, hashLength uint32) []byte {
	return argon2.IDKey(secret, nil, argon2Iterations, argon2Memory, argon2Threads, hashLength)
}

func (publicKey *PublicKey) EncryptKey(secretKey SecretKey) (wrappedKey *WrappedKey, err error) {
	ciphered, err := publicKey.encrypt(secretKey.toBytes())
	if err == nil {
		wrappedKey = &WrappedKey{
			ciphered: ciphered,
		}
	}
	return
}

func (publicKey *PublicKey) EncryptSecret(secret []byte, hashLength *uint32) (sealedSecret *SealedSecret, err error) {
	ciphered, err := publicKey.encrypt(secret)
	if err == nil {
		sealedSecret = &SealedSecret{
			ciphered: ciphered,
		}
		if hashLength != nil && *hashLength > 0 {
			if *hashLength > argon2HashMaxLength {
				*hashLength = argon2HashMaxLength
			}
			hash := publicKey.getHashForSecret(secret, *hashLength)
			sealedSecret.hash = &hash
		}
	}
	return
}

func (publicKey *PublicKey) EncryptSecretString(str string, hashLength *uint32) (sealedSecret *SealedSecret, err error) {
	return publicKey.EncryptSecret([]byte(str), hashLength)
}

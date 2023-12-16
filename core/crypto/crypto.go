package crypto

import (
	"github.com/shibme/slv/core/commons"
	"gopkg.shib.me/gociphers/argon2"
)

func (publicKey *PublicKey) encrypt(data []byte) (*ciphered, error) {
	ciphertext, err := publicKey.pubKey.Encrypt(data, true)
	if err != nil {
		return nil, ErrEncryptionFailed
	}
	return &ciphered{
		version:     publicKey.version,
		keyType:     publicKey.keyType,
		pubKeyBytes: commons.ByteSlicePtr(publicKey.toBytes()),
		ciphertext:  &ciphertext,
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
			hash := argon2.Hash(secret, *hashLength)
			sealedSecret.hash = &hash
		}
	}
	return
}

func (secretKey *SecretKey) decrypt(ciphered *ciphered) (data []byte, err error) {
	publicKey, err := secretKey.PublicKey()
	if err != nil || !ciphered.IsEncryptedBy(publicKey) {
		return nil, ErrSecretKeyMismatch
	}
	data, err = secretKey.privKey.Decrypt(*ciphered.ciphertext)
	if err != nil {
		return nil, ErrDecryptionFailed
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

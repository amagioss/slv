package crypto

import (
	"crypto/rand"

	"github.com/shibme/slv/core/commons"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/box"
)

func nonce() (nonce [24]byte, err error) {
	_, err = rand.Read(nonce[:])
	return
}

func (publicKey *PublicKey) getEncrypter() (*encrypter, error) {
	if publicKey.encrypter == nil {
		if ephPubKey, ephPrivKey, err := box.GenerateKey(rand.Reader); err == nil {
			var sharedKey *[32]byte = new([32]byte)
			box.Precompute(sharedKey, publicKey.key, ephPrivKey)
			publicKey.encrypter = &encrypter{
				encryptionKeyId:    publicKey.Id(),
				ephemeralPublicKey: ephPubKey,
				sharedKey:          sharedKey,
			}
		} else {
			return nil, err
		}
	}
	return publicKey.encrypter, nil
}

func (publicKey *PublicKey) encrypt(data []byte) (ciphertext []byte, err error) {
	if publicKey.encrypter == nil {
		if publicKey.encrypter, err = publicKey.getEncrypter(); err != nil {
			return nil, err
		}
	}
	nonce, err := nonce()
	if err == nil {
		compressedData, err := commons.Compress(data)
		if err == nil {
			encrypted := append(publicKey.encrypter.ephemeralPublicKey[:], nonce[:]...)
			encrypted = append(encrypted, box.SealAfterPrecomputation(nil, compressedData, &nonce, publicKey.encrypter.sharedKey)...)
			ciphertext = encrypted
		}
	}
	return
}

func (publicKey *PublicKey) getHashForSecret(secret []byte, hashLength uint32) []byte {
	return argon2.IDKey(secret, nil, secretHashTime, secretHashMemory, secretHashThreads, hashLength)
}

func (publicKey *PublicKey) EncryptSecret(secret []byte, hashLength uint32) (sealedSecret *SealedSecret, err error) {
	ciphertext, err := publicKey.encrypt(secret)
	if err == nil {
		sealedSecret = &SealedSecret{
			ciphered: &ciphered{
				ciphertext: &ciphertext,
				keyId:      publicKey.Id(),
				keyType:    publicKey.keyType,
			},
		}
		if hashLength > 0 {
			if hashLength > secretHashMaxLength {
				hashLength = secretHashMaxLength
			}
			*sealedSecret.hash = publicKey.getHashForSecret(secret, hashLength)
		}
	}
	return
}

func (publicKey *PublicKey) EncryptSecretString(str string, hashLength uint32) (sealedSecret *SealedSecret, err error) {
	return publicKey.EncryptSecret([]byte(str), hashLength)
}

func (publicKey *PublicKey) EncryptKey(secretKey SecretKey) (wrappedKey *WrappedKey, err error) {
	encrypted, err := publicKey.encrypt(secretKey.toBytes())
	if err == nil {
		wrappedKey = &WrappedKey{
			ciphered: &ciphered{
				ciphertext: &encrypted,
				keyId:      publicKey.Id(),
				keyType:    publicKey.keyType,
			},
		}
	}
	return
}

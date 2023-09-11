package crypto

import (
	"crypto/rand"

	"github.com/shibme/slv/core/commons"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/box"
)

func nonce() (nonce [nonceLength]byte, err error) {
	_, err = rand.Read(nonce[:])
	return
}

func (publicKey *PublicKey) getEncrypter() (*encrypter, error) {
	if publicKey.encrypter == nil {
		if ephPubKey, ephPrivKey, err := box.GenerateKey(rand.Reader); err == nil {
			sharedKey := new([keyLength]byte)
			box.Precompute(sharedKey, publicKey.key, ephPrivKey)
			publicKey.encrypter = &encrypter{
				encryptionKeyId:    publicKey.ShortId(),
				ephemeralPublicKey: ephPubKey,
				sharedKey:          sharedKey,
			}
		} else {
			return nil, err
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
	nonce, err := nonce()
	if err == nil {
		compressedData, err := commons.Compress(data)
		if err == nil {
			ciphertext := append(publicKey.encrypter.ephemeralPublicKey[:], nonce[:]...)
			ciphertext = append(ciphertext, box.SealAfterPrecomputation(nil, compressedData, &nonce, publicKey.encrypter.sharedKey)...)
			ciphd = &ciphered{
				version:    publicKey.version,
				keyType:    publicKey.keyType,
				ciphertext: &ciphertext,
				shortKeyId: publicKey.ShortId(),
			}
		}
	}
	return
}

func (publicKey *PublicKey) getHashForSecret(secret []byte, hashLength uint32) []byte {
	return argon2.IDKey(secret, nil, secretHashTime, secretHashMemory, secretHashThreads, hashLength)
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
			if *hashLength > secretHashMaxLength {
				*hashLength = secretHashMaxLength
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

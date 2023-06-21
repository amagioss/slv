package crypto

import (
	"crypto/rand"
	"crypto/sha1"

	"github.com/shibme/slv/slv/commons"
	"golang.org/x/crypto/nacl/box"
)

type Encrypter struct {
	encryptionKeyId    [8]byte
	ephemeralPublicKey [32]byte
	sharedKey          *[32]byte
}

func (encrypter *Encrypter) nonce() (nonce [24]byte, err error) {
	_, err = rand.Read(nonce[0:24])
	if err != nil {
		return
	}
	return
}

func (encrypter *Encrypter) encrypt(data []byte) (encrypted []byte, err error) {
	nonce, err := encrypter.nonce()
	if err == nil {
		compressedData, err := commons.Compress(data)
		if err == nil {
			encrypted = append(encrypter.ephemeralPublicKey[:], nonce[:]...)
			encrypted = append(encrypted, box.SealAfterPrecomputation(nil, compressedData, &nonce, encrypter.sharedKey)...)
		}
	}
	return
}

func (encrypter *Encrypter) Encrypt(data []byte) (sealedData SealedData, err error) {
	cipherData, err := encrypter.encrypt(data)
	if err == nil {
		sumBytes := sha1.Sum(data)
		sealedData.checksum = [8]byte(sumBytes[len(sumBytes)-8:])
		sealedData.encryptionKeyId = encrypter.encryptionKeyId
		sealedData.data = cipherData
	}
	return
}

func (encrypter *Encrypter) EncryptString(str string) (sealedData SealedData, err error) {
	return encrypter.Encrypt([]byte(str))
}

func (encrypter *Encrypter) EncryptKey(privateKey PrivateKey) (sealedKey SealedKey, err error) {
	encrypted, err := encrypter.encrypt(privateKey.toBytes())
	if err == nil {
		sealedKey.data = encrypted
		sealedKey.encryptionKeyId = encrypter.encryptionKeyId
	}
	return
}

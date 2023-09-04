package crypto

import (
	"github.com/shibme/slv/core/commons"
	"golang.org/x/crypto/nacl/box"
)

type Decrypter struct {
	privateKey *PrivateKey
}

func (decrypter *Decrypter) decrypt(encrypted []byte) (data []byte, err error) {
	ephemeralPublicKey := [32]byte(encrypted[0:32])
	nonce := [24]byte(encrypted[32:56])
	encryptedBytes := encrypted[56:]
	decryptedData, success := box.Open(nil, encryptedBytes, &nonce, &ephemeralPublicKey, &decrypter.privateKey.keyData)
	if !success {
		return nil, ErrDecryptionFailed
	}
	return commons.Decompress(decryptedData)
}

func (decrypter *Decrypter) Decrypt(cipheredData SealedData) (decryptedData []byte, err error) {
	if cipheredData.encryptionKeyId != decrypter.privateKey.id {
		return nil, ErrAccessKeyMismatch
	}
	return decrypter.decrypt(cipheredData.data)
}

func (decrypter *Decrypter) DecrypToString(cipheredData SealedData) (decryptedStr string, err error) {
	var decryptedData []byte
	decryptedData, err = decrypter.Decrypt(cipheredData)
	if err == nil {
		decryptedStr = string(decryptedData)
	}
	return
}

func (decrypter *Decrypter) DecryptKey(sealedKey SealedKey) (privateKey PrivateKey, err error) {
	if sealedKey.encryptionKeyId != decrypter.privateKey.id {
		err = ErrAccessKeyMismatch
		return
	}
	var decryptedBytes []byte
	decryptedBytes, err = decrypter.decrypt(sealedKey.data)
	if err == nil {
		var key *key
		key, err = keyFromBytes(decryptedBytes)
		privateKey.key = key
	}
	return
}

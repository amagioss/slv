package crypto

import (
	"crypto/rand"

	"golang.org/x/crypto/nacl/box"
)

func NewKeyPair(keyType KeyType) (publicKey *PublicKey, secretKey *SecretKey, err error) {
	pubKey, privKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, ErrGeneratingKeyPair
	}
	version := cryptoVersion
	publicKey = &PublicKey{
		key:     pubKey,
		keyType: &keyType,
		version: &version,
	}
	secretKey = &SecretKey{
		key:       privKey,
		PublicKey: publicKey,
	}
	return
}

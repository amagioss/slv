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
	publicKey = &PublicKey{
		key:     pubKey,
		keyType: &keyType,
	}
	secretKey = &SecretKey{
		key:       privKey,
		PublicKey: publicKey,
	}
	return
}

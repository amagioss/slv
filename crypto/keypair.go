package crypto

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/json"

	"golang.org/x/crypto/nacl/box"
	"gopkg.in/yaml.v3"
)

type KeyPair struct {
	publicKey  *PublicKey
	privateKey *PrivateKey
}

func (keyPair *KeyPair) PublicKey() PublicKey {
	return *keyPair.publicKey
}

func (keyPair *KeyPair) PrivateKey() PrivateKey {
	return *keyPair.privateKey
}

func NewKeyPair(keyType KeyType) (keyPair *KeyPair, err error) {
	pubKey, privKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, ErrGeneratingKeyPair
	}
	pubSumBytes := sha1.Sum(pubKey[:])
	privSumBytes := sha1.Sum(privKey[:])
	id := [8]byte(append(pubSumBytes[len(pubSumBytes)-4:], privSumBytes[len(privSumBytes)-4:]...))
	keyPair = &KeyPair{
		publicKey: &PublicKey{
			key: key{
				id:      id,
				public:  true,
				keyType: keyType,
				keyData: *pubKey,
			},
		},
		privateKey: &PrivateKey{
			key: key{
				id:      id,
				public:  false,
				keyType: keyType,
				keyData: *privKey,
			},
		},
	}
	return
}

type PrivateKey struct {
	key
	decrypter *Decrypter
}

func (privateKey *PrivateKey) GetDecrypter() Decrypter {
	if privateKey.decrypter == nil {
		privateKey.decrypter = new(Decrypter)
		privateKey.decrypter.privateKey = privateKey
	}
	return *privateKey.decrypter
}

type PublicKey struct {
	key
	encrypter *Encrypter
}

func (publicKey PublicKey) MarshalYAML() (interface{}, error) {
	return publicKey.String(), nil

}

func (publicKey *PublicKey) UnmarshalYAML(value *yaml.Node) (err error) {
	var slvPublicKey string
	err = value.Decode(&slvPublicKey)
	if err == nil {
		return publicKey.FromString(slvPublicKey)
	}
	return
}

func (publicKey PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(publicKey.String())
}

func (publicKey *PublicKey) UnmarshalJSON(data []byte) (err error) {
	var pubKeyStr string
	err = json.Unmarshal(data, &pubKeyStr)
	if err == nil {
		return publicKey.FromString(pubKeyStr)
	}
	return
}

func (publicKey *PublicKey) GetEncrypter() (encrypter *Encrypter, err error) {
	if publicKey.encrypter == nil {
		if ephemeralKeyPair, err := NewKeyPair(0); err == nil {
			encrypter = new(Encrypter)
			encrypter.encryptionKeyId = publicKey.id
			encrypter.ephemeralPublicKey = ephemeralKeyPair.publicKey.keyData
			encrypter.sharedKey = new([32]byte)
			box.Precompute(encrypter.sharedKey, &publicKey.keyData, &ephemeralKeyPair.privateKey.keyData)
			publicKey.encrypter = encrypter
		}
	}
	return publicKey.encrypter, err
}

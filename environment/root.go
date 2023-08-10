package environment

import "github.com/shibme/slv/crypto"

type Root struct {
	PublicKey        crypto.PublicKey `yaml:"public_key"`
	SealedPrivateKey crypto.SealedKey `yaml:"sealed_key"`
}

func newRoot() (root *Root, rootKey *crypto.PrivateKey, err error) {
	root = &Root{}
	var rootKeyPair, sealingKeyPair *crypto.KeyPair
	rootKeyPair, err = crypto.NewKeyPair(RootKey)
	if err != nil {
		return nil, nil, err
	}
	sealingKeyPair, err = crypto.NewKeyPair(RootKey)
	if err != nil {
		return nil, nil, err
	}
	sealingPublicKey := sealingKeyPair.PublicKey()
	encrypter, err := sealingPublicKey.GetEncrypter()
	if err != nil {
		return nil, nil, err
	}
	root.SealedPrivateKey, err = encrypter.EncryptKey(rootKeyPair.PrivateKey())
	if err != nil {
		return nil, nil, err
	}
	root.PublicKey = rootKeyPair.PublicKey()
	rootPrivKey := rootKeyPair.PrivateKey()
	rootKey = &rootPrivKey
	return
}

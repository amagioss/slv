package environments

import "github.com/shibme/slv/core/crypto"

type Root struct {
	PublicKey  crypto.PublicKey  `yaml:"publicKey"`
	WrappedKey crypto.WrappedKey `yaml:"wrappedKey"`
}

func newRoot() (*Root, *crypto.SecretKey, error) {
	rootSKey, err := crypto.NewSecretKey(RootKey)
	if err != nil {
		return nil, nil, err
	}
	rootPKey, err := rootSKey.PublicKey()
	if err != nil {
		return nil, nil, err
	}
	sealingSKey, err := crypto.NewSecretKey(RootKey)
	if err != nil {
		return nil, nil, err
	}
	sealingPKey, err := sealingSKey.PublicKey()
	if err != nil {
		return nil, nil, err
	}
	rootWrappedKey, err := sealingPKey.EncryptKey(*rootSKey)
	if err != nil {
		return nil, nil, err
	}
	return &Root{
		PublicKey:  *rootPKey,
		WrappedKey: *rootWrappedKey,
	}, sealingSKey, nil
}

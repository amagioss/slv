package providers

import "github.com/amagimedia/slv/core/crypto"

func bindWithLocal(inputs map[string][]byte) (publicKey *crypto.PublicKey, ref map[string][]byte, err error) {
	return nil, nil, nil
}

func unBindFromLocal(ref map[string][]byte) (secretKeyBytes []byte, err error) {
	return nil, nil
}

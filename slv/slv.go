package slv

import "github.com/shibme/slv/slv/cipher"

func KeyGen() (public_key, private_key string, err error) {
	var key_pair cipher.KeyPair
	public_key = key_pair.GetPublicKey()
	private_key = key_pair.GetPrivateKey()
	err = key_pair.New()
	return
}

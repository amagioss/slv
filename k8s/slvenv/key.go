package slvenv

import (
	"oss.amagi.com/slv"
	"oss.amagi.com/slv/core/crypto"
)

var SecretKey *crypto.SecretKey

func InitSLVSecretKey() error {
	if SecretKey == nil {
		sk, err := slv.GetSecretKey()
		if err != nil {
			return err
		}
		SecretKey = sk
	}
	return nil
}

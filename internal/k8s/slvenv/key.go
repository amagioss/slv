package slvenv

import (
	"oss.amagi.com/slv/internal/core/crypto"
	"oss.amagi.com/slv/internal/core/secretkey"
)

var SecretKey *crypto.SecretKey

func InitSLVSecretKey() error {
	if SecretKey == nil {
		sk, err := secretkey.Get()
		if err != nil {
			return err
		}
		SecretKey = sk
	}
	return nil
}

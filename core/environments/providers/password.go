package providers

import (
	"dev.shib.me/xipher"
	"oss.amagi.com/slv/core/config"
	"oss.amagi.com/slv/core/input"
)

func bindWithPassword(skBytes []byte, inputs map[string][]byte) (ref map[string][]byte, err error) {
	password := inputs["password"]
	if len(password) == 0 {
		return nil, err
	}
	xipherKey, err := xipher.NewPrivateKeyForPassword(password)
	if err != nil {
		return nil, err
	}
	sealedSecretKeyBytes, err := xipherKey.Encrypt(skBytes, false)
	if err != nil {
		return nil, err
	}
	ref = make(map[string][]byte)
	ref["ssk"] = sealedSecretKeyBytes
	return
}

func unBindWithPassword(ref map[string][]byte) (secretKeyBytes []byte, err error) {
	sealedSecretKeyBytes := ref["ssk"]
	if len(sealedSecretKeyBytes) == 0 {
		return nil, errSealedSecretKeyRef
	}
	var password []byte
	passwordStr := config.GetEnvSecretPassword()
	if passwordStr == "" {
		if input.IsInteractive() == nil {
			if password, err = input.GetHiddenInput("Enter Password: "); err != nil {
				return nil, err
			}
		}
		if password == nil {
			return nil, errPasswordNotSet
		}
	} else {
		password = []byte(passwordStr)
	}
	xipherKey, err := xipher.NewPrivateKeyForPassword([]byte(password))
	if err != nil {
		return nil, err
	}
	secretKeyBytes, err = xipherKey.Decrypt(sealedSecretKeyBytes)
	if err != nil {
		return nil, errInvalidPassword
	}
	return
}

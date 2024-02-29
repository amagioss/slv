package providers

import (
	"dev.shib.me/xipher"
	"savesecrets.org/slv/core/config"
	"savesecrets.org/slv/core/crypto"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/input"
)

func bindWithPassword(inputs map[string][]byte) (publicKey *crypto.PublicKey, ref map[string][]byte, err error) {
	password := inputs["password"]
	if len(password) == 0 {
		return nil, nil, err
	}
	xipherKey, err := xipher.NewPrivateKeyForPassword(password)
	if err != nil {
		return nil, nil, err
	}
	var secretKey *crypto.SecretKey
	if secretKey, err = crypto.NewSecretKey(environments.EnvironmentKey); err != nil {
		return nil, nil, err
	}
	sealedSecretKeyBytes, err := xipherKey.Encrypt(secretKey.Bytes(), false)
	if err != nil {
		return nil, nil, err
	}
	if publicKey, err = secretKey.PublicKey(); err == nil {
		ref = make(map[string][]byte)
		ref["ssk"] = sealedSecretKeyBytes
		return
	}
	return nil, nil, err
}

func unBindWithPassword(ref map[string][]byte) (secretKeyBytes []byte, err error) {
	sealedSecretKeyBytes := ref["ssk"]
	if len(sealedSecretKeyBytes) == 0 {
		return nil, errSealedSecretKeyRef
	}
	var password []byte
	passwordStr := config.GetEnvSecretPassword()
	if passwordStr == "" {
		if input.IsAllowed() == nil {
			if password, err = input.GetPasswordFromUser(false, nil); err != nil {
				return nil, err
			}
		}
		if password == nil {
			return nil, errPasswordNotSet
		}
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

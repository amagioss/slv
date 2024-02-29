package slv

import (
	"savesecrets.org/slv/core/config"
	"savesecrets.org/slv/core/crypto"
	"savesecrets.org/slv/core/environments"
)

func GetSecretKey() (*crypto.SecretKey, error) {
	if secretKey != nil {
		return secretKey, nil
	}
	var err error
	// Read direct secret key from environment variable
	secretKeyStr := config.GetEnvSecretKey()
	if secretKeyStr != "" {
		secretKey, err = crypto.SecretKeyFromString(secretKeyStr)
		if err != nil {
			return nil, err
		} else {
			return secretKey, nil
		}
	}
	// Read secret key from secret binding
	envSecretBindingStr := config.GetEnvSecretBinding()
	if envSecretBindingStr == "" {
		selfEnv := environments.GetSelf()
		if selfEnv != nil {
			envSecretBindingStr = selfEnv.SecretBinding
		}
	}
	if envSecretBindingStr != "" {
		secretKey, err = environments.GetSecretKeyFromSecretBinding(envSecretBindingStr)
	}
	if secretKey == nil && err == nil {
		err = errEnvironmentAccessNotFound
	}
	return secretKey, err
}

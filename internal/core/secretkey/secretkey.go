package secretkey

import (
	"errors"

	"oss.amagi.com/slv/internal/core/config"
	"oss.amagi.com/slv/internal/core/crypto"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/environments/providers"
)

var (
	secretKey                    *crypto.SecretKey
	errEnvironmentAccessNotFound = errors.New("environment doesn't have access. please set the required environment variables")
)

func Get() (*crypto.SecretKey, error) {
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
		secretKey, err = providers.GetSecretKeyFromSecretBinding(envSecretBindingStr)
	}
	if secretKey == nil && err == nil {
		err = errEnvironmentAccessNotFound
	}
	return secretKey, err
}

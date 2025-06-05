package secretkey

import (
	"errors"

	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/environments/envproviders"
)

var (
	secretKey          *crypto.SecretKey
	ErrEnvSecretNotSet = errors.New("environment secret key or secret binding not set")
)

func Get() (*crypto.SecretKey, error) {
	if secretKey != nil {
		return secretKey, nil
	}
	var err error
	// Read direct secret key from environment variable
	if secretKeyStr := config.GetEnvSecretKey(); secretKeyStr != "" {
		return crypto.SecretKeyFromString(secretKeyStr)
	}
	// Read secret key from secret binding
	envSecretBindingStr := config.GetEnvSecretBinding()
	if envSecretBindingStr == "" {
		if selfEnv := environments.GetSelf(); selfEnv != nil {
			envSecretBindingStr = selfEnv.SecretBinding
		}
	}
	if envSecretBindingStr != "" {
		secretKey, err = envproviders.GetSecretKeyFromSecretBinding(envSecretBindingStr)
	}
	if secretKey == nil && err == nil {
		err = ErrEnvSecretNotSet
	}
	return secretKey, err
}

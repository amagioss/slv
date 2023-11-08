package secretkeystore

import (
	"os"

	"github.com/shibme/slv/core/crypto"
)

var secretKey *crypto.SecretKey

func GetSecretKey() (*crypto.SecretKey, error) {
	if secretKey != nil {
		return secretKey, nil
	}
	var err error
	secretKey, err = getSecretKeyFromEnvar()
	if err == nil && secretKey == nil {
		envProviderContextData := os.Getenv(slvProviderContextEnvarName)
		if envProviderContextData != "" {
			secretKey, err = getSecretKeyFromEnvProviderContext(envProviderContextData)
			if err != nil {
				return nil, err
			}
		}
	}
	if secretKey == nil && err == nil {
		err = ErrEnvironmentAccessNotFound
	}
	return secretKey, err
}

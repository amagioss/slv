package secretkeystore

import (
	"os"

	"github.com/amagimedia/slv/core/crypto"
	"github.com/amagimedia/slv/core/environments"
)

var secretKey *crypto.SecretKey

func GetSecretKey() (*crypto.SecretKey, error) {
	if secretKey != nil {
		return secretKey, nil
	}
	var err error
	secretKey, err = getSecretKeyFromEnvar()
	if err == nil && secretKey == nil {
		envProviderBindingStr := os.Getenv(slvAccessBindingEnvarName)
		if envProviderBindingStr != "" {
			secretKey, err = environments.GetSecretKeyFromAccessBinding(envProviderBindingStr)
			if err != nil {
				return nil, err
			}
		}
	}
	if secretKey == nil && err == nil {
		err = errEnvironmentAccessNotFound
	}
	return secretKey, err
}

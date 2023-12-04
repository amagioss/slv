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
		envAccessBindingStr := os.Getenv(slvAccessBindingEnvarName)
		if envAccessBindingStr != "" {
			secretKey, err = getSecretKeyFromAccessBindingString(envAccessBindingStr)
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

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
		accessKeyDef := os.Getenv(slvAccessKeyDefEnvarName)
		if accessKeyDef != "" {
			secretKey, err = getSecretKeyFromAccessKeyDef(accessKeyDef)
			if err != nil {
				return nil, err
			}
		}

	}
	return secretKey, err
}

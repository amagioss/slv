package secretkeystore

import (
	"os"

	"github.com/shibme/slv/core/crypto"
)

func getSecretKeyFromEnvar() (*crypto.SecretKey, error) {
	secretKeyStr := os.Getenv(slvSecreKeyEnvarName)
	if secretKeyStr == "" {
		return nil, nil
	}
	return crypto.SecretKeyFromString(secretKeyStr)
}

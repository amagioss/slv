package secretkeystore

import (
	"os"

	"github.com/shibme/slv/core/crypto"
)

const (
	slvSecreKeyEnvarName = "SLV_SECRET_KEY"
)

func getSecretKeyFromEnvar() (*crypto.SecretKey, error) {
	secretKeyStr := os.Getenv(slvSecreKeyEnvarName)
	if secretKeyStr == "" {
		return nil, nil
	}
	return crypto.SecretKeyFromString(secretKeyStr)
}

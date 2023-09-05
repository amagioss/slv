package keyreader

import "os"

func GetFromEnvar() (string, error) {
	secretKey := os.Getenv("SLV_SECRET_KEY")
	if secretKey == "" {
		return "", ErrEnvSecretNotSet
	}
	return secretKey, nil
}

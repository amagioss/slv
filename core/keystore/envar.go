package keystore

import "os"

func GetSecretKeyFromEnvar() (string, error) {
	secretKey := os.Getenv(slvSecreKeyEnvarName)
	if secretKey == "" {
		return "", ErrEnvSecretNotSet
	}
	return secretKey, nil
}

func getPassphraseFromEnvar() (string, error) {
	passphrase := os.Getenv(slvPassphraseEnvarName)
	if passphrase == "" {
		return "", ErrEnvSecretNotSet
	}
	return passphrase, nil
}

package crypto

import (
	"crypto/rand"
	"io"

	"github.com/btcsuite/btcutil/base58"
)

func randomStringGen(length int) (string, error) {
	randBytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}
	return base58.Encode(randBytes), nil
}

func SecretRefString() (string, error) {
	randomString, err := randomStringGen(32)
	if err != nil {
		return "", err
	}
	return secretRefPrefix + randomString, nil
}

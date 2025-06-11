package envproviders

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"slv.sh/slv/internal/core/commons"
)

const (
	sealedSecretKeyRefName = "ssk"
	rsaPubKeyRefName       = "rsa-pubkey"
)

var (
	errInvalidRSAPublicKey = errors.New("invalid RSA public key")
	errSealedSecretKeyRef  = errors.New("invalid sealed secret key from provider binding")

	rsaArg = arg{
		id:          rsaPubKeyRefName,
		name:        "RSA public key",
		required:    false,
		description: "RSA public key (file path or PEM encoded content) for offline binding [Only RSA OAEP SHA-256 will be supported - recommended 4096 bits key]",
	}
)

func rsaEncrypt(plain, rsaPublicKey []byte) (encrypted []byte, err error) {
	if commons.FileExists(string(rsaPublicKey)) {
		keyFilePath := string(rsaPublicKey)
		if rsaPublicKey, err = os.ReadFile(keyFilePath); err != nil || len(rsaPublicKey) == 0 {
			return nil, errors.New("failed to read RSA public key from file: " + keyFilePath)
		}
	}
	block, _ := pem.Decode(rsaPublicKey)
	if block == nil {
		return nil, errInvalidRSAPublicKey
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errInvalidRSAPublicKey
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, plain, []byte(""))
}

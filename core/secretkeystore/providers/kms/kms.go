package kms

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"github.com/shibme/slv/core/secretkeystore/providers"
)

func encryptSecretKeyBytesWithKMS(secretKeyBytes, rsaPublicKey []byte) (sealedSecretKeyBytes []byte, err error) {
	// Encrypting Environment Secret Key with RSA OAEP SHA256
	block, _ := pem.Decode(rsaPublicKey)
	if block == nil {
		return nil, ErrInvalidRSAPublicKey
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidRSAPublicKey
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, secretKeyBytes, []byte(""))
}

func LoadKMSProviders() {
	providers.Register("awskms", BindWithAWSKMS, UnBindFromAWSKMS, true)
}

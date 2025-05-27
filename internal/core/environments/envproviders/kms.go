package envproviders

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

func rsaEncrypt(plain, rsaPublicKey []byte) (encrypted []byte, err error) {
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

package environments

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
)

type AccessKey struct {
	*accessKey
}

type accessKey struct {
	Accessor        string `json:"type"`
	Ref             string `json:"ref"`
	SealedSecretKey []byte `json:"ssk"`
}

func AccessKeyFromDefString(accessKeyDef string) (ak *AccessKey, err error) {
	sliced := strings.Split(accessKeyDef, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV || sliced[1] != envAccessKeyDefAbbrev {
		return nil, ErrInvalidAccessKeyDef
	}
	accessKey := new(accessKey)
	err = commons.Deserialize(sliced[2], &accessKey)
	if err == nil {
		ak = &AccessKey{accessKey}
	}
	return
}

func newAccessKeyForSecretKey(accessor string, ref string, secretKey *crypto.SecretKey, rsaPublicKey []byte) (ak *AccessKey, err error) {
	var sealedSecretKey []byte

	//Encrypting Environment Secret Key with RSA OAEP SHA256
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
	sealedSecretKey, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, secretKey.Bytes(), []byte(""))

	if err == nil {
		ak = &AccessKey{
			accessKey: &accessKey{
				Accessor:        accessor,
				Ref:             ref,
				SealedSecretKey: sealedSecretKey,
			},
		}
	}
	return
}

func (ak AccessKey) String() (string, error) {
	data, err := commons.Serialize(*ak.accessKey)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s_%s", commons.SLV, envAccessKeyDefAbbrev, data), nil
}

func (ak *AccessKey) Accessor() string {
	return ak.accessKey.Accessor
}

func (ak *AccessKey) Ref() string {
	return ak.accessKey.Ref
}

func (ak *AccessKey) SealedSecretKey() []byte {
	return ak.accessKey.SealedSecretKey
}

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

type EnvAccessBinding struct {
	*envAccessBinding
}

type envAccessBinding struct {
	Provider        string `json:"t"`
	Ref             string `json:"r"`
	SealedSecretKey []byte `json:"k"`
}

func EnvAccessBindingFromString(envAccessBindingStr string) (epc *EnvAccessBinding, err error) {
	sliced := strings.Split(envAccessBindingStr, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV || sliced[1] != envAccessBindingAbbrev {
		return nil, ErrInvalidEnvAccessBinding
	}
	pd := new(envAccessBinding)
	err = commons.Deserialize(sliced[2], &pd)
	if err == nil {
		epc = &EnvAccessBinding{pd}
	}
	return
}

func newEnvAccessBindingForSecretKey(provider string, ref string, secretKey *crypto.SecretKey, rsaPublicKey []byte) (eab *EnvAccessBinding, err error) {
	//Encrypting Environment Secret Key with RSA OAEP SHA256
	var sealedSecretKey []byte
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
		eab = &EnvAccessBinding{
			envAccessBinding: &envAccessBinding{
				Provider:        provider,
				Ref:             ref,
				SealedSecretKey: sealedSecretKey,
			},
		}
	}
	return
}

func (eab EnvAccessBinding) String() (string, error) {
	data, err := commons.Serialize(*eab.envAccessBinding)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s_%s", commons.SLV, envAccessBindingAbbrev, data), nil
}

func (eab *EnvAccessBinding) Provider() string {
	return eab.envAccessBinding.Provider
}

func (eab *EnvAccessBinding) Ref() string {
	return eab.envAccessBinding.Ref
}

func (eab *EnvAccessBinding) SealedKey() []byte {
	return eab.envAccessBinding.SealedSecretKey
}

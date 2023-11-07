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

type EnvProviderContext struct {
	*envProviderContext
}

type envProviderContext struct {
	ProviderType    string `json:"type"`
	Id              string `json:"id"`
	SealedSecretKey []byte `json:"ssk"`
}

func EnvProviderContextFromStringData(envProviderContextStringData string) (epc *EnvProviderContext, err error) {
	sliced := strings.Split(envProviderContextStringData, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV || sliced[1] != envProviderContextAbbrev {
		return nil, ErrInvalidEnvProviderContextData
	}
	envProviderContext := new(envProviderContext)
	err = commons.Deserialize(sliced[2], &envProviderContext)
	if err == nil {
		epc = &EnvProviderContext{envProviderContext}
	}
	return
}

func newEnvProviderContextForSecretKey(providerType string, id string, secretKey *crypto.SecretKey, rsaPublicKey []byte) (epc *EnvProviderContext, err error) {
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
		epc = &EnvProviderContext{
			envProviderContext: &envProviderContext{
				ProviderType:    providerType,
				Id:              id,
				SealedSecretKey: sealedSecretKey,
			},
		}
	}
	return
}

func (epc EnvProviderContext) String() (string, error) {
	data, err := commons.Serialize(*epc.envProviderContext)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s_%s", commons.SLV, envProviderContextAbbrev, data), nil
}

func (epc *EnvProviderContext) Type() string {
	return epc.envProviderContext.ProviderType
}

func (epc *EnvProviderContext) Id() string {
	return epc.envProviderContext.Id
}

func (epc *EnvProviderContext) SealedSecretKey() []byte {
	return epc.envProviderContext.SealedSecretKey
}

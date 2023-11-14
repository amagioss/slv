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

type ProviderData struct {
	*providerData
}

type providerData struct {
	ProviderType      string `json:"t"`
	ProviderRef       string `json:"r"`
	ProviderSealedKey []byte `json:"k"`
}

func ProviderDataFromString(providerDataString string) (epc *ProviderData, err error) {
	sliced := strings.Split(providerDataString, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV || sliced[1] != providerDataStringAbbrev {
		return nil, ErrInvalidEnvProviderContextData
	}
	pd := new(providerData)
	err = commons.Deserialize(sliced[2], &pd)
	if err == nil {
		epc = &ProviderData{pd}
	}
	return
}

func newProviderDataForSecretKey(providerType string, providerRef string, secretKey *crypto.SecretKey, rsaPublicKey []byte) (epc *ProviderData, err error) {
	//Encrypting Environment Secret Key with RSA OAEP SHA256
	var providerSealedKey []byte
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
	providerSealedKey, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, secretKey.Bytes(), []byte(""))
	if err == nil {
		epc = &ProviderData{
			providerData: &providerData{
				ProviderType:      providerType,
				ProviderRef:       providerRef,
				ProviderSealedKey: providerSealedKey,
			},
		}
	}
	return
}

func (pd ProviderData) String() (string, error) {
	data, err := commons.Serialize(*pd.providerData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s_%s", commons.SLV, providerDataStringAbbrev, data), nil
}

func (pd *ProviderData) Type() string {
	return pd.providerData.ProviderType
}

func (pd *ProviderData) Ref() string {
	return pd.providerData.ProviderRef
}

func (pd *ProviderData) SealedKey() []byte {
	return pd.providerData.ProviderSealedKey
}

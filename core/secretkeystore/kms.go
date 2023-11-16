package secretkeystore

import (
	"strings"

	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/environments"
	"github.com/shibme/slv/core/secretkeystore/awskms"
)

func NewEnvForKMS(name, email string, envType environments.EnvType, kmsType, kmsRef string, rsa4096PublicKey []byte) (env *environments.Environment, err error) {
	kmsType = strings.ToUpper(kmsType)
	switch kmsType {
	case awskms.AccessSourceAWS:
		return awskms.NewEnvironment(name, email, envType, kmsRef, rsa4096PublicKey)
	default:
		err = ErrInvalidEnvProviderType
	}
	return
}

func getSecretKeyFromProviderDataString(providerDataString string) (secretKey *crypto.SecretKey, err error) {
	providerData, err := environments.ProviderDataFromString(providerDataString)
	if err != nil {
		return nil, err
	}
	switch providerData.Type() {
	case awskms.AccessSourceAWS:
		secretKey, err = awskms.GetSecretKeyUsingAWSKMS(providerData)
	default:
		err = ErrInvalidEnvProviderType
	}
	return
}

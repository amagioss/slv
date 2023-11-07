package secretkeystore

import (
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/environments"
	"github.com/shibme/slv/core/secretkeystore/awskms"
)

func NewEnvForKMS(name, email string, envType environments.EnvType, kmsType, kmsRef string, rsa4096PublicKey []byte) (env *environments.Environment, err error) {
	switch kmsType {
	case awskms.AccessSourceAWS:
		return awskms.NewEnvironment(name, email, envType, kmsRef, rsa4096PublicKey)
	default:
		err = ErrInvalidEnvProviderType
	}
	return
}

func getSecretKeyFromEnvProviderContext(envProviderContextData string) (secretKey *crypto.SecretKey, err error) {
	envProviderContext, err := environments.EnvProviderContextFromStringData(envProviderContextData)
	if err != nil {
		return nil, err
	}
	switch envProviderContext.Type() {
	case awskms.AccessSourceAWS:
		secretKey, err = awskms.GetSecretKeyUsingAWSKMS(envProviderContext)
	default:
		err = ErrInvalidEnvProviderType
	}
	return
}

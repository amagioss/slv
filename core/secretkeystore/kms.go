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
		err = ErrInvalidEnvAccessProvider
	}
	return
}

func getSecretKeyFromAccessBindingString(accessBindingStr string) (secretKey *crypto.SecretKey, err error) {
	accessBinding, err := environments.EnvAccessBindingFromString(accessBindingStr)
	if err != nil {
		return nil, err
	}
	switch accessBinding.Provider() {
	case awskms.AccessSourceAWS:
		secretKey, err = awskms.GetSecretKeyUsingAWSKMS(accessBinding)
	default:
		err = ErrInvalidEnvAccessProvider
	}
	return
}

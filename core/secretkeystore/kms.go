package secretkeystore

import (
	"regexp"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/environments"
)

const (
	accessKeyTypeAWS   = "AWS"
	accessKeyTypeGCP   = "GCP"
	accessKeyTypeAzure = "Azure"
)

func NewEnvForKMS(name, email string, envType environments.EnvType, kmsType, kmsRef string, rsa4096PublicKey []byte) (env *environments.Environment, accessKey *environments.AccessKey, err error) {
	if kmsType == accessKeyTypeAWS {
		return NewEnvForAWSKMS(name, email, envType, kmsRef, rsa4096PublicKey)
	}
	return nil, nil, ErrInvalidAccessKeyType
}

func NewEnvForAWSKMS(name, email string, envType environments.EnvType, arn string, rsa4096PublicKey []byte) (env *environments.Environment, accessKey *environments.AccessKey, err error) {
	validARN, _ := regexp.MatchString(awsKMSARNPattern, arn)
	if !validARN {
		return nil, nil, ErrInvalidAWSKMSARN
	}
	return environments.NewEnvironmentWithAccessKey(name, email, envType, accessKeyTypeAWS, arn, rsa4096PublicKey)
}

func retrieveSecretKeyFromAWSKMS(accessKey *environments.AccessKey) (secretKey *crypto.SecretKey, err error) {
	awsSession, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	kmsClient := kms.New(awsSession)
	input := &kms.DecryptInput{
		CiphertextBlob:      accessKey.SealedSecretKey(),
		KeyId:               commons.String(accessKey.Ref()),
		EncryptionAlgorithm: commons.String(awsKMSEncryptionAlgorithm),
	}
	result, err := kmsClient.Decrypt(input)
	if err != nil {
		return nil, err
	}
	return crypto.SecretKeyFromBytes(result.Plaintext)
}

func getSecretKeyFromAccessKeyDef(accessKeyDef string) (secretKey *crypto.SecretKey, err error) {
	accessKey, err := environments.AccessKeyFromDefString(accessKeyDef)
	if err != nil {
		return nil, err
	}
	switch accessKey.Accessor() {
	case accessKeyTypeAWS:
		secretKey, err = retrieveSecretKeyFromAWSKMS(accessKey)
	default:
		err = ErrInvalidAccessKeyType
	}
	return
}

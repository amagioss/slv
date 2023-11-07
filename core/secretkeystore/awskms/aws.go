package awskms

import (
	"errors"
	"regexp"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/environments"
)

const (
	AccessSourceAWS           = "AWS"
	awsKMSEncryptionAlgorithm = "RSAES_OAEP_SHA_256"
	awsKMSARNPattern          = `^arn:aws:kms:[a-z0-9-]+:[0-9]+:key/[a-f0-9-]+$`
)

var ErrInvalidAWSKMSARN = errors.New("invalid AWS KMS ARN")

func NewEnvironment(name, email string, envType environments.EnvType, arn string, rsa4096PublicKey []byte) (env *environments.Environment, err error) {
	validARN, _ := regexp.MatchString(awsKMSARNPattern, arn)
	if !validARN {
		return nil, ErrInvalidAWSKMSARN
	}
	return environments.NewEnvironmentWithProvider(name, email, envType, AccessSourceAWS, arn, rsa4096PublicKey)
}

func GetSecretKeyUsingAWSKMS(envProviderContext *environments.EnvProviderContext) (secretKey *crypto.SecretKey, err error) {
	awsSession, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	kmsClient := kms.New(awsSession)
	input := &kms.DecryptInput{
		CiphertextBlob:      envProviderContext.SealedSecretKey(),
		KeyId:               commons.String(envProviderContext.Id()),
		EncryptionAlgorithm: commons.String(awsKMSEncryptionAlgorithm),
	}
	result, err := kmsClient.Decrypt(input)
	if err != nil {
		return nil, err
	}
	return crypto.SecretKeyFromBytes(result.Plaintext)
}

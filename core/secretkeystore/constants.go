package secretkeystore

import "errors"

const (
	awsKMSEncryptionAlgorithm = "RSAES_OAEP_SHA_256"
	awsKMSARNPattern          = `^arn:aws:kms:[a-z0-9-]+:[0-9]+:key/[a-f0-9-]+$`
	slvAccessKeyDefEnvarName  = "SLV_ACCESS_KEY"
)

var ErrInvalidAccessKeyType = errors.New("invalid access key type")
var ErrSecretKeyNotAccessible = errors.New("secret key not accessible. please set one of the environment variables: " + slvSecreKeyEnvarName + " or " + slvAccessKeyDefEnvarName)
var ErrInvalidAWSKMSARN = errors.New("invalid AWS KMS ARN")

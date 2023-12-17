package providers

import "errors"

const (

	// AWS Provider Constants
	awsKMSAsymmetricEncryptionAlgorithm = "RSAES_OAEP_SHA_256"
	awsKMSARNPattern                    = `^arn:aws:kms:[a-z0-9-]+:[0-9]+:key/[a-f0-9-]+$`
)

var (
	defaultProvidersRegistered = false

	// KMS Provider Errors
	errInvalidRSAPublicKey = errors.New("invalid RSA public key")

	// AWS Provider Errors
	errAWSConfiguration       = errors.New("please configure AWS access")
	errInvalidAWSKMSARN       = errors.New("invalid AWS KMS ARN")
	errInvalidAWSKMSAlgorithm = errors.New("invalid AWS KMS algorithm")
	errSealedSecretKeyRef     = errors.New("invalid sealed secret key from provider binding")
)

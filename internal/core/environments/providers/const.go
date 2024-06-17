package providers

import (
	"errors"

	"oss.amagi.com/slv/internal/core/config"
)

const (
	envSecretBindingAbbrev = "ESB" // Environment Secret Binding
	slvPrefix              = config.AppNameUpperCase

	// Password Provider Constants
	passwordProviderName = "password"
	keyringServiceName   = config.AppNameLowerCase

	// KMS Provider Constants
	rsaPubKeyRefName = "rsa-pubkey"

	// AWS Provider Constants
	awsProviderName                     = "aws"
	awsARNRefName                       = "arn"
	awsKMSAsymmetricEncryptionAlgorithm = "RSAES_OAEP_SHA_256"
	awsKMSARNPattern                    = `^arn:aws:kms:[a-z0-9-]+:[0-9]+:key/[a-f0-9-]+$`

	// GCP Provider Constants
	gcpProviderName    = "gcp"
	gcpResourceNameRef = "resource-name"
	gcpSymmAlgoRef     = "sym"
)

var (
	defaultProvidersRegistered = false

	// Provider Base Errors
	errProviderUnknown               = errors.New("unknown provider")
	errInvalidEnvSecretBindingFormat = errors.New("invalid environment secret binding format")
	errEnvSecretBindingUnspecified   = errors.New("environment secret binding unspecified")
	errProviderRegisteredAlready     = errors.New("env secret provider registered already")

	// KMS Provider Errors
	errInvalidRSAPublicKey = errors.New("invalid RSA public key")
	errSealedSecretKeyRef  = errors.New("invalid sealed secret key from provider binding")

	// AWS Provider Errors
	errAWSConfiguration       = errors.New("please configure AWS access")
	errInvalidAWSKMSARN       = errors.New("invalid AWS KMS ARN")
	errInvalidAWSKMSAlgorithm = errors.New("invalid AWS KMS algorithm")

	// GCP Provider Errors
	errInvalidGCPKMSResourceName = errors.New("invalid GCP KMS resource name")
	errInvalidGCPKMSAlgorithm    = errors.New("invalid GCP KMS algorithm")

	// Password Provider Errors
	errPasswordNotSet  = errors.New("password not set: please set password through the environment variable or use the interactive terminal to enter the password")
	errInvalidPassword = errors.New("invalid password")
)

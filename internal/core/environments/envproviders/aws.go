package envproviders

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"slv.sh/slv/internal/core/commons"
)

const (
	awsProviderId                       = "aws"
	awsProviderName                     = "AWS KMS"
	awsProviderDesc                     = "AWS Key Management Service"
	awsARNRefName                       = "arn"
	awsAlgRefName                       = "alg"
	awsKMSAsymmetricEncryptionAlgorithm = "RSAES_OAEP_SHA_256"
	awsKMSARNPattern                    = `^arn:aws:kms:[a-z0-9-]+:[0-9]+:key/[a-z0-9-]+$`
)

var (
	errAWSConfiguration       = errors.New("please configure AWS access")
	errInvalidAWSKMSARN       = errors.New("invalid AWS KMS ARN")
	errInvalidAWSKMSAlgorithm = errors.New("invalid AWS KMS algorithm")

	awsArgs = []arg{
		{
			id:          awsARNRefName,
			name:        "ARN",
			required:    true,
			description: "ARN of the AWS KMS key to use",
		},
		rsaArg,
	}
)

func isValidARN(arn string) bool {
	validARN, _ := regexp.MatchString(awsKMSARNPattern, arn)
	return validARN
}

func isAWSConfigured(ctx context.Context, cfg aws.Config) bool {
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err == nil && creds.AccessKeyID != "" && (creds.SecretAccessKey != "" || creds.SessionToken != "") {
		stsClient := sts.NewFromConfig(cfg)
		_, err = stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		return err == nil
	}
	return false
}

func encryptWithAWSKMSAPI(secretKeyBytes []byte, arn string) (sealedSecretKeyBytes []byte, algorithm *string, err error) {
	if !isValidARN(arn) {
		return nil, nil, errInvalidAWSKMSARN
	}
	ctx := context.Background()
	arnParts := strings.Split(arn, ":")
	region := arnParts[3]

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, nil, err
	}

	kmsClient := kms.NewFromConfig(cfg)
	input := &kms.DescribeKeyInput{
		KeyId: commons.StringPtr(arn),
	}
	keyDesc, err := kmsClient.DescribeKey(ctx, input)
	if err != nil {
		if !isAWSConfigured(ctx, cfg) {
			return nil, nil, errAWSConfiguration
		}
		return nil, nil, err
	}
	encryptionAlgos := keyDesc.KeyMetadata.EncryptionAlgorithms
	selectedAlgo := encryptionAlgos[len(encryptionAlgos)-1]
	algoString := string(selectedAlgo)
	algorithm = &algoString
	kmsInput := &kms.EncryptInput{
		KeyId:               commons.StringPtr(arn),
		Plaintext:           secretKeyBytes,
		EncryptionAlgorithm: selectedAlgo,
	}
	result, err := kmsClient.Encrypt(ctx, kmsInput)
	if err != nil {
		return nil, nil, err
	}
	return result.CiphertextBlob, algorithm, err
}

func bindWithAWSKMS(skBytes []byte, inputs map[string]string) (ref map[string][]byte, err error) {
	if arn := string(inputs[awsARNRefName]); isValidARN(arn) {
		var sealedSecretKeyBytes []byte
		rsaPublicKey, ok := inputs[rsaPubKeyRefName]
		ref = make(map[string][]byte)
		ref[awsARNRefName] = []byte(arn)
		if !ok || len(rsaPublicKey) == 0 {
			var algorithm *string
			if sealedSecretKeyBytes, algorithm, err = encryptWithAWSKMSAPI(skBytes, arn); err != nil {
				return nil, err
			}
			ref[awsAlgRefName] = []byte(*algorithm)
		} else if sealedSecretKeyBytes, err = rsaEncrypt(skBytes, []byte(rsaPublicKey)); err == nil {
			ref[awsAlgRefName] = []byte(awsKMSAsymmetricEncryptionAlgorithm)
		} else {
			return nil, err
		}
		ref[sealedSecretKeyRefName] = sealedSecretKeyBytes
		return
	}
	return nil, errInvalidAWSKMSARN
}

func unBindFromAWSKMS(ref map[string][]byte) (secretKeyBytes []byte, err error) {
	arn := string(ref[awsARNRefName])
	if !isValidARN(arn) {
		return nil, errInvalidAWSKMSARN
	}
	ctx := context.Background()
	arnParts := strings.Split(arn, ":")
	region := arnParts[3]

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	sealedSecretKeyBytes := ref[sealedSecretKeyRefName]
	if len(sealedSecretKeyBytes) == 0 {
		return nil, errSealedSecretKeyRef
	}
	algorithm := ref[awsAlgRefName]
	if len(algorithm) == 0 {
		return nil, errInvalidAWSKMSAlgorithm
	}

	kmsClient := kms.NewFromConfig(cfg)
	kmsInput := &kms.DecryptInput{
		CiphertextBlob:      sealedSecretKeyBytes,
		KeyId:               commons.StringPtr(arn),
		EncryptionAlgorithm: types.EncryptionAlgorithmSpec(string(algorithm)),
	}
	result, err := kmsClient.Decrypt(ctx, kmsInput)
	if err != nil {
		if !isAWSConfigured(ctx, cfg) {
			return nil, errAWSConfiguration
		}
		return nil, err
	}
	return result.Plaintext, nil
}

package kms

import (
	"errors"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
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
var ErrSealedSecretKeyRef = errors.New("invalid sealed secret key ref from provider binding")

func isValidARN(arn string) bool {
	validARN, _ := regexp.MatchString(awsKMSARNPattern, arn)
	return validARN
}

func BindWithAWSKMS(inputs map[string]any) (publicKey *crypto.PublicKey, ref map[string]any, err error) {
	arn, ok := inputs["arn"].(string)
	if !ok || !isValidARN(arn) {
		return nil, nil, ErrInvalidAWSKMSARN
	}
	rsaPublicKey, ok := inputs["pubkey"].([]byte)
	if !ok || len(rsaPublicKey) == 0 {
		return nil, nil, ErrInvalidRSAPublicKey
	}
	secretKey, err := crypto.NewSecretKey(environments.EnvironmentKey)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err = secretKey.PublicKey()
	if err != nil {
		return
	}
	sealedSecretKeyBytes, err := encryptSecretKeyBytesWithKMS(secretKey.Bytes(), rsaPublicKey)
	if err != nil {
		return nil, nil, err
	}
	ref = make(map[string]any)
	ref["arn"] = arn
	ref["sealedSecretKeyBytes"] = sealedSecretKeyBytes
	return publicKey, ref, nil
}

func UnBindFromAWSKMS(ref map[string]any) (secretKeyBytes []byte, err error) {
	arn, ok := ref["arn"].(string)
	if !ok || !isValidARN(arn) {
		return nil, ErrInvalidAWSKMSARN
	}
	arnParts := strings.Split(arn, ":")
	region := arnParts[3]
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	sealedSecretKeyBytes, ok := ref["arn"].([]byte)
	if !ok || len(sealedSecretKeyBytes) == 0 {
		return nil, ErrSealedSecretKeyRef
	}
	kmsClient := kms.New(awsSession)
	kmsInput := &kms.DecryptInput{
		CiphertextBlob:      sealedSecretKeyBytes,
		KeyId:               commons.StringPtr(arn),
		EncryptionAlgorithm: commons.StringPtr(awsKMSEncryptionAlgorithm),
	}
	result, err := kmsClient.Decrypt(kmsInput)
	if err != nil {
		return nil, err
	}
	return result.Plaintext, nil
}

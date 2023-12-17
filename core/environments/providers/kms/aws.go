package kms

import (
	"errors"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/environments"
)

const (
	AccessSourceAWS                     = "AWS"
	awsKMSAsymmetricEncryptionAlgorithm = "RSAES_OAEP_SHA_256"
	awsKMSARNPattern                    = `^arn:aws:kms:[a-z0-9-]+:[0-9]+:key/[a-f0-9-]+$`
)

var (
	ErrAWSConfiguration   = errors.New("please configure AWS access")
	ErrInvalidAWSKMSARN   = errors.New("invalid AWS KMS ARN")
	ErrSealedSecretKeyRef = errors.New("invalid sealed secret key ref from provider binding")
)

func isValidARN(arn string) bool {
	validARN, _ := regexp.MatchString(awsKMSARNPattern, arn)
	return validARN
}

func isAWSConfigured(session *session.Session) bool {
	if creds, err := session.Config.Credentials.Get(); err == nil && creds.AccessKeyID != "" && (creds.SecretAccessKey != "" || creds.SessionToken != "") {
		_, err = sts.New(session).GetCallerIdentity(&sts.GetCallerIdentityInput{})
		return err == nil
	}
	return false
}

func encryptWithAWSKMSAPI(secretKeyBytes []byte, arn string) (sealedSecretKeyBytes []byte, err error) {
	if !isValidARN(arn) {
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
	kmsClient := kms.New(awsSession)
	kmsInput := &kms.EncryptInput{
		KeyId:     commons.StringPtr(arn),
		Plaintext: secretKeyBytes,
		// EncryptionAlgorithm: commons.StringPtr(awsKMSEncryptionAlgorithm),
	}
	result, err := kmsClient.Encrypt(kmsInput)
	if err != nil {
		if !isAWSConfigured(awsSession) {
			return nil, ErrAWSConfiguration
		}
		return nil, err
	}
	return result.CiphertextBlob, nil
}

func BindWithAWSKMS(inputs map[string][]byte) (publicKey *crypto.PublicKey, ref map[string][]byte, err error) {
	if arn := string(inputs["arn"]); isValidARN(arn) {
		var secretKey *crypto.SecretKey
		if secretKey, err = crypto.NewSecretKey(environments.EnvironmentKey); err != nil {
			return nil, nil, err
		}
		var sealedSecretKeyBytes []byte
		rsaPublicKey, ok := inputs["rsa-pubkey"]
		ref = make(map[string][]byte)
		ref["arn"] = []byte(arn)
		if !ok || len(rsaPublicKey) == 0 {
			if sealedSecretKeyBytes, err = encryptWithAWSKMSAPI(secretKey.Bytes(), arn); err != nil {
				return nil, nil, err
			}
		} else if sealedSecretKeyBytes, err = rsaEncrypt(secretKey.Bytes(), rsaPublicKey); err != nil {
			return nil, nil, err
		} else {
			ref["alg"] = []byte(awsKMSAsymmetricEncryptionAlgorithm)
		}
		if publicKey, err = secretKey.PublicKey(); err != nil {
			return nil, nil, err
		}
		ref["ssk"] = sealedSecretKeyBytes
		return
	}
	return nil, nil, ErrInvalidAWSKMSARN
}

func UnBindFromAWSKMS(ref map[string][]byte) (secretKeyBytes []byte, err error) {
	arn := string(ref["arn"])
	if !isValidARN(arn) {
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
	sealedSecretKeyBytes := ref["ssk"]
	if len(sealedSecretKeyBytes) == 0 {
		return nil, ErrSealedSecretKeyRef
	}
	kmsClient := kms.New(awsSession)
	kmsInput := &kms.DecryptInput{
		CiphertextBlob: sealedSecretKeyBytes,
		KeyId:          commons.StringPtr(arn),
	}
	if alg, ok := ref["alg"]; ok && len(alg) > 0 {
		kmsInput.EncryptionAlgorithm = commons.StringPtr(string(alg))
	}
	result, err := kmsClient.Decrypt(kmsInput)
	if err != nil {
		if !isAWSConfigured(awsSession) {
			return nil, ErrAWSConfiguration
		}
		return nil, err
	}
	return result.Plaintext, nil
}

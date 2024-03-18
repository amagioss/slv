package providers

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/sts"
	"savesecrets.org/slv/core/commons"
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

func encryptWithAWSKMSAPI(secretKeyBytes []byte, arn string) (sealedSecretKeyBytes []byte, algorithm *string, err error) {
	if !isValidARN(arn) {
		return nil, nil, errInvalidAWSKMSARN
	}
	arnParts := strings.Split(arn, ":")
	region := arnParts[3]
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, nil, err
	}
	kmsClient := kms.New(awsSession)
	input := &kms.DescribeKeyInput{
		KeyId: aws.String(arn),
	}
	keyDesc, err := kmsClient.DescribeKey(input)
	if err != nil {
		if !isAWSConfigured(awsSession) {
			return nil, nil, errAWSConfiguration
		}
		return nil, nil, err
	}
	encryptionAlgos := keyDesc.KeyMetadata.EncryptionAlgorithms
	algorithm = encryptionAlgos[len(encryptionAlgos)-1]
	kmsInput := &kms.EncryptInput{
		KeyId:               commons.StringPtr(arn),
		Plaintext:           secretKeyBytes,
		EncryptionAlgorithm: algorithm,
	}
	result, err := kmsClient.Encrypt(kmsInput)
	if err != nil {
		return nil, nil, err
	}
	return result.CiphertextBlob, algorithm, err
}

func bindWithAWSKMS(skBytes []byte, inputs map[string][]byte) (ref map[string][]byte, err error) {
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
			ref["alg"] = []byte(*algorithm)
		} else if sealedSecretKeyBytes, err = rsaEncrypt(skBytes, rsaPublicKey); err == nil {
			ref["alg"] = []byte(awsKMSAsymmetricEncryptionAlgorithm)
		} else {
			return nil, err
		}
		ref["ssk"] = sealedSecretKeyBytes
		return
	}
	return nil, errInvalidAWSKMSARN
}

func unBindFromAWSKMS(ref map[string][]byte) (secretKeyBytes []byte, err error) {
	arn := string(ref["arn"])
	if !isValidARN(arn) {
		return nil, errInvalidAWSKMSARN
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
		return nil, errSealedSecretKeyRef
	}
	algorithm := ref["alg"]
	if len(algorithm) == 0 {
		return nil, errInvalidAWSKMSAlgorithm
	}
	kmsClient := kms.New(awsSession)
	kmsInput := &kms.DecryptInput{
		CiphertextBlob:      sealedSecretKeyBytes,
		KeyId:               commons.StringPtr(arn),
		EncryptionAlgorithm: commons.StringPtr(string(algorithm)),
	}
	result, err := kmsClient.Decrypt(kmsInput)
	if err != nil {
		if !isAWSConfigured(awsSession) {
			return nil, errAWSConfiguration
		}
		return nil, err
	}
	return result.Plaintext, nil
}

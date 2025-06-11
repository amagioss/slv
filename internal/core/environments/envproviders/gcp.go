package envproviders

import (
	"context"
	"errors"
	"strings"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
)

const (
	gcpProviderId      = "gcp"
	gcpProviderName    = "GCP KMS"
	gcpProviderDesc    = "Google Cloud Key Management Service"
	gcpResourceNameRef = "resource-name"
	gcpSymmAlgoRef     = "sym"
)

var (
	errInvalidGCPKMSResourceName = errors.New("invalid Cloud KMS resource name")
	errInvalidGCPKMSAlgorithm    = errors.New("invalid Cloud KMS algorithm")

	gcpArgs = []arg{
		{
			id:          gcpResourceNameRef,
			name:        "Resource Name",
			required:    true,
			description: "Resource name of the GCP KMS key to use",
		},
		rsaArg,
	}
)

func isValidGCPResourceName(resourcePath string, symmetricAlgo bool) bool {
	return strings.HasPrefix(resourcePath, "projects/") && strings.Contains(resourcePath, "locations/") &&
		strings.Contains(resourcePath, "keyRings/") && strings.Contains(resourcePath, "cryptoKeys/") &&
		(symmetricAlgo != strings.Contains(resourcePath, "cryptoKeyVersions/"))
}

func bindWithGCP(skBytes []byte, inputs map[string]string) (ref map[string][]byte, err error) {
	rsaPublicKey, ok := inputs[rsaPubKeyRefName]
	symmetricAlgo := !ok || len(rsaPublicKey) == 0
	if resourceName := string(inputs[gcpResourceNameRef]); isValidGCPResourceName(resourceName, symmetricAlgo) {
		var sealedSecretKeyBytes []byte
		ref = make(map[string][]byte)
		ref[gcpResourceNameRef] = []byte(resourceName)
		if symmetricAlgo {
			ctx := context.Background()
			client, err := kms.NewKeyManagementClient(ctx)
			if err != nil {
				return nil, err
			}
			req := &kmspb.EncryptRequest{
				Name:      resourceName,
				Plaintext: skBytes,
			}
			result, err := client.Encrypt(ctx, req)
			if err != nil {
				return nil, err
			}
			sealedSecretKeyBytes = result.Ciphertext
			symmetricAlgo = true
		} else if sealedSecretKeyBytes, err = rsaEncrypt(skBytes, []byte(rsaPublicKey)); err != nil {
			return nil, err
		}
		var symmetricAlgoByte byte
		if symmetricAlgo {
			symmetricAlgoByte = 1
		} else {
			symmetricAlgoByte = 0
		}
		ref[gcpSymmAlgoRef] = []byte{symmetricAlgoByte}
		ref[sealedSecretKeyRefName] = sealedSecretKeyBytes
		return
	}
	return nil, errInvalidGCPKMSResourceName
}

func unBindWithGCP(ref map[string][]byte) (secretKeyBytes []byte, err error) {
	if ref[gcpSymmAlgoRef] == nil {
		return nil, errInvalidGCPKMSAlgorithm
	}
	symmetricAlgo := ref[gcpSymmAlgoRef][0] == 1
	resourceName := string(ref[gcpResourceNameRef])
	if !isValidGCPResourceName(resourceName, symmetricAlgo) {
		return nil, errInvalidGCPKMSResourceName
	}
	sealedSecretKeyBytes := ref[sealedSecretKeyRefName]
	if len(sealedSecretKeyBytes) == 0 {
		return nil, errSealedSecretKeyRef
	}
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	var plaintext []byte
	if symmetricAlgo {
		req := &kmspb.DecryptRequest{
			Name:       resourceName,
			Ciphertext: sealedSecretKeyBytes,
		}
		result, err := client.Decrypt(ctx, req)
		if err != nil {
			return nil, err
		}
		plaintext = result.Plaintext
	} else {
		req := &kmspb.AsymmetricDecryptRequest{
			Name:       resourceName,
			Ciphertext: sealedSecretKeyBytes,
		}
		result, err := client.AsymmetricDecrypt(ctx, req)
		if err != nil {
			return nil, err
		}
		plaintext = result.Plaintext
	}
	return plaintext, nil
}

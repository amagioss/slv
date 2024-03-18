package providers

import (
	"context"
	"strings"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
)

func isValidGCPResourceName(resourcePath string, symmetricAlgo bool) bool {
	return strings.HasPrefix(resourcePath, "projects/") && strings.Contains(resourcePath, "locations/") &&
		strings.Contains(resourcePath, "keyRings/") && strings.Contains(resourcePath, "cryptoKeys/") &&
		(symmetricAlgo != strings.Contains(resourcePath, "cryptoKeyVersions/"))
}

func bindWithGCP(skBytes []byte, inputs map[string][]byte) (ref map[string][]byte, err error) {
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
		} else if sealedSecretKeyBytes, err = rsaEncrypt(skBytes, rsaPublicKey); err != nil {
			return nil, err
		}
		var symmetricAlgoByte byte
		if symmetricAlgo {
			symmetricAlgoByte = 1
		} else {
			symmetricAlgoByte = 0
		}
		ref[gcpSymmAlgoRef] = []byte{symmetricAlgoByte}
		ref["ssk"] = sealedSecretKeyBytes
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
	sealedSecretKeyBytes := ref["ssk"]
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

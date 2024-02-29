package providers

import (
	"context"
	"strings"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
	"savesecrets.org/slv/core/crypto"
	"savesecrets.org/slv/core/environments"
)

func isValidGCPResourcePath(resourcePath string, symmetricAlgo bool) bool {
	return strings.HasPrefix(resourcePath, "projects/") && strings.Contains(resourcePath, "locations/") &&
		strings.Contains(resourcePath, "keyRings/") && strings.Contains(resourcePath, "cryptoKeys/") &&
		(symmetricAlgo != strings.Contains(resourcePath, "cryptoKeyVersions/"))
}

func bindWithGCP(inputs map[string][]byte) (publicKey *crypto.PublicKey, ref map[string][]byte, err error) {
	if resourceName := string(inputs[gcpResourceNameRef]); isValidGCPResourcePath(resourceName, true) {
		var secretKey *crypto.SecretKey
		if secretKey, err = crypto.NewSecretKey(environments.EnvironmentKey); err != nil {
			return nil, nil, err
		}
		var sealedSecretKeyBytes []byte
		rsaPublicKey, ok := inputs[rsaPubKeyRefName]
		ref = make(map[string][]byte)
		ref[gcpResourceNameRef] = []byte(resourceName)
		symmetricAlgo := false
		if !ok || len(rsaPublicKey) == 0 {
			ctx := context.Background()
			client, err := kms.NewKeyManagementClient(ctx)
			if err != nil {
				return nil, nil, err
			}
			req := &kmspb.EncryptRequest{
				Name:      resourceName,
				Plaintext: secretKey.Bytes(),
			}
			result, err := client.Encrypt(ctx, req)
			if err != nil {
				return nil, nil, err
			}
			sealedSecretKeyBytes = result.Ciphertext
			symmetricAlgo = true
		} else if sealedSecretKeyBytes, err = rsaEncrypt(secretKey.Bytes(), rsaPublicKey); err != nil {
			return nil, nil, err
		}
		if publicKey, err = secretKey.PublicKey(); err != nil {
			return nil, nil, err
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
	return nil, nil, errInvalidGCPKMSResourceName
}

func unBindWithGCP(ref map[string][]byte) (secretKeyBytes []byte, err error) {
	if ref[gcpSymmAlgoRef] == nil {
		return nil, errInvalidGCPKMSAlgorithm
	}
	symmetricAlgo := ref[gcpSymmAlgoRef][0] == 1
	resourceName := string(ref[gcpResourceNameRef])
	if !isValidGCPResourcePath(resourceName, symmetricAlgo) {
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

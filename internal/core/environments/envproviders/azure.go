package envproviders

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

const (
	azureProviderId    = "azure"
	azureProviderName  = "Azure Key Vault"
	azureProviderDesc  = "Azure Key Vault (only RSA 4096 with AES 256 GCM supported)"
	azureVaultUrlRef   = "vault-url"
	azureKeyNameRef    = "key-name"
	azureKeyVersionRef = "key-version"

	azureKeyVaultUrlPattern = `^https://[a-z0-9-]+\.vault\.azure\.net/$`
	akvAlgorithm            = azkeys.EncryptionAlgorithmRSAOAEP256
)

var (
	azureArgs = []arg{
		{
			id:          azureVaultUrlRef,
			name:        "Vault URL",
			required:    true,
			description: "URL of the Azure Key Vault to use (e.g., https://myvault.vault.azure.net/)",
		},
		{
			id:          azureKeyNameRef,
			name:        "Key Name",
			required:    true,
			description: "Name of the key in the Azure Key Vault to use",
		},
		{
			id:          azureKeyVersionRef,
			name:        "Key Version",
			required:    false,
			description: "Version of the key in Azure Key Vault to use (optional, latest version will be used if not specified)",
		},
		rsaArg,
	}
)

func isValidKeyVaultUrl(url string) bool {
	validUrl, _ := regexp.MatchString(azureKeyVaultUrlPattern, url)
	return validUrl
}

func bindWithAzure(skBytes []byte, inputs map[string]string) (ref map[string][]byte, err error) {
	akvUrl := string(inputs[azureVaultUrlRef])
	if !isValidKeyVaultUrl(akvUrl) {
		return nil, fmt.Errorf("invalid Azure Key Vault URL: %s", akvUrl)
	}
	akvKeyName := string(inputs[azureKeyNameRef])
	akvKeyVersion := string(inputs[azureKeyVersionRef])
	if akvUrl == "" || akvKeyName == "" {
		return nil, fmt.Errorf("missing required inputs: %s and %s", azureVaultUrlRef, azureKeyNameRef)
	}
	var sealedSecretKeyBytes []byte
	if rsaPublicKey, ok := inputs[rsaPubKeyRefName]; !ok || len(rsaPublicKey) == 0 {
		alg := akvAlgorithm
		if cred, err := azidentity.NewDefaultAzureCredential(nil); err != nil {
			return nil, fmt.Errorf("failed to create Azure credential: %w", err)
		} else if client, err := azkeys.NewClient(akvUrl, cred, nil); err != nil {
			return nil, fmt.Errorf("failed to create Azure Key Vault client: %w", err)
		} else if encResp, err := client.Encrypt(context.Background(),
			akvKeyName, akvKeyVersion,
			azkeys.KeyOperationParameters{
				Algorithm: &alg,
				Value:     skBytes,
			}, nil); err != nil {
			return nil, fmt.Errorf("failed to encrypt with Azure Key Vault: %w", err)
		} else {
			sealedSecretKeyBytes = encResp.Result
		}
	} else if sealedSecretKeyBytes, err = rsaEncrypt(skBytes, []byte(rsaPublicKey)); err != nil {
		return nil, err
	}
	ref = make(map[string][]byte)
	ref[azureVaultUrlRef] = []byte(akvUrl)
	ref[azureKeyNameRef] = []byte(akvKeyName)
	if akvKeyVersion != "" {
		ref[azureKeyVersionRef] = []byte(akvKeyVersion)
	}
	ref[sealedSecretKeyRefName] = sealedSecretKeyBytes
	return
}

func unBindFromAzure(ref map[string][]byte) (secretKeyBytes []byte, err error) {
	akvUrl := string(ref[azureVaultUrlRef])
	akvKeyName := string(ref[azureKeyNameRef])
	akvKeyVersion := string(ref[azureKeyVersionRef])
	sealedSecretKeyBytes := ref[sealedSecretKeyRefName]
	if len(sealedSecretKeyBytes) == 0 {
		return nil, fmt.Errorf("sealed secret key is missing in the reference")
	}
	alg := akvAlgorithm
	if cred, err := azidentity.NewDefaultAzureCredential(nil); err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %w", err)
	} else if client, err := azkeys.NewClient(akvUrl, cred, nil); err != nil {
		return nil, fmt.Errorf("failed to create Azure Key Vault client: %w", err)
	} else if decResp, err := client.Decrypt(context.Background(),
		akvKeyName, akvKeyVersion,
		azkeys.KeyOperationParameters{
			Algorithm: &alg,
			Value:     sealedSecretKeyBytes,
		}, nil); err != nil {
		return nil, fmt.Errorf("failed to decrypt with Azure Key Vault: %w", err)
	} else {
		return decResp.Result, nil
	}
}

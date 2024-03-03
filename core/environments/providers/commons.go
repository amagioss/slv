package providers

import "savesecrets.org/slv/core/environments"

func LoadDefaultProviders() {
	if !defaultProvidersRegistered {
		environments.RegisterEnvSecretProvider(awsProviderName, bindWithAWSKMS, unBindFromAWSKMS, true)
		environments.RegisterEnvSecretProvider(passwordProviderName, bindWithPassword, unBindWithPassword, true)
		environments.RegisterEnvSecretProvider(gcpProviderName, bindWithGCP, unBindWithGCP, true)
		defaultProvidersRegistered = true
	}
}

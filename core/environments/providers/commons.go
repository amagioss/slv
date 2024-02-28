package providers

import "savesecrets.org/slv/core/environments"

func LoadDefaults() {
	if !defaultProvidersRegistered {
		environments.RegisterEnvSecretProvider(awsProviderName, bindWithAWSKMS, unBindFromAWSKMS, true)
		environments.RegisterEnvSecretProvider(passwordProviderName, bindWithPassword, unBindWithPassword, true)
		defaultProvidersRegistered = true
	}
}

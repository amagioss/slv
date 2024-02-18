package providers

import "savesecrets.org/slv/core/environments"

func LoadDefaults() {
	if !defaultProvidersRegistered {
		environments.RegisterEnvSecretProvider("kms-aws", bindWithAWSKMS, unBindFromAWSKMS, true)
		defaultProvidersRegistered = true
	}
}

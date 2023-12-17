package providers

import "github.com/shibme/slv/core/environments"

func RegisterDefaultProviders() {
	if !defaultProvidersRegistered {
		environments.RegisterAccessProvider("kms-aws", bindWithAWSKMS, unBindFromAWSKMS, true)
		defaultProvidersRegistered = true
	}
}

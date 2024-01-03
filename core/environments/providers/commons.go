package providers

import "github.com/amagimedia/slv/core/environments"

func LoadDefaults() {
	if !defaultProvidersRegistered {
		environments.RegisterAccessProvider("kms-aws", bindWithAWSKMS, unBindFromAWSKMS, true)
		defaultProvidersRegistered = true
	}
}

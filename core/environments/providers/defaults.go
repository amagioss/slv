package providers

func loadDefaultProviders() {
	if !defaultProvidersRegistered {
		registerProvider(passwordProviderName, bindWithPassword, unBindWithPassword, true)
		registerProvider(awsProviderName, bindWithAWSKMS, unBindFromAWSKMS, true)
		registerProvider(gcpProviderName, bindWithGCP, unBindWithGCP, true)
		defaultProvidersRegistered = true
	}
}

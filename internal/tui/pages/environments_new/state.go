package environments_new

// SaveNavigationState implements the Page interface
func (nep *NewEnvironmentPage) SaveNavigationState() {
	// State is preserved in the page struct
}

// RestoreNavigationState implements the Page interface
func (nep *NewEnvironmentPage) RestoreNavigationState() {
	// Update help text when navigating back
	if nep.navigation != nil {
		nep.navigation.updateHelpText()
	}

	// Set focus to current component
	if component := nep.GetCurrentComponent(); component != nil {
		nep.GetTUI().GetApplication().SetFocus(component)
	}
}

// ClearNavigationState implements the Page interface
func (nep *NewEnvironmentPage) ClearNavigationState() {
	// Reset to initial state when clearing
	nep.currentStep = StepProviderSelection
	nep.selectedType = ""
	nep.envName = ""
	nep.envEmail = ""
	nep.envTags = nil
	nep.quantumSafe = false
	nep.selectedProvider = ""
	nep.providerInputs = make(map[string]string)
	nep.addToProfile = false
	nep.createdEnv = nil
	nep.secretKey = ""
}

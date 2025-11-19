package environments

// SaveNavigationState saves the current navigation state
func (ep *EnvironmentsPage) SaveNavigationState() {
	// Environments page doesn't need to save state yet
}

// RestoreNavigationState restores the saved navigation state
func (ep *EnvironmentsPage) RestoreNavigationState() {
	// Environments page doesn't need to restore state yet
	ep.navigation.updateHelpText()
}

// ClearNavigationState clears the navigation state
func (ep *EnvironmentsPage) ClearNavigationState() {
	// Environments page doesn't need to clear state yet
}

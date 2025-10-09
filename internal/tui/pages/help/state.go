package help

// SaveNavigationState implements the Page interface (empty for help page)
func (hp *HelpPage) SaveNavigationState() {
	// Help page doesn't need state management
}

// RestoreNavigationState implements the Page interface (empty for help page)
func (hp *HelpPage) RestoreNavigationState() {
	// Help page doesn't need state management
}

// ClearNavigationState implements the Page interface (empty for help page)
func (hp *HelpPage) ClearNavigationState() {
	// Help page doesn't need state management
}

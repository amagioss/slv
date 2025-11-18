package mainpage

// SaveNavigationState implements the Page interface (empty for main page)
func (mp *MainPage) SaveNavigationState() {
	// Main page doesn't need state management
}

// RestoreNavigationState implements the Page interface (empty for main page)
func (mp *MainPage) RestoreNavigationState() {
	// Update help text when navigating back
	if mp.navigation != nil {
		mp.navigation.updateHelpText()
	}
}

// ClearNavigationState implements the Page interface (empty for main page)
func (mp *MainPage) ClearNavigationState() {
	// Main page doesn't need state management
}

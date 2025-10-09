package navigation

// SavePageState saves a state value for a specific page and key
func (n *Navigation) SavePageState(pageName, stateKey string, stateValue interface{}) {
	if n.pageStates[pageName] == nil {
		n.pageStates[pageName] = make(map[string]interface{})
	}
	n.pageStates[pageName][stateKey] = stateValue
}

// GetPageState retrieves a state value for a specific page and key
func (n *Navigation) GetPageState(pageName, stateKey string) (interface{}, bool) {
	if pageState, exists := n.pageStates[pageName]; exists {
		value, exists := pageState[stateKey]
		return value, exists
	}
	return nil, false
}

// ClearPageState clears all state for a specific page
func (n *Navigation) ClearPageState(pageName string) {
	delete(n.pageStates, pageName)
}

// ClearPageStateKey clears a specific state key for a page
func (n *Navigation) ClearPageStateKey(pageName, stateKey string) {
	if pageState, exists := n.pageStates[pageName]; exists {
		delete(pageState, stateKey)
	}
}

// HasPageState checks if a page has any saved state
func (n *Navigation) HasPageState(pageName string) bool {
	pageState, exists := n.pageStates[pageName]
	return exists && len(pageState) > 0
}

// saveCurrentPageState saves the state of the currently active page
func (n *Navigation) saveCurrentPageState() {
	currentPage := n.GetCurrentPage()
	if page, exists := n.pageInstances[currentPage]; exists && page != nil {
		// Check if the page has a SaveNavigationState method
		// We'll use type assertion to check for pages that support state saving
		switch p := page.(type) {
		case interface{ SaveNavigationState() }:
			p.SaveNavigationState()
		}
	}
}

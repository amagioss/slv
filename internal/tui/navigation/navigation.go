package navigation

import (
	"os"

	"slv.sh/slv/internal/tui/interfaces"
)

// Navigation handles page navigation and routing
type Navigation struct {
	app           interfaces.TUIInterface
	vaultDir      string                     // Store current vault directory
	customHelp    string                     // Store custom help text for current page
	pageInstances map[string]interfaces.Page // Store actual page instances for refresh

	// General state management for all pages
	pageStates map[string]map[string]interface{} // pageName -> stateKey -> stateValue
}

// NewNavigation creates a new navigation controller
func NewNavigation(app interfaces.TUIInterface) *Navigation {
	// Get user's home directory
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}

	nav := &Navigation{
		app:           app,
		vaultDir:      currentDir,
		pageInstances: make(map[string]interfaces.Page),
		pageStates:    make(map[string]map[string]interface{}),
	}

	nav.UpdateStatus()       // Initialize status bar with content
	nav.setupInputHandling() // Setup global input handling
	return nav
}

// GetApp returns the TUI interface
func (n *Navigation) GetApp() interfaces.TUIInterface {
	return n.app
}

// GetVaultDir returns the current vault directory
func (n *Navigation) GetVaultDir() string {
	return n.vaultDir
}

// SetVaultDir sets the current vault directory
func (n *Navigation) SetVaultDir(dir string) {
	n.vaultDir = dir
}

// GetCustomHelp returns the custom help text
func (n *Navigation) GetCustomHelp() string {
	return n.customHelp
}

// StorePageInstance stores a page instance for later refresh
func (n *Navigation) StorePageInstance(pageName string, page interfaces.Page) {
	n.pageInstances[pageName] = page
}

// GetPageInstance retrieves a stored page instance
func (n *Navigation) GetPageInstance(pageName string) (interfaces.Page, bool) {
	page, exists := n.pageInstances[pageName]
	return page, exists
}

// RemovePageInstance removes a page instance from storage
func (n *Navigation) RemovePageInstance(pageName string) {
	delete(n.pageInstances, pageName)
}

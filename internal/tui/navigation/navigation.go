package navigation

import (
	"os"

	"slv.sh/slv/internal/tui/interfaces"
)

// Navigation handles page navigation and routing
type Navigation struct {
	app        interfaces.TUIInterface
	vaultDir   string // Store current vault directory
	customHelp string // Store custom help text for current page
}

// NewNavigation creates a new navigation controller
func NewNavigation(app interfaces.TUIInterface) *Navigation {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	nav := &Navigation{
		app:      app,
		vaultDir: homeDir,
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

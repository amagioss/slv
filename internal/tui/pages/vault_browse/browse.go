package vault_browse

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
)

// VaultFile represents a file or directory in the vault browser
type VaultFile struct {
	Name   string
	Path   string
	IsFile bool
}

// VaultBrowsePage handles the vault browsing functionality
type VaultBrowsePage struct {
	pages.BasePage
	currentDir string
	vault      *vaults.Vault // Store the current vault instance
	vaultPath  string        // Store the current vault path

	// UI components
	mainContent *tview.Grid
	pwdTextView *tview.TextView
	fileList    *tview.List

	navigation *FormNavigation
}

// NewVaultBrowsePage creates a new VaultBrowsePage instance
func NewVaultBrowsePage(tui interfaces.TUIInterface, currentDir string) *VaultBrowsePage {
	vbp := &VaultBrowsePage{
		BasePage:   *pages.NewBasePage(tui, "Vault Management"),
		currentDir: currentDir,
		vault:      nil,
		vaultPath:  "",
	}
	vbp.mainContent = vbp.createMainSection()
	vbp.navigation = (&FormNavigation{}).NewFormNavigation(vbp)
	vbp.updateFileList() // Initial population of the list
	vbp.navigation.SetupNavigation()
	return vbp
}

// Create implements the Page interface
func (vbp *VaultBrowsePage) Create() tview.Primitive {
	// Update status bar with help text
	vbp.GetTUI().UpdateStatusBar("[yellow]←/→: Move between directories | ↑/↓: Navigate | Enter: open vault/directory | Ctrl+N: New vault[white]")

	// Create layout using BasePage method
	vbp.SetTitle("Vault Management")
	return vbp.CreateLayout(vbp.mainContent)
}

// Refresh implements the Page interface
func (vbp *VaultBrowsePage) Refresh() {
	vbp.updateFileList()
	// Recreate page using navigation system
	vbp.GetTUI().GetNavigation().ShowVaultsWithDir(vbp.currentDir, true)

	// Update help text for the current focus
	if vbp.navigation != nil {
		vbp.navigation.updateHelpText()
	}

}

// HandleInput implements the Page interface
func (vbp *VaultBrowsePage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// TODO: Implement vault browsing page input handling
	return event
}

// GetTitle implements the Page interface
func (vbp *VaultBrowsePage) GetTitle() string {
	return vbp.BasePage.GetTitle()
}

// GetCurrentDir returns the current directory
func (vbp *VaultBrowsePage) GetCurrentDir() string {
	return vbp.currentDir
}

// SetCurrentDir sets the current directory
func (vbp *VaultBrowsePage) SetCurrentDir(dir string) {
	vbp.currentDir = dir
}

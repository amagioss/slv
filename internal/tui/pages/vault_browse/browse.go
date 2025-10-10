package vault_browse

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
)

// VaultFile represents a file or directory in the vault browser
type VaultFile struct {
	Name         string
	Path         string
	IsFile       bool
	IsAccessible bool
}

// VaultBrowsePage handles the vault browsing functionality
type VaultBrowsePage struct {
	pages.BasePage
	currentDir string
	vaultPath  string // Store the current vault path

	// UI components
	mainContent   *tview.Grid
	directoryList *tview.List
	fileList      *tview.List

	navigation *FormNavigation

	// Pre-loaded data for performance (lazy loading - only loaded when accessed)
	directoryMap map[string][]VaultFile // dirPath -> []VaultFile
	vaultFileMap map[string][]VaultFile // dirPath -> []VaultFile
}

// NewVaultBrowsePage creates a new VaultBrowsePage instance
func NewVaultBrowsePage(tui interfaces.TUIInterface, currentDir string) *VaultBrowsePage {
	vbp := &VaultBrowsePage{
		BasePage:     *pages.NewBasePage(tui, "Vault Management"),
		currentDir:   currentDir,
		vaultPath:    "",
		directoryMap: make(map[string][]VaultFile),
		vaultFileMap: make(map[string][]VaultFile),
	}

	// Check if we have saved state and use that directory instead
	nav := tui.GetNavigation()
	if lastViewedDir, hasViewedDir := nav.GetPageState("vaults", "lastViewedDir"); hasViewedDir {
		if viewedDir, ok := lastViewedDir.(string); ok {
			vbp.currentDir = viewedDir
		}
	}

	// Load current directory (lazy loading for subdirectories)
	vbp.ensureDirectoryLoaded(vbp.currentDir)

	vbp.mainContent = vbp.createMainSection()
	vbp.navigation = (&FormNavigation{}).NewFormNavigation(vbp)
	vbp.updateFileList()
	vbp.RestoreNavigationState()
	// Initial population of the list
	vbp.navigation.SetupNavigation()
	return vbp
}

// Create implements the Page interface
func (vbp *VaultBrowsePage) Create() tview.Primitive {
	// Update status bar with help text

	// Create layout using BasePage method
	vbp.SetTitle("Vault Management")
	return vbp.CreateLayout(vbp.mainContent)
}

// Refresh implements the Page interface
func (vbp *VaultBrowsePage) Refresh() {
	// Save current state before refreshing
	vbp.SaveNavigationState()

	vbp.updateFileList()
	// Recreate page using navigation system
	vbp.GetTUI().GetNavigation().ShowVaultsWithDir(vbp.currentDir, true)
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

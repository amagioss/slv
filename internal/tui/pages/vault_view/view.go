package vault_view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
)

// VaultViewPage handles the vault details viewing functionality
type VaultViewPage struct {
	pages.BasePage
	vault      *vaults.Vault
	filePath   string
	navigation *FormNavigation

	currentPage tview.Primitive

	accessorsTable    *tview.Table
	itemsTable        *tview.Table
	vaultDetailsTable *tview.Table
	mainFlex          *tview.Flex
}

// NewVaultViewPage creates a new VaultViewPage instance
func NewVaultViewPage(tui interfaces.TUIInterface, vault *vaults.Vault, filePath string) *VaultViewPage {
	vvp := &VaultViewPage{
		BasePage: *pages.NewBasePage(tui, "Vault Details"),
		vault:    vault,
		filePath: filePath,
	}
	vvp.currentPage = vvp.createMainSection()
	vvp.navigation = (&FormNavigation{}).NewFormNavigation(vvp)
	vvp.navigation.SetupNavigation()

	// Restore navigation state will be called after setCurrentPage in ShowVaultDetailsWithVault

	return vvp
}

// Create implements the Page interface
func (vvp *VaultViewPage) Create() tview.Primitive { // Create a flex layout to hold the three tables
	flex := vvp.currentPage
	// Update status bar with help text
	// vvp.GetTUI().UpdateStatusBar("[yellow]q: close | u: unlock | l: lock | r: reload | Tab: switch tables[white]")

	vvp.SetTitle("Vault Details")
	return vvp.CreateLayout(flex)
}

// Refresh implements the Page interface
func (vvp *VaultViewPage) Refresh() {
	// Save current state before refreshing
	vvp.SaveNavigationState()

	// Reload vault data from disk to get fresh state
	vvp.reloadVaultData()

	// Refresh content by recreating the page primitive
	vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vvp.vault, vvp.filePath, true)

	// Update help text for the current focus
	if vvp.navigation != nil {
		vvp.navigation.updateHelpText()
		// vvp.navigation.currentFocus = oldFocus
	}
}

// reloadVaultData reloads the vault data without recreating the page
func (vvp *VaultViewPage) reloadVaultData() {
	if vvp.filePath == "" {
		return
	}

	// Load fresh vault instance
	vault, err := vaults.Get(vvp.filePath)
	if err != nil {
		// Don't show error during refresh, just keep existing vault
		return
	}

	// Update stored instance
	vvp.vault = vault
}

// HandleInput implements the Page interface
func (vvp *VaultViewPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// TODO: Implement vault view page input handling
	return event
}

// GetTitle implements the Page interface
func (vvp *VaultViewPage) GetTitle() string {
	return vvp.BasePage.GetTitle()
}

// SetVault sets the vault
func (vvp *VaultViewPage) SetVault(vault *vaults.Vault) {
	vvp.vault = vault
}

// SetFilePath sets the file path
func (vvp *VaultViewPage) SetFilePath(filePath string) {
	vvp.filePath = filePath
}

// GetVault returns the vault
func (vvp *VaultViewPage) GetVault() *vaults.Vault {
	return vvp.vault
}

// GetFilePath returns the file path
func (vvp *VaultViewPage) GetFilePath() string {
	return vvp.filePath
}

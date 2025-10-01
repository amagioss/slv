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
	return vvp
}

// Create implements the Page interface
func (vvp *VaultViewPage) Create() tview.Primitive { // Create a flex layout to hold the three tables
	flex := vvp.currentPage
	// Update status bar with help text
	vvp.GetTUI().UpdateStatusBar("[yellow]q: close | u: unlock | l: lock | r: reload | Tab: switch tables[white]")

	vvp.SetTitle("Vault Details")
	return vvp.CreateLayout(flex)
}

// Refresh implements the Page interface
func (vvp *VaultViewPage) Refresh() {
	// Reload vault data from disk to get fresh state
	vvp.reloadVaultData()

	// Refresh content by recreating the page primitive
	vvp.refreshContent()

	// Update help text for the current focus
	if vvp.navigation != nil {
		vvp.navigation.updateHelpText()
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

// refreshContent safely refreshes the page content by recreating the primitive
func (vvp *VaultViewPage) refreshContent() {
	// Recreate the main content with fresh tables
	vvp.currentPage = vvp.createMainSection()

	// Set up navigation for the new components
	if vvp.navigation != nil {
		vvp.navigation.SetupNavigation()
	}

	// Get the current page name from the router
	currentPageName := vvp.GetTUI().GetRouter().GetCurrentPage()
	if currentPageName != "" {
		// Replace the page in the main content pages
		vvp.GetTUI().GetComponents().GetMainContentPages().AddPage(currentPageName, vvp.CreateLayout(vvp.currentPage), true, true)

		// Ensure focus is set to the first table after page replacement
		if vvp.navigation != nil {
			vvp.GetTUI().GetApplication().SetFocus(vvp.vaultDetailsTable)
		}
	}
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

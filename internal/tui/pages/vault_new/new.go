package vault_new

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
)

// VaultNewPage handles the new vault creation functionality
type VaultNewPage struct {
	pages.BasePage
	currentDir   string
	publicKeys   []string
	grantedEnvs  []*environments.Environment
	searchEnvMap map[string]*environments.Environment // Map environment names to environment structs for search results
	currentQuery string                               // Store the current search query for refreshing

	// Form references
	vaultConfigForm   *tview.Form   // Vault Configuration form
	optionsForm       *tview.Form   // Options form
	grantAccessForm   *tview.Form   // Grant Access form
	shareWithSelfForm *tview.Form   // Share with Self form
	shareWithK8sForm  *tview.Form   // Share with K8s Context form
	submitButton      *tview.Button // Submit button

	// Lists
	searchResults *tview.List
	grantedAccess *tview.List

	// K8s environment
	k8sEnv *environments.Environment

	navigation *FormNavigation
	// Checkbox references
	// shareWithSelfCheckbox *tview.Checkbox // Reference to the Share with Self checkbox
	// shareWithK8sCheckbox  *tview.Checkbox // Reference to the Share with K8s Context checkbox

	currentPage tview.Primitive // Store reference to current page for modal navigation
}

// NewVaultNewPage creates a new VaultNewPage instance
func NewVaultNewPage(tui interfaces.TUIInterface, currentDir string) *VaultNewPage {
	vnp := &VaultNewPage{
		BasePage:     *pages.NewBasePage(tui, "New Vault"),
		currentDir:   currentDir,
		publicKeys:   []string{},
		grantedEnvs:  []*environments.Environment{},
		searchEnvMap: make(map[string]*environments.Environment),
	}
	vnp.currentPage = vnp.createMainSection()
	vnp.navigation = (&FormNavigation{}).NewFormNavigation(vnp)
	vnp.navigation.SetupNavigation()
	return vnp
}

// Create implements the Page interface
func (vnp *VaultNewPage) Create() tview.Primitive {
	// Create a single comprehensive form
	form := vnp.currentPage

	// Update status bar

	// Create the page layout and show it
	vnp.SetTitle("New Vault at " + vnp.currentDir)
	return vnp.CreateLayout(form)
}

// Refresh implements the Page interface
func (vnp *VaultNewPage) Refresh() {
	// Recreate page using navigation system
	vnp.GetTUI().GetNavigation().ShowNewVaultWithDir(vnp.currentDir, true)

	// Update help text for the current focus
	if vnp.navigation != nil {
		vnp.navigation.updateHelpText()
	}
}

// HandleInput implements the Page interface
func (vnp *VaultNewPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// TODO: Implement new vault page input handling
	return event
}

// GetTitle implements the Page interface
func (vnp *VaultNewPage) GetTitle() string {
	return vnp.BasePage.GetTitle()
}

// SetCurrentDir sets the current directory
func (vnp *VaultNewPage) SetCurrentDir(dir string) {
	// TODO: Implement set current directory
}

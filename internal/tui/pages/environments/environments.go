package environments

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
)

// EnvironmentsPage handles the environments page functionality
type EnvironmentsPage struct {
	pages.BasePage

	// UI components
	mainContent *tview.Flex

	// Environment tables (first row)
	sessionEnvTable *tview.Table
	selfEnvTable    *tview.Table

	// Browse environments section (second row)
	browseEnvsSearch   *tview.InputField
	browseEnvsList     *tview.List
	browseEnvsSection  *tview.Flex
	browseEnvsDetails  *tview.Table // Environment details table (shown when viewing details)
	browseEnvsEDSTable *tview.Table // EDS table (shown when parsing SLV_EDS_)
	editForm           *tview.Form  // Form for editing environment fields

	// Search data
	currentQuery      string                               // Store the current search query for refreshing
	searchEnvMap      map[string]*environments.Environment // Map environment names to environment structs for details view
	currentDetailsEnv *environments.Environment            // Currently displayed environment in details view
	editingField      string                               // Currently editing field name (Name, Email, Tags)

	navigation *FormNavigation
}

// NewEnvironmentsPage creates a new EnvironmentsPage instance
func NewEnvironmentsPage(tui interfaces.TUIInterface) *EnvironmentsPage {
	ep := &EnvironmentsPage{
		BasePage:     *pages.NewBasePage(tui, "Environments"),
		searchEnvMap: make(map[string]*environments.Environment),
	}

	// Create UI structure
	ep.mainContent = ep.createMainSection()

	// Initialize navigation
	ep.navigation = (&FormNavigation{}).NewFormNavigation(ep)
	ep.navigation.SetupNavigation()

	return ep
}

// Create implements the Page interface
func (ep *EnvironmentsPage) Create() tview.Primitive {
	ep.SetTitle("Environments")
	// ep.UpdateStatus("Use Tab to switch between environment tables | ↑↓ to navigate fields | c to copy field | d to copy EDS")
	return ep.CreateLayout(ep.mainContent)
}

// Refresh implements the Page interface
func (ep *EnvironmentsPage) Refresh() {
	// Recreate page using navigation system
	ep.GetTUI().GetNavigation().ShowEnvironments(true)
}

// HandleInput implements the Page interface
func (ep *EnvironmentsPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	return event
}

// GetTitle implements the Page interface
func (ep *EnvironmentsPage) GetTitle() string {
	return ep.BasePage.GetTitle()
}

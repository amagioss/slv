package environments_new

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/theme"
)

// EnvironmentNewPage handles the new environment creation functionality
type EnvironmentNewPage struct {
	pages.BasePage
}

// NewEnvironmentNewPage creates a new EnvironmentNewPage instance
func NewEnvironmentNewPage(tui interfaces.TUIInterface) *EnvironmentNewPage {
	return &EnvironmentNewPage{
		BasePage: *pages.NewBasePage(tui, "New Environment"),
	}
}

// Create implements the Page interface
func (enp *EnvironmentNewPage) Create() tview.Primitive {
	enp.SetTitle("New Environment")

	// Create a simple empty page with just a title
	colors := theme.GetCurrentPalette()
	content := tview.NewTextView().
		SetText("New Environment Page - Coming Soon").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(false).
		SetTextColor(colors.TextPrimary)

	// Update status bar
	enp.UpdateStatus("New Environment - Press ESC to go back")

	return enp.CreateLayout(content)
}

// Refresh implements the Page interface
func (enp *EnvironmentNewPage) Refresh() {
	// Recreate page using navigation system
	enp.GetTUI().GetNavigation().ShowNewEnvironment(true)
}

// HandleInput implements the Page interface
func (enp *EnvironmentNewPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	return event
}

// GetTitle implements the Page interface
func (enp *EnvironmentNewPage) GetTitle() string {
	return enp.BasePage.GetTitle()
}

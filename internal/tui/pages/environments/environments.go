package environments

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/theme"
)

// EnvironmentsPage handles the environments page functionality
type EnvironmentsPage struct {
	pages.BasePage
}

// NewEnvironmentsPage creates a new EnvironmentsPage instance
func NewEnvironmentsPage(tui interfaces.TUIInterface) *EnvironmentsPage {
	return &EnvironmentsPage{
		BasePage: *pages.NewBasePage(tui, "Environments"),
	}
}

// Create implements the Page interface
func (ep *EnvironmentsPage) Create() tview.Primitive {
	// Create content
	text := tview.NewTextView().
		SetText("Environments Page\n\nThis page will show environment management options.").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	// Style the text
	colors := theme.GetCurrentPalette()
	text.SetTextColor(colors.TextPrimary)

	// Update status bar
	ep.UpdateStatus("Environment management - Coming soon")

	// Create layout using BasePage method
	return ep.CreateLayout(text)
}

// Refresh implements the Page interface
func (ep *EnvironmentsPage) Refresh() {
	// Environments page doesn't need refresh logic yet
}

// HandleInput implements the Page interface
func (ep *EnvironmentsPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// Environments page doesn't handle specific input yet
	return event
}

// GetTitle implements the Page interface
func (ep *EnvironmentsPage) GetTitle() string {
	return ep.BasePage.GetTitle()
}

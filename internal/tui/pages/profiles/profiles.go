package profiles

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/theme"
)

// ProfilesPage handles the profiles page functionality
type ProfilesPage struct {
	pages.BasePage
}

// NewProfilesPage creates a new ProfilesPage instance
func NewProfilesPage(tui interfaces.TUIInterface) *ProfilesPage {
	return &ProfilesPage{
		BasePage: *pages.NewBasePage(tui, "Profiles"),
	}
}

// Create implements the Page interface
func (pp *ProfilesPage) Create() tview.Primitive {
	// Create content
	text := tview.NewTextView().
		SetText("Profile management is coming soon.\n\nFor now, use the 'slv profile' CLI commands.").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	// Style the text
	colors := theme.GetCurrentPalette()
	text.SetTextColor(colors.TextPrimary)

	// Update status bar
	pp.UpdateStatus("Profiles - Coming soon (use 'slv profile' on the CLI)")

	// Create layout using BasePage method
	return pp.CreateLayout(text)
}

// Refresh implements the Page interface
func (pp *ProfilesPage) Refresh() {
	// Profiles page doesn't need refresh logic yet
}

// HandleInput implements the Page interface
func (pp *ProfilesPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// Profiles page doesn't handle specific input yet
	return event
}

// GetTitle implements the Page interface
func (pp *ProfilesPage) GetTitle() string {
	return pp.BasePage.GetTitle()
}

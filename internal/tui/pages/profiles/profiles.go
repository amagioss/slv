package profiles

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
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
		SetText("Profiles Page\n\nThis page will show profile management options.").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	// Style the text
	text.SetTextColor(tcell.ColorWhite)

	// Update status bar
	pp.UpdateStatus("[yellow]Profiles management - Coming soon[white]")

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

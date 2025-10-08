package help

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/theme"
)

// HelpPage handles the help page functionality
type HelpPage struct {
	pages.BasePage
}

// NewHelpPage creates a new HelpPage instance
func NewHelpPage(tui interfaces.TUIInterface) *HelpPage {
	return &HelpPage{
		BasePage: *pages.NewBasePage(tui, "Help"),
	}
}

// Create implements the Page interface
func (hp *HelpPage) Create() tview.Primitive {
	// Create help content
	helpText := `SLV TUI Help

[yellow]Navigation:[white]
- Arrow keys: Navigate
- Enter: Select
- Esc: Back
- Ctrl+C: Quit

[yellow]Shortcuts:[white]
- v: Vaults
- p: Profiles
- e: Environments
- h: Help

[yellow]Features:[white]
- Secure vault management
- Profile-based configuration
- Environment-specific settings
- Terminal-based interface

[yellow]For more information:[white]
Visit the documentation or use the CLI commands.`

	text := tview.NewTextView().
		SetText(helpText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	// Style the text
	colors := theme.GetCurrentPalette()
	text.SetTextColor(colors.TextPrimary)

	// Update status bar
	hp.UpdateStatus("Help - Press ESC to go back")

	// Create layout using BasePage method
	return hp.CreateLayout(text)
}

// Refresh implements the Page interface
func (hp *HelpPage) Refresh() {
	// Help page doesn't need refresh logic
}

// HandleInput implements the Page interface
func (hp *HelpPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// Help page doesn't handle specific input yet
	return event
}

// GetTitle implements the Page interface
func (hp *HelpPage) GetTitle() string {
	return hp.BasePage.GetTitle()
}

package mainpage

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/theme"
)

// MainPage handles the main menu page functionality
type MainPage struct {
	pages.BasePage
}

// NewMainPage creates a new MainPage instance
func NewMainPage(tui interfaces.TUIInterface) *MainPage {
	return &MainPage{
		BasePage: *pages.NewBasePage(tui, "Main Menu"),
	}
}

// Create implements the Page interface
func (mp *MainPage) Create() tview.Primitive {
	// Create welcome message parts with subtle color variations
	colors := theme.GetCurrentPalette()

	// Title - brightest white
	titleText := tview.NewTextView().
		SetText("Welcome to SLV - Secure Local Vault").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(colors.TextPrimary) // Brightest white

	// Subtitle - slightly muted
	subtitleText := tview.NewTextView().
		SetText("Your decentralized secrets management solution").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(colors.TextSecondary) // Medium gray

	// Create a flex container for the welcome section
	welcomeFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 2, 0, false).
		AddItem(titleText, 2, 0, false).
		AddItem(subtitleText, 2, 0, false)

	// Create the main menu list with enhanced styling
	list := tview.NewList().
		AddItem("Vaults", "Manage and organize your vaults", 'v', func() {
			mp.NavigateTo("vaults", false)
		}).
		AddItem("Profiles", "View Profile settings and Environments", 'p', func() {
			mp.NavigateTo("profiles", false)
		}).
		AddItem("Environments", "Manage Environments", 'e', func() {
			mp.NavigateTo("environments", false)
		}).
		AddItem("Help", "View documentation and help", 'h', func() {
			mp.NavigateTo("help", false)
		})

	// Style the list
	list.SetSelectedTextColor(colors.ListSelectedText).
		SetSelectedBackgroundColor(colors.ListSelectedBg).
		SetSecondaryTextColor(colors.ListSecondaryText).
		SetMainTextColor(colors.ListMainText)

	// Create a centered layout using grid
	content := tview.NewGrid().
		SetRows(10, 0).         // Two flexible rows for equal centering
		SetColumns(-1, 50, -1). // Single column
		SetBorders(false)

	// Center the welcome text
	content.AddItem(welcomeFlex, 0, 1, 1, 1, 0, 0, false)

	// Center the menu list
	content.AddItem(list, 1, 1, 1, 1, 0, 0, true) // Add padding for centering

	// Update status bar with help text using BasePage method
	mp.UpdateStatus("↑/↓: Navigate | Enter: select")

	// Create layout using BasePage method
	return mp.CreateLayout(content)
}

// Refresh implements the Page interface
func (mp *MainPage) Refresh() {
	// Main page doesn't need refresh logic
}

// HandleInput implements the Page interface
func (mp *MainPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// Main page doesn't handle specific input
	return event
}

// GetTitle implements the Page interface
func (mp *MainPage) GetTitle() string {
	return mp.BasePage.GetTitle()
}

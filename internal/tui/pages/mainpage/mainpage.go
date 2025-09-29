package mainpage

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
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
	// Create a welcome message
	welcomeText := "\n\n" + `[white::b]Welcome to SLV - Secure Local Vault[white::-]

[yellow]Your decentralized secrets management solution[yellow::-]

[cyan]Navigate using arrow keys or press the highlighted letter[cyan::-]
[gray]Press ESC to go back, Ctrl+C to exit[gray::-]` + "\n\n"

	welcome := tview.NewTextView().
		SetText(welcomeText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetWrap(true)

	// Create the main menu list with enhanced styling
	list := tview.NewList().
		AddItem("🔐 Vaults", "Manage and organize your vaults", 'v', func() {
			mp.NavigateTo("vaults")
		}).
		AddItem("👤 Profiles", "View Profile settings and Environments", 'p', func() {
			mp.NavigateTo("profiles")
		}).
		AddItem("🌍 Environments", "Manage Environments", 'e', func() {
			mp.NavigateTo("environments")
		}).
		AddItem("❓ Help", "View documentation and help", 'h', func() {
			mp.NavigateTo("help")
		})

	// Style the list
	list.SetSelectedTextColor(tcell.ColorYellow).
		SetSelectedBackgroundColor(tcell.ColorNavy).
		SetSecondaryTextColor(tcell.ColorGray).
		SetMainTextColor(tcell.ColorWhite)

	// Create a centered layout using grid
	content := tview.NewGrid().
		SetRows(10, 0).         // Two flexible rows for equal centering
		SetColumns(-1, 50, -1). // Single column
		SetBorders(false)

	// Center the welcome text
	content.AddItem(welcome, 0, 1, 1, 1, 0, 0, false)

	// Center the menu list
	content.AddItem(list, 1, 1, 1, 1, 0, 0, true) // Add padding for centering

	// Update status bar with help text using BasePage method
	mp.UpdateStatus("[yellow]↑/↓: Navigate | Enter: select[white]")

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

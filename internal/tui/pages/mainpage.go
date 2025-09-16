package pages

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
)

// MainPage handles the main menu page functionality
type MainPage struct {
	tui interfaces.TUIInterface
}

// NewMainPage creates a new MainPage instance
func NewMainPage(tui interfaces.TUIInterface) *MainPage {
	return &MainPage{
		tui: tui,
	}
}

// CreateMainPage creates the main menu page
func (mp *MainPage) CreateMainPage() tview.Primitive {
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

	// Create the main menu list wit`h enhanced styling
	list := tview.NewList().
		AddItem("🔐 Vaults", "Manage and organize your vaults", 'v', func() {
			mp.tui.GetNavigation().ShowVaults()
		}).
		AddItem("👤 Profiles", "View Profile settings and Environments", 'p', func() {
			mp.tui.GetNavigation().ShowProfiles()
		}).
		AddItem("🌍 Environments", "Manage Environments", 'e', func() {
			mp.tui.GetNavigation().ShowEnvironments()
		}).
		AddItem("❓ Help", "View documentation and help", 'h', func() {
			mp.tui.GetNavigation().ShowHelp()
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

	// Update status bar with help text
	mp.tui.UpdateStatusBar("[yellow]↑/↓: Navigate | Enter: select[white]")
	return mp.tui.CreatePageLayout("Main Menu", content)
}

package pages

import (
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
	list := tview.NewList().
		AddItem("Vaults", "Manage secret vaults", 'v', func() {
			mp.tui.GetNavigation().ShowVaults()
		}).
		AddItem("Profiles", "Manage user profiles", 'p', func() {
			mp.tui.GetNavigation().ShowProfiles()
		}).
		AddItem("Environments", "Manage environments", 'e', func() {
			mp.tui.GetNavigation().ShowEnvironments()
		}).
		AddItem("Help", "Show help", 'h', func() {
			mp.tui.GetNavigation().ShowHelp()
		}).
		AddItem("Exit", "Exit application", 'q', func() {
			mp.tui.Quit()
		})

	return mp.tui.CreatePageLayout("Main Menu", list)
}

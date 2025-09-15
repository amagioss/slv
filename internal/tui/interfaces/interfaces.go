package interfaces

import (
	"github.com/rivo/tview"
)

// TUIInterface defines the interface that pages can use to interact with the TUI
type TUIInterface interface {
	// Navigation methods
	GetNavigation() NavigationInterface

	// UI components
	GetInfoBar() tview.Primitive

	// Page layout method
	CreatePageLayout(title string, content tview.Primitive) tview.Primitive

	// Application control
	Quit()
	GetApplication() *tview.Application
}

// NavigationInterface defines the interface for navigation functionality
type NavigationInterface interface {
	ShowMainMenu()
	ShowVaults()
	ShowVaultsReplace()
	ShowProfiles()
	ShowEnvironments()
	ShowHelp()
	ShowVaultDetails(vaultDetailsPage tview.Primitive)
	UpdateStatus()
	SetVaultDir(dir string)
	GetVaultDir() string
}

package interfaces

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NavigationInterface defines the interface for navigation functionality
type NavigationInterface interface {
	// Page navigation
	ShowMainMenu(replace bool)
	ShowVaults(replace bool)
	ShowProfiles(replace bool)
	ShowEnvironments(replace bool)
	ShowHelp(replace bool)
	ShowVaultDetails(replace bool)
	ShowNewVault(replace bool)

	// Routing and stack management
	GoBack()
	NavigateTo(pageName string)
	GetCurrentPage() string
	GetPageStack() []string
	ClearStack()

	// Input handling
	HandleEscape() *tcell.EventKey

	// Status and help
	UpdateStatus()
	GetStatusBar() tview.Primitive
	SetCustomHelp(helpText string)
	ClearCustomHelp()

	// Vault directory management
	GetVaultDir() string
	SetVaultDir(dir string)

	// App access
	GetApp() TUIInterface
}

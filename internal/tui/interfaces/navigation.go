package interfaces

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/vaults"
)

// NavigationInterface defines the interface for navigation functionality
type NavigationInterface interface {
	// Page navigation
	ShowMainMenu(replace bool)
	ShowVaults(replace bool)
	ShowProfiles(replace bool)
	ShowEnvironments(replace bool)
	ShowNewEnvironment(replace bool)
	ShowHelp(replace bool)
	ShowVaultDetails(replace bool)
	ShowNewVault(replace bool)

	// Parameterized page navigation
	ShowVaultsWithDir(dir string, replace bool)
	ShowVaultDetailsWithVault(vault *vaults.Vault, filePath string, replace bool)
	ShowVaultEditWithVault(vault *vaults.Vault, filePath string, replace bool)
	ShowNewVaultWithDir(dir string, replace bool)

	// Routing and stack management
	GoBack()
	NavigateTo(pageName string, replace bool)
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
	GetPageInstance(pageName string) (Page, bool)
	StorePageInstance(pageName string, page Page)
	RemovePageInstance(pageName string)

	// App access
	GetApp() TUIInterface

	// Page state management
	SavePageState(pageName, stateKey string, stateValue interface{})
	GetPageState(pageName, stateKey string) (interface{}, bool)
	ClearPageState(pageName string)
	ClearPageStateKey(pageName, stateKey string)
	HasPageState(pageName string) bool
}

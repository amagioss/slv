package interfaces

import (
	"github.com/rivo/tview"
)

// RouterInterface defines the interface for router functionality
type RouterInterface interface {
	// Basic tview.Pages operations
	AddPage(name string, page tview.Primitive, resize, visible bool)
	RemovePage(name string)
	HasPage(name string) bool
	GetPages() *tview.Pages

	// Page stack management
	PushPage(page string)
	PopPage() string
	ClearStack()
	GetCurrentPage() string
	SetCurrentPage(page string)
	GetPageStack() []string

	// Page interface support (legacy - for backward compatibility)
	RegisterPage(name string, page Page)
	GetRegisteredPage(name string) Page
	HasRegisteredPage(name string) bool
	GetRegisteredPageNames() []string

	// Page factory support
	RegisterPageFactory(name string, factory PageFactory)
	CreatePage(tui TUIInterface, name string, params ...interface{}) Page
	HasPageFactory(name string) bool
	GetPageFactoryNames() []string

	// Infrastructure methods (to avoid duplication in Navigation)
	AddPageToMainComponent(name string, page tview.Primitive, components ComponentManagerInterface)
	NavigateToPage(name string, components ComponentManagerInterface, replace bool)
	GoBackWithComponents(components ComponentManagerInterface) error
}

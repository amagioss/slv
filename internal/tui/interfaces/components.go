package interfaces

import (
	"github.com/rivo/tview"
)

// Component defines the interface for all UI components
type Component interface {
	Render() tview.Primitive
	Refresh()
}

// InfoBarComponent defines the interface for the info bar
type InfoBarComponent interface {
	Component
	UpdateProfile(profileName string)
	UpdateEnvironment(envName, envEmail, envType, publicKey string)
	ShowNoEnvironment()
}

// StatusBarComponent defines the interface for the status bar
type StatusBarComponent interface {
	Component
	UpdateStatus(pageName string)
	SetCustomHelp(helpText string)
	ClearCustomHelp()
	SetPageName(pageName string)
}

// MainContentComponent defines the interface for the main content area
type MainContentComponent interface {
	Component
	SetContent(content tview.Primitive)
	GetContent() tview.Primitive
}

// ComponentManagerInterface defines the interface for component management
type ComponentManagerInterface interface {
	InitializeComponents()
	GetInfoBar() tview.Primitive
	GetStatusBar() tview.Primitive
	GetMainContent() tview.Primitive
	GetMainContentPages() *tview.Pages
	UpdateStatus(page string)
	UpdateStatusBar(text string)
	ClearStatusBar()
	RefreshInfoBar()
}

// // RouterInterface defines the interface for router functionality
// type RouterInterface interface {
// 	AddPage(name string, page tview.Primitive, resize, visible bool)
// 	RemovePage(name string)
// 	SwitchToPage(name string)
// 	PushPage(page string)
// 	PopPage() string
// 	ClearStack()
// 	GetCurrentPage() string
// 	SetCurrentPage(page string)
// 	GetPageStack() []string
// }

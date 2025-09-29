package interfaces

import (
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/theme"
)

// TUIInterface defines the interface that pages can use to interact with the TUI
type TUIInterface interface {
	// Navigation methods
	GetNavigation() NavigationInterface

	// UI components
	GetInfoBar() tview.Primitive
	UpdateStatusBar(text string)
	ClearStatusBar()

	// Page layout method
	CreatePageLayout(title string, content tview.Primitive) tview.Primitive

	// Modal methods
	ShowError(message string)
	ShowInfo(message string)

	// Application control
	Quit()
	GetApplication() *tview.Application

	// Theme access
	GetTheme() *theme.Theme

	// Core components access
	GetComponents() ComponentManagerInterface
	GetRouter() RouterInterface
}

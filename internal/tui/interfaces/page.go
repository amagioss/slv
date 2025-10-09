package interfaces

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Page defines the interface that all pages must implement
type Page interface {
	// Create returns the tview.Primitive for this page
	Create() tview.Primitive

	// Refresh refreshes the page content (useful for dynamic content)
	Refresh()

	// HandleInput handles input events specific to this page
	HandleInput(event *tcell.EventKey) *tcell.EventKey

	// GetTitle returns the page title
	GetTitle() string

	// SaveNavigationState saves the current navigation state for restoration
	SaveNavigationState()

	// RestoreNavigationState restores the saved navigation state
	RestoreNavigationState()

	// ClearNavigationState clears the saved navigation state
	ClearNavigationState()
}

// PageFactory defines the interface for creating page instances
type PageFactory interface {
	// CreatePage creates a new page instance with the given parameters
	CreatePage(params ...interface{}) Page
}

// PageFactoryFunc is a function type that implements PageFactory
type PageFactoryFunc func(params ...interface{}) Page

// CreatePage implements the PageFactory interface for PageFactoryFunc
func (f PageFactoryFunc) CreatePage(params ...interface{}) Page {
	return f(params...)
}

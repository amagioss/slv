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
}

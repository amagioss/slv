package navigation

import (
	"github.com/rivo/tview"
)

// UpdateStatus updates the status bar
func (n *Navigation) UpdateStatus() {
	n.app.GetComponents().UpdateStatus(n.app.GetRouter().GetCurrentPage())
}

// GetStatusBar returns the status bar primitive
func (n *Navigation) GetStatusBar() tview.Primitive {
	return n.app.GetComponents().GetStatusBar()
}

// SetCustomHelp sets the custom help text for the current page
func (n *Navigation) SetCustomHelp(helpText string) {
	n.customHelp = helpText
	n.app.GetComponents().UpdateStatusBar(helpText)
}

// ClearCustomHelp clears the custom help text
func (n *Navigation) ClearCustomHelp() {
	n.customHelp = ""
	n.app.GetComponents().ClearStatusBar()
}

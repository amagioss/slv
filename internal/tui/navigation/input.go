package navigation

import (
	"github.com/gdamore/tcell/v2"
)

// HandleEscape handles the escape key based on current context
func (n *Navigation) HandleEscape() *tcell.EventKey {
	switch n.app.GetRouter().GetCurrentPage() {
	case "main":
		// On main page, escape quits the app
		n.app.Quit()
		return nil
	case "help":
		// On help page, escape goes back
		n.GoBack()
		return nil
	default:
		// On other pages, escape goes back
		n.GoBack()
		return nil
	}
}

// setupInputHandling sets up global input handling for the navigation
func (n *Navigation) setupInputHandling() {
	n.app.GetApplication().SetInputCapture(n.handleGlobalInput)
}

// handleGlobalInput handles global input events
func (n *Navigation) handleGlobalInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlC:
		n.app.Quit()
		return nil
	case tcell.KeyF1:
		n.ShowHelp(false)
		return nil
	case tcell.KeyEsc:
		return n.HandleEscape()
	}
	return event
}

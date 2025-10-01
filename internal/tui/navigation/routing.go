package navigation

import (
	"fmt"

	"github.com/rivo/tview"
)

// addPage adds a page to the navigation (now uses Router infrastructure)
func (n *Navigation) addPage(name string, page tview.Primitive) {
	n.app.GetRouter().AddPageToMainComponent(name, page, n.app.GetComponents())
}

// setCurrentPage sets the current active page (now uses Router infrastructure)
func (n *Navigation) setCurrentPage(name string, replace bool) {
	n.app.GetRouter().NavigateToPage(name, n.app.GetComponents(), replace)
	n.UpdateStatus()
}

// GoBack navigates to the previous page (now uses Router infrastructure)
func (n *Navigation) GoBack() {
	if err := n.app.GetRouter().GoBackWithComponents(n.app.GetComponents()); err != nil {
		// If no pages in stack, show error
		n.app.ShowError("No previous page to go back to")
		return
	}
	n.UpdateStatus()
}

// NavigateTo navigates to a specific page
func (n *Navigation) NavigateTo(pageName string, replace bool) {
	switch pageName {
	case "main":
		n.ShowMainMenu(replace)
	case "vaults":
		n.ShowVaults(replace)
	case "profiles":
		n.ShowProfiles(replace)
	case "environments":
		n.ShowEnvironments(replace)
	case "help":
		n.ShowHelp(replace)
	case "new-vault":
		n.ShowNewVault(replace)
	default:
		n.app.ShowError(fmt.Sprintf("Unknown page: %s", pageName))
	}
}

// GetCurrentPage returns the current page name
func (n *Navigation) GetCurrentPage() string {
	return n.app.GetRouter().GetCurrentPage()
}

// GetPageStack returns the page stack
func (n *Navigation) GetPageStack() []string {
	return n.app.GetRouter().GetPageStack()
}

// ClearStack clears the page stack
func (n *Navigation) ClearStack() {
	n.app.GetRouter().ClearStack()
}

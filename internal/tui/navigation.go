package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/pages"
)

// Navigation handles page navigation and routing
type Navigation struct {
	app         *TUI
	currentPage string
	pageStack   []string
	statusBar   *tview.TextView
}

// NewNavigation creates a new navigation controller
func NewNavigation(app interface{}) *Navigation {
	nav := &Navigation{
		app:       app.(*TUI),
		pageStack: make([]string, 0),
	}

	nav.createStatusBar()
	nav.UpdateStatus() // Initialize status bar with content
	return nav
}

// createStatusBar creates the status bar
func (n *Navigation) createStatusBar() {
	n.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(false).
		SetTextAlign(tview.AlignLeft).
		SetScrollable(false) // Prevent scrolling

	n.statusBar.SetBorder(true). // Add border for status bar
					SetBorderColor(tcell.ColorAqua).
					SetTitle("Status").
					SetTitleAlign(tview.AlignLeft)
}

// ShowMainMenu displays the main menu
func (n *Navigation) ShowMainMenu() {
	// Create MainPage using the pages package
	mainPage := pages.NewMainPage(n.app)
	menu := mainPage.CreateMainPage()
	n.addPage("main", menu)
	n.setCurrentPage("main")
}

// ShowVaults displays the vaults page
func (n *Navigation) ShowVaults() {
	vaults := n.app.createVaultsPage()
	n.addPage("vaults", vaults)
	n.setCurrentPage("vaults")
	n.UpdateStatus()
}

// ShowProfiles displays the profiles page
func (n *Navigation) ShowProfiles() {
	profiles := n.app.createProfilesPage()
	n.addPage("profiles", profiles)
	n.setCurrentPage("profiles")
	n.UpdateStatus()
}

// ShowEnvironments displays the environments page
func (n *Navigation) ShowEnvironments() {
	environments := n.app.createEnvironmentsPage()
	n.addPage("environments", environments)
	n.setCurrentPage("environments")
	n.UpdateStatus()
}

// ShowHelp displays the help page
func (n *Navigation) ShowHelp() {
	help := n.app.createHelpPage()
	n.addPage("help", help)
	n.setCurrentPage("help")
	n.UpdateStatus()
}

// addPage adds a page to the navigation
func (n *Navigation) addPage(name string, page tview.Primitive) {
	n.app.GetPages().AddPage(name, page, true, false)
}

// setCurrentPage sets the current active page
func (n *Navigation) setCurrentPage(name string) {
	if n.currentPage != "" {
		n.pageStack = append(n.pageStack, n.currentPage)
	}

	n.currentPage = name
	n.app.GetPages().SwitchToPage(name)
	n.UpdateStatus()
}

// GoBack navigates to the previous page
func (n *Navigation) GoBack() {
	if len(n.pageStack) > 0 {
		previousPage := n.pageStack[len(n.pageStack)-1]
		n.pageStack = n.pageStack[:len(n.pageStack)-1]
		n.currentPage = previousPage
		n.app.GetPages().SwitchToPage(previousPage)
		n.UpdateStatus()
	}
}

// HandleEscape handles the escape key based on current context
func (n *Navigation) HandleEscape() *tcell.EventKey {
	switch n.currentPage {
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

// UpdateStatus updates the status bar
func (n *Navigation) UpdateStatus() {
	if n.statusBar == nil {
		return
	}

	status := fmt.Sprintf("[white]Page: [cyan]%s[white] | Stack: [cyan]%d[white] | F1: Help | Esc: Back | Ctrl+C: Quit",
		n.currentPage, len(n.pageStack))

	n.statusBar.SetText(status)
}

// GetStatusBar returns the status bar primitive
func (n *Navigation) GetStatusBar() tview.Primitive {
	return n.statusBar
}

// GetCurrentPage returns the current page name
func (n *Navigation) GetCurrentPage() string {
	return n.currentPage
}

// GetPageStack returns the page stack
func (n *Navigation) GetPageStack() []string {
	return n.pageStack
}

// ClearStack clears the page stack
func (n *Navigation) ClearStack() {
	n.pageStack = make([]string, 0)
}

// NavigateTo navigates to a specific page
func (n *Navigation) NavigateTo(pageName string) {
	switch pageName {
	case "main":
		n.ShowMainMenu()
	case "vaults":
		n.ShowVaults()
	case "profiles":
		n.ShowProfiles()
	case "environments":
		n.ShowEnvironments()
	case "help":
		n.ShowHelp()
	default:
		n.app.ShowError(fmt.Sprintf("Unknown page: %s", pageName))
	}
}

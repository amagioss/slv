package app

import (
	"context"
	"log"
	"os"

	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/components"
	"slv.sh/slv/internal/tui/core"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/navigation"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/pages/environments"
	"slv.sh/slv/internal/tui/pages/help"
	"slv.sh/slv/internal/tui/pages/mainpage"
	"slv.sh/slv/internal/tui/pages/profiles"
	"slv.sh/slv/internal/tui/pages/vault_browse"
	"slv.sh/slv/internal/tui/pages/vault_new"
	"slv.sh/slv/internal/tui/pages/vault_view"
	"slv.sh/slv/internal/tui/theme"
)

// TUI represents the main TUI application
type TUI struct {
	app        *core.Application
	navigation *navigation.Navigation
	components *components.ComponentManager
}

// NewTUI creates a new TUI instance
func NewTUI() *TUI {
	themeInstance := theme.NewTheme()
	tui := &TUI{
		app: core.NewApplication(themeInstance),
	}

	tui.setup()

	return tui
}

func (t *TUI) setup() {
	t.app.GetTheme().ApplyTheme(t.app.GetApplication())
	t.app.GetRouter().GetPages().SetBackgroundColor(t.app.GetTheme().GetBackground())

	// Initialize components
	t.components = components.NewComponentManager(t)
	t.components.InitializeComponents()

	t.app.GetApplication().SetTitle("SLV - Terminal User Interface")
	t.setNavigator()

	// Setup pages with the router (after navigation is initialized)
	t.setupPages()

	// Set up layout manager with components
	t.app.GetLayoutManager().SetInfoBar(t.components.GetInfoBar())
	t.app.GetLayoutManager().SetContent(t.components.GetMainContent())
	t.app.GetLayoutManager().SetStatusBar(t.components.GetStatusBar())
	rootLayout := t.app.GetLayoutManager().BuildLayout()

	t.app.GetApplication().SetRoot(rootLayout, true)
	t.navigation.ShowMainMenu(false)
}

// setNavigator initializes the navigation after theme is applied
func (t *TUI) setNavigator() {
	t.navigation = navigation.NewNavigation(t)
}

// setupPages registers all pages with the router at startup
func (t *TUI) setupPages() {
	// Register main page
	mainPage := mainpage.NewMainPage(t)
	t.app.GetRouter().RegisterPage("main", mainPage)

	// Register profiles page
	profilesPage := profiles.NewProfilesPage(t)
	t.app.GetRouter().RegisterPage("profiles", profilesPage)

	// Register environments page
	environmentsPage := environments.NewEnvironmentsPage(t)
	t.app.GetRouter().RegisterPage("environments", environmentsPage)

	// Register help page
	helpPage := help.NewHelpPage(t)
	t.app.GetRouter().RegisterPage("help", helpPage)

	// Register vault browse page (will be created with current directory when needed)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.LogError(err, true)
		homeDir = "."
	}

	vaultBrowsePage := vault_browse.NewVaultBrowsePage(t, homeDir)
	t.app.GetRouter().RegisterPage("vaults_browse", vaultBrowsePage)

	// Register vault new page (will be created with current directory when needed)
	vaultNewPage := vault_new.NewVaultNewPage(t, homeDir)
	t.app.GetRouter().RegisterPage("vaults_new", vaultNewPage)

	// Register vault view page (will be created with current directory when needed)
	vaultViewPage := vault_view.NewVaultViewPage(t, nil, "")
	t.app.GetRouter().RegisterPage("vaults_view", vaultViewPage)

}

// Run starts the TUI
func (t *TUI) Run() error {
	return t.app.Run()
}

// Quit exits the TUI
func (t *TUI) Quit() {
	t.app.Stop()
}

// GetApplication returns the tview application
func (t *TUI) GetApplication() *tview.Application {
	return t.app.GetApplication()
}

// GetPages returns the pages container
func (t *TUI) GetPages() *tview.Pages {
	return t.app.GetRouter().GetPages()
}

// GetTheme returns the theme
func (t *TUI) GetTheme() *theme.Theme {
	return t.app.GetTheme().(*theme.Theme)
}

// GetInfoBar returns the info bar primitive
func (t *TUI) GetInfoBar() tview.Primitive {
	return t.components.GetInfoBar()
}

// GetComponents returns the component manager
func (t *TUI) GetComponents() interfaces.ComponentManagerInterface {
	return t.components
}

// GetNavigation returns the navigation interface
func (t *TUI) GetNavigation() interfaces.NavigationInterface {
	return t.navigation
}

// GetContext returns the context
func (t *TUI) GetContext() context.Context {
	return t.app.GetContext()
}

// GetRouter returns the router
func (t *TUI) GetRouter() interfaces.RouterInterface {
	return t.app.GetRouter()
}

// ShowError shows an error modal
func (t *TUI) ShowError(message string) {
	modal := tview.NewModal().
		SetText("Error: " + message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(int, string) {
			t.components.GetMainContentPages().RemovePage("error")
		})

	// Add modal to the main content pages
	t.components.GetMainContentPages().AddPage("error", modal, true, true)
}

// ShowInfo shows an info modal
func (t *TUI) ShowInfo(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(int, string) {
			t.components.GetMainContentPages().RemovePage("info")
		})

	// Add modal to the main content pages
	t.components.GetMainContentPages().AddPage("info", modal, true, true)
}

// LogError logs an error
func (t *TUI) LogError(err error, showToUser bool) {
	log.Printf("TUI Error: %v", err)
	if showToUser {
		t.ShowError(err.Error())
	}
}

// CreatePageLayout creates a page layout (renamed from createPageLayout for interface)
func (t *TUI) CreatePageLayout(title string, content tview.Primitive) tview.Primitive {
	return pages.CreatePageLayout(t, title, content)
}

func (t *TUI) UpdateStatusBar(helpText string) {
	// Use the component manager to update status bar
	t.components.UpdateStatusBar(helpText)
}

func (t *TUI) ClearStatusBar() {
	// Use the component manager to clear status bar
	t.components.ClearStatusBar()
}

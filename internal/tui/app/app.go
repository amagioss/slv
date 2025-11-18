package app

import (
	"context"
	"log"

	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/tui/components"
	"slv.sh/slv/internal/tui/core"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/navigation"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/pages/environments"
	environments_new "slv.sh/slv/internal/tui/pages/environments_new"
	"slv.sh/slv/internal/tui/pages/help"
	"slv.sh/slv/internal/tui/pages/mainpage"
	"slv.sh/slv/internal/tui/pages/profiles"
	"slv.sh/slv/internal/tui/pages/vault_browse"
	"slv.sh/slv/internal/tui/pages/vault_edit"
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

// setupPages registers all page factories with the router at startup
func (t *TUI) setupPages() {
	// Register main page factory
	t.app.GetRouter().RegisterPageFactory("main", interfaces.PageFactoryFunc(func(params ...interface{}) interfaces.Page {
		tui := params[0].(interfaces.TUIInterface)
		return mainpage.NewMainPage(tui)
	}))

	// Register profiles page factory
	t.app.GetRouter().RegisterPageFactory("profiles", interfaces.PageFactoryFunc(func(params ...interface{}) interfaces.Page {
		tui := params[0].(interfaces.TUIInterface)
		return profiles.NewProfilesPage(tui)
	}))

	// Register environments page factory
	t.app.GetRouter().RegisterPageFactory("environments", interfaces.PageFactoryFunc(func(params ...interface{}) interfaces.Page {
		tui := params[0].(interfaces.TUIInterface)
		return environments.NewEnvironmentsPage(tui)
	}))

	// Register new environment page factory
	t.app.GetRouter().RegisterPageFactory("environments_new", interfaces.PageFactoryFunc(func(params ...interface{}) interfaces.Page {
		tui := params[0].(interfaces.TUIInterface)
		return environments_new.NewEnvironmentNewPage(tui)
	}))

	// Register help page factory
	t.app.GetRouter().RegisterPageFactory("help", interfaces.PageFactoryFunc(func(params ...interface{}) interfaces.Page {
		tui := params[0].(interfaces.TUIInterface)
		return help.NewHelpPage(tui)
	}))

	// Register vault browse page factory (takes directory as parameter)
	t.app.GetRouter().RegisterPageFactory("vaults_browse", interfaces.PageFactoryFunc(func(params ...interface{}) interfaces.Page {
		tui := params[0].(interfaces.TUIInterface)
		dir := params[1].(string)
		return vault_browse.NewVaultBrowsePage(tui, dir)
	}))

	// Register vault new page factory (takes directory as parameter)
	t.app.GetRouter().RegisterPageFactory("vaults_new", interfaces.PageFactoryFunc(func(params ...interface{}) interfaces.Page {
		tui := params[0].(interfaces.TUIInterface)
		dir := params[1].(string)
		return vault_new.NewVaultNewPage(tui, dir)
	}))

	// Register vault view page factory (takes vault and filepath as parameters)
	t.app.GetRouter().RegisterPageFactory("vaults_view", interfaces.PageFactoryFunc(func(params ...interface{}) interfaces.Page {
		tui := params[0].(interfaces.TUIInterface)
		var vault *vaults.Vault = nil
		if len(params) > 1 && params[1] != nil {
			vault = params[1].(*vaults.Vault)
		}
		filePath := ""
		if len(params) > 2 {
			filePath = params[2].(string)
		}
		return vault_view.NewVaultViewPage(tui, vault, filePath)
	}))

	// Register vault edit page factory (takes vault and filepath as parameters)
	t.app.GetRouter().RegisterPageFactory("vaults_edit", interfaces.PageFactoryFunc(func(params ...interface{}) interfaces.Page {
		tui := params[0].(interfaces.TUIInterface)
		var vault *vaults.Vault = nil
		if len(params) > 1 && params[1] != nil {
			vault = params[1].(*vaults.Vault)
		}
		filePath := ""
		if len(params) > 2 {
			filePath = params[2].(string)
		}
		return vault_edit.NewVaultEditPage(tui, vault, filePath)
	}))
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

// ShowConfirmation shows a confirmation modal
func (t *TUI) ShowConfirmation(message string, onConfirm func(), onCancel func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Cancel", "Yes"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			t.components.GetMainContentPages().RemovePage("confirmation")
			if buttonLabel == "Yes" && onConfirm != nil {
				onConfirm()
			} else if buttonLabel == "Cancel" && onCancel != nil {
				onCancel()
			}
		})

	// Add modal to the main content pages
	t.components.GetMainContentPages().AddPage("confirmation", modal, true, true)
}

// ShowConfirmationWithFocus shows a confirmation modal with focus restoration
func (t *TUI) ShowConfirmationWithFocus(message string, confirmButtonText string, cancelButtonText string, onConfirm func(), onCancel func(), restoreFocus func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{cancelButtonText, confirmButtonText}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			t.components.GetMainContentPages().RemovePage("confirmation")

			// Restore focus after modal is removed
			if restoreFocus != nil {
				restoreFocus()
			}

			if buttonLabel == confirmButtonText && onConfirm != nil {
				onConfirm()
			} else if buttonLabel == cancelButtonText && onCancel != nil {
				onCancel()
			}
		})

	// Add modal to the main content pages
	t.components.GetMainContentPages().AddPage("confirmation", modal, true, true)
}

// ShowModalForm shows a modal form
func (t *TUI) ShowModalForm(title string, form *tview.Form, confirmButtonText string, cancelButtonText string, onConfirm func(), onCancel func(), restoreFocus func()) {
	// Style the form for better button alignment
	form.SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignCenter)

	// Add buttons to the form for Cancel and Add
	form.AddButton(cancelButtonText, func() {
		t.components.GetMainContentPages().RemovePage("modal-form")

		// Restore focus after modal is removed
		if restoreFocus != nil {
			restoreFocus()
		}

		if onCancel != nil {
			onCancel()
		}
	})

	form.AddButton(confirmButtonText, func() {
		t.components.GetMainContentPages().RemovePage("modal-form")

		// Restore focus after modal is removed
		if restoreFocus != nil {
			restoreFocus()
		}

		if onConfirm != nil {
			onConfirm()
		}
	})

	// Set button alignment to center after adding buttons
	form.SetButtonsAlign(tview.AlignCenter)

	// Create a centered modal-like container with more generous spacing
	modalContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false). // Top spacer
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false). // Left spacer
			AddItem(form, 0, 1, true). // Form in center
			AddItem(nil, 0, 1, false), // Right spacer
						0, 1, true). // Form row - let form determine its own height
		AddItem(nil, 0, 1, false) // Bottom spacer

	// Add modal to the main content pages
	t.components.GetMainContentPages().AddPage("modal-form", modalContainer, true, true)
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

package app

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/tui/components"
	"slv.sh/slv/internal/tui/core"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/navigation"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/pages/environments"
	"slv.sh/slv/internal/tui/pages/environments_new"
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
	rootLayout tview.Primitive // Store root layout for splash screen transition
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
	// Initialize clipboard
	if err := clipboard.Init(); err != nil {
		log.Printf("Failed to initialize clipboard: %v", err)
	}

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

	// Build the root layout now so it's ready when we need it
	t.rootLayout = t.app.GetLayoutManager().BuildLayout()

	// Show splash screen first, then main menu after delay
	// The splash screen will set itself as root, then transition to main menu
	t.showSplashScreen()
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
		return environments_new.NewNewEnvironmentPage(tui)
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

// ShowModal shows a generic modal with any content
func (t *TUI) ShowModal(title string, content tview.Primitive, restoreFocus func()) {
	// Create a centered modal-like container
	modalContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false). // Top spacer
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false).    // Left spacer
			AddItem(content, 0, 1, true). // Content in center
			AddItem(nil, 0, 1, false),    // Right spacer
						0, 1, true). // Content row - let content determine its own height
		AddItem(nil, 0, 1, false) // Bottom spacer

	// Handle Escape key to close modal
	modalContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			t.components.GetMainContentPages().RemovePage("modal")
			if restoreFocus != nil {
				restoreFocus()
			}
			return nil
		}
		return event
	})

	// Add modal to the main content pages
	t.components.GetMainContentPages().AddPage("modal", modalContainer, true, true)
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

// showSplashScreen displays the SLV logo for a couple of seconds before showing the main menu
func (t *TUI) showSplashScreen() {
	colors := theme.GetCurrentPalette()

	art := config.Art()
	coloredArt := strings.ReplaceAll(art, "▓", "[#9d3a4f]▓[-]")
	coloredArt = strings.ReplaceAll(coloredArt, "░", "[#4f5559]░[-]")
	coloredArt = strings.ReplaceAll(coloredArt, "▒", "[#4f5559]▒[-]")
	// Split colored art into lines for animation
	coloredArtLines := strings.Split(coloredArt, "\n")

	// Determine art dimensions to size the splash grid appropriately.
	artLines := strings.Split(art, "\n")
	artHeight := len(artLines)
	maxWidth := 0
	for _, line := range artLines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	logoText := tview.NewTextView()
	logoText.SetText("") // Start empty for animation
	logoText.SetDynamicColors(true)
	logoText.SetTextAlign(tview.AlignCenter)
	logoText.SetWrap(false)
	logoText.SetTextColor(colors.Primary)
	logoText.SetBackgroundColor(colors.Background)

	subtitleText := tview.NewTextView()
	subtitleText.SetText("[#f2f2f2]Secure Local Vault[-]")
	subtitleText.SetDynamicColors(true)
	subtitleText.SetTextAlign(tview.AlignCenter)
	subtitleText.SetTextColor(colors.Primary)
	subtitleText.SetBackgroundColor(colors.Background)

	// Use a grid layout to center the logo both horizontally and vertically.
	// Rows/columns with size 0 grow to fill remaining space, while the middle
	// row/column reserves enough space for the logo itself.
	splashGrid := tview.NewGrid()
	splashGrid.SetRows(0, artHeight+2, 1, 0)
	splashGrid.SetColumns(0, maxWidth+4, 0)
	splashGrid.SetBorder(true)
	splashGrid.SetBorderColor(colors.Border)
	splashGrid.SetBackgroundColor(colors.Background)

	splashGrid.AddItem(logoText, 1, 1, 1, 1, artHeight+2, maxWidth+4, false)
	splashGrid.AddItem(subtitleText, 2, 1, 1, 1, 1, maxWidth+4, false)

	app := t.app.GetApplication()

	// Set application background to ensure full screen background
	// tview.Styles.PrimitiveBackgroundColor = colors.Background

	// Fill the entire screen with black before every draw while splash is active.
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		style := tcell.StyleDefault.Background(colors.Background)
		screen.Fill(' ', style)
		return false
	})

	// Helper to transition from splash to main menu.
	finishSplash := func() {
		app.SetBeforeDrawFunc(nil)
		if t.rootLayout == nil {
			t.rootLayout = t.app.GetLayoutManager().BuildLayout()
		}
		app.SetRoot(t.rootLayout, true)
		t.navigation.ShowMainMenu(false)
	}

	// 1) Set splash as root
	app.SetRoot(splashGrid, true)

	// 2) Animate logo line by line, then switch to main layout
	go func() {
		// Animation duration
		lineDelay := 40 * time.Millisecond
		currentText := ""

		for _, line := range coloredArtLines {
			time.Sleep(lineDelay)
			currentText += line + "\n"

			// Capture current text for the closure
			textToUpdate := currentText
			app.QueueUpdateDraw(func() {
				logoText.SetText(textToUpdate)
			})
		}

		// Hold for a bit after animation completes
		time.Sleep(500 * time.Millisecond)
		app.QueueUpdateDraw(finishSplash)
	}()

	// 3) Allow any key to skip splash immediately
	splashGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		app.QueueUpdateDraw(finishSplash)
		return nil // consume key
	})
}

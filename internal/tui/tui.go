package tui

import (
	"context"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TUI represents the main TUI application
type TUI struct {
	app        *tview.Application
	pages      *tview.Pages
	theme      *Theme
	navigation *Navigation
	infoBar    tview.Primitive
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewTUI creates a new TUI instance
func NewTUI() *TUI {
	ctx, cancel := context.WithCancel(context.Background())

	tui := &TUI{
		app:    tview.NewApplication(),
		pages:  tview.NewPages(),
		theme:  NewTheme(),
		ctx:    ctx,
		cancel: cancel,
	}

	tui.setup()

	return tui
}

func (t *TUI) setup() {
	t.theme.ApplyTheme(t.app)
	t.pages.SetBackgroundColor(t.theme.Background)

	t.createInfoBar()
	t.app.SetTitle("SLV - Terminal User Interface")
	t.setNavigator()

	statusBar := t.navigation.GetStatusBar()
	if statusTextView, ok := statusBar.(*tview.TextView); ok {
		statusTextView.SetWrap(false) // Prevent text wrapping
	}

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.infoBar, 8, 1, false). // Info bar
		AddItem(t.pages, 0, 1, true).    // Pages
		AddItem(statusBar, 3, 1, false)  // Status bar

	t.app.SetRoot(flex, true)
	// t.app.SetRoot(t.pages, true)
	t.app.SetInputCapture(t.handleInput)
	t.navigation.ShowMainMenu()
}

// setNavigator initializes the navigation after theme is applied
func (t *TUI) setNavigator() {
	t.navigation = NewNavigation(t)
}

// handleInput handles global input
func (t *TUI) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlC:
		t.Quit()
		return nil
	case tcell.KeyF1:
		t.navigation.ShowHelp()
		return nil
	case tcell.KeyEsc:
		return t.navigation.HandleEscape()
	}
	return event
}

// Run starts the TUI
func (t *TUI) Run() error {
	defer t.cancel()
	return t.app.Run()
}

// Quit exits the TUI
func (t *TUI) Quit() {
	t.cancel()
	t.app.Stop()
}

// GetApplication returns the tview application
func (t *TUI) GetApplication() *tview.Application {
	return t.app
}

// GetPages returns the pages container
func (t *TUI) GetPages() *tview.Pages {
	return t.pages
}

// GetTheme returns the theme
func (t *TUI) GetTheme() *Theme {
	return t.theme
}

// GetInfoBar returns the info bar primitive
func (t *TUI) GetInfoBar() tview.Primitive {
	return t.infoBar
}

// GetNavigation returns the navigation interface
func (t *TUI) GetNavigation() NavigationInterface {
	return t.navigation
}

// GetContext returns the context
func (t *TUI) GetContext() context.Context {
	return t.ctx
}

// ShowError shows an error modal
func (t *TUI) ShowError(message string) {
	modal := tview.NewModal().
		SetText("Error: " + message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(int, string) {
			t.pages.RemovePage("error")
		})
	t.pages.AddPage("error", modal, true, true)
}

// ShowInfo shows an info modal
func (t *TUI) ShowInfo(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(int, string) {
			t.pages.RemovePage("info")
		})
	t.pages.AddPage("info", modal, true, true)
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
	return t.createPageLayout(title, content)
}

func (t *TUI) UpdateStatusBar(helpText string) {
	// Use the navigation's SetCustomHelp method which handles the status bar update
	t.navigation.SetCustomHelp(helpText)
}

func (t *TUI) ClearStatusBar() {
	// Use the navigation's ClearCustomHelp method which handles the status bar update
	t.navigation.ClearCustomHelp()
}

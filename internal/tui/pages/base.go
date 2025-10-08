package pages

import (
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/theme"
)

// BasePage provides common functionality for all pages
type BasePage struct {
	tui   interfaces.TUIInterface
	title string
}

// NewBasePage creates a new BasePage instance
func NewBasePage(tui interfaces.TUIInterface, title string) *BasePage {
	return &BasePage{
		tui:   tui,
		title: title,
	}
}

// CreateLayout creates a common layout for the page content
func (bp *BasePage) CreateLayout(content tview.Primitive) tview.Primitive {
	return CreatePageLayout(bp.tui, bp.title, content)
}

// CreatePageLayout creates a common page layout with title and border
func CreatePageLayout(tui interfaces.TUIInterface, title string, content tview.Primitive) tview.Primitive {
	colors := theme.GetCurrentPalette()

	// Create a flex container
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Set border properties
	flex.SetBorder(true).
		SetBorderColor(colors.Border).
		SetTitle(title).
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(colors.MainContentTitle)

	// Add the content to the flex
	flex.AddItem(content, 0, 1, true)

	return flex
}

// ShowError displays an error modal
func (bp *BasePage) ShowError(message string) {
	bp.tui.ShowError(message)
}

// ShowInfo displays an info modal
func (bp *BasePage) ShowInfo(message string) {
	bp.tui.ShowInfo(message)
}

// UpdateStatus updates the status bar
func (bp *BasePage) UpdateStatus(text string) {
	bp.tui.UpdateStatusBar(text)
}

// ClearStatus clears the status bar
func (bp *BasePage) ClearStatus() {
	bp.tui.ClearStatusBar()
}

// GetTheme returns the current theme
func (bp *BasePage) GetTheme() *theme.Theme {
	return bp.tui.GetTheme()
}

// GetTitle returns the page title
func (bp *BasePage) GetTitle() string {
	return bp.title
}

// SetTitle sets the page title
func (bp *BasePage) SetTitle(title string) {
	bp.title = title
}

// GetTUI returns the TUI interface
func (bp *BasePage) GetTUI() interfaces.TUIInterface {
	return bp.tui
}

// NavigateTo navigates to another page
func (bp *BasePage) NavigateTo(pageName string, replace bool) {
	bp.tui.GetNavigation().NavigateTo(pageName, replace)
}

// GoBack goes back to the previous page
func (bp *BasePage) GoBack() {
	bp.tui.GetNavigation().GoBack()
}

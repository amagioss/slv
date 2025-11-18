package mainpage

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FormNavigation struct {
	mp           *MainPage
	currentFocus int
	focusGroup   []tview.Primitive
	helpTexts    map[tview.Primitive]string
}

func (fn *FormNavigation) NewFormNavigation(mp *MainPage) *FormNavigation {
	focusGroup := []tview.Primitive{
		mp.list, // Main menu list
	}

	return &FormNavigation{
		mp:           mp,
		currentFocus: 0,
		focusGroup:   focusGroup,
		helpTexts:    make(map[tview.Primitive]string),
	}
}

func (fn *FormNavigation) SetupNavigation() {
	// Set up help texts
	fn.setupHelpTexts()

	// Set up input capture for the list
	fn.mp.list.SetInputCapture(fn.handleListInputCapture)

	// Set initial focus to list
	fn.mp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])

	// Set initial help text
	fn.updateHelpText()
}

// setupHelpTexts sets up help text for each component
func (fn *FormNavigation) setupHelpTexts() {
	fn.helpTexts[fn.mp.list] = "Main Menu: ↑/↓ Navigate | Enter: Select | v: Vaults | p: Profiles | e: Environments | h: Help"
}

// updateHelpText updates the status bar with help text for the currently focused component
func (fn *FormNavigation) updateHelpText() {
	if fn.currentFocus >= 0 && fn.currentFocus < len(fn.focusGroup) {
		currentComponent := fn.focusGroup[fn.currentFocus]
		if helpText, exists := fn.helpTexts[currentComponent]; exists {
			fn.mp.GetTUI().UpdateStatusBar(helpText)
		}
	}
}

// handleListInputCapture handles input for the main menu list
func (fn *FormNavigation) handleListInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'v', 'V':
			// Navigate to vaults
			fn.mp.NavigateTo("vaults", false)
			return nil
		case 'p', 'P':
			// Navigate to profiles
			fn.mp.NavigateTo("profiles", false)
			return nil
		case 'e', 'E':
			// Navigate to environments
			fn.mp.NavigateTo("environments", false)
			return nil
		case 'h', 'H':
			// Navigate to help
			fn.mp.NavigateTo("help", false)
			return nil
		}
	}
	// Let other keys pass through (arrow keys, enter, etc. are handled by tview.List)
	return event
}

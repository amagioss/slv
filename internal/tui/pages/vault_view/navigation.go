package vault_view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FormNavigation struct {
	vvp          *VaultViewPage
	currentFocus int
	focusGroup   []tview.Primitive
	helpTexts    map[tview.Primitive]string // Component-specific help texts
}

func (fn *FormNavigation) NewFormNavigation(vvp *VaultViewPage) *FormNavigation {
	focusGroup := []tview.Primitive{
		vvp.vaultDetailsTable,
		vvp.accessorsTable,
		vvp.itemsTable,
	}
	intialFocus := 0

	return &FormNavigation{
		vvp:          vvp,
		currentFocus: intialFocus,
		focusGroup:   focusGroup,
		helpTexts:    make(map[tview.Primitive]string),
	}
}

func (fn *FormNavigation) resetSelectable() {
	for _, table := range fn.focusGroup {
		table.(*tview.Table).SetSelectable(false, false)
	}
}

func (fn *FormNavigation) SetupNavigation() {
	// Set up help texts for each component
	fn.setupHelpTexts()

	fn.vvp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
	fn.resetSelectable()
	fn.focusGroup[fn.currentFocus].(*tview.Table).SetSelectable(true, false)
	fn.vvp.mainFlex.SetInputCapture(fn.handleInputCapture)

	// Set initial help text
	fn.updateHelpText()
}

func (fn *FormNavigation) ShiftFocusForward() {
	fn.currentFocus = (fn.currentFocus + 1) % len(fn.focusGroup)
	fn.resetSelectable()
	fn.focusGroup[fn.currentFocus].(*tview.Table).SetSelectable(true, false)
	fn.vvp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
	fn.updateHelpText()
}

func (fn *FormNavigation) ShiftFocusBackward() {
	fn.currentFocus = (fn.currentFocus - 1 + len(fn.focusGroup)) % len(fn.focusGroup)
	fn.resetSelectable()
	fn.focusGroup[fn.currentFocus].(*tview.Table).SetSelectable(true, false)
	fn.vvp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
	fn.updateHelpText()
}

func (fn *FormNavigation) handleInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	} else {
		switch event.Key() {
		case tcell.KeyTab:
			// Switch focus between tables
			fn.ShiftFocusForward()
			return nil
		case tcell.KeyBacktab:
			// Switch focus between tables
			fn.ShiftFocusBackward()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				fn.vvp.GetTUI().GetNavigation().ShowVaults(false)
				return nil
			case 'u', 'U':
				// Unlock vault
				if fn.vvp.filePath != "" {
					fn.vvp.unlockVault()
				}
				return nil
			case 'l', 'L':
				// Lock vault
				if fn.vvp.filePath != "" {
					fn.vvp.lockVault()
				}
				return nil
			case 'r', 'R':
				// Reload vault
				if fn.vvp.filePath != "" {
					fn.vvp.reloadVault()
					fn.vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(fn.vvp.vault, fn.vvp.filePath, false)
				}
				return nil
			}
		case tcell.KeyEsc:
			fn.vvp.GetTUI().GetNavigation().ShowVaults(false)
			return nil
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight, tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyHome, tcell.KeyEnd:
			// Allow arrow keys and page keys to scroll
			return event
		}
		return event
	}
}

// setupHelpTexts sets up help text for each component
func (fn *FormNavigation) setupHelpTexts() {
	fn.helpTexts[fn.vvp.vaultDetailsTable] = "[yellow]Vault Details: ↑/↓: Navigate rows | Tab: Next table | q: Close | u: Unlock | l: Lock | r: Reload[white]"
	fn.helpTexts[fn.vvp.accessorsTable] = "[yellow]Accessors: ↑/↓: Navigate rows | Tab: Next table | q: Close | u: Unlock | l: Lock | r: Reload[white]"
	fn.helpTexts[fn.vvp.itemsTable] = "[yellow]Items: ↑/↓: Navigate rows | Tab: Next table | q: Close | u: Unlock | l: Lock | r: Reload[white]"
}

// updateHelpText updates the status bar with help text for the currently focused component
func (fn *FormNavigation) updateHelpText() {
	if fn.currentFocus >= 0 && fn.currentFocus < len(fn.focusGroup) {
		currentComponent := fn.focusGroup[fn.currentFocus]
		if helpText, exists := fn.helpTexts[currentComponent]; exists {
			fn.vvp.GetTUI().UpdateStatusBar(helpText)
		}
	}
}

// SetComponentHelpText sets help text for a specific component
func (fn *FormNavigation) SetComponentHelpText(component tview.Primitive, helpText string) {
	fn.helpTexts[component] = helpText
}

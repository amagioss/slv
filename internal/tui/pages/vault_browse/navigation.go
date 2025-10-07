package vault_browse

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FormNavigation struct {
	vbp          *VaultBrowsePage
	currentFocus int
	focusGroup   []tview.Primitive
	helpTexts    map[tview.Primitive]string // Component-specific help texts
}

func (fn *FormNavigation) NewFormNavigation(vbp *VaultBrowsePage) *FormNavigation {
	focusGroup := []tview.Primitive{
		vbp.fileList,
	}

	return &FormNavigation{
		vbp:          vbp,
		currentFocus: 0,
		focusGroup:   focusGroup,
		helpTexts:    make(map[tview.Primitive]string),
	}
}

func (fn *FormNavigation) SetupNavigation() {
	// Set up help texts for each component
	fn.setupHelpTexts()

	// Set up input capture for the file list
	fn.vbp.fileList.SetInputCapture(fn.handleInputCapture)
	fn.vbp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])

	// Set initial help text
	fn.updateHelpText()
}

// setupHelpTexts sets up help text for each component
func (fn *FormNavigation) setupHelpTexts() {
	fn.helpTexts[fn.vbp.fileList] = "[yellow]File Browser: ↑/↓: Navigate | →: Open vault/directory | ←: Go back | Ctrl+N: New vault | Ctrl+E: Edit vault | Ctrl+R: Rename vault | Ctrl+D: Delete vault[white]"
}

// updateHelpText updates the status bar with help text for the currently focused component
func (fn *FormNavigation) updateHelpText() {
	if fn.currentFocus >= 0 && fn.currentFocus < len(fn.focusGroup) {
		currentComponent := fn.focusGroup[fn.currentFocus]
		if helpText, exists := fn.helpTexts[currentComponent]; exists {
			fn.vbp.GetTUI().UpdateStatusBar(helpText)
		}
	}
}

// SetComponentHelpText sets help text for a specific component
func (fn *FormNavigation) SetComponentHelpText(component tview.Primitive, helpText string) {
	fn.helpTexts[component] = helpText
}

// handleInputCapture handles input for the file list
func (fn *FormNavigation) handleInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRight:
		// Load selected directory
		fn.vbp.loadSelectedItem(fn.vbp.fileList)
		return nil
	case tcell.KeyLeft:
		// Go back to previous directory
		fn.vbp.goBackDirectory()
		return nil
	case tcell.KeyCtrlN:
		// Create new vault
		fn.vbp.GetTUI().GetNavigation().ShowNewVaultWithDir(fn.vbp.currentDir, false)
		return nil
	case tcell.KeyCtrlE:
		// Edit selected vault
		fn.vbp.editSelectedVault()
		return nil
	case tcell.KeyCtrlR:
		// Rename selected vault
		fn.vbp.renameSelectedVault()
		return nil
	case tcell.KeyCtrlD:
		// Delete selected vault
		fn.vbp.deleteSelectedVault()
		return nil
	}
	return event
}

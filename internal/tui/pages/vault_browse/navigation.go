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
		vbp.directoryList, // Left column - directories
		vbp.fileList,      // Right column - vault files
	}

	return &FormNavigation{
		vbp:          vbp,
		currentFocus: 0, // Start with directory list focused
		focusGroup:   focusGroup,
		helpTexts:    make(map[tview.Primitive]string),
	}
}

func (fn *FormNavigation) SetupNavigation() {
	// Set up help texts for each component
	fn.setupHelpTexts()

	// Set up input capture for both lists
	fn.vbp.directoryList.SetInputCapture(fn.handleDirectoryInputCapture)
	fn.vbp.fileList.SetInputCapture(fn.handleFileInputCapture)

	// Set initial focus to directory list
	fn.vbp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])

	// Set initial help text
	fn.updateHelpText()
}

// setupHelpTexts sets up help text for each component
func (fn *FormNavigation) setupHelpTexts() {
	fn.helpTexts[fn.vbp.directoryList] = "Directories: Enter/→: Open directory | Backspace/←: Go back | Tab: Switch to files | Ctrl+N: New vault"
	fn.helpTexts[fn.vbp.fileList] = "Vault Files: Enter/→: Open vault | Tab: Switch to directories | Ctrl+E: Edit vault | Ctrl+R: Rename vault | Ctrl+D: Delete vault"
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

// handleDirectoryInputCapture handles input for the directory list
func (fn *FormNavigation) handleDirectoryInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyTab, tcell.KeyBacktab:
		// Switch to file list
		fn.currentFocus = 1
		fn.vbp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
		fn.updateHelpText()
		return nil

	case tcell.KeyEnter, tcell.KeyRight:
		// Load selected directory (navigate into it)
		fn.vbp.loadSelectedDirectory()
		return nil
	case tcell.KeyCtrlN:
		// Create new vault
		dir := fn.vbp.getCurrentDisplayedDirectory()
		fn.vbp.GetTUI().GetNavigation().ShowNewVaultWithDir(dir, false)
		return nil
	case tcell.KeyLeft, tcell.KeyBackspace2, tcell.KeyBackspace:
		// Go back to parent directory
		fn.vbp.goBackDirectory()
		return nil
	}
	return event
}

// handleFileInputCapture handles input for the file list
func (fn *FormNavigation) handleFileInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyTab, tcell.KeyBacktab:
		// Switch to directory list
		fn.currentFocus = 0
		fn.vbp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
		fn.updateHelpText()
		return nil
	case tcell.KeyEnter, tcell.KeyRight:
		// Load selected file (only if it's a real vault file)
		fn.vbp.loadSelectedFile()
		return nil
	case tcell.KeyCtrlE:
		// Edit selected vault (only if it's a real vault file)
		fn.vbp.editSelectedVault()
		return nil
	case tcell.KeyCtrlR:
		// Rename selected vault (only if it's a real vault file)
		fn.vbp.renameSelectedVault()
		return nil
	case tcell.KeyCtrlD:
		// Delete selected vault (only if it's a real vault file)
		fn.vbp.deleteSelectedVault()
		return nil
	}
	return event
}

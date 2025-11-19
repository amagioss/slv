package vault_new

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FormNavigation struct {
	vnp          *VaultNewPage
	currentFocus int
	focusGroup   []tview.Primitive
	helpTexts    map[tview.Primitive]string // Component-specific help texts
}

func (fn *FormNavigation) NewFormNavigation(vnp *VaultNewPage) *FormNavigation {
	focusGroup := []tview.Primitive{
		vnp.vaultConfigForm,
		vnp.optionsForm,
		vnp.grantAccessForm,
		vnp.shareWithSelfForm,
		vnp.shareWithK8sForm,
		vnp.searchResults,
		vnp.grantedAccess,
		vnp.submitButton,
	}
	intialFocus := 0

	return &FormNavigation{
		vnp:          vnp,
		currentFocus: intialFocus,
		focusGroup:   focusGroup,
		helpTexts:    make(map[tview.Primitive]string),
	}
}

func (fn *FormNavigation) SetupNavigation() {
	// Set up help texts for each component
	fn.setupHelpTexts()

	fn.setInputCaptureForConfigForm()
	fn.setInputCaptureForOptionsForm()
	fn.setInputCaptureForShareWithSelfForm()
	fn.setInputCaptureForShareWithK8sForm()
	fn.setInputCaptureForSearchBarForm()
	fn.setInputCaptureForSearchResultsForm()
	fn.setInputCaptureForGrantedAccessForm()
	fn.setInputCaptureForSubmitButton()

	fn.vnp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])

	// Set initial help text
	fn.updateHelpText()
}

func (fn *FormNavigation) ShiftFocusForward() {
	fn.currentFocus = (fn.currentFocus + 1) % len(fn.focusGroup)
	fn.vnp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
	fn.updateHelpText()
}

func (fn *FormNavigation) ShiftFocusBackward() {
	fn.currentFocus = (fn.currentFocus - 1 + len(fn.focusGroup)) % len(fn.focusGroup)
	fn.vnp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
	fn.updateHelpText()
}

func (fn *FormNavigation) defaultFormInputCapture(event *tcell.EventKey) *tcell.EventKey {

	switch event.Key() {
	case tcell.KeyTab:
		// Move to next form/primitive
		fn.ShiftFocusForward()
		return nil
	case tcell.KeyBacktab:
		// Move to previous form/primitive
		fn.ShiftFocusBackward()
		return nil
	case tcell.KeyCtrlS:
		// Create vault with Ctrl+S
		fn.vnp.createVaultFromForm()
		return nil
	}
	// Let all other keys pass through to the primitive for within-form navigation
	return event
}

func (fn *FormNavigation) setInputCaptureForConfigForm() {
	fn.vnp.vaultConfigForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		if event == nil {
			return event
		} else {
			switch event.Key() {
			case tcell.KeyDown:
				// Move to next field
				i, _ := fn.vnp.vaultConfigForm.GetFocusedItemIndex()
				next := (i + 1) % fn.vnp.vaultConfigForm.GetFormItemCount()
				fn.vnp.GetTUI().GetApplication().SetFocus(fn.vnp.vaultConfigForm.GetFormItem(next))
				return nil
			case tcell.KeyUp:
				// Move to previous field
				i, _ := fn.vnp.vaultConfigForm.GetFocusedItemIndex()
				prev := (i - 1 + fn.vnp.vaultConfigForm.GetFormItemCount()) % fn.vnp.vaultConfigForm.GetFormItemCount()
				fn.vnp.GetTUI().GetApplication().SetFocus(fn.vnp.vaultConfigForm.GetFormItem(prev))
				return nil
			}

		}
		return event
	})
}

func (fn *FormNavigation) setInputCaptureForOptionsForm() {
	fn.vnp.optionsForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		if event == nil {
			return event
		} else {
			switch event.Key() {
			case tcell.KeyDown:
				// Move to next field
				i, _ := fn.vnp.optionsForm.GetFocusedItemIndex()
				next := (i + 1) % fn.vnp.optionsForm.GetFormItemCount()
				fn.vnp.GetTUI().GetApplication().SetFocus(fn.vnp.optionsForm.GetFormItem(next))
				return nil
			case tcell.KeyUp:
				// Move to previous field
				i, _ := fn.vnp.optionsForm.GetFocusedItemIndex()
				prev := (i - 1 + fn.vnp.optionsForm.GetFormItemCount()) % fn.vnp.optionsForm.GetFormItemCount()
				fn.vnp.GetTUI().GetApplication().SetFocus(fn.vnp.optionsForm.GetFormItem(prev))
				return nil
			}
			// Let all other keys pass through
			return event
		}
	})
}

func (fn *FormNavigation) setInputCaptureForShareWithSelfForm() {
	fn.vnp.shareWithSelfForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		return event
	})
}

func (fn *FormNavigation) setInputCaptureForShareWithK8sForm() {
	fn.vnp.shareWithK8sForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		return event
	})
}

func (fn *FormNavigation) setInputCaptureForSearchBarForm() {
	fn.vnp.grantAccessForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		if event == nil {
			return event
		} else {
			switch event.Key() {
			case tcell.KeyEnter:
				// Handle Enter key for SLV_EPK inputs
				fn.vnp.handleSearchBarEnter()
				return nil
			case tcell.KeyDown:
				// Move to next field
				fn.vnp.GetTUI().GetApplication().SetFocus(fn.focusGroup[5])
				return nil
			}
			return event
		}
	})
}

func (fn *FormNavigation) setInputCaptureForSearchResultsForm() {
	fn.vnp.searchResults.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		if event == nil {
			return event
		} else {
			switch event.Key() {
			case tcell.KeyEnter:
				// Add selected item to granted access
				selected := fn.vnp.searchResults.GetCurrentItem()
				if selected >= 0 && selected < fn.vnp.searchResults.GetItemCount() {
					mainText, _ := fn.vnp.searchResults.GetItemText(selected)
					// Extract environment name from the formatted text and find the environment
					if strings.HasPrefix(mainText, "🔍 ") {
						envName := strings.TrimPrefix(mainText, "🔍 ")
						// Find the environment in the search results map
						if env, exists := fn.vnp.searchEnvMap[envName]; exists {
							fn.vnp.addToGrantedAccess(env)
						}
					}
				}
				return nil
			}
			return event
		}
	})
}

func (fn *FormNavigation) setInputCaptureForGrantedAccessForm() {
	fn.vnp.grantedAccess.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		if event == nil {
			return event
		} else {
			switch event.Key() {
			case tcell.KeyCtrlD:
				// Remove selected environment from granted access
				selected := fn.vnp.grantedAccess.GetCurrentItem()
				if selected >= 0 && selected < fn.vnp.grantedAccess.GetItemCount() {
					_, secondaryText := fn.vnp.grantedAccess.GetItemText(selected)
					// Extract public key from secondary text
					// Format: "Email: xxx | PK: full_public_key"
					if strings.Contains(secondaryText, " | Key: ") {
						parts := strings.Split(secondaryText, " | Key: ")
						if len(parts) == 2 {
							publicKey := parts[1]
							fn.vnp.removeFromGrantedAccess(publicKey)
						}
					}
				}
				return nil
			}
		}
		return event
	})
}

func (fn *FormNavigation) setInputCaptureForSubmitButton() {
	fn.vnp.submitButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		return event
	})
}

// setupHelpTexts sets up help text for each component
func (fn *FormNavigation) setupHelpTexts() {
	fn.helpTexts[fn.vnp.vaultConfigForm] = "Vault Config: ↑/↓: Navigate fields | Tab: Next section | Enter: Select field"
	fn.helpTexts[fn.vnp.optionsForm] = "Options: ↑/↓: Navigate fields | Space: Toggle checkbox | Tab: Next section"
	fn.helpTexts[fn.vnp.grantAccessForm] = "Grant Access: Type to search | Enter: Add environment | ↓: View results | Tab: Next section"
	fn.helpTexts[fn.vnp.shareWithSelfForm] = "Share With Self: Space: Toggle checkbox | Tab: Next section"
	fn.helpTexts[fn.vnp.shareWithK8sForm] = "Share With K8s: Space: Toggle checkbox | Tab: Next section"
	fn.helpTexts[fn.vnp.searchResults] = "Search Results: ↑/↓: Navigate | Enter: Add to granted access | Tab: Next section"
	fn.helpTexts[fn.vnp.grantedAccess] = "Granted Access: ↑/↓: Navigate | Ctrl+D: Remove environment | Tab: Next section"
	fn.helpTexts[fn.vnp.submitButton] = "Submit: Enter: Create vault | Ctrl+S: Create vault | Tab: Previous section"
}

// updateHelpText updates the status bar with help text for the currently focused component
func (fn *FormNavigation) updateHelpText() {
	if fn.currentFocus >= 0 && fn.currentFocus < len(fn.focusGroup) {
		currentComponent := fn.focusGroup[fn.currentFocus]
		if helpText, exists := fn.helpTexts[currentComponent]; exists {
			fn.vnp.GetTUI().UpdateStatusBar(helpText)
		}
	}
}

// SetComponentHelpText sets help text for a specific component
func (fn *FormNavigation) SetComponentHelpText(component tview.Primitive, helpText string) {
	fn.helpTexts[component] = helpText
}

package vault_edit

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FormNavigation struct {
	vep          *VaultEditPage
	currentFocus int
	focusGroup   []tview.Primitive
	helpTexts    map[tview.Primitive]string // Component-specific help texts
}

func (fn *FormNavigation) NewFormNavigation(vep *VaultEditPage) *FormNavigation {
	// Build focus group based on vault access state
	focusGroup := []tview.Primitive{
		vep.vaultConfigForm,
		// vep.optionsForm, // Commented out by user
	}

	// Only include access-related components if vault is unlocked
	if vep.IsVaultUnlocked() {
		focusGroup = append(focusGroup,
			vep.grantAccessForm,
			vep.shareWithSelfForm,
			vep.shareWithK8sForm,
			vep.searchResults,
			vep.grantedAccess,
		)
	}

	// Always include submit button
	focusGroup = append(focusGroup, vep.submitButton)

	intialFocus := 0

	return &FormNavigation{
		vep:          vep,
		currentFocus: intialFocus,
		focusGroup:   focusGroup,
		helpTexts:    make(map[tview.Primitive]string),
	}
}

func (fn *FormNavigation) SetupNavigation() {
	// Set up help texts for each component
	fn.setupHelpTexts()

	// Always set up input capture for basic components
	fn.setInputCaptureForConfigForm()

	// Only set up input capture for access-related components if vault is unlocked
	if fn.vep.IsVaultUnlocked() {
		fn.setInputCaptureForShareWithSelfForm()
		fn.setInputCaptureForShareWithK8sForm()
		fn.setInputCaptureForSearchBarForm()
		fn.setInputCaptureForSearchResultsForm()
		fn.setInputCaptureForGrantedAccessForm()
	}

	// Always set up submit button input capture
	fn.setInputCaptureForSubmitButton()

	fn.vep.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])

	// Set initial help text
	fn.updateHelpText()
}

func (fn *FormNavigation) ShiftFocusForward() {
	fn.currentFocus = (fn.currentFocus + 1) % len(fn.focusGroup)
	fn.vep.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
	fn.updateHelpText()
}

func (fn *FormNavigation) ShiftFocusBackward() {
	fn.currentFocus = (fn.currentFocus - 1 + len(fn.focusGroup)) % len(fn.focusGroup)
	fn.vep.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
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
		fn.vep.editVaultFromForm()
		return nil
	}
	// Let all other keys pass through to the primitive for within-form navigation
	return event
}

func (fn *FormNavigation) setInputCaptureForConfigForm() {
	fn.vep.vaultConfigForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		if event == nil {
			return event
		} else {
			switch event.Key() {
			case tcell.KeyDown:
				// Move to next field
				i, _ := fn.vep.vaultConfigForm.GetFocusedItemIndex()
				next := (i + 1) % fn.vep.vaultConfigForm.GetFormItemCount()
				fn.vep.GetTUI().GetApplication().SetFocus(fn.vep.vaultConfigForm.GetFormItem(next))
				return nil
			case tcell.KeyUp:
				// Move to previous field
				i, _ := fn.vep.vaultConfigForm.GetFocusedItemIndex()
				prev := (i - 1 + fn.vep.vaultConfigForm.GetFormItemCount()) % fn.vep.vaultConfigForm.GetFormItemCount()
				fn.vep.GetTUI().GetApplication().SetFocus(fn.vep.vaultConfigForm.GetFormItem(prev))
				return nil
			}

		}
		return event
	})
}

// func (fn *FormNavigation) setInputCaptureForOptionsForm() {
// 	fn.vep.optionsForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		event = fn.defaultFormInputCapture(event)
// 		if event == nil {
// 			return event
// 		} else {
// 			switch event.Key() {
// 			case tcell.KeyDown:
// 				// Move to next field
// 				i, _ := fn.vep.optionsForm.GetFocusedItemIndex()
// 				next := (i + 1) % fn.vep.optionsForm.GetFormItemCount()
// 				fn.vep.GetTUI().GetApplication().SetFocus(fn.vep.optionsForm.GetFormItem(next))
// 				return nil
// 			case tcell.KeyUp:
// 				// Move to previous field
// 				i, _ := fn.vep.optionsForm.GetFocusedItemIndex()
// 				prev := (i - 1 + fn.vep.optionsForm.GetFormItemCount()) % fn.vep.optionsForm.GetFormItemCount()
// 				fn.vep.GetTUI().GetApplication().SetFocus(fn.vep.optionsForm.GetFormItem(prev))
// 				return nil
// 			}
// 			// Let all other keys pass through
// 			return event
// 		}
// 	})
// }

func (fn *FormNavigation) setInputCaptureForShareWithSelfForm() {
	fn.vep.shareWithSelfForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		return event
	})
}

func (fn *FormNavigation) setInputCaptureForShareWithK8sForm() {
	fn.vep.shareWithK8sForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		return event
	})
}

func (fn *FormNavigation) setInputCaptureForSearchBarForm() {
	fn.vep.grantAccessForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		if event == nil {
			return event
		} else {
			switch event.Key() {
			case tcell.KeyEnter:
				// Handle Enter key for SLV_EPK inputs
				fn.vep.handleSearchBarEnter()
				return nil
			case tcell.KeyDown:
				// Move to next field
				fn.vep.GetTUI().GetApplication().SetFocus(fn.focusGroup[5])
				return nil
			}
			return event
		}
	})
}

func (fn *FormNavigation) setInputCaptureForSearchResultsForm() {
	fn.vep.searchResults.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		if event == nil {
			return event
		} else {
			switch event.Key() {
			case tcell.KeyEnter:
				// Add selected item to granted access
				selected := fn.vep.searchResults.GetCurrentItem()
				if selected >= 0 && selected < fn.vep.searchResults.GetItemCount() {
					mainText, _ := fn.vep.searchResults.GetItemText(selected)
					// Extract environment name from the formatted text and find the environment
					if strings.HasPrefix(mainText, "🔍 ") {
						envName := strings.TrimPrefix(mainText, "🔍 ")
						// Find the environment in the search results map
						if env, exists := fn.vep.searchEnvMap[envName]; exists {
							fn.vep.addToGrantedAccess(env)
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
	fn.vep.grantedAccess.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		if event == nil {
			return event
		} else {
			switch event.Key() {
			case tcell.KeyCtrlD:
				// Remove selected environment from granted access
				selected := fn.vep.grantedAccess.GetCurrentItem()
				if selected >= 0 && selected < fn.vep.grantedAccess.GetItemCount() {
					_, secondaryText := fn.vep.grantedAccess.GetItemText(selected)
					// Extract public key from secondary text
					// Format: "Email: xxx | PK: full_public_key"
					if strings.Contains(secondaryText, " | Public Key: ") {
						parts := strings.Split(secondaryText, " | Public Key: ")
						if len(parts) == 2 {
							publicKey := parts[1]
							fn.vep.removeFromGrantedAccess(publicKey)
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
	fn.vep.submitButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = fn.defaultFormInputCapture(event)
		return event
	})
}

// setupHelpTexts sets up help text for each component
func (fn *FormNavigation) setupHelpTexts() {
	// Always set up help texts for basic components
	fn.helpTexts[fn.vep.vaultConfigForm] = "Vault Config (Read-Only): ↑/↓: Navigate fields | Tab: Next section"

	// Set up help texts based on vault access state
	if fn.vep.IsVaultUnlocked() {
		// Unlocked state - all access components are interactive
		fn.helpTexts[fn.vep.grantAccessForm] = "Grant Access: Type to search | Enter: Add environment | ↓: View results | Tab: Next section"
		fn.helpTexts[fn.vep.shareWithSelfForm] = "Share With Self: Space: Toggle checkbox | Tab: Next section"
		fn.helpTexts[fn.vep.shareWithK8sForm] = "Share With K8s: Space: Toggle checkbox | Tab: Next section"
		fn.helpTexts[fn.vep.searchResults] = "Search Results: ↑/↓: Navigate | Enter: Add to granted access | Tab: Next section"
		fn.helpTexts[fn.vep.grantedAccess] = "Granted Access: ↑/↓: Navigate | Ctrl+D: Remove environment | Tab: Next section"
		fn.helpTexts[fn.vep.submitButton] = "Submit: Enter: Update vault | Ctrl+S: Update vault | Tab: Previous section"
	} else {
		// Locked state - submit button is disabled
		fn.helpTexts[fn.vep.submitButton] = "Submit (Locked): No access to update vault | Tab: Previous section"
	}
}

// updateHelpText updates the status bar with help text for the currently focused component
func (fn *FormNavigation) updateHelpText() {
	if fn.currentFocus >= 0 && fn.currentFocus < len(fn.focusGroup) {
		currentComponent := fn.focusGroup[fn.currentFocus]
		if helpText, exists := fn.helpTexts[currentComponent]; exists {
			fn.vep.GetTUI().UpdateStatusBar(helpText)
		}
	}
}

// SetComponentHelpText sets help text for a specific component
func (fn *FormNavigation) SetComponentHelpText(component tview.Primitive, helpText string) {
	fn.helpTexts[component] = helpText
}

package vault_view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
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
	fn.vvp.itemsTable.SetInputCapture(fn.handleSecretItemsInputCapture)
	fn.vvp.accessorsTable.SetInputCapture(fn.handleAccessorsInputCapture)

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
		case tcell.KeyCtrlE:
			if fn.currentFocus == 0 || fn.currentFocus == 1 {
				// Save state before navigating to edit
				fn.vvp.SaveNavigationState()
				fn.vvp.GetTUI().GetNavigation().ShowVaultEditWithVault(fn.vvp.vault, fn.vvp.filePath, false)
				return nil
			}
			return event
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				// Clear state when going back to vault browse
				fn.vvp.ClearNavigationState()
				fn.vvp.GetTUI().GetNavigation().ShowVaults(false)
				return nil
			case 'x', 'X':
				// Toggle vault lock state
				if fn.vvp.filePath != "" {
					if fn.vvp.vault.IsLocked() {
						fn.vvp.unlockVault()
					} else {
						fn.vvp.lockVault()
					}
					// Update help texts after state change
					fn.setupHelpTexts()
					fn.updateHelpText()
				}
				return nil
			case 'r', 'R':
				// Reload vault
				if fn.vvp.filePath != "" {
					fn.vvp.reloadVault()
					fn.vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(fn.vvp.vault, fn.vvp.filePath, true)
				}
				return nil
			}
		case tcell.KeyEsc:
			// Clear state when going back to vault browse
			fn.vvp.ClearNavigationState()
			fn.vvp.GetTUI().GetNavigation().ShowVaults(false)
			return event
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight, tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyHome, tcell.KeyEnd:
			// Allow arrow keys and page keys to scroll
			return event
		}
		return event
	}
}

func (fn *FormNavigation) handleSecretItemsInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	} else {
		switch event.Key() {
		case tcell.KeyCtrlD:
			fn.vvp.removeSecretItem()
			return nil
		case tcell.KeyCtrlN:
			// Show add item modal form
			fn.showAddItemModal()
			return nil
		case tcell.KeyCtrlE:
			fn.vvp.showEditSecretItemModal()
			return nil
		case tcell.KeyEnter:
			// Show item details modal
			row, _ := fn.vvp.itemsTable.GetSelection()
			if row >= 1 { // Skip header
				nameCell := fn.vvp.itemsTable.GetCell(row, 0)
				if nameCell != nil {
					itemName := nameCell.Text
					// Get item from vault to verify access
					item, err := fn.vvp.vault.Get(itemName)
					if err == nil {
						// Check if we can access the value (unlocked or plaintext)
						if !fn.vvp.vault.IsLocked() || item.IsPlaintext() {
							value, err := item.ValueString()
							if err == nil {
								itemType := "Secret"
								if item.IsPlaintext() {
									itemType = "Plaintext"
								}
								fn.vvp.showItemDetailsModal(itemName, itemType, value)
							} else {
								fn.vvp.ShowError(fmt.Sprintf("Error getting value: %v", err))
							}
						} else {
							fn.vvp.ShowError("Cannot view secret value while vault is locked")
						}
					}
				}
			}
			return nil
		case tcell.KeyEsc:
			// Clear state when going back to vault browse
			fn.vvp.ClearNavigationState()
			fn.vvp.GetTUI().GetNavigation().ShowVaults(false)
			return event
		case tcell.KeyRune:
			switch event.Rune() {
			case 'c', 'C':
				// Copy value to clipboard
				row, _ := fn.vvp.itemsTable.GetSelection()
				if row >= 1 { // Skip header
					nameCell := fn.vvp.itemsTable.GetCell(row, 0)
					if nameCell != nil {
						itemName := nameCell.Text
						// Get item from vault to verify access
						item, err := fn.vvp.vault.Get(itemName)
						if err == nil {
							// Check if we can access the value (unlocked or plaintext)
							if !fn.vvp.vault.IsLocked() || item.IsPlaintext() {
								value, err := item.ValueString()
								if err == nil {
									clipboard.Write(clipboard.FmtText, []byte(value))
									fn.vvp.GetTUI().UpdateStatusBar(fmt.Sprintf("Copied value of '%s' to clipboard", itemName))
								} else {
									fn.vvp.ShowError(fmt.Sprintf("Error getting value: %v", err))
								}
							} else {
								fn.vvp.ShowError("Cannot copy secret value while vault is locked")
							}
						}
					}
				}
				return nil
			}
		}

		return event
	}
}

// handleAccessorsInputCapture handles input for the accessors table
func (fn *FormNavigation) handleAccessorsInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyEnter:
		// Show accessor details modal
		row, _ := fn.vvp.accessorsTable.GetSelection()
		if row >= 1 { // Skip header
			typeCell := fn.vvp.accessorsTable.GetCell(row, 0)
			nameCell := fn.vvp.accessorsTable.GetCell(row, 1)
			emailCell := fn.vvp.accessorsTable.GetCell(row, 2)
			publicKeyCell := fn.vvp.accessorsTable.GetCell(row, 3)

			if publicKeyCell != nil {
				accessorType := ""
				if typeCell != nil {
					accessorType = typeCell.Text
				}
				name := ""
				if nameCell != nil {
					name = nameCell.Text
				}
				email := ""
				if emailCell != nil {
					email = emailCell.Text
				}
				publicKey := publicKeyCell.Text

				fn.vvp.showAccessorDetailsModal(accessorType, name, email, publicKey)
			}
		}
		return nil
	}

	return event
}

// setupHelpTexts sets up help text for each component
func (fn *FormNavigation) setupHelpTexts() {
	lockAction := "Lock"
	if fn.vvp.vault.IsLocked() {
		lockAction = "Unlock"
	}

	fn.helpTexts[fn.vvp.vaultDetailsTable] = fmt.Sprintf("Vault Details: ↑/↓: Navigate rows | Tab: Next table | x: %s | r: Reload | Ctrl+E: Edit vault", lockAction)
	fn.helpTexts[fn.vvp.accessorsTable] = fmt.Sprintf("Accessors: ↑/↓: Navigate rows | Tab: Next table | Enter: View Details | x: %s | r: Reload | Ctrl+E: Edit vault", lockAction)
	fn.helpTexts[fn.vvp.itemsTable] = fmt.Sprintf("Items: ↑/↓: Navigate rows | Tab: Next table | Enter: View Details | c: Copy value | x: %s | r: Reload | Ctrl+D: Delete | Ctrl+N: Add | Ctrl+E: Edit", lockAction)
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

func (fn *FormNavigation) SetCurrentFocus(focus int) {
	fn.currentFocus = focus
	fn.resetSelectable()
	fn.focusGroup[fn.currentFocus].(*tview.Table).SetSelectable(true, false)
	fn.vvp.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
	fn.updateHelpText()
}

package vault_view

import (
	"github.com/rivo/tview"
)

// SaveNavigationState saves the current navigation state for restoration
func (vvp *VaultViewPage) SaveNavigationState() {
	var currentFocus, vaultDetailsRow, accessorsRow, itemsRow int
	nav := vvp.GetTUI().GetNavigation()

	// Save current focus
	if vvp.navigation != nil {
		currentFocus = vvp.navigation.currentFocus
		nav.SavePageState("vault_details", "currentFocus", currentFocus)
	}

	// Save table selections
	if vvp.vaultDetailsTable != nil && currentFocus == 0 {
		row, _ := vvp.vaultDetailsTable.GetSelection()
		vaultDetailsRow = row
		nav.SavePageState("vault_details", "vaultDetailsRow", vaultDetailsRow)

	}
	if vvp.accessorsTable != nil && currentFocus == 1 {
		row, _ := vvp.accessorsTable.GetSelection()
		accessorsRow = row
		nav.SavePageState("vault_details", "accessorsRow", accessorsRow)
	}
	if vvp.itemsTable != nil && currentFocus == 2 {
		row, _ := vvp.itemsTable.GetSelection()
		itemsRow = row
		nav.SavePageState("vault_details", "itemsRow", itemsRow)
	}

	// Save to general page state management
	// nav.SavePageState("vault_details", "filePath", vvp.filePath)
}

// RestoreNavigationState restores the saved navigation state
func (vvp *VaultViewPage) RestoreNavigationState() {
	// Get saved state from general page state management
	nav := vvp.GetTUI().GetNavigation()

	currentFocusValue, hasCurrentFocus := nav.GetPageState("vault_details", "currentFocus")
	vaultDetailsRowValue, hasVaultDetailsRow := nav.GetPageState("vault_details", "vaultDetailsRow")
	accessorsRowValue, hasAccessorsRow := nav.GetPageState("vault_details", "accessorsRow")
	itemsRowValue, hasItemsRow := nav.GetPageState("vault_details", "itemsRow")

	if !hasCurrentFocus || (hasVaultDetailsRow && !hasAccessorsRow && !hasItemsRow) {
		return // No saved state to restore
	}

	// Type assert the values
	currentFocus, ok1 := currentFocusValue.(int)
	vaultDetailsRow, ok2 := vaultDetailsRowValue.(int)
	accessorsRow, ok3 := accessorsRowValue.(int)
	itemsRow, ok4 := itemsRowValue.(int)

	if !ok1 && (!ok2 || !ok3 || !ok4) {
		return // Invalid state data
	}

	// Restore table selections
	if hasVaultDetailsRow && vvp.vaultDetailsTable != nil && vaultDetailsRow >= 0 {
		totalRows := vvp.vaultDetailsTable.GetRowCount()
		if vaultDetailsRow < totalRows {
			vvp.vaultDetailsTable.Select(vaultDetailsRow, 0)
		}
	}

	if hasAccessorsRow && vvp.accessorsTable != nil && accessorsRow >= 0 {
		totalRows := vvp.accessorsTable.GetRowCount()
		if accessorsRow < totalRows {
			vvp.accessorsTable.Select(accessorsRow, 0)
		}
	}

	if hasItemsRow && vvp.itemsTable != nil && itemsRow >= 0 {
		totalRows := vvp.itemsTable.GetRowCount()
		if itemsRow < totalRows {
			vvp.itemsTable.Select(itemsRow, 0)
		}
	}

	// Restore focus to the appropriate table
	if hasCurrentFocus && vvp.navigation != nil && currentFocus >= 0 && currentFocus < len(vvp.navigation.focusGroup) {
		vvp.navigation.currentFocus = currentFocus
		vvp.navigation.resetSelectable()
		vvp.navigation.focusGroup[vvp.navigation.currentFocus].(*tview.Table).SetSelectable(true, false)
		vvp.GetTUI().GetApplication().SetFocus(vvp.navigation.focusGroup[vvp.navigation.currentFocus])
		vvp.navigation.updateHelpText()
	}
}

// ClearNavigationState clears the saved navigation state
func (vvp *VaultViewPage) ClearNavigationState() {
	nav := vvp.GetTUI().GetNavigation()
	nav.ClearPageState("vault_details")
}

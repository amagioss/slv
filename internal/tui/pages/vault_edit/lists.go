package vault_edit

import (
	"github.com/rivo/tview"
)

func (vep *VaultEditPage) createSearchResultsList() *tview.List {
	searchResults := tview.NewList()
	searchResults.SetBorder(true)
	searchResults.AddItem("", "", 0, nil)

	// Set title based on vault access state
	title := "Environment Results From Profile"
	if !vep.IsVaultUnlocked() {
		title = "Environment Results (Locked - No Access)"
		// Apply disabled styling when vault is locked
		vep.applyDisabledStyling(searchResults)
	}

	searchResults.SetTitle(title).SetTitleAlign(tview.AlignLeft)
	searchResults.SetWrapAround(false) // Disable looping behavior
	vep.searchResults = searchResults
	return searchResults
}

func (vep *VaultEditPage) createGrantedAccessList() *tview.List {
	grantedAccess := tview.NewList()
	grantedAccess.SetBorder(true)

	// Set title based on vault access state
	title := "Environments With Access"
	if !vep.IsVaultUnlocked() {
		title = "Environments With Access (Locked - No Access)"
		// Apply disabled styling when vault is locked
		vep.applyDisabledStyling(grantedAccess)
	}

	grantedAccess.SetTitle(title).SetTitleAlign(tview.AlignLeft)
	grantedAccess.SetWrapAround(false) // Disable looping behavior
	vep.grantedAccess = grantedAccess
	return grantedAccess
}

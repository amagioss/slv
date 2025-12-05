package vault_new

import "github.com/rivo/tview"

func (vnp *VaultNewPage) createSearchResultsList() *tview.List {
	searchResults := tview.NewList()
	searchResults.SetBorder(true)
	searchResults.AddItem("", "", 0, nil)
	searchResults.SetTitle("Environment Results From Profile").SetTitleAlign(tview.AlignLeft)
	searchResults.SetWrapAround(false) // Disable looping behavior
	vnp.searchResults = searchResults
	return searchResults
}

func (vnp *VaultNewPage) createGrantedAccessList() *tview.List {
	grantedAccess := tview.NewList()
	grantedAccess.SetBorder(true).SetTitle("Environments With Access").SetTitleAlign(tview.AlignLeft)
	grantedAccess.AddItem("📝 No access granted yet", "Add public keys or environments to grant access", 0, nil)
	grantedAccess.SetWrapAround(false) // Disable looping behavior
	vnp.grantedAccess = grantedAccess
	return grantedAccess
}

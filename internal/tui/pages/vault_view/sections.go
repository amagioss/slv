package vault_view

import "github.com/rivo/tview"

func (vvp *VaultViewPage) createMainSection() tview.Primitive {
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	mainFlex.AddItem(vvp.createVaultDetailsTable(), 0, 30, true)
	mainFlex.AddItem(vvp.createAccessorsTable(), 0, 30, false)
	mainFlex.AddItem(vvp.createVaultItemsTable(), 0, 40, false)
	vvp.mainFlex = mainFlex
	return mainFlex
}

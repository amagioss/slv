package vault_new

import "github.com/rivo/tview"

func (vnp *VaultNewPage) createMainSection() tview.Primitive {
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	mainFlex.AddItem(vnp.createVaultConfigAndOptionsSection(), 9, 1, true)
	mainFlex.AddItem(vnp.createAccessSection(), 0, 1, false)
	mainFlex.AddItem(vnp.createSubmitButtonSection(), 3, 1, false)
	return mainFlex
}

func (vnp *VaultNewPage) createVaultConfigAndOptionsSection() tview.Primitive {

	if vnp.vaultConfigForm == nil {
		vnp.createVaultConfigForm()
	}
	if vnp.optionsForm == nil {
		vnp.createVaultOptionsForm()
	}

	mainGrid := tview.NewGrid().
		SetRows(0).           // Single row
		SetColumns(-70, -30). // Two columns: Config, Options
		SetBorders(false)

	mainGrid.AddItem(vnp.vaultConfigForm, 0, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(vnp.optionsForm, 0, 1, 1, 1, 0, 0, false)

	return mainGrid
}

func (vnp *VaultNewPage) createAccessSection() tview.Primitive {

	if vnp.grantAccessForm == nil {
		vnp.createVaultGrantAccessForm()
	}
	if vnp.shareWithSelfForm == nil {
		vnp.createVaultShareWithSelfForm()
	}
	if vnp.shareWithK8sForm == nil {
		vnp.createVaultShareWithK8sForm()
	}
	if vnp.searchResults == nil {
		vnp.createSearchResultsList()
	}
	if vnp.grantedAccess == nil {
		vnp.createGrantedAccessList()
	}

	accessFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	accessRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	accessCheckboxesFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	accessCheckboxesFlex.AddItem(vnp.shareWithSelfForm, 0, 1, false)
	accessCheckboxesFlex.AddItem(vnp.shareWithK8sForm, 0, 1, false)
	accessCheckboxesFlex.SetBorder(true).
		SetTitle("Quick Access").
		SetTitleAlign(tview.AlignLeft)

	accessRow.AddItem(vnp.grantAccessForm, 0, 70, true)
	accessRow.AddItem(accessCheckboxesFlex, 0, 30, false)

	resultsFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	resultsFlex.AddItem(vnp.searchResults, 0, 1, true)
	resultsFlex.AddItem(vnp.grantedAccess, 0, 1, false)

	accessFlex.AddItem(accessRow, 5, 1, false)
	accessFlex.AddItem(resultsFlex, 0, 1, false)

	return accessFlex
}

func (vnp *VaultNewPage) createSubmitButtonSection() tview.Primitive {

	if vnp.submitButton == nil {
		vnp.createSubmitButton()
	}
	submitFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	centeredFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centeredFlex.AddItem(nil, 0, 1, false)              // Left spacer
	centeredFlex.AddItem(vnp.submitButton, 20, 1, true) // Button (20 chars wide)
	centeredFlex.AddItem(nil, 0, 1, false)

	submitFlex.AddItem(centeredFlex, 3, 1, true) // Button row
	submitFlex.AddItem(nil, 1, 1, false)         // Right spacer

	return submitFlex
}

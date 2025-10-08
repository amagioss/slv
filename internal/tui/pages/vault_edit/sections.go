package vault_edit

import "github.com/rivo/tview"

func (vep *VaultEditPage) createMainSection() tview.Primitive {
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	mainFlex.AddItem(vep.createVaultConfigAndOptionsSection(), 9, 1, true)
	mainFlex.AddItem(vep.createAccessSection(), 0, 1, false)
	mainFlex.AddItem(vep.createSubmitButtonSection(), 3, 1, false)
	return mainFlex
}

func (vep *VaultEditPage) createVaultConfigAndOptionsSection() tview.Primitive {

	if vep.vaultConfigForm == nil {
		vep.createVaultConfigForm()
	}
	if vep.optionsForm == nil {
		vep.createVaultOptionsForm()
	}

	mainGrid := tview.NewGrid().
		SetRows(0).           // Single row
		SetColumns(-70, -30). // Two columns: Config, Options
		SetBorders(false)

	mainGrid.AddItem(vep.vaultConfigForm, 0, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(vep.optionsForm, 0, 1, 1, 1, 0, 0, false)

	return mainGrid
}

func (vep *VaultEditPage) createAccessSection() tview.Primitive {

	if vep.grantAccessForm == nil {
		vep.createVaultGrantAccessForm()
	}
	if vep.shareWithSelfForm == nil {
		vep.createVaultShareWithSelfForm()
	}
	if vep.shareWithK8sForm == nil {
		vep.createVaultShareWithK8sForm()
	}
	if vep.searchResults == nil {
		vep.createSearchResultsList()
	}
	if vep.grantedAccess == nil {
		vep.createGrantedAccessList()
	}

	accessFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	accessRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	accessCheckboxesFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	accessCheckboxesFlex.AddItem(vep.shareWithSelfForm, 0, 1, false)
	accessCheckboxesFlex.AddItem(vep.shareWithK8sForm, 0, 1, false)

	if vep.IsVaultUnlocked() {
		accessCheckboxesFlex.SetBorder(true).
			SetTitle("Quick Access").
			SetTitleAlign(tview.AlignLeft)
	} else {
		accessCheckboxesFlex.SetBorder(true).
			SetTitle("Quick Access (Locked - No Access)").
			SetTitleAlign(tview.AlignLeft)
	}

	accessRow.AddItem(vep.grantAccessForm, 0, 70, true)
	accessRow.AddItem(accessCheckboxesFlex, 0, 30, false)

	resultsFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	resultsFlex.AddItem(vep.searchResults, 0, 1, true)
	resultsFlex.AddItem(vep.grantedAccess, 0, 1, false)

	accessFlex.AddItem(accessRow, 5, 1, false)
	accessFlex.AddItem(resultsFlex, 0, 1, false)

	if !vep.IsVaultUnlocked() {
		vep.applyDisabledStyling(accessCheckboxesFlex)
	}

	return accessFlex
}

func (vep *VaultEditPage) createSubmitButtonSection() tview.Primitive {

	if vep.submitButton == nil {
		vep.createSubmitButton()
	}
	submitFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	centeredFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centeredFlex.AddItem(nil, 0, 1, false)              // Left spacer
	centeredFlex.AddItem(vep.submitButton, 20, 1, true) // Button (20 chars wide)
	centeredFlex.AddItem(nil, 0, 1, false)

	submitFlex.AddItem(centeredFlex, 3, 1, true) // Button row
	submitFlex.AddItem(nil, 1, 1, false)         // Right spacer

	return submitFlex
}

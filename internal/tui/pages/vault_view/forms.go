package vault_view

import "github.com/rivo/tview"

func (vvp *VaultViewPage) createAddItemForm() *tview.Form {
	form := tview.NewForm().
		AddInputField("Name", "", 40, nil, nil).
		AddInputField("Value", "", 40, nil, nil).
		AddCheckbox("Plain Text", false, nil)

	form.GetFormItem(2).(*tview.Checkbox).SetCheckedString("✓")

	// Style the form with larger dimensions
	form.SetBorder(true).
		SetTitle("Add New Item").
		SetTitleAlign(tview.AlignCenter)

	return form
}

func (vvp *VaultViewPage) createEditItemForm(itemKey string, itemValue string, isPlaintext bool) *tview.Form {
	form := tview.NewForm().
		AddInputField("Name", itemKey, 40, nil, nil).
		AddInputField("Value", itemValue, 40, nil, nil).
		AddCheckbox("Plain Text", isPlaintext, nil)

	// Disable the name field to prevent changes
	form.GetFormItem(0).(*tview.InputField).SetDisabled(true)

	form.GetFormItem(2).(*tview.Checkbox).SetCheckedString("✓")

	// Style the form with larger dimensions
	form.SetBorder(true).
		SetTitle("Edit Item").
		SetTitleAlign(tview.AlignCenter)

	return form
}

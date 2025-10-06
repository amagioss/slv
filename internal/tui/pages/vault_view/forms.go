package vault_view

import "github.com/rivo/tview"

func (vvp *VaultViewPage) createAddItemForm() *tview.Form {
	form := tview.NewForm().
		AddInputField("Name", "", 40, nil, nil).
		AddInputField("Value", "", 40, nil, nil).
		AddCheckbox("Encrypted", true, nil)

	form.GetFormItem(2).(*tview.Checkbox).SetCheckedString("✓")

	// Don't set border/title here - will be handled by ShowModalForm
	return form
}

func (vvp *VaultViewPage) createEditItemForm(itemKey string, itemValue string, isEncrypted bool) *tview.Form {
	form := tview.NewForm().
		AddInputField("Name", itemKey, 40, nil, nil).
		AddInputField("Value", itemValue, 40, nil, nil).
		AddCheckbox("Encrypted", isEncrypted, nil)

	// Disable the name field to prevent changes
	form.GetFormItem(0).(*tview.InputField).SetDisabled(true)

	form.GetFormItem(2).(*tview.Checkbox).SetCheckedString("✓")

	// Don't set border/title here - will be handled by ShowModalForm
	return form
}

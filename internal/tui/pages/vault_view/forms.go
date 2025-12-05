package vault_view

import (
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/utils"
)

func (vvp *VaultViewPage) createAddItemForm() *tview.Form {
	form := tview.NewForm().
		AddInputField("Name", "", 40, nil, nil)

	// Add manual paste handler to Name field
	nameField := form.GetFormItem(0).(*tview.InputField)
	utils.AttachPasteHandler(nameField)

	// Add TextArea for Value field to support unlimited text length
	valueTextArea := tview.NewTextArea()
	valueTextArea.SetLabel("Value: ").
		SetLabelWidth(10).
		SetText("", false).
		SetBorder(false)

	// Add manual paste handler to bypass tcell paste event limitations
	utils.AttachPasteHandler(valueTextArea)

	form.AddFormItem(valueTextArea)

	form.AddCheckbox("Encrypted", true, nil)

	form.GetFormItem(2).(*tview.Checkbox).SetCheckedString("✓")

	// Don't set border/title here - will be handled by ShowModalForm
	return form
}

func (vvp *VaultViewPage) createEditItemForm(itemKey string, itemValue string, isEncrypted bool) *tview.Form {
	form := tview.NewForm().
		AddInputField("Name", itemKey, 40, nil, nil)

	// Disable the name field to prevent changes
	form.GetFormItem(0).(*tview.InputField).SetDisabled(true)

	// Add TextArea for Value field to support unlimited text length
	valueTextArea := tview.NewTextArea()
	valueTextArea.SetLabel("Value: ").
		SetLabelWidth(10).
		SetText(itemValue, false).
		SetBorder(false)

	// Add manual paste handler to bypass tcell paste event limitations
	utils.AttachPasteHandler(valueTextArea)

	form.AddFormItem(valueTextArea)

	form.AddCheckbox("Encrypted", isEncrypted, nil)

	form.GetFormItem(2).(*tview.Checkbox).SetCheckedString("✓")

	// Don't set border/title here - will be handled by ShowModalForm
	return form
}

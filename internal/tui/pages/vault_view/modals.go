package vault_view

import (
	"fmt"

	"github.com/rivo/tview"
)

func (vvp *VaultViewPage) showEditSecretItemModal() {

	var itemKey, itemValue, itemType string
	var isEncrypted bool
	selected, _ := vvp.itemsTable.GetSelection()
	if selected >= 0 && selected < vvp.itemsTable.GetRowCount() {
		itemKey = vvp.itemsTable.GetCell(selected, 0).Text
		itemType = vvp.itemsTable.GetCell(selected, 1).Text
		if !vvp.vault.IsLocked() {
			itemValue = vvp.itemsTable.GetCell(selected, 2).Text
		} else {
			itemValue = ""
		}
		if itemType == "Secret" {
			isEncrypted = true
		} else {
			isEncrypted = false
		}
	}

	form := vvp.createEditItemForm(itemKey, itemValue, isEncrypted)

	// Get the TextArea and ensure it can handle paste properly
	valueTextArea := form.GetFormItem(1).(*tview.TextArea)

	vvp.GetTUI().ShowModalForm("Edit Item", form, "Edit", "Cancel", func() {
		nameField := form.GetFormItem(0).(*tview.InputField)
		valueTextArea := form.GetFormItem(1).(*tview.TextArea)
		encryptedCheckbox := form.GetFormItem(2).(*tview.Checkbox)

		name := nameField.GetText()
		value := valueTextArea.GetText()
		encrypted := encryptedCheckbox.IsChecked()

		if err := vvp.vault.Put(name, []byte(value), encrypted); err != nil {
			vvp.ShowError(fmt.Sprintf("Error updating item: %v", err))
			return
		}

		vvp.reloadVault()
		// Show success message
		// vvp.GetTUI().ShowInfo(fmt.Sprintf("Item updated: Name='%s', Type='%s'", name, map[bool]string{true: "Plaintext", false: "Secret"}[plainText]))

		// Show the vault details with the current vault instance (which has the latest changes)
		vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vvp.vault, vvp.filePath, true)
	}, func() {
		vvp.GetTUI().GetApplication().SetFocus(vvp.itemsTable)
	}, func() {
		// Set focus to TextArea when modal opens to ensure paste works
		vvp.GetTUI().GetApplication().SetFocus(valueTextArea)
	})
}

// showAddItemModal shows the modal form for adding a new item
func (fn *FormNavigation) showAddItemModal() {
	// Create the form with larger input fields
	form := fn.vvp.createAddItemForm()

	// Get the TextArea and ensure it can handle paste properly
	valueTextArea := form.GetFormItem(1).(*tview.TextArea)

	// Show the modal form
	fn.vvp.GetTUI().ShowModalForm("Add New Item", form, "Add", "Cancel", func() {
		// Confirm callback - TODO: implement item addition logic
		// Get form values
		nameField := form.GetFormItem(0).(*tview.InputField)
		valueTextArea := form.GetFormItem(1).(*tview.TextArea)
		encryptedCheckbox := form.GetFormItem(2).(*tview.Checkbox)

		name := nameField.GetText()
		value := valueTextArea.GetText()
		encrypted := encryptedCheckbox.IsChecked()

		if name == "" || value == "" {
			fn.vvp.ShowError("Name and value are required")
			return
		}

		if err := fn.vvp.vault.Put(name, []byte(value), encrypted); err != nil {
			fn.vvp.ShowError(err.Error())
			return
		}

		// TODO: Add item to vault using name, value, and plainText
		fn.vvp.GetTUI().ShowInfo(fmt.Sprintf("Item added: Name='%s', Value='%s', PlainText=%v", name, value, encrypted))
		fn.vvp.SaveNavigationState()
		fn.vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(fn.vvp.vault, fn.vvp.filePath, true)
	}, func() {
		// Cancel callback - do nothing
		fn.vvp.GetTUI().GetApplication().SetFocus(fn.vvp.itemsTable)
	}, func() {
		// Set focus to TextArea when modal opens to ensure paste works
		fn.vvp.GetTUI().GetApplication().SetFocus(valueTextArea)
	})
}

package vault_view

import (
	"fmt"

	"github.com/rivo/tview"
)

func (vvp *VaultViewPage) showEditSecretItemModal() {

	var itemKey, itemValue, itemType string
	var isPlaintext bool
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
			isPlaintext = false
		} else {
			isPlaintext = true
		}
	}

	form := vvp.createEditItemForm(itemKey, itemValue, isPlaintext)
	vvp.GetTUI().ShowModalForm("Edit Item", form, "Edit", "Cancel", func() {
		nameField := form.GetFormItem(0).(*tview.InputField)
		valueField := form.GetFormItem(1).(*tview.InputField)
		plainTextCheckbox := form.GetFormItem(2).(*tview.Checkbox)

		name := nameField.GetText()
		value := valueField.GetText()
		plainText := plainTextCheckbox.IsChecked()

		if err := vvp.vault.Put(name, []byte(value), !plainText); err != nil {
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
		// Restore focus callback - focus will be handled by the new page
	})
}

// showAddItemModal shows the modal form for adding a new item
func (fn *FormNavigation) showAddItemModal() {
	// Create the form with larger input fields
	form := fn.vvp.createAddItemForm()
	// Show the modal form
	fn.vvp.GetTUI().ShowModalForm("Add New Item", form, "Add", "Cancel", func() {
		// Confirm callback - TODO: implement item addition logic
		// Get form values
		nameField := form.GetFormItem(0).(*tview.InputField)
		valueField := form.GetFormItem(1).(*tview.InputField)
		plainTextCheckbox := form.GetFormItem(2).(*tview.Checkbox)

		name := nameField.GetText()
		value := valueField.GetText()
		plainText := plainTextCheckbox.IsChecked()

		if name == "" || value == "" {
			fn.vvp.ShowError("Name and value are required")
			return
		}

		if err := fn.vvp.vault.Put(name, []byte(value), !plainText); err != nil {
			fn.vvp.ShowError(err.Error())
			return
		}

		// TODO: Add item to vault using name, value, and plainText
		fn.vvp.GetTUI().ShowInfo(fmt.Sprintf("Item added: Name='%s', Value='%s', PlainText=%v", name, value, plainText))
		fn.vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(fn.vvp.vault, fn.vvp.filePath, true)
	}, func() {
		// Cancel callback - do nothing
	}, func() {
		// Restore focus to items table
		fn.vvp.GetTUI().GetApplication().SetFocus(fn.vvp.itemsTable)
	})
}

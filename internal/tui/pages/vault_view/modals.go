package vault_view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/theme"
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

// showItemDetailsModal shows a modal with full item details
func (vvp *VaultViewPage) showItemDetailsModal(name, itemType, value string) {
	colors := theme.GetCurrentPalette()

	// Handle empty value
	if value == "" {
		value = "(empty)"
	}

	// Create content container
	content := tview.NewFlex().SetDirection(tview.FlexRow)

	// Name field
	content.AddItem(tview.NewTextView().
		SetText(fmt.Sprintf("[::b]Name[::-]: %s", name)).
		SetDynamicColors(true).
		SetTextColor(colors.TextPrimary), 1, 1, false)

	// Type field
	content.AddItem(tview.NewTextView().
		SetText(fmt.Sprintf("[::b]Type[::-]: %s", itemType)).
		SetDynamicColors(true).
		SetTextColor(colors.TextPrimary), 1, 1, false)

	// Value label
	content.AddItem(tview.NewTextView().
		SetText("[::b]Value[::-]:").
		SetDynamicColors(true).
		SetTextColor(colors.TextPrimary), 1, 1, false)

	// Value text area (scrollable - TextView is scrollable by default when focused)
	valueArea := tview.NewTextView()
	valueArea.SetText(value)
	valueArea.SetTextColor(colors.TextPrimary)
	valueArea.SetBackgroundColor(colors.BackgroundDark)
	valueArea.SetBorder(true)
	valueArea.SetBorderColor(colors.Border)
	valueArea.SetDynamicColors(false) // Disable dynamic colors to ensure plain text display
	valueArea.SetWordWrap(true)       // Enable word wrap for long text

	// Add value area with minimum height to ensure it's visible
	content.AddItem(valueArea, 0, 1, true)

	// Close button
	btn := tview.NewButton("Close").
		SetSelectedFunc(func() {
			vvp.GetTUI().GetComponents().GetMainContentPages().RemovePage("modal")
			vvp.GetTUI().GetApplication().SetFocus(vvp.itemsTable)
		})
	btn.SetBackgroundColor(colors.Primary)
	btn.SetLabelColor(colors.TextPrimary)
	btn.SetBackgroundColorActivated(colors.Accent)
	btn.SetLabelColorActivated(colors.Background)

	content.AddItem(btn, 1, 0, false)

	// Create bordered container for the modal content with fixed width for larger modal
	innerContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(content, 0, 1, true)
	innerContent.SetBorder(true).
		SetTitle(" Item Details ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(colors.Border).
		SetBackgroundColor(colors.Background)

	// Wrap in a fixed-width container to make modal bigger (80 columns wide)
	modalContent := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).          // Left spacer
		AddItem(innerContent, 80, 0, true). // Fixed width of 80 columns
		AddItem(nil, 0, 1, false)           // Right spacer

	// Handle Tab to switch between text area and button
	modalContent.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
			currentFocus := vvp.GetTUI().GetApplication().GetFocus()
			if currentFocus == valueArea {
				vvp.GetTUI().GetApplication().SetFocus(btn)
			} else {
				vvp.GetTUI().GetApplication().SetFocus(valueArea)
			}
			return nil
		}
		// Let Escape be handled by ShowModal's input capture
		return event
	})

	// Use the new ShowModal function
	vvp.GetTUI().ShowModal("Item Details", modalContent, func() {
		vvp.GetTUI().GetApplication().SetFocus(vvp.itemsTable)
	})

	// Set focus directly to the text area (no QueueUpdate)
	vvp.GetTUI().GetApplication().SetFocus(valueArea)
}

// showAccessorDetailsModal shows a modal with full accessor details
func (vvp *VaultViewPage) showAccessorDetailsModal(accessorType, name, email, publicKey string) {
	colors := theme.GetCurrentPalette()

	// Create content container
	content := tview.NewFlex().SetDirection(tview.FlexRow)

	// Type field
	content.AddItem(tview.NewTextView().
		SetText(fmt.Sprintf("[::b]Type[::-]: %s", accessorType)).
		SetDynamicColors(true).
		SetTextColor(colors.TextPrimary), 1, 1, false)

	// Name field
	if name != "" {
		content.AddItem(tview.NewTextView().
			SetText(fmt.Sprintf("[::b]Name[::-]: %s", name)).
			SetDynamicColors(true).
			SetTextColor(colors.TextPrimary), 1, 1, false)
	}

	// Email field
	if email != "" {
		content.AddItem(tview.NewTextView().
			SetText(fmt.Sprintf("[::b]Email[::-]: %s", email)).
			SetDynamicColors(true).
			SetTextColor(colors.TextPrimary), 1, 1, false)
	}

	// Public Key label
	content.AddItem(tview.NewTextView().
		SetText("[::b]Public Key[::-]:").
		SetDynamicColors(true).
		SetTextColor(colors.TextPrimary), 1, 1, false)

	// Public Key text area (scrollable)
	publicKeyArea := tview.NewTextView()
	publicKeyArea.SetText(publicKey)
	publicKeyArea.SetTextColor(colors.TextPrimary)
	publicKeyArea.SetBackgroundColor(colors.BackgroundDark)
	publicKeyArea.SetBorder(true)
	publicKeyArea.SetBorderColor(colors.Border)
	publicKeyArea.SetDynamicColors(false) // Disable dynamic colors to ensure plain text display
	publicKeyArea.SetWordWrap(true)       // Enable word wrap for long text

	// Add public key area with minimum height to ensure it's visible
	content.AddItem(publicKeyArea, 0, 1, true)

	// Close button
	btn := tview.NewButton("Close").
		SetSelectedFunc(func() {
			vvp.GetTUI().GetComponents().GetMainContentPages().RemovePage("modal")
			vvp.GetTUI().GetApplication().SetFocus(vvp.accessorsTable)
		})
	btn.SetBackgroundColor(colors.Primary)
	btn.SetLabelColor(colors.TextPrimary)
	btn.SetBackgroundColorActivated(colors.Accent)
	btn.SetLabelColorActivated(colors.Background)

	content.AddItem(btn, 1, 0, false)

	// Create bordered container for the modal content with fixed width for larger modal
	innerContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(content, 0, 1, true)
	innerContent.SetBorder(true).
		SetTitle(" Accessor Details ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(colors.Border).
		SetBackgroundColor(colors.Background)

	// Wrap in a fixed-width container to make modal bigger (80 columns wide - same as item details)
	modalContent := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).          // Left spacer
		AddItem(innerContent, 80, 0, true). // Fixed width of 80 columns
		AddItem(nil, 0, 1, false)           // Right spacer

	// Handle Tab to switch between text area and button
	modalContent.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
			currentFocus := vvp.GetTUI().GetApplication().GetFocus()
			if currentFocus == publicKeyArea {
				vvp.GetTUI().GetApplication().SetFocus(btn)
			} else {
				vvp.GetTUI().GetApplication().SetFocus(publicKeyArea)
			}
			return nil
		}
		// Let Escape be handled by ShowModal's input capture
		return event
	})

	// Use the new ShowModal function
	vvp.GetTUI().ShowModal("Accessor Details", modalContent, func() {
		vvp.GetTUI().GetApplication().SetFocus(vvp.accessorsTable)
	})

	// Set focus directly to the text area (no QueueUpdate)
	vvp.GetTUI().GetApplication().SetFocus(publicKeyArea)
}

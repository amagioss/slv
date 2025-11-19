package vault_view

import (
	"fmt"

	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
)

// reloadVault reloads the vault
func (vvp *VaultViewPage) reloadVault() {
	// Save current state before reloading
	vvp.SaveNavigationState()

	if vvp.filePath == "" {
		return
	}

	// Load fresh vault instance
	vault, err := vaults.Get(vvp.filePath)
	if err != nil {
		vvp.ShowError(fmt.Sprintf("Error reloading vault: %v", err))
		return
	}

	// Update stored instance
	vvp.vault = vault

	// Refresh the vault details page using the stored instance
	vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vault, vvp.filePath, true)
}

// unlockVault unlocks the vault
func (vvp *VaultViewPage) unlockVault() {
	// Save current state before unlocking
	vvp.SaveNavigationState()

	// Check if we have the vault loaded
	if vvp.vault == nil || vvp.filePath == "" {
		vvp.ShowError("Vault not loaded. Please reopen the vault.")
		return
	}

	// If already unlocked, just refresh the display
	if !vvp.vault.IsLocked() {
		vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vvp.vault, vvp.filePath, true)
		return
	}

	// Attempt to unlock the vault
	secretKey, err := session.GetSecretKey()
	if err != nil {
		vvp.ShowError(fmt.Sprintf("Error getting secret key: %v", err))
		return
	}

	err = vvp.vault.Unlock(secretKey)
	if err != nil {
		vvp.ShowError(fmt.Sprintf("Error unlocking vault: %v", err))
		return
	}

	vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vvp.vault, vvp.filePath, true)
}

// lockVault locks the vault
func (vvp *VaultViewPage) lockVault() {
	// Save current state before locking
	vvp.SaveNavigationState()

	// Check if we have the vault loaded
	if vvp.vault == nil || vvp.filePath == "" {
		vvp.ShowError("Vault not loaded. Please reopen the vault.")
		return
	}

	// If already locked, just refresh the display
	if vvp.vault.IsLocked() {
		vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vvp.vault, vvp.filePath, true)
		return
	}

	// Lock the vault
	vvp.vault.Lock()

	// Refresh the vault details page using the stored instance
	vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vvp.vault, vvp.filePath, true)
}

func (vvp *VaultViewPage) removeSecretItem() {

	var itemKey string
	selected, _ := vvp.itemsTable.GetSelection()
	if selected >= 0 && selected < vvp.itemsTable.GetRowCount() {
		itemKey = vvp.itemsTable.GetCell(selected, 0).Text
	}

	// Show confirmation modal with focus restoration
	vvp.GetTUI().ShowConfirmationWithFocus(
		fmt.Sprintf("Are you sure you want to delete the item '%s'?\n\nThis action cannot be undone.", itemKey),
		"Delete",
		"Cancel",
		func() {
			if vvp.vault == nil || vvp.filePath == "" {
				vvp.ShowError("Vault not loaded. Please reopen the vault.")
				return
			}

			vvp.vault.DeleteItem(itemKey)
			vvp.SaveNavigationState()
			vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vvp.vault, vvp.filePath, true)
		},
		func() {
			// User cancelled - do nothing
		},
		func() {
			// Restore focus to items table after modal is dismissed
			vvp.GetTUI().GetApplication().SetFocus(vvp.itemsTable)
		},
	)
}

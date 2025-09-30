package vault_view

import (
	"fmt"

	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
)

// reloadVault reloads the vault
func (vvp *VaultViewPage) reloadVault() {
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
	vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vault, vvp.filePath, false)
}

// unlockVault unlocks the vault
func (vvp *VaultViewPage) unlockVault() {
	// Check if we have the vault loaded
	if vvp.vault == nil || vvp.filePath == "" {
		vvp.ShowError("Vault not loaded. Please reopen the vault.")
		return
	}

	// If already unlocked, just refresh the display
	if !vvp.vault.IsLocked() {
		vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vvp.vault, vvp.filePath, false)
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

	vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vvp.vault, vvp.filePath, false)
}

// lockVault locks the vault
func (vvp *VaultViewPage) lockVault() {
	// Check if we have the vault loaded
	if vvp.vault == nil || vvp.filePath == "" {
		vvp.ShowError("Vault not loaded. Please reopen the vault.")
		return
	}

	// If already locked, just refresh the display
	if vvp.vault.IsLocked() {
		vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vvp.vault, vvp.filePath, false)
		return
	}

	// Lock the vault
	vvp.vault.Lock()

	// Refresh the vault details page using the stored instance
	vvp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vvp.vault, vvp.filePath, false)
}

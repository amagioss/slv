package vault_browse

import "strings"

// SaveNavigationState saves the current navigation state for restoration
func (vbp *VaultBrowsePage) SaveNavigationState() {
	var dirIndex, fileIndex int
	if vbp.directoryList != nil {
		dirIndex = vbp.directoryList.GetCurrentItem()
	}
	if vbp.fileList != nil {
		fileIndex = vbp.fileList.GetCurrentItem()
	}

	// Save to general page state management
	nav := vbp.GetTUI().GetNavigation()
	nav.SavePageState("vaults", "lastSelectedDirIndex", dirIndex)
	nav.SavePageState("vaults", "lastSelectedFileIndex", fileIndex)
	nav.SavePageState("vaults", "lastViewedDir", vbp.currentDir)
}

// RestoreNavigationState restores the saved navigation state
func (vbp *VaultBrowsePage) RestoreNavigationState() {
	// Get saved state from general page state management
	nav := vbp.GetTUI().GetNavigation()

	dirIndexValue, hasDirIndex := nav.GetPageState("vaults", "lastSelectedDirIndex")
	fileIndexValue, hasFileIndex := nav.GetPageState("vaults", "lastSelectedFileIndex")
	lastViewedDir, hasViewedDir := nav.GetPageState("vaults", "lastViewedDir")

	if !hasDirIndex || !hasFileIndex || !hasViewedDir {
		return // No saved state to restore
	}

	// Type assert the values
	dirIndex, ok1 := dirIndexValue.(int)
	fileIndex, ok2 := fileIndexValue.(int)
	lastViewedDirValue, ok3 := lastViewedDir.(string)
	if !ok1 || !ok2 || !ok3 {
		return // Invalid state data
	}

	vbp.currentDir = lastViewedDirValue

	// Restore directory selection
	if vbp.directoryList != nil {
		// Check if we need to focus on a specific directory (from goBackDirectory)
		if focusDirValue, hasFocusDir := nav.GetPageState("vaults", "focusOnDirectory"); hasFocusDir {
			if focusDirName, ok := focusDirValue.(string); ok {
				// Find the directory with the matching name
				for i := 1; i < vbp.directoryList.GetItemCount(); i++ { // Skip index 0 (current directory)
					itemText, _ := vbp.directoryList.GetItemText(i)
					if strings.Contains(itemText, focusDirName) {
						vbp.directoryList.SetCurrentItem(i)
						// Clear the focus state after using it
						nav.ClearPageStateKey("vaults", "focusOnDirectory")
						break
					}
				}
			}
		} else if dirIndex >= 0 {
			// Use the saved directory index
			totalDirs := vbp.directoryList.GetItemCount()
			if dirIndex < totalDirs {
				vbp.directoryList.SetCurrentItem(dirIndex)
			}
		}
	}

	// Restore file selection
	if vbp.fileList != nil && fileIndex >= 0 {
		totalFiles := vbp.fileList.GetItemCount()
		if fileIndex < totalFiles {
			vbp.fileList.SetCurrentItem(fileIndex)
		}
	}

	// Restore focus to the appropriate list
	if vbp.navigation != nil {
		vbp.navigation.currentFocus = 0 // Default to directory list
		vbp.GetTUI().GetApplication().SetFocus(vbp.directoryList)
		vbp.navigation.updateHelpText()
	}
}

// ClearNavigationState clears the saved navigation state
func (vbp *VaultBrowsePage) ClearNavigationState() {
	nav := vbp.GetTUI().GetNavigation()
	nav.ClearPageState("vaults")
}

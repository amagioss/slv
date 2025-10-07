package vault_browse

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/vaults"
)

// getVaultFiles returns a list of directories and .slv files in the current directory
func (vbp *VaultBrowsePage) getVaultFiles() []VaultFile {
	var items []VaultFile

	// Read the directory
	entries, err := os.ReadDir(vbp.currentDir)
	if err != nil {
		return items
	}

	// Filter and collect items
	for _, entry := range entries {
		// Skip hidden files (except .slv files)
		if strings.HasPrefix(entry.Name(), ".") && !strings.HasSuffix(entry.Name(), ".slv.yaml") && !strings.HasSuffix(entry.Name(), ".slv.yml") {
			continue
		}

		// Check if it's a directory
		if entry.IsDir() {
			items = append(items, VaultFile{
				Name:   entry.Name(),
				Path:   filepath.Join(vbp.currentDir, entry.Name()),
				IsFile: false,
			})
		} else {
			// Check if it's a .slv file
			if strings.HasSuffix(entry.Name(), ".slv.yaml") || strings.HasSuffix(entry.Name(), ".slv.yml") {
				items = append(items, VaultFile{
					Name:   entry.Name(),
					Path:   filepath.Join(vbp.currentDir, entry.Name()),
					IsFile: true,
				})
			}
		}
	}

	return items
}

// handleItemSelection handles selection of a file or directory
func (vbp *VaultBrowsePage) handleItemSelection(item VaultFile) {
	if item.IsFile {
		// Handle .slv file selection - open for viewing
		vbp.openVaultFile(item.Path)
	} else {
		// Handle directory selection - navigate into the directory
		vbp.currentDir = item.Path
		// Replace the current vault page with new directory
		vbp.GetTUI().GetNavigation().ShowVaultsWithDir(item.Path, true)
	}
}

// loadSelectedItem loads the currently selected item
func (vbp *VaultBrowsePage) loadSelectedItem(list *tview.List) {
	// Get the current selection index
	selectedIndex := list.GetCurrentItem()

	// Skip the "Go Back" option (index 0)
	if selectedIndex == 0 {
		vbp.goBackDirectory()
		return
	}

	// Adjust index for the "Go Back" option
	itemIndex := selectedIndex - 1

	// Get the items
	items := vbp.getVaultFiles()

	// Check if the index is valid
	if itemIndex >= 0 && itemIndex < len(items) {
		item := items[itemIndex]
		vbp.handleItemSelection(item)
	}
}

// goBackDirectory navigates back to the parent directory
func (vbp *VaultBrowsePage) goBackDirectory() {
	parentDir := filepath.Dir(vbp.currentDir)
	// Don't go back if we're already at the root
	if parentDir != vbp.currentDir {
		vbp.currentDir = parentDir
		// Replace the current vault page with parent directory
		vbp.GetTUI().GetNavigation().ShowVaultsWithDir(parentDir, true)
	}
}

// updateFileList refreshes the file list displayed in the UI
func (vbp *VaultBrowsePage) updateFileList() {
	vbp.fileList.Clear()
	vbp.pwdTextView.SetText(fmt.Sprintf("\n[white]Browse Vaults[white::-]\n[gray](Use arrow keys [←] and [→] to navigate directories)[gray::-]\n\nCurrent Directory: %s", vbp.currentDir))

	// Add "go back one directory" option at the top
	vbp.fileList.AddItem("⬆️ Go Back", "Go back to parent directory", 'b', func() {
		vbp.goBackDirectory()
	})

	// Get directories and .slv files
	items := vbp.getVaultFiles()

	// Add items to the list
	for _, item := range items {
		icon := "📁"
		if item.IsFile {
			icon = "📄"
		}

		vbp.fileList.AddItem(
			fmt.Sprintf("%s %s", icon, item.Name),
			"",
			0,
			func() {
				vbp.handleItemSelection(item)
			},
		)
	}
}

// openVaultFile opens a vault file for viewing
func (vbp *VaultBrowsePage) openVaultFile(filePath string) {
	// Load the vault using vaults.Get
	vault, err := vaults.Get(filePath)
	if err != nil {
		vbp.ShowError(fmt.Sprintf("Error loading vault: %v", err))
		return
	}

	// Show vault details page with the loaded vault and filepath
	vbp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vault, filePath, false)
}

// editSelectedVault edits the selected vault
func (vbp *VaultBrowsePage) editSelectedVault() {
	// Get the current selection index
	selectedIndex := vbp.fileList.GetCurrentItem()

	// Skip the "Go Back" option (index 0)
	if selectedIndex == 0 {
		vbp.ShowError("Please select a vault file to edit")
		return
	}

	// Adjust index for the "Go Back" option
	itemIndex := selectedIndex - 1

	// Get the items
	items := vbp.getVaultFiles()

	// Check if the index is valid
	if itemIndex >= 0 && itemIndex < len(items) {
		item := items[itemIndex]

		// Check if it's a vault file
		if !item.IsFile {
			vbp.ShowError("Please select a vault file to edit")
			return
		}

		// Load the vault
		vault, err := vaults.Get(item.Path)
		if err != nil {
			vbp.ShowError(fmt.Sprintf("Error loading vault: %v", err))
			return
		}

		// Navigate to vault edit page
		vbp.GetTUI().GetNavigation().ShowVaultEditWithVault(vault, item.Path, false)
	}
}

// renameSelectedVault renames the selected vault file
func (vbp *VaultBrowsePage) renameSelectedVault() {
	// Get the current selection index
	selectedIndex := vbp.fileList.GetCurrentItem()

	// Skip the "Go Back" option (index 0)
	if selectedIndex == 0 {
		vbp.ShowError("Please select a vault file to edit")
		return
	}

	// Adjust index for the "Go Back" option
	itemIndex := selectedIndex - 1

	// Get the items
	items := vbp.getVaultFiles()

	// Check if the index is valid
	if itemIndex >= 0 && itemIndex < len(items) {
		item := items[itemIndex]

		// Check if it's a vault file
		if !item.IsFile {
			vbp.ShowError("Please select a vault file to edit")
			return
		}

		// Create rename form
		filePath := item.Path
		fileName := item.Name
		form := tview.NewForm().
			AddInputField("New File Name", fileName, 40, nil, nil)

		// Show modal form for renaming
		vbp.GetTUI().ShowModalForm("Rename Vault", form, "Rename", "Cancel", func() {
			newFileName := form.GetFormItem(0).(*tview.InputField).GetText()
			if newFileName == "" {
				vbp.ShowError("File name cannot be empty")
				return
			}
			// Ensure the new filename has the correct extension
			if !strings.HasSuffix(newFileName, ".slv.yaml") && !strings.HasSuffix(newFileName, ".slv.yml") {
				vbp.ShowError("File name must end with .slv.yaml or .slv.yml")
				return
			}

			// Construct new file path
			dir := filepath.Dir(filePath)
			newFilePath := filepath.Join(dir, newFileName)

			// Check if new file already exists
			if _, err := os.Stat(newFilePath); err == nil {
				vbp.ShowError(fmt.Sprintf("File '%s' already exists", newFileName))
				return
			}
			// Rename the file
			if err := os.Rename(filePath, newFilePath); err != nil {
				vbp.ShowError(fmt.Sprintf("Error renaming vault: %v", err))
				return
			}

			vbp.ShowInfo(fmt.Sprintf("Vault renamed to '%s'", newFileName))
			// Refresh the file list to show the renamed file
			vbp.Refresh()
		}, func() {
			// Cancel - do nothing
		}, func() {
			// Restore focus to file list
			vbp.GetTUI().GetApplication().SetFocus(vbp.fileList)
		})
	}
}

// deleteSelectedVault deletes the selected vault file after confirmation
func (vbp *VaultBrowsePage) deleteSelectedVault() {
	// Get the current selection index
	selectedIndex := vbp.fileList.GetCurrentItem()

	// Skip the "Go Back" option (index 0)
	if selectedIndex == 0 {
		vbp.ShowError("Please select a vault file to edit")
		return
	}

	// Adjust index for the "Go Back" option
	itemIndex := selectedIndex - 1

	// Get the items
	items := vbp.getVaultFiles()

	// Check if the index is valid
	if itemIndex >= 0 && itemIndex < len(items) {
		item := items[itemIndex]

		// Check if it's a vault file
		if !item.IsFile {
			vbp.ShowError("Please select a vault file to edit")
			return
		}
		// Get just the filename for display
		fileName := filepath.Base(item.Path)

		if !strings.HasSuffix(fileName, ".slv.yaml") && !strings.HasSuffix(fileName, ".slv.yml") {
			vbp.ShowError("File name must end with .slv.yaml or .slv.yml")
			return
		}
		// Show confirmation modal
		vbp.GetTUI().ShowConfirmationWithFocus(
			fmt.Sprintf("Are you sure you want to delete vault '%s'?\n\nThis action cannot be undone.", fileName),
			func() {
				// Confirm deletion
				if err := os.Remove(item.Path); err != nil {
					vbp.ShowError(fmt.Sprintf("Error deleting vault: %v", err))
					return
				}

				vbp.ShowInfo(fmt.Sprintf("Vault '%s' deleted successfully", fileName))
				// Refresh the file list to remove the deleted file
				vbp.Refresh()
			},
			func() {
				// Cancel - do nothing
			},
			func() {
				// Restore focus to file list
				vbp.GetTUI().GetApplication().SetFocus(vbp.fileList)
			},
		)
	}
}

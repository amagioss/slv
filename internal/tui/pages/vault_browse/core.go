package vault_browse

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
)

// ensureDirectoryLoaded ensures a directory is loaded (lazy loading)
func (vbp *VaultBrowsePage) ensureDirectoryLoaded(dirPath string) {
	// Check if directory is already loaded
	if _, exists := vbp.vaultFileMap[dirPath]; exists {
		return // Already loaded
	}

	// Get session and secret key for accessibility checking
	session, err := session.GetSession()
	if err != nil {
		session = nil
	}

	var secretKey *crypto.SecretKey
	if session != nil {
		secretKey = session.SecretKey()
	}

	// Load the directory
	vbp.scanDirectory(dirPath, secretKey)
}

// scanDirectory scans a single directory and stores its contents
func (vbp *VaultBrowsePage) scanDirectory(dirPath string, secretKey *crypto.SecretKey) {
	// Read the directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}

	var directories []VaultFile
	var vaultFiles []VaultFile

	// Process each entry
	for _, entry := range entries {
		// Skip hidden files/directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			// It's a directory - add to directories list
			directories = append(directories, VaultFile{
				Name:         entry.Name(),
				Path:         entryPath,
				IsFile:       false,
				IsAccessible: true, // Directories are always accessible
			})
		} else if strings.HasSuffix(entry.Name(), ".slv.yaml") || strings.HasSuffix(entry.Name(), ".slv.yml") {
			// It's a vault file - check accessibility
			isAccessible := false
			if secretKey != nil {
				vault, err := vaults.Get(entryPath)
				if err == nil {
					isAccessible = vault.IsAccessibleBy(secretKey)
				}
			}

			vaultFiles = append(vaultFiles, VaultFile{
				Name:         entry.Name(),
				Path:         entryPath,
				IsFile:       true,
				IsAccessible: isAccessible,
			})
		}
	}

	// Store the pre-loaded data
	vbp.directoryMap[dirPath] = directories
	vbp.vaultFileMap[dirPath] = vaultFiles
}

// getDirectories returns only directories from the current directory (using pre-loaded data)
func (vbp *VaultBrowsePage) getDirectories() []VaultFile {
	// Use pre-loaded data if available
	if directories, exists := vbp.directoryMap[vbp.currentDir]; exists {
		return directories
	}

	// Fallback to real-time scanning if pre-loaded data is not available
	var directories []VaultFile
	entries, err := os.ReadDir(vbp.currentDir)
	if err != nil {
		return directories
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if entry.IsDir() {
			directories = append(directories, VaultFile{
				Name:         entry.Name(),
				Path:         filepath.Join(vbp.currentDir, entry.Name()),
				IsFile:       false,
				IsAccessible: true,
			})
		}
	}

	return directories
}

// handleItemSelection handles selection of a file or directory
func (vbp *VaultBrowsePage) handleItemSelection(item VaultFile) {
	if item.IsFile {
		// Handle .slv file selection - open for viewing
		vbp.SaveNavigationState()
		vbp.openVaultFile(item.Path)
	} else {
		// Handle directory selection - navigate into the directory
		vbp.currentDir = item.Path
		vbp.directoryList.SetCurrentItem(0)
		vbp.fileList.SetCurrentItem(0)
		vbp.SaveNavigationState()
		// Replace the current vault page with new directory
		vbp.GetTUI().GetNavigation().ShowVaultsWithDir(item.Path, true)
	}
}

// loadSelectedDirectory loads the currently selected directory
func (vbp *VaultBrowsePage) loadSelectedDirectory() {
	// Get the current selection index
	selectedIndex := vbp.directoryList.GetCurrentItem()

	// Skip the "Go Back" option (index 0)
	if selectedIndex == 0 {
		vbp.GetTUI().GetNavigation().ShowVaultsWithDir(vbp.currentDir, true)
		return
	}

	// Adjust index for the "Go Back" option
	itemIndex := selectedIndex - 1

	// Get the directories
	directories := vbp.getDirectories()

	// Check if the index is valid
	if itemIndex >= 0 && itemIndex < len(directories) {
		dir := directories[itemIndex]
		vbp.handleItemSelection(dir)
	}
}

// loadSelectedFile loads the currently selected vault file
func (vbp *VaultBrowsePage) loadSelectedFile() {
	// Save navigation state before opening vault
	vbp.SaveNavigationState()

	// Get the current selection index
	selectedIndex := vbp.fileList.GetCurrentItem()

	// Get the files from the currently displayed directory
	files := vbp.getVaultFilesFromDirectory(vbp.getCurrentDisplayedDirectory())

	// Check if the index is valid and not the "no files found" placeholder
	if selectedIndex >= 0 && selectedIndex < len(files) {
		file := files[selectedIndex]
		vbp.handleItemSelection(file)
	} else {
		// Handle case when "No vault files found" is selected
		vbp.ShowError("No vault files available to open")
	}
}

// goBackDirectory navigates back to the parent directory
func (vbp *VaultBrowsePage) goBackDirectory() {
	parentDir := filepath.Dir(vbp.currentDir)
	// Don't go back if we're already at the root
	if parentDir != vbp.currentDir {
		// Save the current directory name so we can focus on it when we return
		currentDirName := filepath.Base(vbp.currentDir)
		vbp.GetTUI().GetNavigation().SavePageState("vaults", "focusOnDirectory", currentDirName)

		vbp.currentDir = parentDir
		// Replace the current vault page with parent directory
		vbp.SaveNavigationState()
		vbp.GetTUI().GetNavigation().ShowVaultsWithDir(parentDir, true)
	}
}

// updateFileList refreshes both directory and file lists displayed in the UI
func (vbp *VaultBrowsePage) updateFileList() {
	// Clear both lists
	vbp.directoryList.Clear()
	vbp.fileList.Clear()

	// Update the directory list title to show current directory
	vbp.directoryList.SetTitle(fmt.Sprintf("Directories (%s)", vbp.currentDir))

	// Add "." entry to represent current directory
	vbp.directoryList.AddItem("📁 . (current)", "Current directory", 0, func() {
		vbp.updateVaultFilesForCurrentDir()
	})

	// Get directories and populate directory list
	directories := vbp.getDirectories()
	for _, dir := range directories {
		vbp.directoryList.AddItem(
			fmt.Sprintf("📁 %s", dir.Name),
			"",
			0,
			func() {
				vbp.handleItemSelection(dir)
			},
		)
	}

	// Set up dynamic loading for directory list
	vbp.setupDynamicLoading()

	// Initially populate file list with current directory files
	vbp.updateVaultFilesForCurrentDir()
}

// setupDynamicLoading sets up the directory list to dynamically update vault files
func (vbp *VaultBrowsePage) setupDynamicLoading() {
	vbp.directoryList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		// Handle "." entry (index 0) - show current directory files
		if index == 0 {
			vbp.updateVaultFilesForCurrentDir()
			return
		}

		// Adjust index for the "." entry
		itemIndex := index - 1
		directories := vbp.getDirectories()

		// Check if the index is valid
		if itemIndex >= 0 && itemIndex < len(directories) {
			selectedDir := directories[itemIndex]
			vbp.updateVaultFilesForDirectory(selectedDir.Path)
		}
	})
}

// updateVaultFilesForCurrentDir updates the vault files list for the current directory
func (vbp *VaultBrowsePage) updateVaultFilesForCurrentDir() {
	vbp.updateVaultFilesForDirectory(vbp.currentDir)
}

// updateVaultFilesForDirectory updates the vault files list for a specific directory
func (vbp *VaultBrowsePage) updateVaultFilesForDirectory(dirPath string) {
	// Ensure the directory is loaded (lazy loading)
	vbp.ensureDirectoryLoaded(dirPath)

	// Clear the file list
	vbp.fileList.Clear()

	// Update the title to show the directory path
	vbp.fileList.SetTitle(fmt.Sprintf("Vault Files (%s)", dirPath))

	// Get vault files from the specified directory
	files := vbp.getVaultFilesFromDirectory(dirPath)

	if len(files) == 0 {
		// Add "no vault files found" entry when directory is empty
		vbp.fileList.AddItem(
			"📄 No vault files found",
			"",
			0,
			nil, // No action for this placeholder item
		)
	} else {
		// Add actual vault files
		for _, file := range files {
			if !file.IsAccessible {
				vbp.fileList.AddItem(
					fmt.Sprintf("📄 [lightcoral]%s[white]", file.Name),
					"",
					0,
					func() {
						vbp.handleItemSelection(file)
					},
				)
			} else {
				vbp.fileList.AddItem(
					fmt.Sprintf("📄 [lightgreen]%s[white]", file.Name),
					"",
					0,
					func() {
						vbp.handleItemSelection(file)
					},
				)
			}

		}
	}
}

// getCurrentDisplayedDirectory returns the directory path currently being displayed in the vault files list
func (vbp *VaultBrowsePage) getCurrentDisplayedDirectory() string {
	// Get the current selection index from directory list
	selectedIndex := vbp.directoryList.GetCurrentItem()

	// Skip the "Go Back" option (index 0)
	if selectedIndex == 0 {
		return vbp.currentDir
	}

	// Adjust index for the "Go Back" option
	itemIndex := selectedIndex - 1
	directories := vbp.getDirectories()

	// Check if the index is valid
	if itemIndex >= 0 && itemIndex < len(directories) {
		return directories[itemIndex].Path
	}

	// Fallback to current directory
	return vbp.currentDir
}

// getVaultFilesFromDirectory returns only .slv files from a specific directory (using pre-loaded data)
func (vbp *VaultBrowsePage) getVaultFilesFromDirectory(dirPath string) []VaultFile {
	// Use pre-loaded data if available
	if vaultFiles, exists := vbp.vaultFileMap[dirPath]; exists {
		return vaultFiles
	}

	// Fallback to real-time scanning if pre-loaded data is not available
	var files []VaultFile
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return files
	}

	session, _ := session.GetSession()
	secretKey := session.SecretKey()

	// Filter and collect .slv files only
	for _, entry := range entries {
		// Skip hidden files (except .slv files)
		if strings.HasPrefix(entry.Name(), ".") && !strings.HasSuffix(entry.Name(), ".slv.yaml") && !strings.HasSuffix(entry.Name(), ".slv.yml") {
			continue
		}

		// Check if it's a .slv file
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".slv.yaml") || strings.HasSuffix(entry.Name(), ".slv.yml")) {
			vault, err := vaults.Get(filepath.Join(dirPath, entry.Name()))
			isVaultAccessible := false
			if err == nil && secretKey != nil {
				isVaultAccessible = vault.IsAccessibleBy(secretKey)
			}

			files = append(files, VaultFile{
				Name:         entry.Name(),
				Path:         filepath.Join(dirPath, entry.Name()),
				IsFile:       true,
				IsAccessible: isVaultAccessible,
			})
		}
	}

	return files
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
	// Save navigation state before editing vault
	vbp.SaveNavigationState()

	// Get the current selection index
	selectedIndex := vbp.fileList.GetCurrentItem()

	// Get the files from the currently displayed directory
	files := vbp.getVaultFilesFromDirectory(vbp.getCurrentDisplayedDirectory())

	// Check if the index is valid and not the "no files found" placeholder
	if selectedIndex >= 0 && selectedIndex < len(files) {
		file := files[selectedIndex]

		// Load the vault
		vault, err := vaults.Get(file.Path)
		if err != nil {
			vbp.ShowError(fmt.Sprintf("Error loading vault: %v", err))
			return
		}

		// Navigate to vault edit page
		vbp.GetTUI().GetNavigation().ShowVaultEditWithVault(vault, file.Path, false)
	} else {
		vbp.ShowError("No vault files available to edit")
	}
}

// renameSelectedVault renames the selected vault file
func (vbp *VaultBrowsePage) renameSelectedVault() {
	// Get the current selection index
	selectedIndex := vbp.fileList.GetCurrentItem()

	// Get the files from the currently displayed directory
	files := vbp.getVaultFilesFromDirectory(vbp.getCurrentDisplayedDirectory())

	// Check if the index is valid and not the "no files found" placeholder
	if selectedIndex >= 0 && selectedIndex < len(files) {
		file := files[selectedIndex]

		// Create rename form
		filePath := file.Path
		fileName := file.Name
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

			// Construct new file path using the file's actual directory
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
	} else {
		vbp.ShowError("No vault files available to rename")
	}
}

// deleteSelectedVault deletes the selected vault file after confirmation
func (vbp *VaultBrowsePage) deleteSelectedVault() {
	// Get the current selection index
	selectedIndex := vbp.fileList.GetCurrentItem()

	// Get the files from the currently displayed directory
	files := vbp.getVaultFilesFromDirectory(vbp.getCurrentDisplayedDirectory())

	// Check if the index is valid and not the "no files found" placeholder
	if selectedIndex >= 0 && selectedIndex < len(files) {
		file := files[selectedIndex]

		// Get just the filename for display
		fileName := filepath.Base(file.Path)

		if !strings.HasSuffix(fileName, ".slv.yaml") && !strings.HasSuffix(fileName, ".slv.yml") {
			vbp.ShowError("File name must end with .slv.yaml or .slv.yml")
			return
		}
		// Show confirmation modal
		vbp.GetTUI().ShowConfirmationWithFocus(
			fmt.Sprintf("Are you sure you want to delete vault '%s'?\n\nThis action cannot be undone.", fileName),
			"Delete",
			"Cancel",
			func() {
				// Confirm deletion
				if err := os.Remove(file.Path); err != nil {
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
	} else {
		vbp.ShowError("No vault files available to delete")
	}
}

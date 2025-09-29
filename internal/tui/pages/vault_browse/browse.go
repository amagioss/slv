package vault_browse

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/pages/vault_view"
)

// VaultFile represents a file or directory in the vault browser
type VaultFile struct {
	Name   string
	Path   string
	IsFile bool
}

// VaultBrowsePage handles the vault browsing functionality
type VaultBrowsePage struct {
	pages.BasePage
	currentDir string
	vault      *vaults.Vault // Store the current vault instance
	vaultPath  string        // Store the current vault path
}

// NewVaultBrowsePage creates a new VaultBrowsePage instance
func NewVaultBrowsePage(tui interfaces.TUIInterface, currentDir string) *VaultBrowsePage {
	return &VaultBrowsePage{
		BasePage:   *pages.NewBasePage(tui, "Vault Management"),
		currentDir: currentDir,
		vault:      nil,
		vaultPath:  "",
	}
}

// Create implements the Page interface
func (vbp *VaultBrowsePage) Create() tview.Primitive {
	welcomeText := fmt.Sprintf("\n[white]Browse Vaults[white::-]\n[gray](Use arrow keys [←] and [→] to navigate directories)[gray::-]\n\nCurrent Directory: %s", vbp.currentDir)

	pwd := tview.NewTextView().
		SetText(welcomeText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetWrap(true)

	// Get directories and .slv files
	items := vbp.getVaultFiles()

	// Create the list
	list := tview.NewList()

	// Add "go back one directory" option at the top
	list.AddItem("⬆️ Go Back", "Go back to parent directory", 'b', func() {
		vbp.goBackDirectory()
	})

	// Add items to the list
	for _, item := range items {
		icon := "📁"
		if item.IsFile {
			icon = "📄"
		}

		list.AddItem(
			fmt.Sprintf("%s %s", icon, item.Name),
			"",
			0,
			func() {
				vbp.handleItemSelection(item)
			},
		)
	}

	// Style the list
	list.SetSelectedTextColor(tcell.ColorYellow).
		SetSelectedBackgroundColor(tcell.ColorNavy).
		SetSecondaryTextColor(tcell.ColorGray).
		SetMainTextColor(tcell.ColorWhite)

	// Set up keyboard navigation
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			// Load selected directory
			vbp.loadSelectedItem(list)
			return nil
		case tcell.KeyLeft:
			// Go back to previous directory
			vbp.goBackDirectory()
			return nil
		case tcell.KeyCtrlN:
			// Create new vault
			vbp.GetTUI().GetNavigation().ShowNewVault()
			return nil
		}
		return event
	})

	// Create a centered layout using grid
	content := tview.NewGrid().
		SetRows(6, 0). // Two flexible rows
		SetColumns(0). // Single column
		SetBorders(false)

	// Center the welcome text
	content.AddItem(pwd, 0, 0, 1, 1, 0, 0, false)

	// Center the list
	content.AddItem(list, 1, 0, 1, 1, 0, 0, true)

	// Update status bar with help text
	vbp.GetTUI().UpdateStatusBar("[yellow]←/→: Move between directories | ↑/↓: Navigate | Enter: open vault/directory | Ctrl+N: New vault[white]")

	// Create layout using BasePage method
	vbp.SetTitle("Vault Management")
	return vbp.CreateLayout(content)
}

// Refresh implements the Page interface
func (vbp *VaultBrowsePage) Refresh() {
	// TODO: Implement vault browsing page refresh
}

// HandleInput implements the Page interface
func (vbp *VaultBrowsePage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// TODO: Implement vault browsing page input handling
	return event
}

// GetTitle implements the Page interface
func (vbp *VaultBrowsePage) GetTitle() string {
	return vbp.BasePage.GetTitle()
}

// GetCurrentDir returns the current directory
func (vbp *VaultBrowsePage) GetCurrentDir() string {
	return vbp.currentDir
}

// SetCurrentDir sets the current directory
func (vbp *VaultBrowsePage) SetCurrentDir(dir string) {
	vbp.currentDir = dir
}

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
		vbp.GetTUI().GetNavigation().ShowVaultsReplace()
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
		vbp.GetTUI().GetNavigation().ShowVaultsReplace()
	}
}

// openVaultFile opens a vault file for viewing
func (vbp *VaultBrowsePage) openVaultFile(filePath string) { // Check if we already have this vault loaded

	// Load the vault using vaults.Get
	vault, err := vaults.Get(filePath)
	if err != nil {
		vbp.ShowError(fmt.Sprintf("Error loading vault: %v", err))
		return
	}

	// Store the vault instance and path
	vbp.GetTUI().GetRouter().GetRegisteredPage("vaults_view").(*vault_view.VaultViewPage).SetVault(vault)
	vbp.GetTUI().GetRouter().GetRegisteredPage("vaults_view").(*vault_view.VaultViewPage).SetFilePath(filePath)

	// Create and show vault details page
	vbp.GetTUI().GetNavigation().ShowVaultDetails()
}

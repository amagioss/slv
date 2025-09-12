package pages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
)

// VaultPage handles the vault management page functionality
type VaultPage struct {
	tui        interfaces.TUIInterface
	currentDir string
}

// NewVaultPage creates a new VaultPage instance
func NewVaultPage(tui interfaces.TUIInterface, currentDir string) *VaultPage {
	return &VaultPage{
		tui:        tui,
		currentDir: currentDir,
	}
}

// CreateVaultPage creates the vault management page
func (vp *VaultPage) CreateVaultPage() tview.Primitive {
	// Create welcome message
	welcomeText := fmt.Sprintf("\n[white]Browse Vaults[white::-]\n[gray](Use arrow keys [←] and [→] to navigate directories)[gray::-]\n\nCurrent Directory: %s", vp.currentDir)

	pwd := tview.NewTextView().
		SetText(welcomeText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetWrap(true)

	// Get directories and .slv files
	items := vp.getVaultItems()

	// Create the list
	list := tview.NewList()

	// Add "go back one directory" option at the top
	list.AddItem("⬆️ Go Back", "Go back to parent directory", 'b', func() {
		vp.goBackDirectory()
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
				vp.handleItemSelection(item)
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
			vp.loadSelectedItem(list)
			return nil
		case tcell.KeyLeft:
			// Go back to previous directory
			vp.goBackDirectory()
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

	return vp.tui.CreatePageLayout("Vault Management", content)
}

// VaultItem represents a directory or .slv file
type VaultItem struct {
	Name   string
	Path   string
	IsFile bool
}

// getVaultItems scans the home directory for directories and .slv files
func (vp *VaultPage) getVaultItems() []VaultItem {
	var items []VaultItem

	// Read the directory
	entries, err := os.ReadDir(vp.currentDir)
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
			items = append(items, VaultItem{
				Name:   entry.Name(),
				Path:   filepath.Join(vp.currentDir, entry.Name()),
				IsFile: false,
			})
		} else {
			// Check if it's a .slv file
			if strings.HasSuffix(entry.Name(), ".slv.yaml") || strings.HasSuffix(entry.Name(), ".slv.yml") {
				items = append(items, VaultItem{
					Name:   entry.Name(),
					Path:   filepath.Join(vp.currentDir, entry.Name()),
					IsFile: true,
				})
			}
		}
	}

	return items
}

// handleItemSelection handles when a user selects an item
func (vp *VaultPage) handleItemSelection(item VaultItem) {
	if item.IsFile {
		// Handle .slv file selection
		// TODO: Open the .slv file for editing/viewing
		// For now, just update status
		vp.tui.GetNavigation().UpdateStatus()
	} else {
		// Handle directory selection - navigate into the directory
		vp.tui.GetNavigation().SetVaultDir(item.Path)
		// Replace the current vault page with new directory
		vp.tui.GetNavigation().ShowVaultsReplace()
	}
}

// loadSelectedItem loads the currently selected item
func (vp *VaultPage) loadSelectedItem(list *tview.List) {
	// Get the current selection index
	selectedIndex := list.GetCurrentItem()

	// Skip the "Go Back" option (index 0)
	if selectedIndex == 0 {
		vp.goBackDirectory()
		return
	}

	// Adjust index for the "Go Back" option
	itemIndex := selectedIndex - 1

	// Get the items
	items := vp.getVaultItems()

	// Check if the index is valid
	if itemIndex >= 0 && itemIndex < len(items) {
		item := items[itemIndex]
		vp.handleItemSelection(item)
	}
}

// goBackDirectory goes back to the parent directory
func (vp *VaultPage) goBackDirectory() {
	parentDir := filepath.Dir(vp.currentDir)
	// Don't go back if we're already at the root
	if parentDir != vp.currentDir {
		vp.tui.GetNavigation().SetVaultDir(parentDir)
		// Replace the current vault page with parent directory
		vp.tui.GetNavigation().ShowVaultsReplace()
	}
}

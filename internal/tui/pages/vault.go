package pages

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/tui/interfaces"
)

// VaultPage handles the vault management page functionality
type VaultPage struct {
	tui        interfaces.TUIInterface
	currentDir string
	vault      *vaults.Vault // Store the current vault instance
	vaultPath  string        // Store the current vault path
}

// NewVaultPage creates a new VaultPage instance
func NewVaultPage(tui interfaces.TUIInterface, currentDir string) *VaultPage {
	return &VaultPage{
		tui:        tui,
		currentDir: currentDir,
		vault:      nil,
		vaultPath:  "",
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
	items := vp.getVaultFiles()

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

	// Update status bar with help text
	vp.tui.UpdateStatusBar("[yellow]←/→: Move between directories | ↑/↓: Navigate | Enter: open vault/directory[white]")
	return vp.tui.CreatePageLayout("Vault Management", content)
}

// VaultFile represents a directory or .slv file
type VaultFile struct {
	Name   string
	Path   string
	IsFile bool
}

// getVaultFiles scans the home directory for directories and .slv files
func (vp *VaultPage) getVaultFiles() []VaultFile {
	var items []VaultFile

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
			items = append(items, VaultFile{
				Name:   entry.Name(),
				Path:   filepath.Join(vp.currentDir, entry.Name()),
				IsFile: false,
			})
		} else {
			// Check if it's a .slv file
			if strings.HasSuffix(entry.Name(), ".slv.yaml") || strings.HasSuffix(entry.Name(), ".slv.yml") {
				items = append(items, VaultFile{
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
func (vp *VaultPage) handleItemSelection(item VaultFile) {
	if item.IsFile {
		// Handle .slv file selection - open for viewing
		vp.openVaultFile(item.Path)
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
	items := vp.getVaultFiles()

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

// openVaultFile opens and displays an SLV file for viewing
func (vp *VaultPage) openVaultFile(filePath string) {
	// Check if we already have this vault loaded
	if vp.vault != nil && vp.vaultPath == filePath {
		// Use existing vault instance
		vp.showVaultDetails(vp.vault, filePath)
		return
	}

	// Load the vault using vaults.Get
	vault, err := vaults.Get(filePath)
	if err != nil {
		vp.showError(fmt.Sprintf("Error loading vault: %v", err))
		return
	}

	// Store the vault instance and path
	vp.vault = vault
	vp.vaultPath = filePath

	// Create and show vault details page
	vp.showVaultDetails(vault, filePath)
}

// showVaultDetails displays detailed information about a vault
func (vp *VaultPage) showVaultDetails(vault *vaults.Vault, filePath string) {
	// Create a flex layout to hold the three tables
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// 1. Vault Details Table (30% height)
	vaultDetailsTable := vp.createVaultDetailsTable(vault, filePath)
	flex.AddItem(vaultDetailsTable, 0, 30, true) // First table gets focus

	// 2. Accessors Table (30% height)
	accessorsTable := vp.createAccessorsTable(vault)
	flex.AddItem(accessorsTable, 0, 30, false)

	// 3. Vault Items Table (40% height)
	itemsTable := vp.createVaultItemsTable(vault)
	flex.AddItem(itemsTable, 0, 40, false)

	// Set initial focus to vault details table
	vp.tui.GetApplication().SetFocus(vaultDetailsTable)

	// Track current focus index (0 = vault details, 1 = accessors, 2 = items)
	currentFocusIndex := 0

	// Function to switch focus between tables
	switchFocus := func() {
		// Clear focus from all tables first
		vaultDetailsTable.SetSelectable(false, false)
		accessorsTable.SetSelectable(false, false)
		itemsTable.SetSelectable(false, false)

		// Set focus to the next table
		currentFocusIndex = (currentFocusIndex + 1) % 3
		switch currentFocusIndex {
		case 0:
			vaultDetailsTable.SetSelectable(true, false)
			vp.tui.GetApplication().SetFocus(vaultDetailsTable)
		case 1:
			accessorsTable.SetSelectable(true, false)
			vp.tui.GetApplication().SetFocus(accessorsTable)
		case 2:
			itemsTable.SetSelectable(true, false)
			vp.tui.GetApplication().SetFocus(itemsTable)
		}
	}

	// Update status bar with help text
	vp.tui.UpdateStatusBar("[yellow]q: close | u: unlock | l: lock | r: reload | Tab: switch tables[white]")

	// Set up input capture for the flex
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Switch focus between tables
			switchFocus()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				vp.clearVault()
				vp.tui.GetNavigation().ShowVaults()
				return nil
			case 'u', 'U':
				// Unlock vault
				if filePath != "" {
					vp.unlockVault(filePath)
				}
				return nil
			case 'l', 'L':
				// Lock vault
				if filePath != "" {
					vp.lockVault(filePath)
				}
				return nil
			case 'r', 'R':
				// Reload vault
				if filePath != "" {
					vp.reloadVault()
					vp.showVaultDetails(vp.vault, filePath)
				}
				return nil
			}
		case tcell.KeyEsc:
			vp.tui.GetNavigation().ShowVaults()
			return nil
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight, tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyHome, tcell.KeyEnd:
			// Allow arrow keys and page keys to scroll
			return event
		}
		return event
	})

	// Create the page layout and show it
	page := vp.tui.CreatePageLayout("Vault Details", flex)
	vp.tui.GetNavigation().ShowVaultDetails(page)
}

// showError displays an error message
func (vp *VaultPage) showError(message string) {
	modal := tview.NewModal()
	modal.SetText(fmt.Sprintf("[red]Error[white]\n\n%s", message)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			vp.tui.GetNavigation().ShowVaults()
		})

	modal.SetBackgroundColor(tcell.ColorDarkRed).
		SetTextColor(tcell.ColorWhite).
		SetButtonBackgroundColor(tcell.ColorMaroon).
		SetButtonTextColor(tcell.ColorYellow)

	vp.showModalWithText(fmt.Sprintf("[red]Error[white]\n\n%s", message), "error", "")
}

// showModalWithText displays a modal dialog with the given text
func (vp *VaultPage) showModalWithText(text string, pageName string, filePath string) {
	// Create a text view to display the modal content
	textView := tview.NewTextView().
		SetText(text).
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)

	// Create a layout with the text view and close button
	content := tview.NewGrid().
		SetRows(0, 3). // Text view takes most space, button area at bottom
		SetColumns(0).
		SetBorders(false)

	content.AddItem(textView, 0, 0, 1, 1, 0, 0, true)

	// Add close button as text
	closeButton := tview.NewTextView().
		SetText("[yellow]Press 'q' to close, 'u' to unlock, 'l' to lock, 'r' to reload, ESC to go back, arrow/page keys to scroll[white]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	content.AddItem(closeButton, 1, 0, 1, 1, 0, 0, false)

	// Set up input capture for the text view
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				vp.clearVault()
				vp.tui.GetNavigation().ShowVaults()
				return nil
			case 'u', 'U':
				// Unlock vault
				if filePath != "" {
					vp.unlockVault(filePath)
				}
				return nil
			case 'l', 'L':
				// Lock vault
				if filePath != "" {
					vp.lockVault(filePath)
				}
				return nil
			case 'r', 'R':
				// Reload vault
				if filePath != "" {
					vp.reloadVault()
					vp.showVaultDetails(vp.vault, filePath)
				}
				return nil
			}
		case tcell.KeyEsc:
			vp.tui.GetNavigation().ShowVaults()
			return nil
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight, tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyHome, tcell.KeyEnd:
			// Allow arrow keys and page keys to scroll the text view
			return event
		}
		return event
	})

	// Create the page layout and show it
	page := vp.tui.CreatePageLayout("Vault Details", content)
	vp.tui.GetNavigation().ShowVaultDetails(page)
}

// unlockVault attempts to unlock the vault
func (vp *VaultPage) unlockVault(filePath string) {
	// Check if we have the vault loaded
	if vp.vault == nil || vp.vaultPath != filePath {
		vp.showError("Vault not loaded. Please reopen the vault.")
		return
	}

	// If already unlocked, just refresh the display
	if !vp.vault.IsLocked() {
		vp.showVaultDetails(vp.vault, filePath)
		return
	}

	// Attempt to unlock the vault
	secretKey, err := session.GetSecretKey()
	if err != nil {
		vp.showError(fmt.Sprintf("Error getting secret key: %v", err))
		return
	}

	err = vp.vault.Unlock(secretKey)
	if err != nil {
		vp.showError(fmt.Sprintf("Error unlocking vault: %v", err))
		return
	}

	vp.showVaultDetails(vp.vault, filePath)
}

// lockVault locks the vault
func (vp *VaultPage) lockVault(filePath string) {
	// Check if we have the vault loaded
	if vp.vault == nil || vp.vaultPath != filePath {
		vp.showError("Vault not loaded. Please reopen the vault.")
		return
	}

	// If already locked, just refresh the display
	if vp.vault.IsLocked() {
		vp.showVaultDetails(vp.vault, filePath)
		return
	}

	// Lock the vault
	vp.vault.Lock()

	// Refresh the vault details page using the stored instance
	vp.showVaultDetails(vp.vault, filePath)
}

// clearVault clears the stored vault instance
func (vp *VaultPage) clearVault() {
	vp.vault = nil
	vp.vaultPath = ""
}

// reloadVault reloads the vault from disk (useful if file was modified externally)
func (vp *VaultPage) reloadVault() {
	if vp.vaultPath == "" {
		return
	}

	// Load fresh vault instance
	vault, err := vaults.Get(vp.vaultPath)
	if err != nil {
		vp.showError(fmt.Sprintf("Error reloading vault: %v", err))
		return
	}

	// Update stored instance
	vp.vault = vault
}

// createVaultDetailsTable creates a table for vault details
func (vp *VaultPage) createVaultDetailsTable(vault *vaults.Vault, filePath string) *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true).SetTitle("Metadata").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0) // Fix the first row (header) and no columns

	// Set headers (non-selectable) with fixed width for first column
	table.SetCell(0, 0, tview.NewTableCell("Property").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(20))
	table.SetCell(0, 1, tview.NewTableCell("Value").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))

	// Add vault details
	row := 1
	table.SetCell(row, 0, tview.NewTableCell("Vault Path").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(filePath).SetTextColor(tcell.ColorWhite))
	row++

	table.SetCell(row, 0, tview.NewTableCell("Vault Name").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(vault.ObjectMeta.Name).SetTextColor(tcell.ColorWhite))
	row++

	if vault.ObjectMeta.Namespace != "" {
		table.SetCell(row, 0, tview.NewTableCell("Namespace").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
		table.SetCell(row, 1, tview.NewTableCell(vault.ObjectMeta.Namespace).SetTextColor(tcell.ColorWhite))
		row++
	} else {
		table.SetCell(row, 0, tview.NewTableCell("Namespace").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
		table.SetCell(row, 1, tview.NewTableCell("No Namespace").SetTextColor(tcell.ColorWhite))
		row++
	}

	table.SetCell(row, 0, tview.NewTableCell("Public Key").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(vault.Spec.Config.PublicKey).SetTextColor(tcell.ColorWhite))
	row++

	table.SetCell(row, 0, tview.NewTableCell("Number of Accessors").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%d", len(vault.Spec.Config.WrappedKeys))).SetTextColor(tcell.ColorWhite))
	row++
	// Make table focusable for scrolling with custom selection colors
	table.SetSelectable(true, false) // Vault details table is initially selectable
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))

	return table
}

// createAccessorsTable creates a table for accessors
func (vp *VaultPage) createAccessorsTable(vault *vaults.Vault) *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true).SetTitle("Access").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0) // Fix the first row (header) and no columns

	// Set headers (non-selectable) with fixed column widths
	table.SetCell(0, 0, tview.NewTableCell("Type").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(8))
	table.SetCell(0, 1, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(30))
	table.SetCell(0, 2, tview.NewTableCell("Email").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(30))
	table.SetCell(0, 3, tview.NewTableCell("Public Key").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))

	accessors, err := vault.ListAccessors()
	if err != nil || len(accessors) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No accessors found").SetTextColor(tcell.ColorGray).SetAlign(tview.AlignCenter))
		table.SetCell(1, 1, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		table.SetCell(1, 2, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		table.SetCell(1, 3, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		return table
	}

	// Get profile and environment information
	profile, _ := profiles.GetActiveProfile()
	var root *environments.Environment
	if profile != nil {
		root, _ = profile.GetRoot()
	}
	self := environments.GetSelf()

	row := 1
	for _, accessor := range accessors {
		accessorPubKey, err := accessor.String()
		if err != nil {
			continue
		}

		// Determine accessor type and name
		var accessorType, accessorName, accessorEmail string
		if self != nil && self.PublicKey == accessorPubKey {
			accessorType = "Self"
			accessorName = self.Name
			accessorEmail = self.Email
		} else if root != nil && root.PublicKey == accessorPubKey {
			accessorType = "Root"
			accessorName = root.Name
			accessorEmail = root.Email
		} else if profile != nil {
			if env, _ := profile.GetEnv(accessorPubKey); env != nil {
				if env.EnvType == environments.USER {
					accessorType = "User"
				} else {
					accessorType = "Service"
				}
				accessorName = env.Name
				accessorEmail = env.Email
			} else {
				accessorType = "Unknown"
				accessorName = ""
				accessorEmail = ""
			}
		} else {
			accessorType = "Unknown"
			accessorName = ""
			accessorEmail = ""
		}

		table.SetCell(row, 0, tview.NewTableCell(accessorType).SetTextColor(tcell.ColorAqua).SetMaxWidth(8))
		table.SetCell(row, 1, tview.NewTableCell(accessorName).SetTextColor(tcell.ColorGreen).SetMaxWidth(30))
		table.SetCell(row, 2, tview.NewTableCell(accessorEmail).SetTextColor(tcell.ColorWhite).SetMaxWidth(30))
		table.SetCell(row, 3, tview.NewTableCell(accessorPubKey).SetTextColor(tcell.ColorGray))
		row++
	}

	// Make table focusable for scrolling with custom selection colors
	table.SetSelectable(false, false) // Initially not selectable
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))

	return table
}

// createVaultItemsTable creates a table for vault items
func (vp *VaultPage) createVaultItemsTable(vault *vaults.Vault) *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true).SetTitle("Items").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0) // Fix the first row (header) and no columns

	// Set headers (non-selectable) with fixed column widths
	table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(40))
	table.SetCell(0, 1, tview.NewTableCell("Type").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(12))
	table.SetCell(0, 2, tview.NewTableCell("Value").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))

	itemNames := vault.GetItemNames()
	if len(itemNames) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No items found").SetTextColor(tcell.ColorGray).SetAlign(tview.AlignCenter))
		table.SetCell(1, 1, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		table.SetCell(1, 2, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		return table
	}

	// Sort item names for consistent display order
	sort.Strings(itemNames)

	row := 1
	for _, name := range itemNames {
		table.SetCell(row, 0, tview.NewTableCell(name).SetTextColor(tcell.ColorGreen).SetMaxWidth(25))

		if !vault.IsLocked() {
			// Vault is unlocked - show actual item details
			item, err := vault.Get(name)
			if err == nil {
				encryptedStatus := "Secret"
				if item.IsPlaintext() {
					encryptedStatus = "Plaintext"
				}
				table.SetCell(row, 1, tview.NewTableCell(encryptedStatus).SetTextColor(tcell.ColorWhite).SetMaxWidth(12))

				value, err := item.ValueString()
				if err != nil {
					value = "Error loading value"
				}
				table.SetCell(row, 2, tview.NewTableCell(value).SetTextColor(tcell.ColorWhite))
			} else {
				table.SetCell(row, 1, tview.NewTableCell("Error").SetTextColor(tcell.ColorRed).SetMaxWidth(12))
				table.SetCell(row, 2, tview.NewTableCell("Error loading item").SetTextColor(tcell.ColorRed))
			}
		} else {
			// Vault is locked - show masked value
			table.SetCell(row, 1, tview.NewTableCell("***").SetTextColor(tcell.ColorYellow).SetMaxWidth(12))
			table.SetCell(row, 2, tview.NewTableCell("***").SetTextColor(tcell.ColorGray))
		}
		row++
	}

	// Make table focusable for scrolling with custom selection colors
	table.SetSelectable(false, false) // Initially not selectable
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))

	return table
}

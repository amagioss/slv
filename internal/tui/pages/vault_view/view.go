package vault_view

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
)

// VaultViewPage handles the vault details viewing functionality
type VaultViewPage struct {
	pages.BasePage
	vault    *vaults.Vault
	filePath string
}

// NewVaultViewPage creates a new VaultViewPage instance
func NewVaultViewPage(tui interfaces.TUIInterface, vault *vaults.Vault, filePath string) *VaultViewPage {
	return &VaultViewPage{
		BasePage: *pages.NewBasePage(tui, "Vault Details"),
		vault:    vault,
		filePath: filePath,
	}
}

// Create implements the Page interface
func (vvp *VaultViewPage) Create() tview.Primitive { // Create a flex layout to hold the three tables
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// 1. Vault Details Table (30% height)
	vaultDetailsTable := vvp.createVaultDetailsTable(vvp.vault, vvp.filePath)
	flex.AddItem(vaultDetailsTable, 0, 30, true) // First table gets focus

	// 2. Accessors Table (30% height)
	accessorsTable := vvp.createAccessorsTable(vvp.vault)
	flex.AddItem(accessorsTable, 0, 30, false)

	// 3. Vault Items Table (40% height)
	itemsTable := vvp.createVaultItemsTable(vvp.vault)
	flex.AddItem(itemsTable, 0, 40, false)

	// Set initial focus to vault details table
	vvp.GetTUI().GetApplication().SetFocus(vaultDetailsTable)

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
			vvp.GetTUI().GetApplication().SetFocus(vaultDetailsTable)
		case 1:
			accessorsTable.SetSelectable(true, false)
			vvp.GetTUI().GetApplication().SetFocus(accessorsTable)
		case 2:
			itemsTable.SetSelectable(true, false)
			vvp.GetTUI().GetApplication().SetFocus(itemsTable)
		}
	}

	// Update status bar with help text
	vvp.GetTUI().UpdateStatusBar("[yellow]q: close | u: unlock | l: lock | r: reload | Tab: switch tables[white]")

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
				vvp.GetTUI().GetNavigation().ShowVaults()
				return nil
			case 'u', 'U':
				// Unlock vault
				if vvp.filePath != "" {
					vvp.unlockVault()
				}
				return nil
			case 'l', 'L':
				// Lock vault
				if vvp.filePath != "" {
					vvp.lockVault()
				}
				return nil
			case 'r', 'R':
				// Reload vault
				if vvp.filePath != "" {
					vvp.reloadVault()
					vvp.GetTUI().GetNavigation().ShowVaultDetails()
				}
				return nil
			}
		case tcell.KeyEsc:
			vvp.GetTUI().GetNavigation().ShowVaults()
			return nil
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight, tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyHome, tcell.KeyEnd:
			// Allow arrow keys and page keys to scroll
			return event
		}
		return event
	})

	vvp.SetTitle("Vault Details")
	return vvp.CreateLayout(flex)
}

// Refresh implements the Page interface
func (vvp *VaultViewPage) Refresh() {
	// TODO: Implement vault view page refresh
}

// HandleInput implements the Page interface
func (vvp *VaultViewPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// TODO: Implement vault view page input handling
	return event
}

// GetTitle implements the Page interface
func (vvp *VaultViewPage) GetTitle() string {
	return vvp.BasePage.GetTitle()
}

// SetVault sets the vault
func (vvp *VaultViewPage) SetVault(vault *vaults.Vault) {
	vvp.vault = vault
}

// SetFilePath sets the file path
func (vvp *VaultViewPage) SetFilePath(filePath string) {
	vvp.filePath = filePath
}

// GetVault returns the vault
func (vvp *VaultViewPage) GetVault() *vaults.Vault {
	return vvp.vault
}

// GetFilePath returns the file path
func (vvp *VaultViewPage) GetFilePath() string {
	return vvp.filePath
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
		vvp.GetTUI().GetNavigation().ShowVaultDetails()
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

	vvp.GetTUI().GetNavigation().ShowVaultDetails()
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
		vvp.GetTUI().GetNavigation().ShowVaultDetails()
		return
	}

	// Lock the vault
	vvp.vault.Lock()

	// Refresh the vault details page using the stored instance
	vvp.GetTUI().GetNavigation().ShowVaultDetails()
}

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
	vvp.GetTUI().GetNavigation().ShowVaultDetails()
}

// func (vvp *VaultViewPage) clearVault() {
// 	vvp.vault = nil
// 	vvp.filePath = ""
// }

func (vvp *VaultViewPage) createVaultDetailsTable(vault *vaults.Vault, filePath string) *tview.Table {
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
func (vvp *VaultViewPage) createAccessorsTable(vault *vaults.Vault) *tview.Table {
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
func (vvp *VaultViewPage) createVaultItemsTable(vault *vaults.Vault) *tview.Table {
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

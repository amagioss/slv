package vault_view

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
)

func (vvp *VaultViewPage) createVaultDetailsTable() tview.Primitive {
	table := tview.NewTable()
	table.SetBorder(true).SetTitle("Metadata").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0) // Fix the first row (header) and no columns

	// Set headers (non-selectable) with fixed width for first column
	table.SetCell(0, 0, tview.NewTableCell("Property").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(20))
	table.SetCell(0, 1, tview.NewTableCell("Value").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))

	// Add vault details
	row := 1
	table.SetCell(row, 0, tview.NewTableCell("Vault Path").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(vvp.filePath).SetTextColor(tcell.ColorWhite))
	row++

	table.SetCell(row, 0, tview.NewTableCell("Vault Name").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(vvp.vault.ObjectMeta.Name).SetTextColor(tcell.ColorWhite))
	row++

	if vvp.vault.ObjectMeta.Namespace != "" {
		table.SetCell(row, 0, tview.NewTableCell("Namespace").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
		table.SetCell(row, 1, tview.NewTableCell(vvp.vault.ObjectMeta.Namespace).SetTextColor(tcell.ColorWhite))
		row++
	} else {
		table.SetCell(row, 0, tview.NewTableCell("Namespace").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
		table.SetCell(row, 1, tview.NewTableCell("No Namespace").SetTextColor(tcell.ColorWhite))
		row++
	}

	table.SetCell(row, 0, tview.NewTableCell("Public Key").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(vvp.vault.Spec.Config.PublicKey).SetTextColor(tcell.ColorWhite))
	row++

	table.SetCell(row, 0, tview.NewTableCell("Number of Accessors").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%d", len(vvp.vault.Spec.Config.WrappedKeys))).SetTextColor(tcell.ColorWhite))
	row++
	// Make table focusable for scrolling with custom selection colors
	table.SetSelectable(true, false) // Vault details table is initially selectable
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))

	vvp.vaultDetailsTable = table
	return table
}

func (vvp *VaultViewPage) createAccessorsTable() tview.Primitive {
	table := tview.NewTable()
	table.SetBorder(true).SetTitle("Access").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0) // Fix the first row (header) and no columns

	// Set headers (non-selectable) with fixed column widths
	table.SetCell(0, 0, tview.NewTableCell("Type").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(8))
	table.SetCell(0, 1, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(30))
	table.SetCell(0, 2, tview.NewTableCell("Email").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(30))
	table.SetCell(0, 3, tview.NewTableCell("Public Key").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))

	accessors, err := vvp.vault.ListAccessors()
	if err != nil || len(accessors) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No accessors found").SetTextColor(tcell.ColorGray).SetAlign(tview.AlignCenter))
		table.SetCell(1, 1, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		table.SetCell(1, 2, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		table.SetCell(1, 3, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		vvp.accessorsTable = table
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

	vvp.accessorsTable = table
	return table
}

func (vvp *VaultViewPage) createVaultItemsTable() tview.Primitive {
	table := tview.NewTable()
	table.SetBorder(true).SetTitle("Items").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0) // Fix the first row (header) and no columns

	// Set headers (non-selectable) with fixed column widths
	table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(40))
	table.SetCell(0, 1, tview.NewTableCell("Type").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(12))
	table.SetCell(0, 2, tview.NewTableCell("Value").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))

	itemNames := vvp.vault.GetItemNames()
	if len(itemNames) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No items found").SetTextColor(tcell.ColorGray).SetAlign(tview.AlignCenter))
		table.SetCell(1, 1, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		table.SetCell(1, 2, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		vvp.itemsTable = table
		return table
	}

	// Sort item names for consistent display order
	sort.Strings(itemNames)

	row := 1
	for _, name := range itemNames {
		table.SetCell(row, 0, tview.NewTableCell(name).SetTextColor(tcell.ColorGreen).SetMaxWidth(25))

		if !vvp.vault.IsLocked() {
			// Vault is unlocked - show actual item details
			item, err := vvp.vault.Get(name)
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

	vvp.itemsTable = table
	return table
}

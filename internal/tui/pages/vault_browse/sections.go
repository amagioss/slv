package vault_browse

import (
	"fmt"

	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/theme"
)

func (vbp *VaultBrowsePage) createMainSection() *tview.Grid {
	colors := theme.GetCurrentPalette()

	// Create the directory list (left column)
	vbp.directoryList = tview.NewList()
	vbp.directoryList.SetTitle(fmt.Sprintf("Directories (%s)", vbp.currentDir)).SetTitleAlign(tview.AlignLeft)
	vbp.directoryList.SetBorder(true)

	// Style the directory list
	vbp.directoryList.SetSelectedTextColor(colors.ListSelectedText).
		SetSelectedBackgroundColor(colors.ListSelectedBg).
		SetSecondaryTextColor(colors.ListSecondaryText).
		SetMainTextColor(colors.ListMainText).
		SetWrapAround(false) // Disable looping behavior

	// Create the file list (right column)
	vbp.fileList = tview.NewList()
	vbp.fileList.SetTitle("Vault Files").SetTitleAlign(tview.AlignLeft)
	vbp.fileList.SetBorder(true)

	// Style the file list
	vbp.fileList.SetSelectedTextColor(colors.ListSelectedText).
		SetSelectedBackgroundColor(colors.ListSelectedBg).
		SetSecondaryTextColor(colors.ListSecondaryText).
		SetMainTextColor(colors.ListMainText).
		SetWrapAround(false) // Disable looping behavior

	// Create a two-column layout using grid
	mainContent := tview.NewGrid().
		SetRows(0).       // Single row taking full height
		SetColumns(0, 0). // Two equal columns
		SetBorders(false)

	// Add the directory list (left column)
	mainContent.AddItem(vbp.directoryList, 0, 0, 1, 1, 0, 0, true)

	// Add the file list (right column)
	mainContent.AddItem(vbp.fileList, 0, 1, 1, 1, 0, 0, false)

	return mainContent
}

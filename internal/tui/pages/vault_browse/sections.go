package vault_browse

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (vbp *VaultBrowsePage) createMainSection() *tview.Grid {
	welcomeText := fmt.Sprintf("\n[white]Browse Vaults[white::-]\n[gray](Use arrow keys [←] and [→] to navigate directories)[gray::-]\n\nCurrent Directory: %s", vbp.currentDir)

	vbp.pwdTextView = tview.NewTextView().
		SetText(welcomeText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetWrap(true)

	// Create the list
	vbp.fileList = tview.NewList()

	// Style the list
	vbp.fileList.SetSelectedTextColor(tcell.ColorYellow).
		SetSelectedBackgroundColor(tcell.ColorNavy).
		SetSecondaryTextColor(tcell.ColorGray).
		SetMainTextColor(tcell.ColorWhite)

	// Create a centered layout using grid
	mainContent := tview.NewGrid().
		SetRows(6, 0). // Two flexible rows
		SetColumns(0). // Single column
		SetBorders(false)

	// Center the welcome text
	mainContent.AddItem(vbp.pwdTextView, 0, 0, 1, 1, 0, 0, false)

	// Center the list
	mainContent.AddItem(vbp.fileList, 1, 0, 1, 1, 0, 0, true)

	return mainContent
}

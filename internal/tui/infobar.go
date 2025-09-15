package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/session"
)

func (t *TUI) createInfoBar() {
	profileName := getProfileName()
	var selfEnvironment *environments.Environment

	session, err := session.GetSession()
	if err != nil {
		session = nil
	}
	if session != nil {
		selfEnvironment, err = session.Env()
		if err != nil {
			selfEnvironment = nil
		}
	}

	// Create a table for the info (left side)
	infoTable := tview.NewTable()
	infoTable.SetBorder(false) // No borders as requested
	infoTable.SetBorders(false)
	infoTable.SetSelectable(false, false) // Not selectable

	// Add profile and environment info to the table
	row := 0
	infoTable.SetCell(row, 0, tview.NewTableCell("Profile:").SetTextColor(tcell.ColorDarkCyan).SetMaxWidth(20))
	infoTable.SetCell(row, 1, tview.NewTableCell(profileName).SetTextColor(tcell.ColorWhite))
	row++

	if selfEnvironment != nil {
		// Environment exists - show all details
		envName := selfEnvironment.Name
		envEmail := selfEnvironment.Email
		envType := string(selfEnvironment.EnvType)

		var publicKey string
		if pubKey, err := selfEnvironment.GetPublicKey(); err == nil && pubKey != nil {
			if keyStr, err := pubKey.String(); err == nil {
				publicKey = keyStr
			} else {
				publicKey = "Error getting key"
			}
		} else {
			publicKey = "No key available"
		}

		infoTable.SetCell(row, 0, tview.NewTableCell("Environment:").SetTextColor(tcell.ColorDarkCyan).SetMaxWidth(20))
		infoTable.SetCell(row, 1, tview.NewTableCell(envName).SetTextColor(tcell.ColorWhite))
		row++

		infoTable.SetCell(row, 0, tview.NewTableCell("Email:").SetTextColor(tcell.ColorDarkCyan).SetMaxWidth(20))
		infoTable.SetCell(row, 1, tview.NewTableCell(envEmail).SetTextColor(tcell.ColorWhite))
		row++

		infoTable.SetCell(row, 0, tview.NewTableCell("Type:").SetTextColor(tcell.ColorDarkCyan).SetMaxWidth(20))
		infoTable.SetCell(row, 1, tview.NewTableCell(envType).SetTextColor(tcell.ColorWhite))
		row++

		infoTable.SetCell(row, 0, tview.NewTableCell("Public Key:").SetTextColor(tcell.ColorDarkCyan).SetMaxWidth(20))
		infoTable.SetCell(row, 1, tview.NewTableCell(publicKey).SetTextColor(tcell.ColorWhite))
	} else {
		// No environment - show minimal info
		infoTable.SetCell(row, 0, tview.NewTableCell("Status:").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
		infoTable.SetCell(row, 1, tview.NewTableCell("No self environment is set").SetTextColor(tcell.ColorYellow))
	}

	// Create logo content (right side)
	logoContent := ` ____  _ __     __
/ ___|| |\ \   / /
\___ \| | \ \ / / 
 ___) | |__\ V /  
|____/|_____\_/   
`

	logoTextView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetWrap(false).
		SetText(logoContent)

	// Create flex container to hold both the table and logo
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).     // Set direction to column (horizontal)
		AddItem(infoTable, 0, 1, false).    // Info table takes remaining space
		AddItem(logoTextView, 30, 0, false) // Logo has fixed width

	flex.SetBorder(true).
		SetBorderColor(t.theme.Accent).
		SetTitle("Secure Local Vault").
		SetTitleAlign(tview.AlignCenter)

	t.infoBar = flex
}

func getProfileName() string {
	profile, err := profiles.GetActiveProfile()
	if err != nil || profile == nil {
		return "No Profile"
	}
	return profile.Name()
}

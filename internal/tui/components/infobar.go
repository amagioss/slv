package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/tui/interfaces"
)

// InfoBar represents the info bar component
type InfoBar struct {
	tui       interfaces.TUIInterface
	primitive tview.Primitive
	infoTable *tview.Table
	logoView  *tview.TextView
}

// NewInfoBar creates a new InfoBar component
func NewInfoBar(tui interfaces.TUIInterface) *InfoBar {
	ib := &InfoBar{
		tui: tui,
	}
	ib.createComponents()
	ib.Refresh()
	return ib
}

// createComponents creates the underlying UI components
func (ib *InfoBar) createComponents() {
	// Create info table
	ib.infoTable = tview.NewTable()
	ib.infoTable.SetBorder(false)
	ib.infoTable.SetBorders(false)
	ib.infoTable.SetSelectable(false, false)

	// Create logo view
	logoContent := ` ____  _ __     __
/ ___|| |\ \   / /
\___ \| | \ \ / / 
 ___) | |__\ V /  
|____/|_____\_/   
`

	ib.logoView = tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetWrap(false).
		SetText(logoContent)

	// Create flex container
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(ib.infoTable, 0, 1, false).
		AddItem(ib.logoView, 30, 0, false)

	flex.SetBorder(true).
		SetBorderColor(tcell.ColorAqua).
		SetTitle("Secure Local Vault").
		SetTitleAlign(tview.AlignCenter)

	ib.primitive = flex
}

// Render returns the primitive for this component
func (ib *InfoBar) Render() tview.Primitive {
	return ib.primitive
}

// Refresh refreshes the component with current data
func (ib *InfoBar) Refresh() {
	profileName := ib.getProfileName()
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

	// Clear existing cells
	ib.infoTable.Clear()

	// Add profile and environment info to the table
	row := 0
	ib.infoTable.SetCell(row, 0, tview.NewTableCell("Profile:").SetTextColor(tcell.ColorDarkCyan).SetMaxWidth(20))
	ib.infoTable.SetCell(row, 1, tview.NewTableCell(profileName).SetTextColor(tcell.ColorWhite))
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

		ib.infoTable.SetCell(row, 0, tview.NewTableCell("Environment:").SetTextColor(tcell.ColorDarkCyan).SetMaxWidth(20))
		ib.infoTable.SetCell(row, 1, tview.NewTableCell(envName).SetTextColor(tcell.ColorWhite))
		row++

		ib.infoTable.SetCell(row, 0, tview.NewTableCell("Email:").SetTextColor(tcell.ColorDarkCyan).SetMaxWidth(20))
		ib.infoTable.SetCell(row, 1, tview.NewTableCell(envEmail).SetTextColor(tcell.ColorWhite))
		row++

		ib.infoTable.SetCell(row, 0, tview.NewTableCell("Type:").SetTextColor(tcell.ColorDarkCyan).SetMaxWidth(20))
		ib.infoTable.SetCell(row, 1, tview.NewTableCell(envType).SetTextColor(tcell.ColorWhite))
		row++

		ib.infoTable.SetCell(row, 0, tview.NewTableCell("Public Key:").SetTextColor(tcell.ColorDarkCyan).SetMaxWidth(20))
		ib.infoTable.SetCell(row, 1, tview.NewTableCell(publicKey).SetTextColor(tcell.ColorWhite))
	} else {
		// No environment - show minimal info
		ib.infoTable.SetCell(row, 0, tview.NewTableCell("Status:").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
		ib.infoTable.SetCell(row, 1, tview.NewTableCell("No self environment is set").SetTextColor(tcell.ColorYellow))
	}
}

// SetFocus sets focus on the component
func (ib *InfoBar) SetFocus(focus bool) {
	// Info bar is not focusable
}

// UpdateProfile updates the profile information
func (ib *InfoBar) UpdateProfile(profileName string) {
	// This could be used for real-time updates if needed
	ib.Refresh()
}

// UpdateEnvironment updates the environment information
func (ib *InfoBar) UpdateEnvironment(envName, envEmail, envType, publicKey string) {
	// This could be used for real-time updates if needed
	ib.Refresh()
}

// ShowNoEnvironment shows the no environment state
func (ib *InfoBar) ShowNoEnvironment() {
	// This could be used for real-time updates if needed
	ib.Refresh()
}

// getProfileName gets the current profile name
func (ib *InfoBar) getProfileName() string {
	profile, err := profiles.GetActiveProfile()
	if err != nil || profile == nil {
		return "No Profile"
	}
	return profile.Name()
}

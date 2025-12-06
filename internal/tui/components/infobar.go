package components

import (
	"runtime"
	"time"

	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/theme"
)

// InfoBar represents the info bar component
type InfoBar struct {
	tui          interfaces.TUIInterface
	primitive    tview.Primitive
	infoTable    *tview.Table
	versionTable *tview.Table
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
	colors := theme.GetCurrentPalette()
	// Create info table
	ib.infoTable = tview.NewTable()
	ib.infoTable.SetBorder(false)
	ib.infoTable.SetBorders(false)
	ib.infoTable.SetSelectable(false, false)

	// Create version info table
	var versionTableWidth int
	ib.versionTable, versionTableWidth = createVersionTable(colors)

	// Create flex container
	// Add a spacer to push version table to the right, then add version table with fixed width
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(ib.infoTable, 0, 1, false).
		// AddItem(nil, 0, 1, false).                            // Spacer to push version table to the right
		AddItem(ib.versionTable, versionTableWidth, 0, false) // Fixed width, pinned to right

	flex.SetBorder(true).
		SetBorderColor(colors.InfobarBorder).
		SetTitle("Secure Local Vault").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(colors.InfobarTitle)

	ib.primitive = flex
}

// Render returns the primitive for this component
func (ib *InfoBar) Render() tview.Primitive {
	return ib.primitive
}

// Refresh refreshes the component with current data
func (ib *InfoBar) Refresh() {
	colors := theme.GetCurrentPalette()
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
	ib.infoTable.SetCell(row, 0, tview.NewTableCell("Profile:").SetTextColor(colors.InfoTableLabel).SetMaxWidth(20))
	ib.infoTable.SetCell(row, 1, tview.NewTableCell(profileName).SetTextColor(colors.InfoTableValue))
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

		ib.infoTable.SetCell(row, 0, tview.NewTableCell("Environment:").SetTextColor(colors.InfoTableLabel).SetMaxWidth(20))
		ib.infoTable.SetCell(row, 1, tview.NewTableCell(envName).SetTextColor(colors.InfoTableValue))
		row++

		ib.infoTable.SetCell(row, 0, tview.NewTableCell("Email:").SetTextColor(colors.InfoTableLabel).SetMaxWidth(20))
		ib.infoTable.SetCell(row, 1, tview.NewTableCell(envEmail).SetTextColor(colors.InfoTableValue))
		row++

		ib.infoTable.SetCell(row, 0, tview.NewTableCell("Type:").SetTextColor(colors.InfoTableLabel).SetMaxWidth(20))
		ib.infoTable.SetCell(row, 1, tview.NewTableCell(envType).SetTextColor(colors.InfoTableValue))
		row++

		ib.infoTable.SetCell(row, 0, tview.NewTableCell("Public Key:").SetTextColor(colors.InfoTableLabel).SetMaxWidth(20))
		ib.infoTable.SetCell(row, 1, tview.NewTableCell(publicKey).SetTextColor(colors.InfoTableValue))
	} else {
		// No environment - show minimal info
		ib.infoTable.SetCell(row, 0, tview.NewTableCell("Status:").SetTextColor(colors.InfoTableLabel).SetMaxWidth(20))
		ib.infoTable.SetCell(row, 1, tview.NewTableCell("No self environment is set").SetTextColor(colors.InfoTableValue))
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

// createVersionTable creates a table with version information
// Returns the table and its calculated fixed width
func createVersionTable(colors theme.ColorPalette) (*tview.Table, int) {
	var committedAt string
	if builtAtTime, err := time.Parse(time.RFC3339, config.GetCommitDate()); err == nil {
		builtAtLocalTime := builtAtTime.Local()
		committedAt = builtAtLocalTime.Format("02 Jan 2006 03:04:05 PM MST")
	}

	versionTable := tview.NewTable().
		SetBorders(false)

	// Calculate maximum width needed for the table
	// Format: "Label : Value"
	maxWidth := 0
	rows := []struct {
		label string
		value string
	}{
		{"SLV Version", config.GetVersion()},
		{"Built At", committedAt},
		{"Release", config.GetReleaseURL()},
		{"Git Commit", config.GetFullCommit()},
		{"Web", "https://slv.sh"},
		{"Platform", runtime.GOOS + "/" + runtime.GOARCH},
		{"Go Version", runtime.Version()},
	}

	for _, row := range rows {
		// Calculate width: label + " : " + value
		width := len(row.label) + 3 + len(row.value)
		if width > maxWidth {
			maxWidth = width
		}
	}

	// Add padding for borders and spacing (typically 4-6 characters)
	// The Git Commit hash is usually the longest, so we ensure enough space
	fixedWidth := maxWidth + 6

	addVersionRow := func(label, value string) {
		row := versionTable.GetRowCount()
		versionTable.SetCell(row, 0, tview.NewTableCell(label).
			SetTextColor(colors.TextSecondary).
			SetAlign(tview.AlignLeft))
		versionTable.SetCell(row, 1, tview.NewTableCell(" : "+value).
			SetTextColor(colors.TextPrimary).
			SetAlign(tview.AlignLeft))
	}

	for _, row := range rows {
		addVersionRow(row.label, row.value)
	}

	return versionTable, fixedWidth
}

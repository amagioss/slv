package mainpage

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/theme"
)

// MainPage handles the main menu page functionality
type MainPage struct {
	pages.BasePage
	list       *tview.List
	navigation *FormNavigation
}

// NewMainPage creates a new MainPage instance
func NewMainPage(tui interfaces.TUIInterface) *MainPage {
	return &MainPage{
		BasePage: *pages.NewBasePage(tui, "Main Menu"),
	}
}

// Create implements the Page interface
func (mp *MainPage) Create() tview.Primitive {
	// Create welcome message parts with subtle color variations
	colors := theme.GetCurrentPalette()

	leftPanel, leftPanelWidth := createLogoPanel(colors)

	// Create a flex container for the right panel
	rightPanel := tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Create the main menu list with enhanced styling
	mp.list = tview.NewList().
		AddItem("Vaults", "Manage and organize your vaults", 'v', func() {
			mp.NavigateTo("vaults", false)
		}).
		AddItem("Profiles", "View Profile settings and Environments", 'p', func() {
			mp.NavigateTo("profiles", false)
		}).
		AddItem("Environments", "Manage Environments", 'e', func() {
			mp.NavigateTo("environments", false)
		}).
		AddItem("Help", "View documentation and help", 'h', func() {
			mp.NavigateTo("help", false)
		})

	// Style the list
	mp.list.SetSelectedTextColor(colors.ListSelectedText).
		SetSelectedBackgroundColor(colors.ListSelectedBg).
		SetSecondaryTextColor(colors.ListSecondaryText).
		SetMainTextColor(colors.ListMainText).
		SetWrapAround(false) // Disable looping behavior
	listRow := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 20, 0, false).
		AddItem(mp.list, 0, 1, true).
		AddItem(nil, 4, 0, false)

	// Center the list vertically in the right panel
	rightPanel.AddItem(nil, 0, 1, false).
		AddItem(listRow, 0, 2, true).
		AddItem(nil, 0, 1, false)

	rightPanel.SetBorder(true).
		SetBorderColor(colors.Border)

	// Create a two-column layout using Flex.
	// We give the left panel a fixed width (art width + padding).
	content := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftPanel, leftPanelWidth, 0, false).
		AddItem(rightPanel, 0, 1, true)

	// Set up navigation
	mp.navigation = (&FormNavigation{}).NewFormNavigation(mp)
	mp.navigation.SetupNavigation()

	// Create layout using BasePage method
	return mp.CreateLayout(content)
}
func createLogoPanel(colors theme.ColorPalette) (tview.Primitive, int) {
	art := config.Art()
	coloredArt := strings.ReplaceAll(art, "▓", "[#9d3a4f]▓[-]")
	coloredArt = strings.ReplaceAll(coloredArt, "░", "[#4f5559]░[-]")
	coloredArt = strings.ReplaceAll(coloredArt, "▒", "[#4f5559]▒[-]")

	artLines := strings.Split(art, "\n")
	maxWidth := 0
	for _, line := range artLines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	logoText := tview.NewTextView().
		SetText(coloredArt).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetWrap(false)
		// SetBackgroundColor(colors.Background)

	infoLines := "[#9d3a4f]Made with ❤️  from 🇮🇳  for a secure decentralized future.[-]"
	footerLen := len("Made with ❤️  from 🇮🇳  for a secure decentralized future.") + 4 // +4 for emoji width correction/padding

	infoView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(infoLines).
		SetTextColor(colors.TextPrimary)
	// SetBackgroundColor(colors.Background)

	// Add dynamic version info
	var committedAt string
	if builtAtTime, err := time.Parse(time.RFC3339, config.GetCommitDate()); err == nil {
		builtAtLocalTime := builtAtTime.Local()
		committedAt = builtAtLocalTime.Format("02 Jan 2006 03:04:05 PM MST")
	}

	versionTable := tview.NewTable().
		SetBorders(false)

	addVersionRow := func(label, value string) {
		row := versionTable.GetRowCount()
		versionTable.SetCell(row, 0, tview.NewTableCell(label).
			SetTextColor(colors.TextSecondary).
			SetAlign(tview.AlignLeft))
		versionTable.SetCell(row, 1, tview.NewTableCell(" : "+value).
			SetTextColor(colors.TextPrimary).
			SetAlign(tview.AlignLeft))
	}

	addVersionRow("SLV Version", config.GetVersion())
	addVersionRow("Built At", committedAt)
	addVersionRow("Release", config.GetReleaseURL())
	addVersionRow("Git Commit", config.GetFullCommit())
	addVersionRow("Web", "https://slv.sh")
	addVersionRow("Platform", runtime.GOOS+"/"+runtime.GOARCH)
	addVersionRow("Go Version", runtime.Version())

	// Description above the table
	descriptionText := fmt.Sprintf(`[#f2f2f2]%s (Secure Local Vault)[-] : Securely store, share, and access secrets alongside the codebase.`, config.AppNameUpperCase)
	descriptionView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(descriptionText).
		SetTextColor(colors.TextPrimary).
		SetWrap(true)

	// Combine description and table in a vertical flex
	// Use a horizontal flex to center the table horizontally
	tableFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 30, 0, false).
		AddItem(versionTable, 0, 8, false). // Give table 80% width relative to spacers
		AddItem(nil, 0, 1, false)

	infoFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(descriptionView, 3, 1, false). // Fixed height for description
		AddItem(tview.NewBox(), 1, 0, false).  // Spacer
		AddItem(tableFlex, 0, 1, false)        // Table centered horizontally takes remaining vertical space

	// Use a flex container to center the logo vertically
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(logoText, len(artLines)+2, 0, false). // Fixed height for art
		AddItem(infoFlex, 12, 0, false).              // Flexible height for version info
		AddItem(nil, 4, 0, false).                    // Spacer
		AddItem(infoView, 2, 0, false).               // Made with... footer
		AddItem(nil, 0, 1, false)

	flex.SetBorder(true).
		SetBorderColor(colors.Border)
	// SetBackgroundColor(colors.Background)

	// Determine optimal width for the left panel
	finalWidth := maxWidth
	if footerLen > finalWidth {
		finalWidth = footerLen
	}

	// Return flex and calculated width (content width + padding + border)
	return flex, finalWidth
}

// Refresh implements the Page interface
func (mp *MainPage) Refresh() {
	// Main page doesn't need refresh logic
}

// HandleInput implements the Page interface
func (mp *MainPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// Main page doesn't handle specific input
	return event
}

// GetTitle implements the Page interface
func (mp *MainPage) GetTitle() string {
	return mp.BasePage.GetTitle()
}

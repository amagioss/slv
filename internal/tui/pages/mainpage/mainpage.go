package mainpage

import (
	"strings"

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
		AddItem(nil, 10, 0, false).
		AddItem(mp.list, 0, 1, true).
		AddItem(nil, 4, 0, false)

	// Center the list vertically in the right panel
	rightPanel.AddItem(nil, 0, 1, false).
		AddItem(listRow, 10, 0, true).
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

	// Add "Secure Local Vault" text below the logo
	titleText := "[#9d3a4f]Secure Local Vault[-]"
	titleView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(titleText).
		SetTextColor(colors.TextPrimary)

	infoLines := "[#858a8e]Made with ❤️  from 🇮🇳  for a secure decentralized future.[-]"
	footerLen := len("Made with ❤️  from 🇮🇳  for a secure decentralized future.") + 4 // +4 for emoji width correction/padding
	titleLen := len("Secure Local Vault") + 4                                         // +4 for padding

	infoView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(infoLines).
		SetTextColor(colors.TextPrimary)
	// SetBackgroundColor(colors.Background)

	// Use a flex container to center the logo vertically
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(logoText, len(artLines)+1, 0, false). // Fixed height for art
		AddItem(nil, 1, 0, false).                    // Minimal spacer between logo and title
		AddItem(titleView, 1, 0, false).              // "Secure Local Vault" title
		AddItem(nil, 1, 0, false).                    // Minimal spacer between title and footer
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
	if titleLen > finalWidth {
		finalWidth = titleLen
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

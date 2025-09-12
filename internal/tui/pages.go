package tui

import (
	"github.com/rivo/tview"
)

// createVaultsPage creates the vaults page
func (t *TUI) createVaultsPage() tview.Primitive {
	text := tview.NewTextView().
		SetText("Vaults Page\n\nThis page will show vault management options.").
		SetTextAlign(tview.AlignCenter)

	return t.createPageLayout("Vaults", text)
}

// createProfilesPage creates the profiles page
func (t *TUI) createProfilesPage() tview.Primitive {
	text := tview.NewTextView().
		SetText("Profiles Page\n\nThis page will show profile management options.").
		SetTextAlign(tview.AlignCenter)

	return t.createPageLayout("Profiles", text)
}

// createEnvironmentsPage creates the environments page
func (t *TUI) createEnvironmentsPage() tview.Primitive {
	text := tview.NewTextView().
		SetText("Environments Page\n\nThis page will show environment management options.").
		SetTextAlign(tview.AlignCenter)

	return t.createPageLayout("Environments", text)
}

// createHelpPage creates the help page
func (t *TUI) createHelpPage() tview.Primitive {
	text := tview.NewTextView().
		SetText("SLV TUI Help\n\nNavigation:\n- Arrow keys: Navigate\n- Enter: Select\n- Esc: Back\n- Ctrl+C: Quit\n\nShortcuts:\n- v: Vaults\n- p: Profiles\n- e: Environments\n- h: Help")

	return t.createPageLayout("Help", text)
}

// createPageLayout creates a common layout for all pages with border
func (t *TUI) createPageLayout(title string, content tview.Primitive) tview.Primitive {
	theme := t.GetTheme()

	// Create a flex container
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Set border properties
	flex.SetBorder(true).
		SetBorderColor(theme.Accent).
		SetTitle(title).
		SetTitleAlign(tview.AlignCenter)

	// Add the content to the flex
	flex.AddItem(content, 0, 1, true)

	return flex
}

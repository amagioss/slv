package vault_edit

import (
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/theme"
)

func (vep *VaultEditPage) applyDisabledStyling(component tview.Primitive) {
	colors := theme.GetCurrentPalette()
	switch c := component.(type) {
	case *tview.Form:
		// Use very muted colors for disabled appearance
		c.SetLabelColor(colors.FormDisabledText).
			SetTitleColor(colors.FormDisabledText)
	case *tview.List:
		// Apply muted colors to lists
		c.SetSelectedTextColor(colors.FormDisabledText).
			SetSelectedBackgroundColor(colors.FormDisabledBg).
			SetSecondaryTextColor(colors.FormDisabledText).
			SetMainTextColor(colors.FormDisabledText).
			SetBorderColor(colors.FormDisabledBorder).
			SetTitleColor(colors.FormDisabledText)
	case *tview.InputField:
		// Apply muted colors to input fields
		c.SetFieldTextColor(colors.FormDisabledText).
			SetLabelColor(colors.FormDisabledText).
			SetBorderColor(colors.FormDisabledBorder).
			SetTitleColor(colors.FormDisabledText)
	case *tview.Table:
		// Apply muted colors to tables
		c.SetBorderColor(colors.FormDisabledBorder).
			SetTitleColor(colors.FormDisabledText)
		// Note: Table doesn't have SetSelectedTextColor/SetSelectedBackgroundColor methods
	case *tview.Flex:
		c.SetBorderColor(colors.FormDisabledBorder).
			SetTitleColor(colors.FormDisabledText)
	}
}

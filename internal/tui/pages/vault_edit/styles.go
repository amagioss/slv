package vault_edit

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (vep *VaultEditPage) applyDisabledStyling(component tview.Primitive) {
	switch c := component.(type) {
	case *tview.Form:
		// Use very muted colors for disabled appearance
		// c.SetFieldBackgroundColor(tcell.ColorDarkGray).
		// SetFieldTextColor(tcell.ColorGray).
		c.SetLabelColor(tcell.ColorGray)
		// SetButtonBackgroundColor(tcell.ColorDarkGray).
		// SetButtonTextColor(tcell.ColorGray).
		// SetBorderColor(tcell.ColorGray)
	case *tview.List:
		// Apply muted colors to lists
		c.SetSelectedTextColor(tcell.ColorGray).
			SetSelectedBackgroundColor(tcell.Color16).
			SetSecondaryTextColor(tcell.ColorDarkGray).
			SetMainTextColor(tcell.ColorGray).
			SetBorderColor(tcell.ColorGray)
	case *tview.InputField:
		// Apply muted colors to input fields
		// c.SetFieldBackgroundColor(tcell.ColorDarkGray).
		c.SetFieldTextColor(tcell.ColorGray).
			SetLabelColor(tcell.ColorGray).
			SetBorderColor(tcell.ColorGray)
	case *tview.Table:
		// Apply muted colors to tables
		c.SetBorderColor(tcell.ColorGray)
		// Note: Table doesn't have SetSelectedTextColor/SetSelectedBackgroundColor methods
	}
}

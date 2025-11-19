package components

import (
	"fmt"

	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/theme"
)

// StatusBar represents the status bar component
type StatusBar struct {
	tui        interfaces.TUIInterface
	primitive  *tview.TextView
	pageName   string
	customHelp string
}

// NewStatusBar creates a new StatusBar component
func NewStatusBar(tui interfaces.TUIInterface) *StatusBar {
	sb := &StatusBar{
		tui: tui,
	}
	sb.createComponent()
	sb.UpdateStatus("")
	return sb
}

// createComponent creates the underlying UI component
func (sb *StatusBar) createComponent() {
	colors := theme.GetCurrentPalette()
	sb.primitive = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignLeft).
		SetScrollable(true)

	sb.primitive.SetBorder(true).
		SetBorderColor(colors.Border).
		SetTitle("Status").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(colors.MainContentTitle)
}

// Render returns the primitive for this component
func (sb *StatusBar) Render() tview.Primitive {
	return sb.primitive
}

// Refresh refreshes the component with current data
func (sb *StatusBar) Refresh() {
	sb.UpdateStatus(sb.pageName)
}

// SetFocus sets focus on the component
func (sb *StatusBar) SetFocus(focus bool) {
	// Status bar is not focusable
}

// UpdateStatus updates the status bar with the current page name
func (sb *StatusBar) UpdateStatus(pageName string) {
	sb.pageName = pageName
	colors := theme.GetCurrentPalette()

	status := fmt.Sprintf("Page: %s | F1: Help | Esc: Back | Ctrl+C: Quit", pageName)

	// Add custom help text if available
	if sb.customHelp != "" {
		status += " | " + sb.customHelp
	}

	sb.primitive.SetText(status)
	sb.primitive.SetTextColor(colors.TextPrimary)
}

// SetCustomHelp sets the custom help text for the current page
func (sb *StatusBar) SetCustomHelp(helpText string) {
	sb.customHelp = helpText
	sb.Refresh()
}

// ClearCustomHelp clears the custom help text
func (sb *StatusBar) ClearCustomHelp() {
	sb.customHelp = ""
	sb.Refresh()
}

// SetPageName sets the current page name
func (sb *StatusBar) SetPageName(pageName string) {
	sb.UpdateStatus(pageName)
}

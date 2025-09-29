package theme

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Theme defines the color scheme and styling for the TUI
type Theme struct {
	palette ColorPalette
}

// NewTheme creates a new theme with the default color palette
func NewTheme() *Theme {
	return &Theme{
		palette: GetDefaultPalette(),
	}
}

// NewThemeWithPalette creates a new theme with a specific color palette
func NewThemeWithPalette(palette ColorPalette) *Theme {
	return &Theme{
		palette: palette,
	}
}

// GetPalette returns the current color palette
func (t *Theme) GetPalette() ColorPalette {
	return t.palette
}

// SetPalette sets a new color palette
func (t *Theme) SetPalette(palette ColorPalette) {
	t.palette = palette
}

// ApplyTheme applies the theme to a tview application
func (t *Theme) ApplyTheme(app *tview.Application) {
	// Use tview's global Styles to set the theme
	tview.Styles.PrimitiveBackgroundColor = t.palette.Background
	tview.Styles.ContrastBackgroundColor = t.palette.BackgroundDark
	tview.Styles.MoreContrastBackgroundColor = t.palette.BackgroundDark
	tview.Styles.PrimaryTextColor = t.palette.TextPrimary
	tview.Styles.SecondaryTextColor = t.palette.TextSecondary
	tview.Styles.TertiaryTextColor = t.palette.TextMuted
	tview.Styles.InverseTextColor = t.palette.TextPrimary
	tview.Styles.ContrastSecondaryTextColor = t.palette.TextSecondary
	tview.Styles.BorderColor = t.palette.Border
}

// GetBackground returns the background color
func (t *Theme) GetBackground() tcell.Color {
	return t.palette.Background
}

// GetAccent returns the accent color
func (t *Theme) GetAccent() tcell.Color {
	return t.palette.Accent
}

// GetPrimary returns the primary color
func (t *Theme) GetPrimary() tcell.Color {
	return t.palette.Primary
}

// GetSecondary returns the secondary color
func (t *Theme) GetSecondary() tcell.Color {
	return t.palette.Secondary
}

// GetTextPrimary returns the primary text color
func (t *Theme) GetTextPrimary() tcell.Color {
	return t.palette.TextPrimary
}

// GetTextSecondary returns the secondary text color
func (t *Theme) GetTextSecondary() tcell.Color {
	return t.palette.TextSecondary
}

// GetTextMuted returns the muted text color
func (t *Theme) GetTextMuted() tcell.Color {
	return t.palette.TextMuted
}

// GetSuccess returns the success color
func (t *Theme) GetSuccess() tcell.Color {
	return t.palette.Success
}

// GetWarning returns the warning color
func (t *Theme) GetWarning() tcell.Color {
	return t.palette.Warning
}

// GetError returns the error color
func (t *Theme) GetError() tcell.Color {
	return t.palette.Error
}

// GetInfo returns the info color
func (t *Theme) GetInfo() tcell.Color {
	return t.palette.Info
}

// GetBorder returns the border color
func (t *Theme) GetBorder() tcell.Color {
	return t.palette.Border
}

// GetFocus returns the focus color
func (t *Theme) GetFocus() tcell.Color {
	return t.palette.Focus
}

// GetSelection returns the selection color
func (t *Theme) GetSelection() tcell.Color {
	return t.palette.Selection
}

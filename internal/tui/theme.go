package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Theme defines the color scheme and styling for the TUI
type Theme struct {
	// Primary colors
	Primary   tcell.Color
	Secondary tcell.Color
	Accent    tcell.Color

	// Background colors
	Background     tcell.Color
	BackgroundDark tcell.Color

	// Text colors
	TextPrimary   tcell.Color
	TextSecondary tcell.Color
	TextMuted     tcell.Color

	// Status colors
	Success tcell.Color
	Warning tcell.Color
	Error   tcell.Color
	Info    tcell.Color

	// Border and UI colors
	Border     tcell.Color
	BorderDark tcell.Color

	// Focus and selection
	Focus     tcell.Color
	Selection tcell.Color
}

// NewTheme creates a new theme with elegant dark colors
func NewTheme() *Theme {
	return &Theme{
		// Primary colors - elegant purple/blue palette
		Primary:   tcell.Color141, // Soft purple
		Secondary: tcell.Color75,  // Light blue
		Accent:    tcell.Color87,  // Bright cyan

		// Background colors - deep dark theme
		Background:     tcell.Color16,  // Pure black
		BackgroundDark: tcell.Color232, // Very dark gray

		// Text colors - elegant whites and grays
		TextPrimary:   tcell.Color255, // Pure white
		TextSecondary: tcell.Color250, // Off-white
		TextMuted:     tcell.Color244, // Medium gray

		// Status colors - vibrant but elegant
		Success: tcell.Color82,  // Bright green
		Warning: tcell.Color220, // Golden yellow
		Error:   tcell.Color196, // Bright red
		Info:    tcell.Color75,  // Light blue

		// Border and UI colors - subtle grays
		Border:     tcell.Color240, // Medium gray
		BorderDark: tcell.Color236, // Dark gray

		// Focus and selection - elegant highlights
		Focus:     tcell.Color87, // Bright cyan
		Selection: tcell.Color25, // Dark blue
	}
}

// GetColorStyles returns a map of color styles for different UI elements
func (t *Theme) GetColorStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"primary":      tcell.StyleDefault.Foreground(t.Primary),
		"primary-bold": tcell.StyleDefault.Foreground(t.Primary).Bold(true),
		"secondary":    tcell.StyleDefault.Foreground(t.Secondary),
		"accent":       tcell.StyleDefault.Foreground(t.Accent),
		"accent-bold":  tcell.StyleDefault.Foreground(t.Accent).Bold(true),
		"background":   tcell.StyleDefault.Background(t.Background),
		"text":         tcell.StyleDefault.Foreground(t.TextPrimary),
		"text-bold":    tcell.StyleDefault.Foreground(t.TextPrimary).Bold(true),
		"text-muted":   tcell.StyleDefault.Foreground(t.TextMuted),
		"text-dim":     tcell.StyleDefault.Foreground(t.TextSecondary),
		"success":      tcell.StyleDefault.Foreground(t.Success),
		"success-bold": tcell.StyleDefault.Foreground(t.Success).Bold(true),
		"warning":      tcell.StyleDefault.Foreground(t.Warning),
		"warning-bold": tcell.StyleDefault.Foreground(t.Warning).Bold(true),
		"error":        tcell.StyleDefault.Foreground(t.Error),
		"error-bold":   tcell.StyleDefault.Foreground(t.Error).Bold(true),
		"info":         tcell.StyleDefault.Foreground(t.Info),
		"info-bold":    tcell.StyleDefault.Foreground(t.Info).Bold(true),
		"border":       tcell.StyleDefault.Foreground(t.Border),
		"focus":        tcell.StyleDefault.Foreground(t.Focus).Bold(true),
		"selection":    tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.Selection),
		"highlight":    tcell.StyleDefault.Foreground(t.Background).Background(t.Accent),
	}
}

// GetTableStyles returns styles for table elements
func (t *Theme) GetTableStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"header":        tcell.StyleDefault.Foreground(t.Primary).Bold(true),
		"header-alt":    tcell.StyleDefault.Foreground(t.Accent).Bold(true),
		"cell":          tcell.StyleDefault.Foreground(t.TextPrimary),
		"cell-alt":      tcell.StyleDefault.Foreground(t.TextSecondary).Background(t.BackgroundDark),
		"cell-muted":    tcell.StyleDefault.Foreground(t.TextMuted),
		"selected":      tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.Selection),
		"selected-alt":  tcell.StyleDefault.Foreground(t.Background).Background(t.Accent),
		"border":        tcell.StyleDefault.Foreground(t.Border),
		"border-bright": tcell.StyleDefault.Foreground(t.Primary),
	}
}

// GetFormStyles returns styles for form elements
func (t *Theme) GetFormStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"label":            tcell.StyleDefault.Foreground(t.TextPrimary).Bold(true),
		"label-muted":      tcell.StyleDefault.Foreground(t.TextMuted),
		"field":            tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.BackgroundDark),
		"field-focus":      tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.BackgroundDark).Bold(true),
		"field-border":     tcell.StyleDefault.Foreground(t.Border),
		"button":           tcell.StyleDefault.Foreground(t.Background).Background(t.Primary),
		"button-focus":     tcell.StyleDefault.Foreground(t.Background).Background(t.Accent).Bold(true),
		"button-secondary": tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.BackgroundDark),
		"button-danger":    tcell.StyleDefault.Foreground(t.Background).Background(t.Error),
	}
}

// GetListStyles returns styles for list elements
func (t *Theme) GetListStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"item":          tcell.StyleDefault.Foreground(t.TextPrimary),
		"item-focus":    tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.Selection),
		"item-selected": tcell.StyleDefault.Foreground(t.Background).Background(t.Primary),
		"item-muted":    tcell.StyleDefault.Foreground(t.TextMuted),
		"item-accent":   tcell.StyleDefault.Foreground(t.Accent),
		"secondary":     tcell.StyleDefault.Foreground(t.TextSecondary),
		"divider":       tcell.StyleDefault.Foreground(t.Border),
	}
}

// GetModalStyles returns styles for modal dialogs
func (t *Theme) GetModalStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"background":    tcell.StyleDefault.Background(t.Background).Foreground(t.TextPrimary),
		"overlay":       tcell.StyleDefault.Background(t.BackgroundDark),
		"border":        tcell.StyleDefault.Foreground(t.Border),
		"border-bright": tcell.StyleDefault.Foreground(t.Primary),
		"title":         tcell.StyleDefault.Foreground(t.Primary).Bold(true),
		"title-accent":  tcell.StyleDefault.Foreground(t.Accent).Bold(true),
		"text":          tcell.StyleDefault.Foreground(t.TextPrimary),
		"text-muted":    tcell.StyleDefault.Foreground(t.TextMuted),
		"button":        tcell.StyleDefault.Foreground(t.Background).Background(t.Primary),
		"button-focus":  tcell.StyleDefault.Foreground(t.Background).Background(t.Accent).Bold(true),
		"button-danger": tcell.StyleDefault.Foreground(t.Background).Background(t.Error),
		"button-cancel": tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.BackgroundDark),
	}
}

// ApplyTheme applies the theme to a tview.Application
func (t *Theme) ApplyTheme(app *tview.Application) {
	// Use tview's global Styles to set the theme
	tview.Styles.PrimitiveBackgroundColor = t.Background
	tview.Styles.ContrastBackgroundColor = t.BackgroundDark
	tview.Styles.MoreContrastBackgroundColor = t.BackgroundDark
	tview.Styles.PrimaryTextColor = t.TextPrimary
	tview.Styles.SecondaryTextColor = t.TextSecondary
	tview.Styles.TertiaryTextColor = t.TextMuted
	tview.Styles.InverseTextColor = t.TextPrimary
	tview.Styles.ContrastSecondaryTextColor = t.TextSecondary
	tview.Styles.BorderColor = t.Border

	// Note: Background color is set via tview.Styles above
}

// GetStatusStyles returns styles for status indicators
func (t *Theme) GetStatusStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"online":   tcell.StyleDefault.Foreground(t.Success).Bold(true),
		"offline":  tcell.StyleDefault.Foreground(t.Error),
		"pending":  tcell.StyleDefault.Foreground(t.Warning),
		"loading":  tcell.StyleDefault.Foreground(t.Info),
		"locked":   tcell.StyleDefault.Foreground(t.Warning).Bold(true),
		"unlocked": tcell.StyleDefault.Foreground(t.Success).Bold(true),
		"error":    tcell.StyleDefault.Foreground(t.Error).Bold(true),
		"success":  tcell.StyleDefault.Foreground(t.Success).Bold(true),
	}
}

// GetVaultStyles returns specialized styles for vault-related elements
func (t *Theme) GetVaultStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"vault-name":       tcell.StyleDefault.Foreground(t.Primary).Bold(true),
		"vault-path":       tcell.StyleDefault.Foreground(t.TextMuted),
		"vault-locked":     tcell.StyleDefault.Foreground(t.Warning).Bold(true),
		"vault-unlocked":   tcell.StyleDefault.Foreground(t.Success).Bold(true),
		"item-name":        tcell.StyleDefault.Foreground(t.Accent),
		"item-value":       tcell.StyleDefault.Foreground(t.TextPrimary),
		"item-secret":      tcell.StyleDefault.Foreground(t.TextMuted),
		"accessor-self":    tcell.StyleDefault.Foreground(t.Success).Bold(true),
		"accessor-root":    tcell.StyleDefault.Foreground(t.Primary).Bold(true),
		"accessor-user":    tcell.StyleDefault.Foreground(t.Accent),
		"accessor-service": tcell.StyleDefault.Foreground(t.Info),
	}
}

// GetElegantPalette returns a curated set of elegant colors
func (t *Theme) GetElegantPalette() map[string]tcell.Color {
	return map[string]tcell.Color{
		"midnight": tcell.Color16,  // Pure black
		"charcoal": tcell.Color232, // Very dark gray
		"slate":    tcell.Color236, // Dark gray
		"steel":    tcell.Color240, // Medium gray
		"silver":   tcell.Color244, // Light gray
		"pearl":    tcell.Color250, // Off-white
		"snow":     tcell.Color255, // Pure white
		"amethyst": tcell.Color141, // Soft purple
		"azure":    tcell.Color75,  // Light blue
		"cyan":     tcell.Color87,  // Bright cyan
		"emerald":  tcell.Color82,  // Bright green
		"gold":     tcell.Color220, // Golden yellow
		"ruby":     tcell.Color196, // Bright red
		"ocean":    tcell.Color25,  // Dark blue
	}
}

// GetBackground returns the background color
func (t *Theme) GetBackground() tcell.Color {
	return t.Background
}

// GetTextPrimary returns the primary text color
func (t *Theme) GetTextPrimary() tcell.Color {
	return t.TextPrimary
}

// GetTextSecondary returns the secondary text color
func (t *Theme) GetTextSecondary() tcell.Color {
	return t.TextSecondary
}

// GetTextMuted returns the muted text color
func (t *Theme) GetTextMuted() tcell.Color {
	return t.TextMuted
}

// GetPrimary returns the primary color
func (t *Theme) GetPrimary() tcell.Color {
	return t.Primary
}

// GetSecondary returns the secondary color
func (t *Theme) GetSecondary() tcell.Color {
	return t.Secondary
}

// GetAccent returns the accent color
func (t *Theme) GetAccent() tcell.Color {
	return t.Accent
}

// GetSuccess returns the success color
func (t *Theme) GetSuccess() tcell.Color {
	return t.Success
}

// GetWarning returns the warning color
func (t *Theme) GetWarning() tcell.Color {
	return t.Warning
}

// GetError returns the error color
func (t *Theme) GetError() tcell.Color {
	return t.Error
}

// GetInfo returns the info color
func (t *Theme) GetInfo() tcell.Color {
	return t.Info
}

// GetBorder returns the border color
func (t *Theme) GetBorder() tcell.Color {
	return t.Border
}

// GetFocus returns the focus color
func (t *Theme) GetFocus() tcell.Color {
	return t.Focus
}

// GetSelection returns the selection color
func (t *Theme) GetSelection() tcell.Color {
	return t.Selection
}

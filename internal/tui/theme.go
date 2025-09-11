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

// NewTheme creates a new theme with default colors
func NewTheme() *Theme {
	return &Theme{
		// Primary colors - using a modern blue palette
		Primary:   tcell.ColorBlue,
		Secondary: tcell.ColorNavy,
		Accent:    tcell.ColorAqua,

		// Background colors
		Background:     tcell.ColorDefault,
		BackgroundDark: tcell.ColorDefault,

		// Text colors
		TextPrimary:   tcell.ColorWhite,
		TextSecondary: tcell.ColorSilver,
		TextMuted:     tcell.ColorGray,

		// Status colors
		Success: tcell.ColorGreen,
		Warning: tcell.ColorYellow,
		Error:   tcell.ColorRed,
		Info:    tcell.ColorBlue,

		// Border and UI colors
		Border:     tcell.ColorSilver,
		BorderDark: tcell.ColorGray,

		// Focus and selection
		Focus:     tcell.ColorAqua,
		Selection: tcell.ColorBlue,
	}
}

// GetColorStyles returns a map of color styles for different UI elements
func (t *Theme) GetColorStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"primary":    tcell.StyleDefault.Foreground(t.Primary),
		"secondary":  tcell.StyleDefault.Foreground(t.Secondary),
		"accent":     tcell.StyleDefault.Foreground(t.Accent),
		"background": tcell.StyleDefault.Background(t.Background),
		"text":       tcell.StyleDefault.Foreground(t.TextPrimary),
		"text-muted": tcell.StyleDefault.Foreground(t.TextMuted),
		"success":    tcell.StyleDefault.Foreground(t.Success),
		"warning":    tcell.StyleDefault.Foreground(t.Warning),
		"error":      tcell.StyleDefault.Foreground(t.Error),
		"info":       tcell.StyleDefault.Foreground(t.Info),
		"border":     tcell.StyleDefault.Foreground(t.Border),
		"focus":      tcell.StyleDefault.Foreground(t.Focus).Bold(true),
		"selection":  tcell.StyleDefault.Foreground(t.Selection).Background(t.BackgroundDark),
	}
}

// GetTableStyles returns styles for table elements
func (t *Theme) GetTableStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"header":   tcell.StyleDefault.Foreground(t.Primary).Bold(true),
		"cell":     tcell.StyleDefault.Foreground(t.TextPrimary),
		"cell-alt": tcell.StyleDefault.Foreground(t.TextSecondary).Background(t.BackgroundDark),
		"selected": tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.Selection),
		"border":   tcell.StyleDefault.Foreground(t.Border),
	}
}

// GetFormStyles returns styles for form elements
func (t *Theme) GetFormStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"label":        tcell.StyleDefault.Foreground(t.TextPrimary),
		"field":        tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.BackgroundDark),
		"field-focus":  tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.BackgroundDark).Bold(true),
		"button":       tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.Primary),
		"button-focus": tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.Accent).Bold(true),
	}
}

// GetListStyles returns styles for list elements
func (t *Theme) GetListStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"item":       tcell.StyleDefault.Foreground(t.TextPrimary),
		"item-focus": tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.Selection),
		"selected":   tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.Primary),
		"secondary":  tcell.StyleDefault.Foreground(t.TextSecondary),
	}
}

// GetModalStyles returns styles for modal dialogs
func (t *Theme) GetModalStyles() map[string]tcell.Style {
	return map[string]tcell.Style{
		"background":   tcell.StyleDefault.Background(t.Background).Foreground(t.TextPrimary),
		"border":       tcell.StyleDefault.Foreground(t.Border),
		"title":        tcell.StyleDefault.Foreground(t.Primary).Bold(true),
		"text":         tcell.StyleDefault.Foreground(t.TextPrimary),
		"button":       tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.Primary),
		"button-focus": tcell.StyleDefault.Foreground(t.TextPrimary).Background(t.Accent).Bold(true),
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
}

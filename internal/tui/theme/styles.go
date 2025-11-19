package theme

import (
	"github.com/gdamore/tcell/v2"
)

// StyleManager manages color styles for different UI elements
type StyleManager struct {
	theme *Theme
}

// NewStyleManager creates a new style manager
func NewStyleManager(theme *Theme) *StyleManager {
	return &StyleManager{
		theme: theme,
	}
}

// GetColorStyles returns a map of color styles for different UI elements
func (sm *StyleManager) GetColorStyles() map[string]tcell.Style {
	palette := sm.theme.GetPalette()
	return map[string]tcell.Style{
		"primary":      tcell.StyleDefault.Foreground(palette.Primary),
		"primary-bold": tcell.StyleDefault.Foreground(palette.Primary).Bold(true),
		"secondary":    tcell.StyleDefault.Foreground(palette.Secondary),
		"accent":       tcell.StyleDefault.Foreground(palette.Accent),
		"accent-bold":  tcell.StyleDefault.Foreground(palette.Accent).Bold(true),
		"background":   tcell.StyleDefault.Background(palette.Background),
		"text":         tcell.StyleDefault.Foreground(palette.TextPrimary),
		"text-bold":    tcell.StyleDefault.Foreground(palette.TextPrimary).Bold(true),
		"text-muted":   tcell.StyleDefault.Foreground(palette.TextMuted),
		"text-dim":     tcell.StyleDefault.Foreground(palette.TextSecondary),
		"success":      tcell.StyleDefault.Foreground(palette.Success),
		"success-bold": tcell.StyleDefault.Foreground(palette.Success).Bold(true),
		"warning":      tcell.StyleDefault.Foreground(palette.Warning),
		"warning-bold": tcell.StyleDefault.Foreground(palette.Warning).Bold(true),
		"error":        tcell.StyleDefault.Foreground(palette.Error),
		"error-bold":   tcell.StyleDefault.Foreground(palette.Error).Bold(true),
		"info":         tcell.StyleDefault.Foreground(palette.Info),
		"info-bold":    tcell.StyleDefault.Foreground(palette.Info).Bold(true),
		"border":       tcell.StyleDefault.Foreground(palette.Border),
		"focus":        tcell.StyleDefault.Foreground(palette.Focus).Bold(true),
		"selection":    tcell.StyleDefault.Foreground(palette.TextPrimary).Background(palette.Selection),
		"highlight":    tcell.StyleDefault.Foreground(palette.Background).Background(palette.Accent),
	}
}

// GetTableStyles returns styles for table elements
func (sm *StyleManager) GetTableStyles() map[string]tcell.Style {
	palette := sm.theme.GetPalette()
	return map[string]tcell.Style{
		"header":        tcell.StyleDefault.Foreground(palette.Primary).Bold(true),
		"header-alt":    tcell.StyleDefault.Foreground(palette.Accent).Bold(true),
		"cell":          tcell.StyleDefault.Foreground(palette.TextPrimary),
		"cell-alt":      tcell.StyleDefault.Foreground(palette.TextSecondary).Background(palette.BackgroundDark),
		"cell-muted":    tcell.StyleDefault.Foreground(palette.TextMuted),
		"selected":      tcell.StyleDefault.Foreground(palette.TextPrimary).Background(palette.Selection),
		"selected-alt":  tcell.StyleDefault.Foreground(palette.Background).Background(palette.Accent),
		"border":        tcell.StyleDefault.Foreground(palette.Border),
		"border-bright": tcell.StyleDefault.Foreground(palette.Primary),
	}
}

// GetListStyles returns styles for list elements
func (sm *StyleManager) GetListStyles() map[string]tcell.Style {
	palette := sm.theme.GetPalette()
	return map[string]tcell.Style{
		"main-text":         tcell.StyleDefault.Foreground(palette.TextPrimary),
		"secondary-text":    tcell.StyleDefault.Foreground(palette.TextSecondary),
		"shortcut":          tcell.StyleDefault.Foreground(palette.Accent),
		"selected":          tcell.StyleDefault.Foreground(palette.TextPrimary).Background(palette.Selection),
		"selected-shortcut": tcell.StyleDefault.Foreground(palette.Accent).Background(palette.Selection),
	}
}

// GetFormStyles returns styles for form elements
func (sm *StyleManager) GetFormStyles() map[string]tcell.Style {
	palette := sm.theme.GetPalette()
	return map[string]tcell.Style{
		"label":          tcell.StyleDefault.Foreground(palette.TextPrimary),
		"field":          tcell.StyleDefault.Foreground(palette.TextPrimary).Background(palette.BackgroundDark),
		"field-focused":  tcell.StyleDefault.Foreground(palette.TextPrimary).Background(palette.BackgroundDark),
		"button":         tcell.StyleDefault.Foreground(palette.TextPrimary).Background(palette.Primary),
		"button-focused": tcell.StyleDefault.Foreground(palette.TextPrimary).Background(palette.Accent),
	}
}

// GetModalStyles returns styles for modal elements
func (sm *StyleManager) GetModalStyles() map[string]tcell.Style {
	palette := sm.theme.GetPalette()
	return map[string]tcell.Style{
		"background":     tcell.StyleDefault.Background(palette.BackgroundDark),
		"border":         tcell.StyleDefault.Foreground(palette.Border),
		"title":          tcell.StyleDefault.Foreground(palette.TextPrimary).Bold(true),
		"text":           tcell.StyleDefault.Foreground(palette.TextPrimary),
		"button":         tcell.StyleDefault.Foreground(palette.TextPrimary).Background(palette.Primary),
		"button-focused": tcell.StyleDefault.Foreground(palette.TextPrimary).Background(palette.Accent),
	}
}

// GetStatusStyles returns styles for status elements
func (sm *StyleManager) GetStatusStyles() map[string]tcell.Style {
	palette := sm.theme.GetPalette()
	return map[string]tcell.Style{
		"normal":    tcell.StyleDefault.Foreground(palette.TextPrimary),
		"success":   tcell.StyleDefault.Foreground(palette.Success),
		"warning":   tcell.StyleDefault.Foreground(palette.Warning),
		"error":     tcell.StyleDefault.Foreground(palette.Error),
		"info":      tcell.StyleDefault.Foreground(palette.Info),
		"highlight": tcell.StyleDefault.Foreground(palette.Accent),
		"online":    tcell.StyleDefault.Foreground(palette.Success).Bold(true),
		"offline":   tcell.StyleDefault.Foreground(palette.Error),
		"pending":   tcell.StyleDefault.Foreground(palette.Warning),
		"loading":   tcell.StyleDefault.Foreground(palette.Info),
		"locked":    tcell.StyleDefault.Foreground(palette.Warning).Bold(true),
		"unlocked":  tcell.StyleDefault.Foreground(palette.Success).Bold(true),
	}
}

// GetVaultStyles returns specialized styles for vault-related elements
func (sm *StyleManager) GetVaultStyles() map[string]tcell.Style {
	palette := sm.theme.GetPalette()
	return map[string]tcell.Style{
		"vault-name":       tcell.StyleDefault.Foreground(palette.Primary).Bold(true),
		"vault-path":       tcell.StyleDefault.Foreground(palette.TextMuted),
		"vault-locked":     tcell.StyleDefault.Foreground(palette.Warning).Bold(true),
		"vault-unlocked":   tcell.StyleDefault.Foreground(palette.Success).Bold(true),
		"item-name":        tcell.StyleDefault.Foreground(palette.Accent),
		"item-value":       tcell.StyleDefault.Foreground(palette.TextPrimary),
		"item-secret":      tcell.StyleDefault.Foreground(palette.TextMuted),
		"accessor-self":    tcell.StyleDefault.Foreground(palette.Success).Bold(true),
		"accessor-root":    tcell.StyleDefault.Foreground(palette.Primary).Bold(true),
		"accessor-user":    tcell.StyleDefault.Foreground(palette.Accent),
		"accessor-service": tcell.StyleDefault.Foreground(palette.Info),
	}
}

// GetElegantPalette returns a curated set of elegant colors
func (sm *StyleManager) GetElegantPalette() map[string]tcell.Color {
	palette := sm.theme.GetPalette()
	return map[string]tcell.Color{
		"midnight": palette.Background,     // Pure black
		"charcoal": palette.BackgroundDark, // Very dark gray
		"slate":    palette.BorderDark,     // Dark gray
		"steel":    palette.Border,         // Medium gray
		"silver":   palette.TextMuted,      // Light gray
		"pearl":    palette.TextSecondary,  // Off-white
		"snow":     palette.TextPrimary,    // Pure white
		"amethyst": palette.Primary,        // Soft purple
		"azure":    palette.Secondary,      // Light blue
		"cyan":     palette.Accent,         // Bright cyan
		"emerald":  palette.Success,        // Bright green
		"gold":     palette.Warning,        // Golden yellow
		"ruby":     palette.Error,          // Bright red
		"ocean":    palette.Selection,      // Dark blue
	}
}

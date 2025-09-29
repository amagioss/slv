package theme

import (
	"github.com/gdamore/tcell/v2"
)

// ColorPalette defines the color scheme for the theme
type ColorPalette struct {
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

// GetDefaultPalette returns the default color palette
func GetDefaultPalette() ColorPalette {
	return ColorPalette{
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

// GetDarkPalette returns a darker color palette
func GetDarkPalette() ColorPalette {
	return ColorPalette{
		Primary:   tcell.Color90, // Darker purple
		Secondary: tcell.Color67, // Darker blue
		Accent:    tcell.Color81, // Darker cyan

		Background:     tcell.ColorBlack, // Pure black
		BackgroundDark: tcell.Color232,   // Very dark gray

		TextPrimary:   tcell.Color250, // Off-white
		TextSecondary: tcell.Color244, // Medium gray
		TextMuted:     tcell.Color238, // Dark gray

		Success: tcell.Color70,  // Darker green
		Warning: tcell.Color214, // Darker yellow
		Error:   tcell.Color160, // Darker red
		Info:    tcell.Color67,  // Darker blue

		Border:     tcell.Color236, // Dark gray
		BorderDark: tcell.Color232, // Very dark gray

		Focus:     tcell.Color81, // Darker cyan
		Selection: tcell.Color17, // Darker blue
	}
}

// GetLightPalette returns a light color palette
func GetLightPalette() ColorPalette {
	return ColorPalette{
		Primary:   tcell.Color54, // Light purple
		Secondary: tcell.Color39, // Light blue
		Accent:    tcell.Color51, // Light cyan

		Background:     tcell.Color255, // Pure white
		BackgroundDark: tcell.Color250, // Light gray

		TextPrimary:   tcell.Color16,  // Pure black
		TextSecondary: tcell.Color238, // Dark gray
		TextMuted:     tcell.Color244, // Medium gray

		Success: tcell.Color34,  // Light green
		Warning: tcell.Color214, // Light yellow
		Error:   tcell.Color196, // Light red
		Info:    tcell.Color39,  // Light blue

		Border:     tcell.Color244, // Medium gray
		BorderDark: tcell.Color238, // Dark gray

		Focus:     tcell.Color51,  // Light cyan
		Selection: tcell.Color153, // Light blue
	}
}

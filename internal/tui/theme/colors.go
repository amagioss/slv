package theme

import (
	"github.com/gdamore/tcell/v2"
)

// ColorPalette defines the color scheme for the theme
type ColorPalette struct {
	//infobar colors
	InfobarBorder   tcell.Color
	InfobarTitle    tcell.Color
	InfoTableLabel  tcell.Color
	InfoTableValue  tcell.Color
	InfobarASCIIArt tcell.Color

	// main content colors
	MainContentTitle  tcell.Color
	MainContentBorder tcell.Color

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

	// Table and list specific colors
	TableHeader tcell.Color
	TableLabel  tcell.Color
	TableValue  tcell.Color
	TableType   tcell.Color
	TableName   tcell.Color
	TableEmail  tcell.Color
	TableKey    tcell.Color
	TableEmpty  tcell.Color
	TableMasked tcell.Color
	TableError  tcell.Color

	// Component specific colors
	ListSelectedText   tcell.Color
	ListSelectedBg     tcell.Color
	ListSecondaryText  tcell.Color
	ListMainText       tcell.Color
	FormLabel          tcell.Color
	FormFieldText      tcell.Color
	FormBorder         tcell.Color
	FormDisabledText   tcell.Color
	FormDisabledBg     tcell.Color
	FormDisabledBorder tcell.Color
}

// GetDefaultPalette returns the default color palette
func GetDefaultPalette() ColorPalette {
	return ColorPalette{

		//infobar colors - professional grays
		InfobarBorder:   tcell.Color240, // Medium gray
		InfobarTitle:    tcell.Color250, // Off-white
		InfoTableLabel:  tcell.Color244, // Medium gray
		InfoTableValue:  tcell.Color250, // Off-white
		InfobarASCIIArt: tcell.Color244, // Medium gray

		// Main content colors - clean and minimal
		MainContentTitle:  tcell.Color250, // Off-white
		MainContentBorder: tcell.Color240, // Medium gray

		// Primary colors - neutral professional palette
		Primary:   tcell.Color244, // Medium gray
		Secondary: tcell.Color240, // Medium-dark gray
		Accent:    tcell.Color250, // Off-white

		// Background colors - sophisticated dark theme
		Background:     tcell.Color16,  // Pure black
		BackgroundDark: tcell.Color232, // Very dark gray

		// Text colors - high contrast for readability
		TextPrimary:   tcell.Color250, // Off-white
		TextSecondary: tcell.Color244, // Medium gray
		TextMuted:     tcell.Color238, // Dark gray

		// Status colors - professional and clear
		Success: tcell.Color70,  // Soft green
		Warning: tcell.Color214, // Warm yellow
		Error:   tcell.Color160, // Soft red
		Info:    tcell.Color75,  // Light blue

		// Border and UI colors - subtle definition
		Border:     tcell.Color240, // Dark gray
		BorderDark: tcell.Color232, // Very dark gray

		// Focus and selection - clear indication
		Focus:     tcell.Color250, // Off-white
		Selection: tcell.Color240, // Medium gray

		// Table and list specific colors - professional hierarchy
		TableHeader: tcell.Color244, // Medium gray
		TableLabel:  tcell.Color244, // Medium gray
		TableValue:  tcell.Color250, // Off-white
		TableType:   tcell.Color194, // Very light green (closest to white)
		TableName:   tcell.Color195, // Very light blue (closest to white)
		TableEmail:  tcell.Color250, // Off-white
		TableKey:    tcell.Color244, // Medium gray
		TableEmpty:  tcell.Color238, // Dark gray
		TableMasked: tcell.Color230, // Very light yellow (closest to white)
		TableError:  tcell.Color160, // Soft red

		// Component specific colors - consistent styling
		ListSelectedText:   tcell.Color250, // Off-white
		ListSelectedBg:     tcell.Color240, // Medium gray
		ListSecondaryText:  tcell.Color244, // Medium gray
		ListMainText:       tcell.Color250, // Off-white
		FormLabel:          tcell.Color244, // Medium gray
		FormFieldText:      tcell.Color244, // Medium gray
		FormBorder:         tcell.Color244, // Medium gray
		FormDisabledText:   tcell.Color238, // Dark gray
		FormDisabledBg:     tcell.Color16,  // Black
		FormDisabledBorder: tcell.Color238, // Dark gray
	}
}

// GetCurrentPalette returns the currently active color palette
func GetCurrentPalette() ColorPalette {
	// For now, return the default palette
	// This can be extended to support theme switching
	return GetDefaultPalette()
}

package core

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Theme defines the interface for theme functionality
type Theme interface {
	ApplyTheme(app *tview.Application)
	GetBackground() tcell.Color
	GetAccent() tcell.Color
	GetPrimary() tcell.Color
	GetSecondary() tcell.Color
	GetTextPrimary() tcell.Color
	GetTextSecondary() tcell.Color
	GetTextMuted() tcell.Color
	GetSuccess() tcell.Color
	GetWarning() tcell.Color
	GetError() tcell.Color
	GetInfo() tcell.Color
	GetBorder() tcell.Color
	GetFocus() tcell.Color
	GetSelection() tcell.Color
}

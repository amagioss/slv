package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

// AttachPasteHandler attaches a Ctrl+V paste handler to a tview component (InputField or TextArea)
// bypassing tcell's event queue to avoid truncation issues
func AttachPasteHandler(component interface{}) {
	if component == nil {
		return
	}

	switch c := component.(type) {
	case *tview.InputField:
		existingCapture := c.GetInputCapture()
		c.SetInputCapture(createPasteHandler(
			func() string { return c.GetText() },
			func(text string) { c.SetText(text) },
			existingCapture,
		))
	case *tview.TextArea:
		existingCapture := c.GetInputCapture()
		c.SetInputCapture(createPasteHandler(
			func() string { return c.GetText() },
			func(text string) { c.SetText(text, false) },
			existingCapture,
		))
	}
}

// createPasteHandler creates a reusable input capture function for paste handling
func createPasteHandler(getText func() string, setText func(string), existingCapture func(*tcell.EventKey) *tcell.EventKey) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		// Handle Ctrl+V (or Cmd+V on Mac) for manual paste
		if event.Key() == tcell.KeyCtrlV || (event.Key() == tcell.KeyRune && event.Rune() == 'v' && event.Modifiers()&tcell.ModCtrl != 0) {
			// Read from clipboard
			clipboardBytes := clipboard.Read(clipboard.FmtText)
			if len(clipboardBytes) > 0 {
				clipboardText := string(clipboardBytes)
				// Get current text
				currentText := getText()
				// Append clipboard text
				setText(currentText + clipboardText)
			}
			return nil
		}

		// Chain to existing handler if present
		if existingCapture != nil {
			return existingCapture(event)
		}
		return event
	}
}

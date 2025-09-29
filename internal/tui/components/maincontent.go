package components

import (
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
)

// MainContent represents the main content area component
type MainContent struct {
	tui       interfaces.TUIInterface
	primitive *tview.Pages
	content   tview.Primitive
}

// NewMainContent creates a new MainContent component
func NewMainContent(tui interfaces.TUIInterface) *MainContent {
	mc := &MainContent{
		tui:       tui,
		primitive: tview.NewPages(),
	}
	return mc
}

// Render returns the primitive for this component
func (mc *MainContent) Render() tview.Primitive {
	return mc.primitive
}

// Refresh refreshes the component with current data
func (mc *MainContent) Refresh() {
	// Main content refreshes are handled by the router
	// This method is here for interface compliance
}

// SetFocus sets focus on the component
func (mc *MainContent) SetFocus(focus bool) {
	// Pages component doesn't have SetFocus method
	// Focus is handled by the application
}

// SetContent sets the content for the main area
func (mc *MainContent) SetContent(content tview.Primitive) {
	mc.content = content
	mc.primitive.AddPage("main", content, true, true)
}

// GetContent returns the current content
func (mc *MainContent) GetContent() tview.Primitive {
	return mc.content
}

// GetPages returns the pages container for direct manipulation
func (mc *MainContent) GetPages() *tview.Pages {
	return mc.primitive
}

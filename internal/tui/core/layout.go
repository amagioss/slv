package core

import (
	"github.com/rivo/tview"
)

// LayoutManager handles the layout of UI components
type LayoutManager struct {
	infoBar   tview.Primitive
	statusBar tview.Primitive
	content   tview.Primitive
	root      tview.Primitive
}

// NewLayoutManager creates a new LayoutManager instance
func NewLayoutManager() *LayoutManager {
	return &LayoutManager{}
}

// SetInfoBar sets the info bar primitive
func (lm *LayoutManager) SetInfoBar(infoBar tview.Primitive) {
	lm.infoBar = infoBar
}

// SetStatusBar sets the status bar primitive
func (lm *LayoutManager) SetStatusBar(statusBar tview.Primitive) {
	lm.statusBar = statusBar
}

// SetContent sets the content primitive
func (lm *LayoutManager) SetContent(content tview.Primitive) {
	lm.content = content
}

// GetInfoBar returns the info bar primitive
func (lm *LayoutManager) GetInfoBar() tview.Primitive {
	return lm.infoBar
}

// GetStatusBar returns the status bar primitive
func (lm *LayoutManager) GetStatusBar() tview.Primitive {
	return lm.statusBar
}

// GetContent returns the content primitive
func (lm *LayoutManager) GetContent() tview.Primitive {
	return lm.content
}

// GetRoot returns the root layout primitive
func (lm *LayoutManager) GetRoot() tview.Primitive {
	return lm.root
}

// BuildLayout builds the complete layout
func (lm *LayoutManager) BuildLayout() tview.Primitive {
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(lm.infoBar, 8, 1, false).  // Info bar
		AddItem(lm.content, 0, 1, true).   // Content (pages)
		AddItem(lm.statusBar, 3, 1, false) // Status bar

	lm.root = flex
	return lm.root
}

// UpdateLayout updates the layout with new components
func (lm *LayoutManager) UpdateLayout(infoBar, content, statusBar tview.Primitive) {
	lm.infoBar = infoBar
	lm.content = content
	lm.statusBar = statusBar
	lm.BuildLayout()
}

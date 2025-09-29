package components

import (
	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
)

// ComponentManager manages all UI components
type ComponentManager struct {
	tui         interfaces.TUIInterface
	infoBar     *InfoBar
	statusBar   *StatusBar
	mainContent *MainContent
}

// NewComponentManager creates a new component manager
func NewComponentManager(tui interfaces.TUIInterface) *ComponentManager {
	return &ComponentManager{
		tui: tui,
	}
}

// InitializeComponents initializes all components
func (cm *ComponentManager) InitializeComponents() {
	cm.infoBar = NewInfoBar(cm.tui)
	cm.statusBar = NewStatusBar(cm.tui)
	cm.mainContent = NewMainContent(cm.tui)
}

// GetInfoBar returns the info bar component
func (cm *ComponentManager) GetInfoBar() tview.Primitive {
	return cm.infoBar.Render()
}

// GetStatusBar returns the status bar component
func (cm *ComponentManager) GetStatusBar() tview.Primitive {
	return cm.statusBar.Render()
}

// GetMainContent returns the main content component
func (cm *ComponentManager) GetMainContent() tview.Primitive {
	return cm.mainContent.Render()
}

// GetMainContentPages returns the pages container for direct manipulation
func (cm *ComponentManager) GetMainContentPages() *tview.Pages {
	return cm.mainContent.GetPages()
}

// UpdateStatusBar updates the status bar with new text
func (cm *ComponentManager) UpdateStatusBar(text string) {
	cm.statusBar.SetCustomHelp(text)
}

// ClearStatusBar clears the status bar custom help
func (cm *ComponentManager) ClearStatusBar() {
	cm.statusBar.ClearCustomHelp()
}

// RefreshInfoBar refreshes the info bar with current data
func (cm *ComponentManager) RefreshInfoBar() {
	cm.infoBar.Refresh()
}

// UpdateStatus updates the status bar with page information
func (cm *ComponentManager) UpdateStatus(pageName string) {
	cm.statusBar.UpdateStatus(pageName)
}

// GetInfoBarComponent returns the InfoBar component for direct access
func (cm *ComponentManager) GetInfoBarComponent() *InfoBar {
	return cm.infoBar
}

// GetStatusBarComponent returns the StatusBar component for direct access
func (cm *ComponentManager) GetStatusBarComponent() *StatusBar {
	return cm.statusBar
}

// GetMainContentComponent returns the MainContent component for direct access
func (cm *ComponentManager) GetMainContentComponent() *MainContent {
	return cm.mainContent
}

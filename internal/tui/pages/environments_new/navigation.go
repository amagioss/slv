package environments_new

import (
	"github.com/gdamore/tcell/v2"
)

type FormNavigation struct {
	nep       *NewEnvironmentPage
	helpTexts map[int]string // Step-specific help texts
}

func (fn *FormNavigation) NewFormNavigation(nep *NewEnvironmentPage) *FormNavigation {
	return &FormNavigation{
		nep:       nep,
		helpTexts: make(map[int]string),
	}
}

func (fn *FormNavigation) SetupNavigation() {
	// Set up help texts for each step
	fn.setupHelpTexts()

	// Set initial help text
	fn.updateHelpText()
}

// setupHelpTexts sets up help text for each step
func (fn *FormNavigation) setupHelpTexts() {
	fn.helpTexts[StepProviderSelection] = "Select Provider: ↑/↓: Navigate list | Enter: Select | Tab: Switch to Cancel | Esc: Cancel"
	fn.helpTexts[StepMetadata] = "Environment Metadata: Tab: Next field | Enter: Submit | Backspace: Back"
	fn.helpTexts[StepProviderConfig] = "Provider Configuration: Tab: Next field | Enter: Submit | Backspace: Back"
	fn.helpTexts[StepConfirmation] = "Review Environment: Tab: Navigate | Enter: Create | Backspace: Back"
	fn.helpTexts[StepResult] = "Environment Created: ↑/↓: Navigate table | c: Copy field | Tab: To buttons | Shift+Tab: To table"
}

// updateHelpText updates the status bar with help text for the current step
func (fn *FormNavigation) updateHelpText() {
	if helpText, exists := fn.helpTexts[fn.nep.currentStep]; exists {
		fn.nep.GetTUI().UpdateStatusBar(helpText)
	}
}

// HandleInput handles input for the new environment page
func (fn *FormNavigation) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEsc:
		// ESC to go back or cancel
		switch fn.nep.currentStep {
		case StepProviderSelection:
			fn.nep.GetTUI().GetNavigation().GoBack()
			return nil
		case StepMetadata:
			fn.nep.showProviderSelection()
			return nil
		case StepProviderConfig:
			fn.nep.showMetadataForm()
			return nil
		case StepConfirmation:
			if fn.nep.selectedProvider == "direct" {
				fn.nep.showMetadataForm()
			} else {
				fn.nep.showProviderConfig()
			}
			return nil
		}
	}

	// Let the form handle other inputs
	return event
}

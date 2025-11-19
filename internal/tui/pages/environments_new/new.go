package environments_new

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/theme"
)

const (
	StepProviderSelection = iota
	StepMetadata
	StepProviderConfig
	StepConfirmation
	StepResult
)

// NewEnvironmentPage handles the new environment page functionality
type NewEnvironmentPage struct {
	pages.BasePage

	// Navigation
	navigation *FormNavigation

	// UI Components
	mainFlex      *tview.Flex
	contentFlex   *tview.Flex
	stepIndicator *tview.TextView

	// Step 1: Provider Selection
	providerList       *tview.List
	providerCancelForm *tview.Form

	// Step 2: Metadata Form
	metadataForm *tview.Form

	// Step 3: Provider Config
	providerForm *tview.Form

	// Step 4: Confirmation
	confirmText *tview.TextView
	confirmForm *tview.Form

	// Step 5: Result
	resultTable *tview.Table
	resultForm  *tview.Form

	// State
	currentStep      int
	selectedType     environments.EnvType
	envName          string
	envEmail         string
	envTags          []string
	quantumSafe      bool
	selectedProvider string
	providerInputs   map[string]string
	addToProfile     bool
	createdEnv       *environments.Environment
	secretKey        string
}

// NewNewEnvironmentPage creates a new NewEnvironmentPage instance
func NewNewEnvironmentPage(tui interfaces.TUIInterface) *NewEnvironmentPage {
	return &NewEnvironmentPage{
		BasePage:       *pages.NewBasePage(tui, "New Environment"),
		currentStep:    StepProviderSelection,
		providerInputs: make(map[string]string),
	}
}

// Create implements the Page interface
func (nep *NewEnvironmentPage) Create() tview.Primitive {
	colors := theme.GetCurrentPalette()

	// Create step indicator
	nep.stepIndicator = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(colors.TextPrimary)
	nep.updateStepIndicator()

	// Create main content flex
	nep.contentFlex = tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Create main layout
	nep.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nep.stepIndicator, 3, 0, false).
		AddItem(nep.contentFlex, 0, 1, true)

	// Start with provider selection
	nep.showProviderSelection()

	// Set up navigation
	nep.navigation = (&FormNavigation{}).NewFormNavigation(nep)
	nep.navigation.SetupNavigation()

	// Create layout using BasePage method
	return nep.CreateLayout(nep.mainFlex)
}

// Refresh implements the Page interface
func (nep *NewEnvironmentPage) Refresh() {
	// Reset to initial state
	nep.currentStep = StepProviderSelection
	nep.selectedType = ""
	nep.envName = ""
	nep.envEmail = ""
	nep.envTags = nil
	nep.quantumSafe = false
	nep.selectedProvider = ""
	nep.providerInputs = make(map[string]string)
	nep.addToProfile = false
	nep.createdEnv = nil
	nep.secretKey = ""

	nep.showProviderSelection()
}

// HandleInput implements the Page interface
func (nep *NewEnvironmentPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// Handled by navigation
	return event
}

// GetTitle implements the Page interface
func (nep *NewEnvironmentPage) GetTitle() string {
	return nep.BasePage.GetTitle()
}

// updateStepIndicator updates the step indicator text
func (nep *NewEnvironmentPage) updateStepIndicator() {
	colors := theme.GetCurrentPalette()
	steps := []string{"Provider", "Metadata", "Config", "Review", "Result"}
	var parts []string

	for i, step := range steps {
		if i == nep.currentStep {
			parts = append(parts, fmt.Sprintf("[%s]● %s[-]", colors.Accent.String(), step))
		} else if i < nep.currentStep {
			parts = append(parts, fmt.Sprintf("[%s]✓ %s[-]", colors.Success.String(), step))
		} else {
			parts = append(parts, fmt.Sprintf("[%s]○ %s[-]", colors.TextSecondary.String(), step))
		}
	}

	nep.stepIndicator.SetText("\n" + strings.Join(parts, "  →  "))
}

// GetCurrentComponent returns the currently focused component for navigation
func (nep *NewEnvironmentPage) GetCurrentComponent() tview.Primitive {
	switch nep.currentStep {
	case StepProviderSelection:
		return nep.providerList
	case StepMetadata:
		return nep.metadataForm
	case StepProviderConfig:
		return nep.providerForm
	case StepConfirmation:
		return nep.confirmForm
	case StepResult:
		return nep.resultTable
	default:
		return nil
	}
}

package vault_new

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
	"slv.sh/slv/internal/tui/pages/vault_view"
)

// VaultNewPage handles the new vault creation functionality
type VaultNewPage struct {
	pages.BasePage
	currentDir   string
	publicKeys   []string
	grantedEnvs  []*environments.Environment
	searchEnvMap map[string]*environments.Environment // Map environment names to environment structs for search results
	currentQuery string                               // Store the current search query for refreshing

	// Form references
	vaultConfigForm   *tview.Form   // Vault Configuration form
	optionsForm       *tview.Form   // Options form
	grantAccessForm   *tview.Form   // Grant Access form
	shareWithSelfForm *tview.Form   // Share with Self form
	shareWithK8sForm  *tview.Form   // Share with K8s Context form
	submitButton      *tview.Button // Submit button

	// Lists
	searchResults *tview.List
	grantedAccess *tview.List

	// Checkbox references
	shareWithSelfCheckbox *tview.Checkbox // Reference to the Share with Self checkbox
	shareWithK8sCheckbox  *tview.Checkbox // Reference to the Share with K8s Context checkbox

	currentPage tview.Primitive // Store reference to current page for modal navigation
}

// NewVaultNewPage creates a new VaultNewPage instance
func NewVaultNewPage(tui interfaces.TUIInterface, currentDir string) *VaultNewPage {
	return &VaultNewPage{
		BasePage:     *pages.NewBasePage(tui, "New Vault"),
		currentDir:   currentDir,
		publicKeys:   []string{},
		grantedEnvs:  []*environments.Environment{},
		searchEnvMap: make(map[string]*environments.Environment),
	}
}

// Create implements the Page interface
func (vnp *VaultNewPage) Create() tview.Primitive {
	// Create a single comprehensive form
	form := vnp.createComprehensiveVaultForm()

	// Update status bar
	vnp.GetTUI().UpdateStatusBar("[yellow]Tab: Navigate fields | Enter: Submit | Esc: Cancel[white]")

	// Create the page layout and show it
	vnp.SetTitle("New Vault at " + vnp.currentDir)
	return vnp.CreateLayout(form)
}

// Refresh implements the Page interface
func (vnp *VaultNewPage) Refresh() {
	// TODO: Implement new vault page refresh
}

// HandleInput implements the Page interface
func (vnp *VaultNewPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// TODO: Implement new vault page input handling
	return event
}

// GetTitle implements the Page interface
func (vnp *VaultNewPage) GetTitle() string {
	return vnp.BasePage.GetTitle()
}

// SetCurrentDir sets the current directory
func (vnp *VaultNewPage) SetCurrentDir(dir string) {
	// TODO: Implement set current directory
}

// createComprehensiveVaultForm creates a single, well-organized vault creation form
func (vnp *VaultNewPage) createComprehensiveVaultForm() tview.Primitive {
	// Create a two-column layout
	mainGrid := tview.NewGrid().
		SetRows(0).           // Single row
		SetColumns(-70, -30). // Two columns: Config, Options
		SetBorders(false)

	// Left Column: Main form with input fields
	leftForm := tview.NewForm()

	// Vault Metadata Section
	leftForm.AddInputField("Vault Name", "", 30, nil, func(text string) {
		// Auto-update file name based on vault name
		if text != "" {
			fileName := text + ".slv.yaml"
			leftForm.GetFormItem(1).(*tview.InputField).SetText(fileName)
		} else {
			leftForm.GetFormItem(1).(*tview.InputField).SetText("")
		}
	}).
		AddInputField("File Name", "", 40, nil, nil).
		AddInputField("K8s Namespace (optional)", "", 30, nil, nil)

	leftForm.SetBorder(true).
		SetTitle("Vault Configuration").
		SetTitleAlign(tview.AlignLeft)

	// Right Column: Vault Options
	optionsForm := tview.NewForm()

	// Vault Options
	optionsForm.AddCheckbox("Enable Hashing", false, nil).
		AddCheckbox("Quantum Safe", false, nil)

	optionsForm.SetBorder(true).
		SetTitle("Options").
		SetTitleAlign(tview.AlignLeft)

	// Set checkbox character to tick mark
	for i := 0; i < optionsForm.GetFormItemCount(); i++ {
		if checkbox, ok := optionsForm.GetFormItem(i).(*tview.Checkbox); ok {
			checkbox.SetCheckedString("✓")
		}
	}

	// Add forms to grid
	mainGrid.AddItem(leftForm, 0, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(optionsForm, 0, 1, 1, 1, 0, 0, false)

	// Create a flex layout to include search results and granted access
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Add the two-column form
	mainFlex.AddItem(mainGrid, 9, 1, true)

	// Add Grant Access section with input fields and results
	grantAccessFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Grant Access input fields and access checkboxes in same row
	grantAccessRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Left side: Public Key input
	grantAccessForm := tview.NewForm()
	grantAccessForm.AddInputField("Public Key / Search String", "", 80, nil, func(text string) {
		// Handle search string based on prefix
		if text != "" {
			vnp.searchEnvironments(text)
		} else {
			// Clear results when input is empty
			vnp.searchResults.Clear()
			vnp.searchResults.AddItem("", "", 0, nil)
		}
	})
	grantAccessForm.SetBorder(true).SetTitle("Grant Access").SetTitleAlign(tview.AlignLeft)

	// Store reference to the form for Enter key handling
	vnp.grantAccessForm = grantAccessForm

	// Right side: Access checkboxes (in same row)
	accessCheckboxesFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Create individual forms for each checkbox
	shareWithSelfForm := tview.NewForm()
	shareWithSelfForm.AddCheckbox("Share with Self", false, func(checked bool) {
		vnp.handleShareWithSelfChange(checked)
	})
	shareWithSelfForm.SetBorder(false)

	// Set checkbox character to tick mark and store reference
	if checkbox, ok := shareWithSelfForm.GetFormItem(0).(*tview.Checkbox); ok {
		checkbox.SetCheckedString("✓")
		vnp.shareWithSelfCheckbox = checkbox
	}

	shareWithK8sForm := tview.NewForm()
	shareWithK8sForm.AddCheckbox("Share with K8s Context", false, func(checked bool) {
		vnp.handleShareWithK8sChange(checked)
	})
	shareWithK8sForm.SetBorder(false)

	// Set checkbox character to tick mark
	if checkbox, ok := shareWithK8sForm.GetFormItem(0).(*tview.Checkbox); ok {
		checkbox.SetCheckedString("✓")
	}

	// Add both forms to the flex
	accessCheckboxesFlex.AddItem(shareWithSelfForm, 0, 1, false)
	accessCheckboxesFlex.AddItem(shareWithK8sForm, 0, 1, false)
	accessCheckboxesFlex.SetBorder(true).
		SetTitle("Quick Access").
		SetTitleAlign(tview.AlignLeft)

	// Add to the row
	grantAccessRow.AddItem(grantAccessForm, 0, 70, true)
	grantAccessRow.AddItem(accessCheckboxesFlex, 0, 30, false)

	// Results section
	resultsFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Search results list
	searchResults := tview.NewList()
	searchResults.SetBorder(true)
	searchResults.AddItem("", "", 0, nil)
	searchResults.SetTitle("Environment Results From Profile").SetTitleAlign(tview.AlignLeft)

	// Granted Access
	grantedAccess := tview.NewList()
	grantedAccess.SetBorder(true).SetTitle("Environments With Access").SetTitleAlign(tview.AlignLeft)
	grantedAccess.AddItem("📝 No access granted yet", "Add public keys or environments to grant access", 0, nil)

	resultsFlex.AddItem(searchResults, 0, 1, true)
	resultsFlex.AddItem(grantedAccess, 0, 1, false)

	// Add input row and results to Grant Access section
	grantAccessFlex.AddItem(grantAccessRow, 5, 1, false)
	grantAccessFlex.AddItem(resultsFlex, 0, 1, false)

	// Add Grant Access section
	mainFlex.AddItem(grantAccessFlex, 0, 1, false)

	// Add Submit Button section
	submitButtonFlex, submitButton := vnp.createSubmitButtonSection(leftForm, optionsForm)
	mainFlex.AddItem(submitButtonFlex, 3, 1, false) // Increased from 3 to 5 for more space

	// Store references for real-time updates
	vnp.searchResults = searchResults
	vnp.grantedAccess = grantedAccess
	vnp.publicKeys = []string{}
	vnp.grantedEnvs = []*environments.Environment{}
	vnp.searchEnvMap = make(map[string]*environments.Environment)

	// Set up navigation between forms while preserving within-form navigation
	vnp.setupComprehensiveFormNavigation(leftForm, optionsForm, grantAccessForm, shareWithSelfForm, shareWithK8sForm, searchResults, grantedAccess, submitButton)

	// Store reference to current page for modal navigation
	vnp.currentPage = mainFlex

	return mainFlex
}

// createSubmitButtonSection creates the submit button section with border
func (vnp *VaultNewPage) createSubmitButtonSection(leftForm, optionsForm *tview.Form) (tview.Primitive, *tview.Button) {
	// Create a flex container for the submit button section
	submitFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create the submit button
	submitButton := tview.NewButton("Create Vault").
		SetSelectedFunc(func() {
			vnp.createVaultFromForm(leftForm, optionsForm)
		})

	// Style the submit button
	submitButton.SetBorder(true).
		// SetTitle("Actions").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(vnp.GetTheme().GetAccent()).
		SetBackgroundColor(vnp.GetTheme().GetBackground())

	// Set button text and background colors for better visibility
	submitButton.SetLabelColor(vnp.GetTheme().GetTextPrimary()).
		SetLabelColorActivated(vnp.GetTheme().GetBackground()).
		SetBackgroundColorActivated(vnp.GetTheme().GetAccent())

	// Center the button horizontally
	centeredFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centeredFlex.AddItem(nil, 0, 1, false)          // Left spacer
	centeredFlex.AddItem(submitButton, 20, 1, true) // Button (20 chars wide)
	centeredFlex.AddItem(nil, 0, 1, false)          // Right spacer

	// Add some vertical spacing
	// submitFlex.AddItem(nil, 1, 1, false)         // Top spacer
	submitFlex.AddItem(centeredFlex, 3, 1, true) // Button row
	submitFlex.AddItem(nil, 1, 1, false)         // Bottom spacer (increased from 1 to 2)

	return submitFlex, submitButton
}

// setupComprehensiveFormNavigation sets up navigation between all forms
func (vnp *VaultNewPage) setupComprehensiveFormNavigation(leftForm, optionsForm, grantAccessForm, shareWithSelfForm, shareWithK8sForm *tview.Form, searchResults, grantedAccess *tview.List, submitButton *tview.Button) {
	// Create a focus group for inter-form navigation
	focusGroup := []tview.Primitive{
		leftForm,          // Vault Configuration
		optionsForm,       // Options
		grantAccessForm,   // Grant Access
		shareWithSelfForm, // Share with Self
		shareWithK8sForm,  // Share with K8s Context
		searchResults,     // Search Results list
		grantedAccess,     // Granted Access list
		submitButton,      // Submit Button
	}

	currentFocus := 0

	// Create a shared input capture function for inter-form navigation
	createInterFormInputCapture := func() func(*tcell.EventKey) *tcell.EventKey {
		return func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyTab:
				// Move to next form/primitive
				currentFocus = (currentFocus + 1) % len(focusGroup)
				vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
				return nil
			case tcell.KeyBacktab:
				// Move to previous form/primitive
				currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
				vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
				return nil
			case tcell.KeyCtrlS:
				// Create vault with Ctrl+S
				vnp.createVaultFromForm(leftForm, optionsForm)
				return nil
			}
			// Let all other keys pass through to the primitive for within-form navigation
			return event
		}
	}

	// Set up specific navigation for Vault Configuration form (leftForm)
	leftForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Move to next form/primitive
			currentFocus = (currentFocus + 1) % len(focusGroup)
			vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyBacktab:
			// Move to previous form/primitive
			currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
			vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyDown:
			// Move to next field
			i, _ := leftForm.GetFocusedItemIndex()
			next := (i + 1) % leftForm.GetFormItemCount()
			vnp.GetTUI().GetApplication().SetFocus(leftForm.GetFormItem(next))
			return nil
		case tcell.KeyUp:
			// Move to previous field
			i, _ := leftForm.GetFocusedItemIndex()
			prev := (i - 1 + leftForm.GetFormItemCount()) % leftForm.GetFormItemCount()
			vnp.GetTUI().GetApplication().SetFocus(leftForm.GetFormItem(prev))
			return nil
		case tcell.KeyCtrlS:
			// Create vault with Ctrl+S
			vnp.createVaultFromForm(leftForm, optionsForm)
			return nil
		}
		// Let all other keys pass through
		return event
	})

	// Set up specific navigation for Options form
	optionsForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Move to next form/primitive
			currentFocus = (currentFocus + 1) % len(focusGroup)
			vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyBacktab:
			// Move to previous form/primitive
			currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
			vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyDown:
			// Move to next field
			i, _ := optionsForm.GetFocusedItemIndex()
			next := (i + 1) % optionsForm.GetFormItemCount()
			vnp.GetTUI().GetApplication().SetFocus(optionsForm.GetFormItem(next))
			return nil
		case tcell.KeyUp:
			// Move to previous field
			i, _ := optionsForm.GetFocusedItemIndex()
			prev := (i - 1 + optionsForm.GetFormItemCount()) % optionsForm.GetFormItemCount()
			vnp.GetTUI().GetApplication().SetFocus(optionsForm.GetFormItem(prev))
			return nil
		case tcell.KeyCtrlS:
			// Create vault with Ctrl+S
			vnp.createVaultFromForm(leftForm, optionsForm)
			return nil
		}
		// Let all other keys pass through
		return event
	})

	// Set up navigation for individual access control forms
	shareWithSelfForm.SetInputCapture(createInterFormInputCapture())
	shareWithK8sForm.SetInputCapture(createInterFormInputCapture())

	// Set up input capture for Grant Access form with special Enter handling
	grantAccessForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Move to next form/primitive
			currentFocus = (currentFocus + 1) % len(focusGroup)
			vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyBacktab:
			// Move to previous form/primitive
			currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
			vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyEnter:
			// Handle Enter key for SLV_EPK inputs
			vnp.handleGrantAccessEnter()
			return nil
		case tcell.KeyCtrlS:
			// Create vault with Ctrl+S
			vnp.createVaultFromForm(leftForm, optionsForm)
			return nil
		case tcell.KeyDown:
			// Move to next field
			vnp.GetTUI().GetApplication().SetFocus(focusGroup[5])
			return nil
		}
		// Let all other keys pass through to the primitive for within-form navigation
		return event
	})

	// Set up input capture for lists
	searchResults.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Move to next form/primitive
			currentFocus = (currentFocus + 1) % len(focusGroup)
			vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyBacktab:
			// Move to previous form/primitive
			currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
			vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyCtrlS:
			// Create vault with Ctrl+S
			vnp.createVaultFromForm(leftForm, optionsForm)
			return nil
		case tcell.KeyEnter:
			// Add selected item to granted access
			selected := searchResults.GetCurrentItem()
			if selected >= 0 && selected < searchResults.GetItemCount() {
				mainText, _ := searchResults.GetItemText(selected)
				// Extract environment name from the formatted text and find the environment
				if strings.HasPrefix(mainText, "🔍 ") {
					envName := strings.TrimPrefix(mainText, "🔍 ")
					// Find the environment in the search results map
					if env, exists := vnp.searchEnvMap[envName]; exists {
						vnp.addToGrantedAccess(env)
					}
				}
			}
			return nil
		}
		// Let all other keys pass through for list navigation
		return event
	})

	grantedAccess.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Move to next form/primitive
			currentFocus = (currentFocus + 1) % len(focusGroup)
			vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyBacktab:
			// Move to previous form/primitive
			currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
			vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyCtrlS:
			// Create vault with Ctrl+S
			vnp.createVaultFromForm(leftForm, optionsForm)
			return nil
		case tcell.KeyCtrlD:
			// Remove selected environment from granted access
			selected := grantedAccess.GetCurrentItem()
			if selected >= 0 && selected < grantedAccess.GetItemCount() {
				mainText, _ := grantedAccess.GetItemText(selected)
				// Extract environment name from the formatted text
				if strings.HasPrefix(mainText, "🌍 ") {
					envName := strings.TrimPrefix(mainText, "🌍 ")
					vnp.removeFromGrantedAccess(envName)
				}
			}
			return nil
		}
		// Let all other keys pass through for list navigation
		return event
	})

	submitButton.SetInputCapture(createInterFormInputCapture())

	// Set initial focus to the first form
	vnp.GetTUI().GetApplication().SetFocus(focusGroup[currentFocus])
}

// searchEnvironments searches for environments based on query
func (vnp *VaultNewPage) searchEnvironments(query string) {
	vnp.currentQuery = query // Store the current query for refreshing
	vnp.searchResults.Clear()
	vnp.searchEnvMap = make(map[string]*environments.Environment) // Clear previous search results
	profile, err := profiles.GetActiveProfile()
	if err != nil {
		vnp.showError(fmt.Sprintf("Error getting active profile: %v", err))
		vnp.searchResults.AddItem("", "", 0, nil)
		return
	}

	// Helper function to check if environment is already granted access
	isAlreadyGranted := func(env *environments.Environment) bool {
		for _, grantedEnv := range vnp.grantedEnvs {
			if grantedEnv.PublicKey == env.PublicKey {
				return true
			}
		}
		return false
	}

	var matchingEnvs []*environments.Environment

	if strings.HasPrefix(strings.ToUpper(query), "SLV_EPK") {
		envs, err := profile.ListEnvs()
		if err != nil {
			vnp.showError(fmt.Sprintf("Error listing environments: %v", err))
			vnp.searchResults.AddItem("", "", 0, nil)
			return
		}
		for _, env := range envs {
			if env.PublicKey == query {
				matchingEnvs = append(matchingEnvs, env)
			}
		}
		if len(matchingEnvs) == 0 {
			vnp.searchResults.AddItem("❌ Environment not found in the profile", "", 0, nil)
			return
		}
	} else {
		envs, err := profile.SearchEnvs([]string{query})
		if err != nil {
			vnp.showError(fmt.Sprintf("Error searching environments: %v", err))
			vnp.searchResults.AddItem("", "", 0, nil)
			return
		}
		matchingEnvs = envs
	}

	// Sort environments by name for consistent display order
	sort.Slice(matchingEnvs, func(i, j int) bool {
		return matchingEnvs[i].Name < matchingEnvs[j].Name
	})

	// Display sorted environments
	for _, env := range matchingEnvs {
		// Check if already granted access
		if isAlreadyGranted(env) {
			vnp.searchResults.AddItem("✅ Environment already has access", fmt.Sprintf("Name: %s | Type: %s", env.Name, string(env.EnvType)), 0, nil)
		} else {
			// Store environment in map for later retrieval
			vnp.searchEnvMap[env.Name] = env
			// Format with colors and proper spacing
			mainText := fmt.Sprintf("🔍 %s", env.Name)
			secondaryText := fmt.Sprintf("Type: %s | Email: %s", string(env.EnvType), env.Email)
			vnp.searchResults.AddItem(mainText, secondaryText, 0, func() {
				vnp.addToGrantedAccess(env)
			})
		}
	}
}

// addToGrantedAccess adds an environment to the granted access list
func (vnp *VaultNewPage) addToGrantedAccess(env *environments.Environment) {
	// Check if environment is already granted
	for _, existingEnv := range vnp.grantedEnvs {
		if existingEnv.PublicKey == env.PublicKey {
			return // Already granted
		}
	}

	// Add to granted environments
	vnp.grantedEnvs = append(vnp.grantedEnvs, env)

	// Update the display
	vnp.updateGrantedAccessList()

	// Refresh search results to show the newly added environment as "already granted"
	vnp.refreshSearchResults()
}

func (vnp *VaultNewPage) refreshSearchResults() {
	if vnp.currentQuery != "" {
		vnp.searchEnvironments(vnp.currentQuery)
	}
}

// removeFromGrantedAccess removes an environment from granted access
func (vnp *VaultNewPage) removeFromGrantedAccess(envName string) {
	// Find and remove the environment from grantedEnvs
	var removedEnv *environments.Environment
	for i, env := range vnp.grantedEnvs {
		if env.Name == envName {
			removedEnv = env
			// Remove the environment at index i
			vnp.grantedEnvs = append(vnp.grantedEnvs[:i], vnp.grantedEnvs[i+1:]...)
			break
		}
	}

	// Check if the removed environment is the self environment
	if removedEnv != nil {
		selfEnv := environments.GetSelf()
		if selfEnv != nil && selfEnv.PublicKey == removedEnv.PublicKey {
			// Uncheck the "Share with Self" checkbox
			vnp.uncheckShareWithSelf()
		}
	}

	// Update the display
	vnp.updateGrantedAccessList()

	// Refresh search results to show the removed environment as available again
	vnp.refreshSearchResults()
}

// uncheckShareWithSelf unchecks the Share with Self checkbox
func (vnp *VaultNewPage) uncheckShareWithSelf() {
	if vnp.shareWithSelfCheckbox != nil {
		vnp.shareWithSelfCheckbox.SetChecked(false)
	}
}

func (vnp *VaultNewPage) createVaultFromForm(leftForm, optionsForm *tview.Form) {
	// Collect form data
	vaultName := leftForm.GetFormItem(0).(*tview.InputField).GetText()
	fileName := leftForm.GetFormItem(1).(*tview.InputField).GetText()
	namespace := leftForm.GetFormItem(2).(*tview.InputField).GetText()

	// Get checkbox states
	enableHashing := optionsForm.GetFormItem(0).(*tview.Checkbox).IsChecked()
	quantumSafe := optionsForm.GetFormItem(1).(*tview.Checkbox).IsChecked()

	// Validate inputs
	if err := vnp.validateVaultInputsForCreation(vaultName, fileName); err != nil {
		vnp.showError(err.Error())
		return
	}

	// Prepare vault file path
	vaultFilePath := filepath.Join(vnp.currentDir, fileName)

	// Collect public keys for vault access from granted environments
	var publicKeys []*crypto.PublicKey
	for _, env := range vnp.grantedEnvs {
		if pk, err := env.GetPublicKey(); err == nil {
			publicKeys = append(publicKeys, pk)
		}
	}

	// Create the vault
	vault, err := vaults.New(vaultFilePath, vaultName, namespace, enableHashing, quantumSafe, publicKeys...)
	if err != nil {
		vnp.showError(fmt.Sprintf("Failed to create vault: %v", err))
		return
	}

	// Show success message
	vnp.showSuccess(fmt.Sprintf("Vault '%s' created successfully at %s", vaultName, vaultFilePath))

	// Navigate to vault details page and remove new vault page from stack
	vnp.navigateToVaultDetails(vault, vaultFilePath)
}

// navigateToVaultDetails navigates to the vault details page and removes new vault page from stack
func (vnp *VaultNewPage) navigateToVaultDetails(vault *vaults.Vault, vaultFilePath string) {
	// Get the registered vault view page
	vaultViewPage := vnp.GetTUI().GetRouter().GetRegisteredPage("vaults_view").(*vault_view.VaultViewPage)

	// Set the vault and filepath for the registered page
	vaultViewPage.SetVault(vault)
	vaultViewPage.SetFilePath(vaultFilePath)

	// Remove the new vault page from the stack first
	// Use the existing ShowVaultDetails method with replace=true
	vnp.GetTUI().GetNavigation().ShowVaultDetails(true)

	vnp.showSuccess(fmt.Sprintf("Vault '%s' created successfully at %s", vault.Name, vaultFilePath))
}

// validateVaultInputs validates the vault creation inputs
func (vnp *VaultNewPage) validateVaultInputsForCreation(vaultName, fileName string) error {
	// Validate vault name
	if strings.TrimSpace(vaultName) == "" {
		return fmt.Errorf("vault name is required")
	}

	// Validate file name
	if strings.TrimSpace(fileName) == "" {
		return fmt.Errorf("file name is required")
	}

	// Check if file already exists
	vaultFilePath := filepath.Join(vnp.currentDir, fileName)
	if _, err := os.Stat(vaultFilePath); err == nil {
		return fmt.Errorf("file already exists: %s", fileName)
	}

	if len(vnp.grantedEnvs) == 0 {
		return fmt.Errorf("no environments granted access. please grant access to at least one environment")
	}

	return nil
}

// showError displays an error message using the TUI's built-in error modal
func (vnp *VaultNewPage) showError(message string) {
	vnp.GetTUI().ShowError(message)
}

// showSuccess displays a success message using the TUI's built-in info modal
func (vnp *VaultNewPage) showSuccess(message string) {
	vnp.GetTUI().ShowInfo("✅ Success: " + message)
}

// handleShareWithSelfChange handles the "Share with Self" checkbox change
func (vnp *VaultNewPage) handleShareWithSelfChange(checked bool) {
	if checked {
		// Add self environment to granted access
		selfEnv := environments.GetSelf()
		if selfEnv != nil {
			// Check if self environment is already granted
			alreadyGranted := false
			for _, env := range vnp.grantedEnvs {
				if env.PublicKey == selfEnv.PublicKey {
					alreadyGranted = true
					break
				}
			}

			if !alreadyGranted {
				vnp.grantedEnvs = append(vnp.grantedEnvs, selfEnv)
				vnp.updateGrantedAccessList()
				vnp.refreshSearchResults()
			}
		}
	} else {
		// Remove self environment from granted access
		selfEnv := environments.GetSelf()
		if selfEnv != nil {
			for i, env := range vnp.grantedEnvs {
				if env.PublicKey == selfEnv.PublicKey {
					vnp.grantedEnvs = append(vnp.grantedEnvs[:i], vnp.grantedEnvs[i+1:]...)
					break
				}
			}
			vnp.updateGrantedAccessList()
			vnp.refreshSearchResults()
		}
	}
}

// handleShareWithK8sChange handles the "Share with K8s Context" checkbox change
func (vnp *VaultNewPage) handleShareWithK8sChange(checked bool) {
	if checked {
		// Try to get K8s environment details
		k8sEnv, err := vnp.getK8sEnvironment()
		if err != nil {
			// Show error and uncheck the checkbox
			vnp.showError(fmt.Sprintf("Failed to get K8s environment: %v", err))
			// Note: We can't easily uncheck the checkbox here since we don't have a reference
			// The user will need to manually uncheck it after seeing the error
			return
		}

		// Check if K8s environment is already granted
		alreadyGranted := false
		for _, env := range vnp.grantedEnvs {
			if env.PublicKey == k8sEnv.PublicKey {
				alreadyGranted = true
				break
			}
		}

		if !alreadyGranted {
			vnp.grantedEnvs = append(vnp.grantedEnvs, k8sEnv)
			vnp.updateGrantedAccessList()
			vnp.refreshSearchResults()
		}
	}
}

// getK8sEnvironment gets environment details from K8s context
func (vnp *VaultNewPage) getK8sEnvironment() (*environments.Environment, error) {
	// Get current K8s namespace
	namespace := session.GetK8sNamespace()
	if namespace == "" {
		namespace = config.AppNameLowerCase
	}

	// Try to get EC public key first (non-post-quantum)
	publicKeyStr, err := session.GetPublicKeyFromK8s(namespace, false)
	if err != nil {
		// If EC key fails, try post-quantum key
		publicKeyStr, err = session.GetPublicKeyFromK8s(namespace, true)
		if err != nil {
			return nil, fmt.Errorf("no public key found in K8s namespace '%s': %w", namespace, err)
		}
	}

	// Validate the public key
	_, err = crypto.PublicKeyFromString(publicKeyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid public key from K8s: %w", err)
	}

	// Search through profile to see if this public key matches any existing environment
	profile, err := profiles.GetActiveProfile()
	if err != nil {
		// If we can't get the profile, create a "No information" environment
		return vnp.createNoInformationK8sEnvironment(publicKeyStr, namespace), nil
	}

	// Get all environments from the profile
	allEnvs, err := profile.ListEnvs()
	if err != nil {
		// If we can't list environments, create a "No information" environment
		return vnp.createNoInformationK8sEnvironment(publicKeyStr, namespace), nil
	}

	// Search for matching environment by public key
	for _, env := range allEnvs {
		if env.PublicKey == publicKeyStr {
			// Found matching environment, return it
			return env, nil
		}
	}

	// No matching environment found, create a "No information" environment
	return vnp.createNoInformationK8sEnvironment(publicKeyStr, namespace), nil
}

// createNoInformationK8sEnvironment creates a "No information" environment for K8s context
func (vnp *VaultNewPage) createNoInformationK8sEnvironment(publicKeyStr, namespace string) *environments.Environment {
	return &environments.Environment{
		PublicKey: publicKeyStr,
		Name:      "No information",
		Email:     "Unknown",
		EnvType:   environments.SERVICE,
		Tags:      []string{"kubernetes", "context", namespace, "no-information"},
	}
}

// handleGrantAccessEnter handles Enter key press in Grant Access field
func (vnp *VaultNewPage) handleGrantAccessEnter() {
	// Get the current input text from the Grant Access field
	// We need to find the input field in the grant access form
	// Since we know it's the first (and only) form item, we can access it directly
	if vnp.grantAccessForm == nil {
		return
	}

	formItem := vnp.grantAccessForm.GetFormItem(0)
	if inputField, ok := formItem.(*tview.InputField); ok {
		text := inputField.GetText()

		// Check if the text starts with "SLV_EPK"
		if strings.HasPrefix(text, "SLV_EPK") {
			// Extract the public key (remove "SLV_EPK" prefix)
			publicKeyStr := strings.TrimSpace(text)

			if publicKeyStr == "" {
				vnp.GetTUI().ShowError("Invalid public key format. Please provide a valid public key.")
				return
			}

			// Validate the public key format
			_, err := crypto.PublicKeyFromString(publicKeyStr)
			if err != nil {
				vnp.GetTUI().ShowError(fmt.Sprintf("Invalid public key format: %v", err))
				return
			}

			// Search through profile to see if this public key matches any existing environment
			profile, err := profiles.GetActiveProfile()
			if err != nil {
				// If we can't get the profile, create a "No information" environment
				env := vnp.createNoInformationEnvironment(publicKeyStr)
				vnp.addEnvironmentToGrantedAccess(env)
				return
			}

			// Get all environments from the profile
			allEnvs, err := profile.ListEnvs()
			if err != nil {
				// If we can't list environments, create a "No information" environment
				env := vnp.createNoInformationEnvironment(publicKeyStr)
				vnp.addEnvironmentToGrantedAccess(env)
				return
			}

			// Search for matching environment by public key
			var foundEnv *environments.Environment
			for _, env := range allEnvs {
				if env.PublicKey == publicKeyStr {
					foundEnv = env
					break
				}
			}

			if foundEnv != nil {
				// Found matching environment, use it
				vnp.addEnvironmentToGrantedAccess(foundEnv)
			} else {
				// No matching environment found, create a "No information" environment
				env := vnp.createNoInformationEnvironment(publicKeyStr)
				vnp.addEnvironmentToGrantedAccess(env)
			}

			// Clear the input field
			inputField.SetText("")
		}
	}
}

// createNoInformationEnvironment creates a "No information" environment for a public key
func (vnp *VaultNewPage) createNoInformationEnvironment(publicKeyStr string) *environments.Environment {
	return &environments.Environment{
		PublicKey: publicKeyStr,
		Name:      "No information",
		Email:     "Unknown",
		EnvType:   environments.SERVICE,
		Tags:      []string{"no-information", "manual-entry"},
	}
}

// addEnvironmentToGrantedAccess adds an environment to the granted access list if not already present
func (vnp *VaultNewPage) addEnvironmentToGrantedAccess(env *environments.Environment) {
	// Check if environment is already granted
	alreadyGranted := false
	for _, grantedEnv := range vnp.grantedEnvs {
		if grantedEnv.PublicKey == env.PublicKey {
			alreadyGranted = true
			break
		}
	}

	if !alreadyGranted {
		vnp.grantedEnvs = append(vnp.grantedEnvs, env)
		vnp.updateGrantedAccessList()
		vnp.refreshSearchResults()
	}
}

// updateGrantedAccessList updates the granted access display
func (vnp *VaultNewPage) updateGrantedAccessList() {
	vnp.grantedAccess.Clear()

	// Add public keys
	for i, key := range vnp.publicKeys {
		keyDisplay := key
		if len(key) > 25 {
			keyDisplay = key[:25] + "..."
		}
		mainText := fmt.Sprintf("🔑 Public Key %d", i+1)
		secondaryText := fmt.Sprintf("Key: %s", keyDisplay)
		vnp.grantedAccess.AddItem(mainText, secondaryText, 0, nil)
	}

	// Add granted environments (sorted by name for consistent display)
	sortedGrantedEnvs := make([]*environments.Environment, len(vnp.grantedEnvs))
	copy(sortedGrantedEnvs, vnp.grantedEnvs)
	sort.Slice(sortedGrantedEnvs, func(i, j int) bool {
		return sortedGrantedEnvs[i].Name < sortedGrantedEnvs[j].Name
	})

	for _, env := range sortedGrantedEnvs {
		// Handle unknown name and email
		name := env.Name
		if name == "" {
			name = "Unknown"
		}

		email := env.Email
		if email == "" {
			email = "Unknown"
		}

		mainText := fmt.Sprintf("🌍 %s", name)
		secondaryText := fmt.Sprintf("Type: %s | Email: %s | Key: %s...", string(env.EnvType), email, env.PublicKey[:min(15, len(env.PublicKey))])
		vnp.grantedAccess.AddItem(mainText, secondaryText, 0, nil)
	}

	// Show message if no access granted
	if len(vnp.publicKeys) == 0 && len(vnp.grantedEnvs) == 0 {
		vnp.grantedAccess.AddItem("📝 No access granted yet", "Add public keys or environments to grant access", 0, nil)
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

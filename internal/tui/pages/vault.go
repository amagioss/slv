package pages

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
)

// VaultPage handles the vault management page functionality
type VaultPage struct {
	tui                   interfaces.TUIInterface
	currentDir            string
	vault                 *vaults.Vault // Store the current vault instance
	vaultPath             string        // Store the current vault path
	searchResults         *tview.List
	grantedAccess         *tview.List
	publicKeys            []string
	grantedEnvs           []*environments.Environment
	searchEnvMap          map[string]*environments.Environment // Map environment names to environment structs for search results
	currentQuery          string                               // Store the current search query for refreshing
	shareWithSelfCheckbox *tview.Checkbox                      // Reference to the Share with Self checkbox
	grantAccessForm       *tview.Form                          // Reference to the Grant Access form
	currentPage           tview.Primitive                      // Store reference to current page for modal navigation
}

// NewVaultPage creates a new VaultPage instance
func NewVaultPage(tui interfaces.TUIInterface, currentDir string) *VaultPage {
	return &VaultPage{
		tui:        tui,
		currentDir: currentDir,
		vault:      nil,
		vaultPath:  "",
	}
}

// CreateVaultPage creates the vault management page
func (vp *VaultPage) CreateVaultPage() tview.Primitive {
	// Create welcome message
	welcomeText := fmt.Sprintf("\n[white]Browse Vaults[white::-]\n[gray](Use arrow keys [←] and [→] to navigate directories)[gray::-]\n\nCurrent Directory: %s", vp.currentDir)

	pwd := tview.NewTextView().
		SetText(welcomeText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetWrap(true)

	// Get directories and .slv files
	items := vp.getVaultFiles()

	// Create the list
	list := tview.NewList()

	// Add "go back one directory" option at the top
	list.AddItem("⬆️ Go Back", "Go back to parent directory", 'b', func() {
		vp.goBackDirectory()
	})

	// Add items to the list
	for _, item := range items {
		icon := "📁"
		if item.IsFile {
			icon = "📄"
		}

		list.AddItem(
			fmt.Sprintf("%s %s", icon, item.Name),
			"",
			0,
			func() {
				vp.handleItemSelection(item)
			},
		)
	}

	// Style the list
	list.SetSelectedTextColor(tcell.ColorYellow).
		SetSelectedBackgroundColor(tcell.ColorNavy).
		SetSecondaryTextColor(tcell.ColorGray).
		SetMainTextColor(tcell.ColorWhite)

	// Set up keyboard navigation
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			// Load selected directory
			vp.loadSelectedItem(list)
			return nil
		case tcell.KeyLeft:
			// Go back to previous directory
			vp.goBackDirectory()
			return nil
		case tcell.KeyCtrlN:
			// Create new vault
			vp.tui.GetNavigation().ShowNewVault()
			return nil
		}
		return event
	})

	// Create a centered layout using grid
	content := tview.NewGrid().
		SetRows(6, 0). // Two flexible rows
		SetColumns(0). // Single column
		SetBorders(false)

	// Center the welcome text
	content.AddItem(pwd, 0, 0, 1, 1, 0, 0, false)

	// Center the list
	content.AddItem(list, 1, 0, 1, 1, 0, 0, true)

	// Update status bar with help text
	vp.tui.UpdateStatusBar("[yellow]←/→: Move between directories | ↑/↓: Navigate | Enter: open vault/directory | Ctrl+N: New vault[white]")
	return vp.tui.CreatePageLayout("Vault Management", content)
}

// CreateNewVaultPage creates the new vault creation page
func (vp *VaultPage) CreateNewVaultPage() tview.Primitive {
	// Create a single comprehensive form
	form := vp.createComprehensiveVaultForm()

	// Update status bar
	vp.tui.UpdateStatusBar("[yellow]Tab: Navigate fields | Enter: Submit | Esc: Cancel[white]")

	// Create the page layout and show it
	return vp.tui.CreatePageLayout("New Vault at "+vp.currentDir, form)
}

// createComprehensiveVaultForm creates a single, well-organized vault creation form
func (vp *VaultPage) createComprehensiveVaultForm() tview.Primitive {
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
			vp.searchEnvironments(text)
		} else {
			// Clear results when input is empty
			vp.searchResults.Clear()
			vp.searchResults.AddItem("", "", 0, nil)
		}
	})
	grantAccessForm.SetBorder(true).SetTitle("Grant Access").SetTitleAlign(tview.AlignLeft)

	// Store reference to the form for Enter key handling
	vp.grantAccessForm = grantAccessForm

	// Right side: Access checkboxes (in same row)
	accessCheckboxesFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Create individual forms for each checkbox
	shareWithSelfForm := tview.NewForm()
	shareWithSelfForm.AddCheckbox("Share with Self", false, func(checked bool) {
		vp.handleShareWithSelfChange(checked)
	})
	shareWithSelfForm.SetBorder(false)

	// Set checkbox character to tick mark and store reference
	if checkbox, ok := shareWithSelfForm.GetFormItem(0).(*tview.Checkbox); ok {
		checkbox.SetCheckedString("✓")
		vp.shareWithSelfCheckbox = checkbox
	}

	shareWithK8sForm := tview.NewForm()
	shareWithK8sForm.AddCheckbox("Share with K8s Context", false, func(checked bool) {
		vp.handleShareWithK8sChange(checked)
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

	// Store references for real-time updates
	vp.searchResults = searchResults
	vp.grantedAccess = grantedAccess
	vp.publicKeys = []string{}
	vp.grantedEnvs = []*environments.Environment{}
	vp.searchEnvMap = make(map[string]*environments.Environment)

	// Set up navigation between forms while preserving within-form navigation
	vp.setupComprehensiveFormNavigation(leftForm, optionsForm, grantAccessForm, shareWithSelfForm, shareWithK8sForm, searchResults, grantedAccess)

	// Store reference to current page for modal navigation
	vp.currentPage = mainFlex

	return mainFlex
}

// setupComprehensiveFormNavigation sets up navigation between forms while preserving within-form navigation
func (vp *VaultPage) setupComprehensiveFormNavigation(leftForm, optionsForm, grantAccessForm, shareWithSelfForm, shareWithK8sForm *tview.Form, searchResults, grantedAccess *tview.List) {
	// Create a focus group for inter-form navigation
	focusGroup := []tview.Primitive{
		leftForm,          // Vault Configuration
		optionsForm,       // Options
		grantAccessForm,   // Grant Access
		shareWithSelfForm, // Share with Self
		shareWithK8sForm,  // Share with K8s Context
		searchResults,     // Search Results list
		grantedAccess,     // Granted Access list
	}

	currentFocus := 0

	// Create a shared input capture function for inter-form navigation
	createInterFormInputCapture := func() func(*tcell.EventKey) *tcell.EventKey {
		return func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyTab:
				// Move to next form/primitive
				currentFocus = (currentFocus + 1) % len(focusGroup)
				vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
				return nil
			case tcell.KeyBacktab:
				// Move to previous form/primitive
				currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
				vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
				return nil
			case tcell.KeyCtrlS:
				// Create vault with Ctrl+S
				vp.createVaultFromForm(leftForm, optionsForm)
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
			vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyBacktab:
			// Move to previous form/primitive
			currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
			vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyDown:
			// Move to next field
			i, _ := leftForm.GetFocusedItemIndex()
			next := (i + 1) % leftForm.GetFormItemCount()
			vp.tui.GetApplication().SetFocus(leftForm.GetFormItem(next))
			return nil
		case tcell.KeyUp:
			// Move to previous field
			i, _ := leftForm.GetFocusedItemIndex()
			prev := (i - 1 + leftForm.GetFormItemCount()) % leftForm.GetFormItemCount()
			vp.tui.GetApplication().SetFocus(leftForm.GetFormItem(prev))
			return nil
		case tcell.KeyCtrlS:
			// Create vault with Ctrl+S
			vp.createVaultFromForm(leftForm, optionsForm)
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
			vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyBacktab:
			// Move to previous form/primitive
			currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
			vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyDown:
			// Move to next field
			i, _ := optionsForm.GetFocusedItemIndex()
			next := (i + 1) % optionsForm.GetFormItemCount()
			vp.tui.GetApplication().SetFocus(optionsForm.GetFormItem(next))
			return nil
		case tcell.KeyUp:
			// Move to previous field
			i, _ := optionsForm.GetFocusedItemIndex()
			prev := (i - 1 + optionsForm.GetFormItemCount()) % optionsForm.GetFormItemCount()
			vp.tui.GetApplication().SetFocus(optionsForm.GetFormItem(prev))
			return nil
		case tcell.KeyCtrlS:
			// Create vault with Ctrl+S
			vp.createVaultFromForm(leftForm, optionsForm)
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
			vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyBacktab:
			// Move to previous form/primitive
			currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
			vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyEnter:
			// Handle Enter key for SLV_EPK inputs
			vp.handleGrantAccessEnter()
			return nil
		case tcell.KeyCtrlS:
			// Create vault with Ctrl+S
			vp.createVaultFromForm(leftForm, optionsForm)
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
			vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyBacktab:
			// Move to previous form/primitive
			currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
			vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyCtrlS:
			// Create vault with Ctrl+S
			vp.createVaultFromForm(leftForm, optionsForm)
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
					if env, exists := vp.searchEnvMap[envName]; exists {
						vp.addToGrantedAccess(env)
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
			vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyBacktab:
			// Move to previous form/primitive
			currentFocus = (currentFocus - 1 + len(focusGroup)) % len(focusGroup)
			vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
			return nil
		case tcell.KeyCtrlS:
			// Create vault with Ctrl+S
			vp.createVaultFromForm(leftForm, optionsForm)
			return nil
		case tcell.KeyCtrlD:
			// Remove selected environment from granted access
			selected := grantedAccess.GetCurrentItem()
			if selected >= 0 && selected < grantedAccess.GetItemCount() {
				mainText, _ := grantedAccess.GetItemText(selected)
				// Extract environment name from the formatted text
				if strings.HasPrefix(mainText, "🌍 ") {
					envName := strings.TrimPrefix(mainText, "🌍 ")
					vp.removeFromGrantedAccess(envName)
				}
			}
			return nil
		}
		// Let all other keys pass through for list navigation
		return event
	})

	// Set initial focus to the first form
	vp.tui.GetApplication().SetFocus(focusGroup[currentFocus])
}

// searchEnvironments performs real-time environment search
func (vp *VaultPage) searchEnvironments(query string) {
	vp.currentQuery = query // Store the current query for refreshing
	vp.searchResults.Clear()
	vp.searchEnvMap = make(map[string]*environments.Environment) // Clear previous search results
	profile, err := profiles.GetActiveProfile()
	if err != nil {
		vp.showError(fmt.Sprintf("Error getting active profile: %v", err))
		vp.searchResults.AddItem("", "", 0, nil)
		return
	}

	// Helper function to check if environment is already granted access
	isAlreadyGranted := func(env *environments.Environment) bool {
		for _, grantedEnv := range vp.grantedEnvs {
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
			vp.showError(fmt.Sprintf("Error listing environments: %v", err))
			vp.searchResults.AddItem("", "", 0, nil)
			return
		}
		for _, env := range envs {
			if env.PublicKey == query {
				matchingEnvs = append(matchingEnvs, env)
			}
		}
		if len(matchingEnvs) == 0 {
			vp.searchResults.AddItem("❌ Environment not found in the profile", "", 0, nil)
			return
		}
	} else {
		envs, err := profile.SearchEnvs([]string{query})
		if err != nil {
			vp.showError(fmt.Sprintf("Error searching environments: %v", err))
			vp.searchResults.AddItem("", "", 0, nil)
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
			vp.searchResults.AddItem("✅ Environment already has access", fmt.Sprintf("Name: %s | Type: %s", env.Name, string(env.EnvType)), 0, nil)
		} else {
			// Store environment in map for later retrieval
			vp.searchEnvMap[env.Name] = env
			// Format with colors and proper spacing
			mainText := fmt.Sprintf("🔍 %s", env.Name)
			secondaryText := fmt.Sprintf("Type: %s | Email: %s", string(env.EnvType), env.Email)
			vp.searchResults.AddItem(mainText, secondaryText, 0, func() {
				vp.addToGrantedAccess(env)
			})
		}
	}
}

// addToGrantedAccess adds an environment to granted access
func (vp *VaultPage) addToGrantedAccess(env *environments.Environment) {
	// Check if environment is already granted
	for _, existingEnv := range vp.grantedEnvs {
		if existingEnv.PublicKey == env.PublicKey {
			return // Already granted
		}
	}

	// Add to granted environments
	vp.grantedEnvs = append(vp.grantedEnvs, env)

	// Update the display
	vp.updateGrantedAccessList()

	// Refresh search results to show the newly added environment as "already granted"
	vp.refreshSearchResults()
}

// refreshSearchResults refreshes the search results with the current query
func (vp *VaultPage) refreshSearchResults() {
	if vp.currentQuery != "" {
		vp.searchEnvironments(vp.currentQuery)
	}
}

// removeFromGrantedAccess removes an environment from granted access
func (vp *VaultPage) removeFromGrantedAccess(envName string) {
	// Find and remove the environment from grantedEnvs
	var removedEnv *environments.Environment
	for i, env := range vp.grantedEnvs {
		if env.Name == envName {
			removedEnv = env
			// Remove the environment at index i
			vp.grantedEnvs = append(vp.grantedEnvs[:i], vp.grantedEnvs[i+1:]...)
			break
		}
	}

	// Check if the removed environment is the self environment
	if removedEnv != nil {
		selfEnv := environments.GetSelf()
		if selfEnv != nil && selfEnv.PublicKey == removedEnv.PublicKey {
			// Uncheck the "Share with Self" checkbox
			vp.uncheckShareWithSelf()
		}
	}

	// Update the display
	vp.updateGrantedAccessList()

	// Refresh search results to show the removed environment as available again
	vp.refreshSearchResults()
}

// uncheckShareWithSelf unchecks the Share with Self checkbox
func (vp *VaultPage) uncheckShareWithSelf() {
	if vp.shareWithSelfCheckbox != nil {
		vp.shareWithSelfCheckbox.SetChecked(false)
	}
}

// createVaultFromForm creates a vault from the form data
func (vp *VaultPage) createVaultFromForm(leftForm, optionsForm *tview.Form) {
	// Collect form data
	vaultName := leftForm.GetFormItem(0).(*tview.InputField).GetText()
	fileName := leftForm.GetFormItem(1).(*tview.InputField).GetText()
	namespace := leftForm.GetFormItem(2).(*tview.InputField).GetText()

	// Get checkbox states
	enableHashing := optionsForm.GetFormItem(0).(*tview.Checkbox).IsChecked()
	quantumSafe := optionsForm.GetFormItem(1).(*tview.Checkbox).IsChecked()

	// Validate inputs
	if err := vp.validateVaultInputsForCreation(vaultName, fileName); err != nil {
		vp.showError(err.Error())
		return
	}

	// Prepare vault file path
	vaultFilePath := filepath.Join(vp.currentDir, fileName)

	// Collect public keys for vault access from granted environments
	var publicKeys []*crypto.PublicKey
	for _, env := range vp.grantedEnvs {
		if pk, err := env.GetPublicKey(); err == nil {
			publicKeys = append(publicKeys, pk)
		}
	}

	// Create the vault
	_, err := vaults.New(vaultFilePath, vaultName, namespace, enableHashing, quantumSafe, publicKeys...)
	if err != nil {
		vp.showError(fmt.Sprintf("Failed to create vault: %v", err))
		return
	}

	// Show success message
	vp.showSuccess(fmt.Sprintf("Vault '%s' created successfully at %s", vaultName, vaultFilePath))

	// TODO: Navigate back to vault browser or refresh the current view
}

// validateVaultInputs validates the vault creation inputs
func (vp *VaultPage) validateVaultInputsForCreation(vaultName, fileName string) error {
	// Validate vault name
	if strings.TrimSpace(vaultName) == "" {
		return fmt.Errorf("vault name is required")
	}

	// Validate file name
	if strings.TrimSpace(fileName) == "" {
		return fmt.Errorf("file name is required")
	}

	// Check if file already exists
	vaultFilePath := filepath.Join(vp.currentDir, fileName)
	if _, err := os.Stat(vaultFilePath); err == nil {
		return fmt.Errorf("file already exists: %s", fileName)
	}

	if len(vp.grantedEnvs) == 0 {
		return fmt.Errorf("no environments granted access. please grant access to at least one environment")
	}

	return nil
}

// showError displays an error message using the TUI's built-in error modal
func (vp *VaultPage) showError(message string) {
	vp.tui.ShowError(message)
}

// showSuccess displays a success message using the TUI's built-in info modal
func (vp *VaultPage) showSuccess(message string) {
	vp.tui.ShowInfo("✅ Success: " + message)
}

// handleShareWithSelfChange handles the "Share with Self" checkbox change
func (vp *VaultPage) handleShareWithSelfChange(checked bool) {
	if checked {
		// Add self environment to granted access
		selfEnv := environments.GetSelf()
		if selfEnv != nil {
			// Check if self environment is already granted
			alreadyGranted := false
			for _, env := range vp.grantedEnvs {
				if env.PublicKey == selfEnv.PublicKey {
					alreadyGranted = true
					break
				}
			}

			if !alreadyGranted {
				vp.grantedEnvs = append(vp.grantedEnvs, selfEnv)
				vp.updateGrantedAccessList()
				vp.refreshSearchResults()
			}
		}
	} else {
		// Remove self environment from granted access
		selfEnv := environments.GetSelf()
		if selfEnv != nil {
			for i, env := range vp.grantedEnvs {
				if env.PublicKey == selfEnv.PublicKey {
					vp.grantedEnvs = append(vp.grantedEnvs[:i], vp.grantedEnvs[i+1:]...)
					break
				}
			}
			vp.updateGrantedAccessList()
			vp.refreshSearchResults()
		}
	}
}

// handleShareWithK8sChange handles the "Share with K8s Context" checkbox change
func (vp *VaultPage) handleShareWithK8sChange(checked bool) {
	if checked {
		// Try to get K8s environment details
		k8sEnv, err := vp.getK8sEnvironment()
		if err != nil {
			// Show error and uncheck the checkbox
			vp.showError(fmt.Sprintf("Failed to get K8s environment: %v", err))
			// Note: We can't easily uncheck the checkbox here since we don't have a reference
			// The user will need to manually uncheck it after seeing the error
			return
		}

		// Check if K8s environment is already granted
		alreadyGranted := false
		for _, env := range vp.grantedEnvs {
			if env.PublicKey == k8sEnv.PublicKey {
				alreadyGranted = true
				break
			}
		}

		if !alreadyGranted {
			vp.grantedEnvs = append(vp.grantedEnvs, k8sEnv)
			vp.updateGrantedAccessList()
			vp.refreshSearchResults()
		}
	}
}

// getK8sEnvironment gets environment details from K8s context
func (vp *VaultPage) getK8sEnvironment() (*environments.Environment, error) {
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
		return vp.createNoInformationK8sEnvironment(publicKeyStr, namespace), nil
	}

	// Get all environments from the profile
	allEnvs, err := profile.ListEnvs()
	if err != nil {
		// If we can't list environments, create a "No information" environment
		return vp.createNoInformationK8sEnvironment(publicKeyStr, namespace), nil
	}

	// Search for matching environment by public key
	for _, env := range allEnvs {
		if env.PublicKey == publicKeyStr {
			// Found matching environment, return it
			return env, nil
		}
	}

	// No matching environment found, create a "No information" environment
	return vp.createNoInformationK8sEnvironment(publicKeyStr, namespace), nil
}

// createNoInformationK8sEnvironment creates a "No information" environment for K8s context
func (vp *VaultPage) createNoInformationK8sEnvironment(publicKeyStr, namespace string) *environments.Environment {
	return &environments.Environment{
		PublicKey: publicKeyStr,
		Name:      "No information",
		Email:     "Unknown",
		EnvType:   environments.SERVICE,
		Tags:      []string{"kubernetes", "context", namespace, "no-information"},
	}
}

// handleGrantAccessEnter handles Enter key press in Grant Access field
func (vp *VaultPage) handleGrantAccessEnter() {
	// Get the current input text from the Grant Access field
	// We need to find the input field in the grant access form
	// Since we know it's the first (and only) form item, we can access it directly
	if vp.grantAccessForm == nil {
		return
	}

	formItem := vp.grantAccessForm.GetFormItem(0)
	if inputField, ok := formItem.(*tview.InputField); ok {
		text := inputField.GetText()

		// Check if the text starts with "SLV_EPK"
		if strings.HasPrefix(text, "SLV_EPK") {
			// Extract the public key (remove "SLV_EPK" prefix)
			publicKeyStr := strings.TrimSpace(text)

			if publicKeyStr == "" {
				vp.tui.ShowError("Invalid public key format. Please provide a valid public key.")
				return
			}

			// Validate the public key format
			_, err := crypto.PublicKeyFromString(publicKeyStr)
			if err != nil {
				vp.tui.ShowError(fmt.Sprintf("Invalid public key format: %v", err))
				return
			}

			// Search through profile to see if this public key matches any existing environment
			profile, err := profiles.GetActiveProfile()
			if err != nil {
				// If we can't get the profile, create a "No information" environment
				env := vp.createNoInformationEnvironment(publicKeyStr)
				vp.addEnvironmentToGrantedAccess(env)
				return
			}

			// Get all environments from the profile
			allEnvs, err := profile.ListEnvs()
			if err != nil {
				// If we can't list environments, create a "No information" environment
				env := vp.createNoInformationEnvironment(publicKeyStr)
				vp.addEnvironmentToGrantedAccess(env)
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
				vp.addEnvironmentToGrantedAccess(foundEnv)
			} else {
				// No matching environment found, create a "No information" environment
				env := vp.createNoInformationEnvironment(publicKeyStr)
				vp.addEnvironmentToGrantedAccess(env)
			}

			// Clear the input field
			inputField.SetText("")
		}
	}
}

// createNoInformationEnvironment creates a "No information" environment for a public key
func (vp *VaultPage) createNoInformationEnvironment(publicKeyStr string) *environments.Environment {
	return &environments.Environment{
		PublicKey: publicKeyStr,
		Name:      "No information",
		Email:     "Unknown",
		EnvType:   environments.SERVICE,
		Tags:      []string{"no-information", "manual-entry"},
	}
}

// addEnvironmentToGrantedAccess adds an environment to the granted access list if not already present
func (vp *VaultPage) addEnvironmentToGrantedAccess(env *environments.Environment) {
	// Check if environment is already granted
	alreadyGranted := false
	for _, grantedEnv := range vp.grantedEnvs {
		if grantedEnv.PublicKey == env.PublicKey {
			alreadyGranted = true
			break
		}
	}

	if !alreadyGranted {
		vp.grantedEnvs = append(vp.grantedEnvs, env)
		vp.updateGrantedAccessList()
		vp.refreshSearchResults()
	}
}

// updateGrantedAccessList updates the granted access display
func (vp *VaultPage) updateGrantedAccessList() {
	vp.grantedAccess.Clear()

	// Add public keys
	for i, key := range vp.publicKeys {
		keyDisplay := key
		if len(key) > 25 {
			keyDisplay = key[:25] + "..."
		}
		mainText := fmt.Sprintf("🔑 Public Key %d", i+1)
		secondaryText := fmt.Sprintf("Key: %s", keyDisplay)
		vp.grantedAccess.AddItem(mainText, secondaryText, 0, nil)
	}

	// Add granted environments (sorted by name for consistent display)
	sortedGrantedEnvs := make([]*environments.Environment, len(vp.grantedEnvs))
	copy(sortedGrantedEnvs, vp.grantedEnvs)
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
		vp.grantedAccess.AddItem(mainText, secondaryText, 0, nil)
	}

	// Show message if no access granted
	if len(vp.publicKeys) == 0 && len(vp.grantedEnvs) == 0 {
		vp.grantedAccess.AddItem("📝 No access granted yet", "Add public keys or environments to grant access", 0, nil)
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// VaultFile represents a directory or .slv file
type VaultFile struct {
	Name   string
	Path   string
	IsFile bool
}

// getVaultFiles scans the home directory for directories and .slv files
func (vp *VaultPage) getVaultFiles() []VaultFile {
	var items []VaultFile

	// Read the directory
	entries, err := os.ReadDir(vp.currentDir)
	if err != nil {
		return items
	}

	// Filter and collect items
	for _, entry := range entries {
		// Skip hidden files (except .slv files)
		if strings.HasPrefix(entry.Name(), ".") && !strings.HasSuffix(entry.Name(), ".slv.yaml") && !strings.HasSuffix(entry.Name(), ".slv.yml") {
			continue
		}

		// Check if it's a directory
		if entry.IsDir() {
			items = append(items, VaultFile{
				Name:   entry.Name(),
				Path:   filepath.Join(vp.currentDir, entry.Name()),
				IsFile: false,
			})
		} else {
			// Check if it's a .slv file
			if strings.HasSuffix(entry.Name(), ".slv.yaml") || strings.HasSuffix(entry.Name(), ".slv.yml") {
				items = append(items, VaultFile{
					Name:   entry.Name(),
					Path:   filepath.Join(vp.currentDir, entry.Name()),
					IsFile: true,
				})
			}
		}
	}

	return items
}

// handleItemSelection handles when a user selects an item
func (vp *VaultPage) handleItemSelection(item VaultFile) {
	if item.IsFile {
		// Handle .slv file selection - open for viewing
		vp.openVaultFile(item.Path)
	} else {
		// Handle directory selection - navigate into the directory
		vp.tui.GetNavigation().SetVaultDir(item.Path)
		// Replace the current vault page with new directory
		vp.tui.GetNavigation().ShowVaultsReplace()
	}
}

// loadSelectedItem loads the currently selected item
func (vp *VaultPage) loadSelectedItem(list *tview.List) {
	// Get the current selection index
	selectedIndex := list.GetCurrentItem()

	// Skip the "Go Back" option (index 0)
	if selectedIndex == 0 {
		vp.goBackDirectory()
		return
	}

	// Adjust index for the "Go Back" option
	itemIndex := selectedIndex - 1

	// Get the items
	items := vp.getVaultFiles()

	// Check if the index is valid
	if itemIndex >= 0 && itemIndex < len(items) {
		item := items[itemIndex]
		vp.handleItemSelection(item)
	}
}

// goBackDirectory goes back to the parent directory
func (vp *VaultPage) goBackDirectory() {
	parentDir := filepath.Dir(vp.currentDir)
	// Don't go back if we're already at the root
	if parentDir != vp.currentDir {
		vp.tui.GetNavigation().SetVaultDir(parentDir)
		// Replace the current vault page with parent directory
		vp.tui.GetNavigation().ShowVaultsReplace()
	}
}

// openVaultFile opens and displays an SLV file for viewing
func (vp *VaultPage) openVaultFile(filePath string) {
	// Check if we already have this vault loaded
	if vp.vault != nil && vp.vaultPath == filePath {
		// Use existing vault instance
		vp.showVaultDetails(vp.vault, filePath)
		return
	}

	// Load the vault using vaults.Get
	vault, err := vaults.Get(filePath)
	if err != nil {
		vp.showError(fmt.Sprintf("Error loading vault: %v", err))
		return
	}

	// Store the vault instance and path
	vp.vault = vault
	vp.vaultPath = filePath

	// Create and show vault details page
	vp.showVaultDetails(vault, filePath)
}

// showVaultDetails displays detailed information about a vault
func (vp *VaultPage) showVaultDetails(vault *vaults.Vault, filePath string) {
	// Create a flex layout to hold the three tables
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// 1. Vault Details Table (30% height)
	vaultDetailsTable := vp.createVaultDetailsTable(vault, filePath)
	flex.AddItem(vaultDetailsTable, 0, 30, true) // First table gets focus

	// 2. Accessors Table (30% height)
	accessorsTable := vp.createAccessorsTable(vault)
	flex.AddItem(accessorsTable, 0, 30, false)

	// 3. Vault Items Table (40% height)
	itemsTable := vp.createVaultItemsTable(vault)
	flex.AddItem(itemsTable, 0, 40, false)

	// Set initial focus to vault details table
	vp.tui.GetApplication().SetFocus(vaultDetailsTable)

	// Track current focus index (0 = vault details, 1 = accessors, 2 = items)
	currentFocusIndex := 0

	// Function to switch focus between tables
	switchFocus := func() {
		// Clear focus from all tables first
		vaultDetailsTable.SetSelectable(false, false)
		accessorsTable.SetSelectable(false, false)
		itemsTable.SetSelectable(false, false)

		// Set focus to the next table
		currentFocusIndex = (currentFocusIndex + 1) % 3
		switch currentFocusIndex {
		case 0:
			vaultDetailsTable.SetSelectable(true, false)
			vp.tui.GetApplication().SetFocus(vaultDetailsTable)
		case 1:
			accessorsTable.SetSelectable(true, false)
			vp.tui.GetApplication().SetFocus(accessorsTable)
		case 2:
			itemsTable.SetSelectable(true, false)
			vp.tui.GetApplication().SetFocus(itemsTable)
		}
	}

	// Update status bar with help text
	vp.tui.UpdateStatusBar("[yellow]q: close | u: unlock | l: lock | r: reload | Tab: switch tables[white]")

	// Set up input capture for the flex
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Switch focus between tables
			switchFocus()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				vp.clearVault()
				vp.tui.GetNavigation().ShowVaults()
				return nil
			case 'u', 'U':
				// Unlock vault
				if filePath != "" {
					vp.unlockVault(filePath)
				}
				return nil
			case 'l', 'L':
				// Lock vault
				if filePath != "" {
					vp.lockVault(filePath)
				}
				return nil
			case 'r', 'R':
				// Reload vault
				if filePath != "" {
					vp.reloadVault()
					vp.showVaultDetails(vp.vault, filePath)
				}
				return nil
			}
		case tcell.KeyEsc:
			vp.tui.GetNavigation().ShowVaults()
			return nil
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight, tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyHome, tcell.KeyEnd:
			// Allow arrow keys and page keys to scroll
			return event
		}
		return event
	})

	// Create the page layout and show it
	page := vp.tui.CreatePageLayout("Vault Details", flex)
	vp.tui.GetNavigation().ShowVaultDetails(page)
}

// unlockVault attempts to unlock the vault
func (vp *VaultPage) unlockVault(filePath string) {
	// Check if we have the vault loaded
	if vp.vault == nil || vp.vaultPath != filePath {
		vp.showError("Vault not loaded. Please reopen the vault.")
		return
	}

	// If already unlocked, just refresh the display
	if !vp.vault.IsLocked() {
		vp.showVaultDetails(vp.vault, filePath)
		return
	}

	// Attempt to unlock the vault
	secretKey, err := session.GetSecretKey()
	if err != nil {
		vp.showError(fmt.Sprintf("Error getting secret key: %v", err))
		return
	}

	err = vp.vault.Unlock(secretKey)
	if err != nil {
		vp.showError(fmt.Sprintf("Error unlocking vault: %v", err))
		return
	}

	vp.showVaultDetails(vp.vault, filePath)
}

// lockVault locks the vault
func (vp *VaultPage) lockVault(filePath string) {
	// Check if we have the vault loaded
	if vp.vault == nil || vp.vaultPath != filePath {
		vp.showError("Vault not loaded. Please reopen the vault.")
		return
	}

	// If already locked, just refresh the display
	if vp.vault.IsLocked() {
		vp.showVaultDetails(vp.vault, filePath)
		return
	}

	// Lock the vault
	vp.vault.Lock()

	// Refresh the vault details page using the stored instance
	vp.showVaultDetails(vp.vault, filePath)
}

// clearVault clears the stored vault instance
func (vp *VaultPage) clearVault() {
	vp.vault = nil
	vp.vaultPath = ""
}

// reloadVault reloads the vault from disk (useful if file was modified externally)
func (vp *VaultPage) reloadVault() {
	if vp.vaultPath == "" {
		return
	}

	// Load fresh vault instance
	vault, err := vaults.Get(vp.vaultPath)
	if err != nil {
		vp.showError(fmt.Sprintf("Error reloading vault: %v", err))
		return
	}

	// Update stored instance
	vp.vault = vault
}

// createVaultDetailsTable creates a table for vault details
func (vp *VaultPage) createVaultDetailsTable(vault *vaults.Vault, filePath string) *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true).SetTitle("Metadata").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0) // Fix the first row (header) and no columns

	// Set headers (non-selectable) with fixed width for first column
	table.SetCell(0, 0, tview.NewTableCell("Property").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(20))
	table.SetCell(0, 1, tview.NewTableCell("Value").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))

	// Add vault details
	row := 1
	table.SetCell(row, 0, tview.NewTableCell("Vault Path").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(filePath).SetTextColor(tcell.ColorWhite))
	row++

	table.SetCell(row, 0, tview.NewTableCell("Vault Name").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(vault.ObjectMeta.Name).SetTextColor(tcell.ColorWhite))
	row++

	if vault.ObjectMeta.Namespace != "" {
		table.SetCell(row, 0, tview.NewTableCell("Namespace").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
		table.SetCell(row, 1, tview.NewTableCell(vault.ObjectMeta.Namespace).SetTextColor(tcell.ColorWhite))
		row++
	} else {
		table.SetCell(row, 0, tview.NewTableCell("Namespace").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
		table.SetCell(row, 1, tview.NewTableCell("No Namespace").SetTextColor(tcell.ColorWhite))
		row++
	}

	table.SetCell(row, 0, tview.NewTableCell("Public Key").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(vault.Spec.Config.PublicKey).SetTextColor(tcell.ColorWhite))
	row++

	table.SetCell(row, 0, tview.NewTableCell("Number of Accessors").SetTextColor(tcell.ColorAqua).SetMaxWidth(20))
	table.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%d", len(vault.Spec.Config.WrappedKeys))).SetTextColor(tcell.ColorWhite))
	row++
	// Make table focusable for scrolling with custom selection colors
	table.SetSelectable(true, false) // Vault details table is initially selectable
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))

	return table
}

// createAccessorsTable creates a table for accessors
func (vp *VaultPage) createAccessorsTable(vault *vaults.Vault) *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true).SetTitle("Access").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0) // Fix the first row (header) and no columns

	// Set headers (non-selectable) with fixed column widths
	table.SetCell(0, 0, tview.NewTableCell("Type").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(8))
	table.SetCell(0, 1, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(30))
	table.SetCell(0, 2, tview.NewTableCell("Email").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(30))
	table.SetCell(0, 3, tview.NewTableCell("Public Key").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))

	accessors, err := vault.ListAccessors()
	if err != nil || len(accessors) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No accessors found").SetTextColor(tcell.ColorGray).SetAlign(tview.AlignCenter))
		table.SetCell(1, 1, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		table.SetCell(1, 2, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		table.SetCell(1, 3, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		return table
	}

	// Get profile and environment information
	profile, _ := profiles.GetActiveProfile()
	var root *environments.Environment
	if profile != nil {
		root, _ = profile.GetRoot()
	}
	self := environments.GetSelf()

	row := 1
	for _, accessor := range accessors {
		accessorPubKey, err := accessor.String()
		if err != nil {
			continue
		}

		// Determine accessor type and name
		var accessorType, accessorName, accessorEmail string
		if self != nil && self.PublicKey == accessorPubKey {
			accessorType = "Self"
			accessorName = self.Name
			accessorEmail = self.Email
		} else if root != nil && root.PublicKey == accessorPubKey {
			accessorType = "Root"
			accessorName = root.Name
			accessorEmail = root.Email
		} else if profile != nil {
			if env, _ := profile.GetEnv(accessorPubKey); env != nil {
				if env.EnvType == environments.USER {
					accessorType = "User"
				} else {
					accessorType = "Service"
				}
				accessorName = env.Name
				accessorEmail = env.Email
			} else {
				accessorType = "Unknown"
				accessorName = ""
				accessorEmail = ""
			}
		} else {
			accessorType = "Unknown"
			accessorName = ""
			accessorEmail = ""
		}

		table.SetCell(row, 0, tview.NewTableCell(accessorType).SetTextColor(tcell.ColorAqua).SetMaxWidth(8))
		table.SetCell(row, 1, tview.NewTableCell(accessorName).SetTextColor(tcell.ColorGreen).SetMaxWidth(30))
		table.SetCell(row, 2, tview.NewTableCell(accessorEmail).SetTextColor(tcell.ColorWhite).SetMaxWidth(30))
		table.SetCell(row, 3, tview.NewTableCell(accessorPubKey).SetTextColor(tcell.ColorGray))
		row++
	}

	// Make table focusable for scrolling with custom selection colors
	table.SetSelectable(false, false) // Initially not selectable
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))

	return table
}

// createVaultItemsTable creates a table for vault items
func (vp *VaultPage) createVaultItemsTable(vault *vaults.Vault) *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true).SetTitle("Items").SetTitleAlign(tview.AlignLeft)
	table.SetFixed(1, 0) // Fix the first row (header) and no columns

	// Set headers (non-selectable) with fixed column widths
	table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(40))
	table.SetCell(0, 1, tview.NewTableCell("Type").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(12))
	table.SetCell(0, 2, tview.NewTableCell("Value").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false))

	itemNames := vault.GetItemNames()
	if len(itemNames) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No items found").SetTextColor(tcell.ColorGray).SetAlign(tview.AlignCenter))
		table.SetCell(1, 1, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		table.SetCell(1, 2, tview.NewTableCell("").SetTextColor(tcell.ColorGray))
		return table
	}

	// Sort item names for consistent display order
	sort.Strings(itemNames)

	row := 1
	for _, name := range itemNames {
		table.SetCell(row, 0, tview.NewTableCell(name).SetTextColor(tcell.ColorGreen).SetMaxWidth(25))

		if !vault.IsLocked() {
			// Vault is unlocked - show actual item details
			item, err := vault.Get(name)
			if err == nil {
				encryptedStatus := "Secret"
				if item.IsPlaintext() {
					encryptedStatus = "Plaintext"
				}
				table.SetCell(row, 1, tview.NewTableCell(encryptedStatus).SetTextColor(tcell.ColorWhite).SetMaxWidth(12))

				value, err := item.ValueString()
				if err != nil {
					value = "Error loading value"
				}
				table.SetCell(row, 2, tview.NewTableCell(value).SetTextColor(tcell.ColorWhite))
			} else {
				table.SetCell(row, 1, tview.NewTableCell("Error").SetTextColor(tcell.ColorRed).SetMaxWidth(12))
				table.SetCell(row, 2, tview.NewTableCell("Error loading item").SetTextColor(tcell.ColorRed))
			}
		} else {
			// Vault is locked - show masked value
			table.SetCell(row, 1, tview.NewTableCell("***").SetTextColor(tcell.ColorYellow).SetMaxWidth(12))
			table.SetCell(row, 2, tview.NewTableCell("***").SetTextColor(tcell.ColorGray))
		}
		row++
	}

	// Make table focusable for scrolling with custom selection colors
	table.SetSelectable(false, false) // Initially not selectable
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))

	return table
}

// showSuccess displays a success message

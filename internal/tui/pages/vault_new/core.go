package vault_new

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
)

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

// removeFromGrantedAccess removes an environment from granted access by public key
func (vnp *VaultNewPage) removeFromGrantedAccess(publicKey string) {
	// Find and remove the environment from grantedEnvs
	var removedEnv *environments.Environment
	for i, env := range vnp.grantedEnvs {
		if env.PublicKey == publicKey {
			removedEnv = env
			// Remove the environment at index i
			vnp.grantedEnvs = append(vnp.grantedEnvs[:i], vnp.grantedEnvs[i+1:]...)
			break
		}
	}

	// Check if the removed environment is the self environment or K8s environment
	if removedEnv != nil {
		selfEnv := environments.GetSelf()
		if selfEnv != nil && selfEnv.PublicKey == removedEnv.PublicKey {
			// Uncheck the "Share with Self" checkbox
			vnp.uncheckShareWithSelf()
		}

		// Check if the removed environment is the K8s environment
		if vnp.k8sEnv != nil && vnp.k8sEnv.PublicKey == removedEnv.PublicKey {
			// Uncheck the "Share with K8s Context" checkbox
			vnp.uncheckShareWithK8s()
			// Clear the stored K8s environment
			vnp.k8sEnv = nil
		}
	}

	// Update the display
	vnp.updateGrantedAccessList()

	// Refresh search results to show the removed environment as available again
	vnp.refreshSearchResults()
}

// uncheckShareWithSelf unchecks the Share with Self checkbox
func (vnp *VaultNewPage) uncheckShareWithSelf() {
	if vnp.shareWithSelfForm != nil {
		vnp.shareWithSelfForm.GetFormItem(0).(*tview.Checkbox).SetChecked(false)
	}
}

// uncheckShareWithK8s unchecks the Share with K8s Context checkbox
func (vnp *VaultNewPage) uncheckShareWithK8s() {
	if vnp.shareWithK8sForm != nil {
		vnp.shareWithK8sForm.GetFormItem(0).(*tview.Checkbox).SetChecked(false)
	}
}

func (vnp *VaultNewPage) createVaultFromForm() {
	// Collect form data
	vaultName := vnp.vaultConfigForm.GetFormItem(0).(*tview.InputField).GetText()
	fileName := vnp.vaultConfigForm.GetFormItem(1).(*tview.InputField).GetText()
	namespace := vnp.vaultConfigForm.GetFormItem(2).(*tview.InputField).GetText()

	// Get checkbox states
	enableHashing := vnp.optionsForm.GetFormItem(0).(*tview.Checkbox).IsChecked()
	quantumSafe := vnp.optionsForm.GetFormItem(1).(*tview.Checkbox).IsChecked()

	// Validate inputs
	if err := vnp.validateVaultInputsForCreation(vaultName, fileName); err != nil {
		vnp.showError(err.Error())
		return
	}

	if len(vnp.grantedEnvs) == 0 {
		vnp.GetTUI().ShowConfirmationWithFocus(
			"No environments granted access. Please grant access to at least one environment.",
			"Create with Self Environment",
			"Cancel",
			func() {
				vnp.grantedEnvs = append(vnp.grantedEnvs, environments.GetSelf())
				vnp.createVaultFromForm()
			}, nil, nil)
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
	// Show vault details page with the created vault and filepath
	vnp.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vault, vaultFilePath, true)

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

	if !strings.HasSuffix(fileName, ".slv.yaml") && !strings.HasSuffix(fileName, ".slv.yml") {
		return fmt.Errorf("file name must end with .slv.yaml or .slv.yml")
	}

	// Check if file already exists
	vaultFilePath := filepath.Join(vnp.currentDir, fileName)
	if _, err := os.Stat(vaultFilePath); err == nil {
		return fmt.Errorf("file already exists: %s", fileName)
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
			vnp.shareWithK8sForm.GetFormItem(0).(*tview.Checkbox).SetChecked(false)
			return
		}

		// Store the K8s environment
		vnp.k8sEnv = k8sEnv

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
	} else {
		// Remove K8s environment from granted access if it exists
		if vnp.k8sEnv != nil {
			for i, env := range vnp.grantedEnvs {
				if env.PublicKey == vnp.k8sEnv.PublicKey {
					vnp.grantedEnvs = append(vnp.grantedEnvs[:i], vnp.grantedEnvs[i+1:]...)
					break
				}
			}
			vnp.updateGrantedAccessList()
			vnp.refreshSearchResults()
		}
		// Clear the stored K8s environment
		vnp.k8sEnv = nil
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
func (vnp *VaultNewPage) handleSearchBarEnter() {
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
		var mainText string
		if env.EnvType == environments.SERVICE {
			mainText = fmt.Sprintf("💻 %s", name)
		} else if env.EnvType == environments.USER {
			mainText = fmt.Sprintf("👤 %s", name)
		} else {
			mainText = fmt.Sprintf("🌍 %s", name)
		}
		// Store the full public key in the secondary text for later retrieval
		// Format: "Email: xxx | PK: full_public_key"
		secondaryText := fmt.Sprintf("Email: %s | PK: %s", email, env.PublicKey)
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

package vault_edit

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

func (vep *VaultEditPage) searchEnvironments(query string) {
	vep.currentQuery = query // Store the current query for refreshing
	vep.searchResults.Clear()
	vep.searchEnvMap = make(map[string]*environments.Environment) // Clear previous search results
	profile, err := profiles.GetActiveProfile()
	if err != nil {
		vep.showError(fmt.Sprintf("Error getting active profile: %v", err))
		vep.searchResults.AddItem("", "", 0, nil)
		return
	}

	// Helper function to check if environment is already granted access
	isAlreadyGranted := func(env *environments.Environment) bool {
		for _, grantedEnv := range vep.grantedEnvs {
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
			vep.showError(fmt.Sprintf("Error listing environments: %v", err))
			vep.searchResults.AddItem("", "", 0, nil)
			return
		}
		for _, env := range envs {
			if env.PublicKey == query {
				matchingEnvs = append(matchingEnvs, env)
			}
		}
		if len(matchingEnvs) == 0 {
			vep.searchResults.AddItem("❌ Environment not found in the profile", "", 0, nil)
			return
		}
	} else {
		envs, err := profile.SearchEnvs([]string{query})
		if err != nil {
			vep.showError(fmt.Sprintf("Error searching environments: %v", err))
			vep.searchResults.AddItem("", "", 0, nil)
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
			vep.searchResults.AddItem("✅ Environment already has access", fmt.Sprintf("Name: %s | Type: %s", env.Name, string(env.EnvType)), 0, nil)
		} else {
			// Store environment in map for later retrieval
			vep.searchEnvMap[env.Name] = env
			// Format with colors and proper spacing
			mainText := fmt.Sprintf("🔍 %s", env.Name)
			secondaryText := fmt.Sprintf("Type: %s | Email: %s", string(env.EnvType), env.Email)
			vep.searchResults.AddItem(mainText, secondaryText, 0, func() {
				vep.addToGrantedAccess(env)
			})
		}
	}
}

// addToGrantedAccess adds an environment to the granted access list
func (vep *VaultEditPage) addToGrantedAccess(env *environments.Environment) {
	// Check if environment is already granted
	for _, existingEnv := range vep.grantedEnvs {
		if existingEnv.PublicKey == env.PublicKey {
			return // Already granted
		}
	}

	// Add to granted environments
	vep.grantedEnvs = append(vep.grantedEnvs, env)

	// Update the display
	vep.updateGrantedAccessList()

	// Refresh search results to show the newly added environment as "already granted"
	vep.refreshSearchResults()
}

func (vep *VaultEditPage) refreshSearchResults() {
	if vep.currentQuery != "" {
		vep.searchEnvironments(vep.currentQuery)
	}
}

// removeFromGrantedAccess removes an environment from granted access by public key
func (vep *VaultEditPage) removeFromGrantedAccess(publicKey string) {
	// Find and remove the environment from grantedEnvs
	var removedEnv *environments.Environment
	for i, env := range vep.grantedEnvs {
		if env.PublicKey == publicKey {
			removedEnv = env
			// Remove the environment at index i
			vep.grantedEnvs = append(vep.grantedEnvs[:i], vep.grantedEnvs[i+1:]...)
			break
		}
	}

	// Check if the removed environment is the self environment or K8s environment
	if removedEnv != nil {
		selfEnv := environments.GetSelf()
		if selfEnv != nil && selfEnv.PublicKey == removedEnv.PublicKey {
			// Uncheck the "Share with Self" checkbox
			vep.uncheckShareWithSelf()
		}

		// Check if the removed environment is the K8s environment
		if vep.k8sEnv != nil && vep.k8sEnv.PublicKey == removedEnv.PublicKey {
			// Uncheck the "Share with K8s Context" checkbox
			vep.uncheckShareWithK8s()
			// Clear the stored K8s environment
			vep.k8sEnv = nil
		}
	}

	// Update the display
	vep.updateGrantedAccessList()

	// Refresh search results to show the removed environment as available again
	vep.refreshSearchResults()
}

// uncheckShareWithSelf unchecks the Share with Self checkbox
func (vep *VaultEditPage) uncheckShareWithSelf() {
	if vep.shareWithSelfForm != nil {
		vep.shareWithSelfForm.GetFormItem(0).(*tview.Checkbox).SetChecked(false)
	}
}

// uncheckShareWithK8s unchecks the Share with K8s Context checkbox
func (vep *VaultEditPage) uncheckShareWithK8s() {
	if vep.shareWithK8sForm != nil {
		vep.shareWithK8sForm.GetFormItem(0).(*tview.Checkbox).SetChecked(false)
	}
}

func (vep *VaultEditPage) editVaultFromForm() {

	vaultName := vep.vaultConfigForm.GetFormItem(0).(*tview.InputField).GetText()
	fileName := vep.vaultConfigForm.GetFormItem(1).(*tview.InputField).GetText()
	namespace := vep.vaultConfigForm.GetFormItem(2).(*tview.InputField).GetText()

	if err := vep.validateVaultInputsForCreation(vaultName, fileName); err != nil {
		vep.showError(err.Error())
		return
	}

	if vep.vault.Name != vaultName {
		vep.vault.Name = vaultName
	}
	if vep.vault.Namespace != namespace {
		vep.vault.Namespace = namespace
	}

	// Construct new file path
	dir := filepath.Dir(vep.filePath)
	newFilePath := filepath.Join(dir, fileName)

	// Only rename if the filename has changed
	if newFilePath != vep.filePath {
		// Check if new file already exists
		if _, err := os.Stat(newFilePath); err == nil {
			vep.ShowError(fmt.Sprintf("File '%s' already exists", fileName))
			return
		}
		// Rename the file
		if err := os.Rename(vep.filePath, newFilePath); err != nil {
			vep.ShowError(fmt.Sprintf("Error renaming vault: %v", err))
			return
		}

		vep.filePath = newFilePath
	}

	vep.vault.Update(vaultName, namespace, "", nil)

	if !vep.IsVaultUnlocked() {
		// Show success message
		vep.showSuccess(fmt.Sprintf("Vault '%s' edited successfully at %s", vaultName, vep.filePath))

		// Navigate to vault details page and remove new vault page from stack
		vep.navigateToVaultDetails(vep.vault, vep.filePath)

		return
	}
	secretKey, err := session.GetSecretKey()
	if err != nil {
		vep.showError(fmt.Sprintf("error editing vault: %v", err))
		return
	}
	if err := vep.vault.Unlock(secretKey); err != nil {
		vep.showError(fmt.Sprintf("error editing vault: %v", err))
		return
	}

	// Collect public keys for vault access from granted environments
	var publicKeys []crypto.PublicKey
	for _, env := range vep.grantedEnvs {
		if pk, err := env.GetPublicKey(); err == nil {
			publicKeys = append(publicKeys, *pk)
		}
	}
	existingKeys, err := vep.vault.ListAccessors()
	if err != nil {
		vep.showError(fmt.Sprintf("error editing vault: %v", err))
		return
	}
	vep.ShowInfo(fmt.Sprintf("%v,%v", publicKeys, existingKeys))
	for _, key := range publicKeys {
		if !containsPublicKey(existingKeys, key) {
			vep.vault.Share(&key)
		}
	}

	// Revoke keys that are no longer granted
	var keysToRevoke []*crypto.PublicKey
	for _, key := range existingKeys {
		if !containsPublicKey(publicKeys, key) {
			keysToRevoke = append(keysToRevoke, &key)
		}
	}
	if err := vep.vault.Revoke(keysToRevoke, false); err != nil {
		vep.showError(fmt.Sprintf("error editing vault: %v", err))
		return
	}

	// Show success message
	vep.showSuccess(fmt.Sprintf("Vault '%s' edited successfully at %s", vaultName, vep.filePath))

	// Navigate to vault details page and remove new vault page from stack
	vep.navigateToVaultDetails(vep.vault, vep.filePath)
}

// navigateToVaultDetails navigates to the vault details page and removes new vault page from stack
func (vep *VaultEditPage) navigateToVaultDetails(vault *vaults.Vault, vaultFilePath string) {
	// Show vault details page with the created vault and filepath
	vep.GetTUI().GetNavigation().ShowVaultDetailsWithVault(vault, vaultFilePath, true)

	vep.showSuccess(fmt.Sprintf("Vault '%s' created successfully at %s", vault.Name, vaultFilePath))
}

// validateVaultInputs validates the vault creation inputs
func (vep *VaultEditPage) validateVaultInputsForCreation(vaultName, fileName string) error {
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

	if len(vep.grantedEnvs) == 0 {
		return fmt.Errorf("no environments granted access. please grant access to at least one environment")
	}

	return nil
}

// showError displays an error message using the TUI's built-in error modal
func (vep *VaultEditPage) showError(message string) {
	vep.GetTUI().ShowError(message)
}

// showSuccess displays a success message using the TUI's built-in info modal
func (vep *VaultEditPage) showSuccess(message string) {
	vep.GetTUI().ShowInfo("✅ Success: " + message)
}

// handleShareWithSelfChange handles the "Share with Self" checkbox change
func (vep *VaultEditPage) handleShareWithSelfChange(checked bool) {
	if checked {
		// Add self environment to granted access
		selfEnv := environments.GetSelf()
		if selfEnv != nil {
			// Check if self environment is already granted
			alreadyGranted := false
			for _, env := range vep.grantedEnvs {
				if env.PublicKey == selfEnv.PublicKey {
					alreadyGranted = true
					break
				}
			}

			if !alreadyGranted {
				vep.grantedEnvs = append(vep.grantedEnvs, selfEnv)
				vep.updateGrantedAccessList()
				vep.refreshSearchResults()
			}
		}
	} else {
		// Remove self environment from granted access
		selfEnv := environments.GetSelf()
		if selfEnv != nil {
			for i, env := range vep.grantedEnvs {
				if env.PublicKey == selfEnv.PublicKey {
					vep.grantedEnvs = append(vep.grantedEnvs[:i], vep.grantedEnvs[i+1:]...)
					break
				}
			}
			vep.updateGrantedAccessList()
			vep.refreshSearchResults()
		}
	}
}

// handleShareWithK8sChange handles the "Share with K8s Context" checkbox change
func (vep *VaultEditPage) handleShareWithK8sChange(checked bool) {
	if checked {
		// Try to get K8s environment details
		k8sEnv, err := vep.getK8sEnvironment()
		if err != nil {
			// Show error and uncheck the checkbox
			vep.showError(fmt.Sprintf("Failed to get K8s environment: %v", err))
			vep.shareWithK8sForm.GetFormItem(0).(*tview.Checkbox).SetChecked(false)
			return
		}

		// Store the K8s environment
		vep.k8sEnv = k8sEnv

		// Check if K8s environment is already granted
		alreadyGranted := false
		for _, env := range vep.grantedEnvs {
			if env.PublicKey == k8sEnv.PublicKey {
				alreadyGranted = true
				break
			}
		}

		if !alreadyGranted {
			vep.grantedEnvs = append(vep.grantedEnvs, k8sEnv)
			vep.updateGrantedAccessList()
			vep.refreshSearchResults()
		}
	} else {
		// Remove K8s environment from granted access if it exists
		if vep.k8sEnv != nil {
			for i, env := range vep.grantedEnvs {
				if env.PublicKey == vep.k8sEnv.PublicKey {
					vep.grantedEnvs = append(vep.grantedEnvs[:i], vep.grantedEnvs[i+1:]...)
					break
				}
			}
			vep.updateGrantedAccessList()
			vep.refreshSearchResults()
		}
		// Clear the stored K8s environment
		vep.k8sEnv = nil
	}
}

// getK8sEnvironment gets environment details from K8s context
func (vep *VaultEditPage) getK8sEnvironment() (*environments.Environment, error) {
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
		return vep.createNoInformationK8sEnvironment(publicKeyStr, namespace), nil
	}

	// Get all environments from the profile
	allEnvs, err := profile.ListEnvs()
	if err != nil {
		// If we can't list environments, create a "No information" environment
		return vep.createNoInformationK8sEnvironment(publicKeyStr, namespace), nil
	}

	// Search for matching environment by public key
	for _, env := range allEnvs {
		if env.PublicKey == publicKeyStr {
			// Found matching environment, return it
			return env, nil
		}
	}

	// No matching environment found, create a "No information" environment
	return vep.createNoInformationK8sEnvironment(publicKeyStr, namespace), nil
}

// createNoInformationK8sEnvironment creates a "No information" environment for K8s context
func (vep *VaultEditPage) createNoInformationK8sEnvironment(publicKeyStr, namespace string) *environments.Environment {
	return &environments.Environment{
		PublicKey: publicKeyStr,
		Name:      "No information",
		Email:     "Unknown",
		EnvType:   environments.SERVICE,
		Tags:      []string{"kubernetes", "context", namespace, "no-information"},
	}
}

// handleGrantAccessEnter handles Enter key press in Grant Access field
func (vep *VaultEditPage) handleSearchBarEnter() {
	// Get the current input text from the Grant Access field
	// We need to find the input field in the grant access form
	// Since we know it's the first (and only) form item, we can access it directly
	if vep.grantAccessForm == nil {
		return
	}

	formItem := vep.grantAccessForm.GetFormItem(0)
	if inputField, ok := formItem.(*tview.InputField); ok {
		text := inputField.GetText()

		// Check if the text starts with "SLV_EPK"
		if strings.HasPrefix(text, "SLV_EPK") {
			// Extract the public key (remove "SLV_EPK" prefix)
			publicKeyStr := strings.TrimSpace(text)

			if publicKeyStr == "" {
				vep.GetTUI().ShowError("Invalid public key format. Please provide a valid public key.")
				return
			}

			// Validate the public key format
			_, err := crypto.PublicKeyFromString(publicKeyStr)
			if err != nil {
				vep.GetTUI().ShowError(fmt.Sprintf("Invalid public key format: %v", err))
				return
			}

			// Search through profile to see if this public key matches any existing environment
			profile, err := profiles.GetActiveProfile()
			if err != nil {
				// If we can't get the profile, create a "No information" environment
				env := vep.createNoInformationEnvironment(publicKeyStr)
				vep.addEnvironmentToGrantedAccess(env)
				return
			}

			// Get all environments from the profile
			allEnvs, err := profile.ListEnvs()
			if err != nil {
				// If we can't list environments, create a "No information" environment
				env := vep.createNoInformationEnvironment(publicKeyStr)
				vep.addEnvironmentToGrantedAccess(env)
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
				vep.addEnvironmentToGrantedAccess(foundEnv)
			} else {
				// No matching environment found, create a "No information" environment
				env := vep.createNoInformationEnvironment(publicKeyStr)
				vep.addEnvironmentToGrantedAccess(env)
			}

			// Clear the input field
			inputField.SetText("")
		}
	}
}

// createNoInformationEnvironment creates a "No information" environment for a public key
func (vep *VaultEditPage) createNoInformationEnvironment(publicKeyStr string) *environments.Environment {
	return &environments.Environment{
		PublicKey: publicKeyStr,
		Name:      "No information",
		Email:     "Unknown",
		EnvType:   environments.SERVICE,
		Tags:      []string{"no-information", "manual-entry"},
	}
}

// addEnvironmentToGrantedAccess adds an environment to the granted access list if not already present
func (vep *VaultEditPage) addEnvironmentToGrantedAccess(env *environments.Environment) {
	// Check if environment is already granted
	alreadyGranted := false
	for _, grantedEnv := range vep.grantedEnvs {
		if grantedEnv.PublicKey == env.PublicKey {
			alreadyGranted = true
			break
		}
	}

	if !alreadyGranted {
		vep.grantedEnvs = append(vep.grantedEnvs, env)
		vep.updateGrantedAccessList()
		vep.refreshSearchResults()
	}
}

// updateGrantedAccessList updates the granted access display
func (vep *VaultEditPage) updateGrantedAccessList() {
	vep.grantedAccess.Clear()

	// Add public keys
	for i, key := range vep.publicKeys {
		keyDisplay := key
		if len(key) > 25 {
			keyDisplay = key[:25] + "..."
		}
		mainText := fmt.Sprintf("🔑 Public Key %d", i+1)
		secondaryText := fmt.Sprintf("Key: %s", keyDisplay)
		vep.grantedAccess.AddItem(mainText, secondaryText, 0, nil)
	}

	for _, env := range vep.grantedEnvs {
		// Handle unknown name and email

		name := env.Name
		if name == "" {
			name = "Unknown"
			env.Name = name
		}

		email := env.Email
		if email == "" {
			email = "Unknown"
			env.Email = email
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
		secondaryText := fmt.Sprintf("Email: %s | Public Key: %s", email, env.PublicKey)
		vep.grantedAccess.AddItem(mainText, secondaryText, 0, nil)
	}

	// Add granted environments (sorted by name for consistent display)
	sortedGrantedEnvs := make([]*environments.Environment, len(vep.grantedEnvs))
	copy(sortedGrantedEnvs, vep.grantedEnvs)
	sort.Slice(sortedGrantedEnvs, func(i, j int) bool {
		return sortedGrantedEnvs[i].Name < sortedGrantedEnvs[j].Name
	})

	// Show message if no access granted
	if len(vep.publicKeys) == 0 && len(vep.grantedEnvs) == 0 {
		vep.grantedAccess.AddItem("📝 No access granted yet", "Add public keys or environments to grant access", 0, nil)
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func containsPublicKey(keys []crypto.PublicKey, target crypto.PublicKey) bool {
	targetStr, err := target.String()
	if err != nil {
		return false
	}
	for _, key := range keys {
		keyStr, err := key.String()
		if err != nil {
			continue
		}
		if keyStr == targetStr {
			return true
		}
	}
	return false
}

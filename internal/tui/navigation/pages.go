package navigation

import (
	"slv.sh/slv/internal/core/vaults"
)

// ShowMainMenu displays the main menu
func (n *Navigation) ShowMainMenu(replace bool) {
	// Create fresh MainPage instance using factory
	mainPage := n.app.GetRouter().CreatePage(n.app, "main")
	n.StorePageInstance("main", mainPage) // Store page instance for refresh
	menu := mainPage.Create()
	n.addPage("main", menu)
	n.setCurrentPage("main", replace)

	// Restore navigation state after setCurrentPage
	mainPage.RestoreNavigationState()
}

// ShowVaults displays the vaults page
func (n *Navigation) ShowVaults(replace bool) {
	// Create fresh VaultPage instance using factory
	vaultPage := n.app.GetRouter().CreatePage(n.app, "vaults_browse", n.GetVaultDir())
	n.StorePageInstance("vaults", vaultPage) // Store page instance for refresh

	// // Check if we have saved state for this page
	// if n.HasPageState("vaults") {
	// 	// If we have saved state, restore it after creating the page
	// 	if vbp, ok := vaultPage.(*vault_browse.VaultBrowsePage); ok {
	// 		vbp.RestoreNavigationState()
	// 	}
	// }

	vaults := vaultPage.Create()
	n.addPage("vaults", vaults)
	n.setCurrentPage("vaults", replace)

	// Restore navigation state after setCurrentPage
	vaultPage.RestoreNavigationState()

	n.UpdateStatus()
}

// ShowProfiles displays the profiles page
func (n *Navigation) ShowProfiles(replace bool) {
	// Create fresh ProfilesPage instance using factory
	profilesPage := n.app.GetRouter().CreatePage(n.app, "profiles")
	n.StorePageInstance("profiles", profilesPage) // Store page instance for refresh
	page := profilesPage.Create()
	n.addPage("profiles", page)
	n.setCurrentPage("profiles", replace)

	// Restore navigation state after setCurrentPage
	profilesPage.RestoreNavigationState()

	n.UpdateStatus()
}

// ShowEnvironments displays the environments page
func (n *Navigation) ShowEnvironments(replace bool) {
	// Create fresh EnvironmentsPage instance using factory
	environmentsPage := n.app.GetRouter().CreatePage(n.app, "environments")
	n.StorePageInstance("environments", environmentsPage) // Store page instance for refresh
	page := environmentsPage.Create()
	n.addPage("environments", page)
	n.setCurrentPage("environments", replace)

	// Restore navigation state after setCurrentPage
	environmentsPage.RestoreNavigationState()

	n.UpdateStatus()
}

// ShowNewEnvironment displays the new environment creation page
func (n *Navigation) ShowNewEnvironment(replace bool) {
	// Create fresh EnvironmentNewPage instance using factory
	newEnvPage := n.app.GetRouter().CreatePage(n.app, "environments_new")
	n.StorePageInstance("environments_new", newEnvPage) // Store page instance for refresh
	page := newEnvPage.Create()
	n.addPage("environments_new", page)
	n.setCurrentPage("environments_new", replace)

	// Restore navigation state after setCurrentPage
	newEnvPage.RestoreNavigationState()

	n.UpdateStatus()
}

// ShowHelp displays the help page
func (n *Navigation) ShowHelp(replace bool) {
	// Create fresh HelpPage instance using factory
	helpPage := n.app.GetRouter().CreatePage(n.app, "help")
	n.StorePageInstance("help", helpPage) // Store page instance for refresh
	page := helpPage.Create()
	n.addPage("help", page)
	n.setCurrentPage("help", replace)

	// Restore navigation state after setCurrentPage
	helpPage.RestoreNavigationState()

	n.UpdateStatus()
}

// ShowVaultDetails shows a vault details page
func (n *Navigation) ShowVaultDetails(replace bool) {
	// Create fresh VaultViewPage instance using factory
	// Note: This method needs to be updated to accept vault and filepath parameters
	// For now, we'll create with nil vault and empty filepath
	vaultViewPage := n.app.GetRouter().CreatePage(n.app, "vaults_view", nil, "")
	n.StorePageInstance("vault-details", vaultViewPage) // Store page instance for refresh
	vaultDetailsPage := vaultViewPage.Create()
	n.addPage("vault-details", vaultDetailsPage)
	n.setCurrentPage("vault-details", replace)

	// Restore navigation state after setCurrentPage
	vaultViewPage.RestoreNavigationState()

	n.UpdateStatus()
}

// ShowNewVault shows the new vault creation page
func (n *Navigation) ShowNewVault(replace bool) {
	// Create fresh VaultNewPage instance using factory with current directory
	newVaultPage := n.app.GetRouter().CreatePage(n.app, "vaults_new", n.GetVaultDir())
	n.StorePageInstance("new-vault", newVaultPage) // Store page instance for refresh
	page := newVaultPage.Create()
	n.addPage("new-vault", page)
	n.setCurrentPage("new-vault", replace)

	// Restore navigation state after setCurrentPage
	newVaultPage.RestoreNavigationState()

	n.UpdateStatus()
}

// ShowVaultsWithDir shows the vaults page with a specific directory
func (n *Navigation) ShowVaultsWithDir(dir string, replace bool) {
	// Create fresh VaultPage instance using factory with specified directory
	vaultPage := n.app.GetRouter().CreatePage(n.app, "vaults_browse", dir)
	n.StorePageInstance("vaults", vaultPage) // Store page instance for refresh

	// // Check if we have saved state for this page
	// if n.HasPageState("vaults") {
	// 	// If we have saved state, restore it after creating the page
	// 	if vbp, ok := vaultPage.(*vault_browse.VaultBrowsePage); ok {
	// 		vbp.RestoreNavigationState()
	// 	}
	// }

	vaults := vaultPage.Create()
	n.addPage("vaults", vaults)
	n.setCurrentPage("vaults", replace)

	// Restore navigation state after setCurrentPage
	vaultPage.RestoreNavigationState()

	n.UpdateStatus()
}

// ShowVaultDetailsWithVault shows vault details page with specific vault and filepath
func (n *Navigation) ShowVaultDetailsWithVault(vault *vaults.Vault, filePath string, replace bool) {
	// Save current page state before navigating away
	n.saveCurrentPageState()

	// Create fresh VaultViewPage instance using factory with vault and filepath
	vaultViewPage := n.app.GetRouter().CreatePage(n.app, "vaults_view", vault, filePath)
	n.StorePageInstance("vault-details", vaultViewPage) // Store page instance for refresh
	vaultDetailsPage := vaultViewPage.Create()
	n.addPage("vault-details", vaultDetailsPage)
	n.setCurrentPage("vault-details", replace)

	// Restore navigation state after setCurrentPage (which resets focus)
	vaultViewPage.RestoreNavigationState()

	n.UpdateStatus()

}

// ShowNewVaultWithDir shows the new vault creation page with a specific directory
func (n *Navigation) ShowNewVaultWithDir(dir string, replace bool) {
	// Save current page state before navigating away
	n.saveCurrentPageState()

	// Create fresh VaultNewPage instance using factory with specified directory
	newVaultPage := n.app.GetRouter().CreatePage(n.app, "vaults_new", dir)
	n.StorePageInstance("new-vault", newVaultPage) // Store page instance for refresh
	page := newVaultPage.Create()
	n.addPage("new-vault", page)
	n.setCurrentPage("new-vault", replace)

	// Restore navigation state after setCurrentPage
	newVaultPage.RestoreNavigationState()

	n.UpdateStatus()
}

// ShowVaultEditWithVault shows vault edit page with specific vault and filepath
func (n *Navigation) ShowVaultEditWithVault(vault *vaults.Vault, filePath string, replace bool) {
	// Save current page state before navigating away
	n.saveCurrentPageState()

	// Create fresh VaultEditPage instance using factory with vault and filepath
	vaultEditPage := n.app.GetRouter().CreatePage(n.app, "vaults_edit", vault, filePath)
	n.StorePageInstance("vault-edit", vaultEditPage) // Store page instance for refresh
	vaultEditDetailsPage := vaultEditPage.Create()
	n.addPage("vault-edit", vaultEditDetailsPage)
	n.setCurrentPage("vault-edit", replace)

	// Restore navigation state after setCurrentPage
	vaultEditPage.RestoreNavigationState()

	n.UpdateStatus()
}

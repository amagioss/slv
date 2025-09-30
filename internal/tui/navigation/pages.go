package navigation

import (
	"slv.sh/slv/internal/core/vaults"
)

// ShowMainMenu displays the main menu
func (n *Navigation) ShowMainMenu(replace bool) {
	// Create fresh MainPage instance using factory
	mainPage := n.app.GetRouter().CreatePage(n.app, "main")
	menu := mainPage.Create()
	n.addPage("main", menu)
	n.setCurrentPage("main", replace)
}

// ShowVaults displays the vaults page
func (n *Navigation) ShowVaults(replace bool) {
	// Create fresh VaultPage instance using factory
	vaultPage := n.app.GetRouter().CreatePage(n.app, "vaults_browse", n.GetVaultDir())
	vaults := vaultPage.Create()
	n.addPage("vaults", vaults)
	n.setCurrentPage("vaults", replace)
	n.UpdateStatus()
}

// ShowProfiles displays the profiles page
func (n *Navigation) ShowProfiles(replace bool) {
	// Create fresh ProfilesPage instance using factory
	profilesPage := n.app.GetRouter().CreatePage(n.app, "profiles")
	page := profilesPage.Create()
	n.addPage("profiles", page)
	n.setCurrentPage("profiles", replace)
	n.UpdateStatus()
}

// ShowEnvironments displays the environments page
func (n *Navigation) ShowEnvironments(replace bool) {
	// Create fresh EnvironmentsPage instance using factory
	environmentsPage := n.app.GetRouter().CreatePage(n.app, "environments")
	page := environmentsPage.Create()
	n.addPage("environments", page)
	n.setCurrentPage("environments", replace)
	n.UpdateStatus()
}

// ShowHelp displays the help page
func (n *Navigation) ShowHelp(replace bool) {
	// Create fresh HelpPage instance using factory
	helpPage := n.app.GetRouter().CreatePage(n.app, "help")
	page := helpPage.Create()
	n.addPage("help", page)
	n.setCurrentPage("help", replace)
	n.UpdateStatus()
}

// ShowVaultDetails shows a vault details page
func (n *Navigation) ShowVaultDetails(replace bool) {
	// Create fresh VaultViewPage instance using factory
	// Note: This method needs to be updated to accept vault and filepath parameters
	// For now, we'll create with nil vault and empty filepath
	vaultViewPage := n.app.GetRouter().CreatePage(n.app, "vaults_view", nil, "")
	vaultDetailsPage := vaultViewPage.Create()
	n.addPage("vault-details", vaultDetailsPage)
	n.setCurrentPage("vault-details", replace)
	n.UpdateStatus()
}

// ShowNewVault shows the new vault creation page
func (n *Navigation) ShowNewVault(replace bool) {
	// Create fresh VaultNewPage instance using factory with current directory
	newVaultPage := n.app.GetRouter().CreatePage(n.app, "vaults_new", n.GetVaultDir())
	page := newVaultPage.Create()
	n.addPage("new-vault", page)
	n.setCurrentPage("new-vault", replace)
	n.UpdateStatus()
}

// ShowVaultsWithDir shows the vaults page with a specific directory
func (n *Navigation) ShowVaultsWithDir(dir string, replace bool) {
	// Create fresh VaultPage instance using factory with specified directory
	vaultPage := n.app.GetRouter().CreatePage(n.app, "vaults_browse", dir)
	vaults := vaultPage.Create()
	n.addPage("vaults", vaults)
	n.setCurrentPage("vaults", replace)
	n.UpdateStatus()
}

// ShowVaultDetailsWithVault shows vault details page with specific vault and filepath
func (n *Navigation) ShowVaultDetailsWithVault(vault *vaults.Vault, filePath string, replace bool) {
	// Create fresh VaultViewPage instance using factory with vault and filepath
	vaultViewPage := n.app.GetRouter().CreatePage(n.app, "vaults_view", vault, filePath)
	vaultDetailsPage := vaultViewPage.Create()
	n.addPage("vault-details", vaultDetailsPage)
	n.setCurrentPage("vault-details", replace)
	n.UpdateStatus()
}

// ShowNewVaultWithDir shows the new vault creation page with a specific directory
func (n *Navigation) ShowNewVaultWithDir(dir string, replace bool) {
	// Create fresh VaultNewPage instance using factory with specified directory
	newVaultPage := n.app.GetRouter().CreatePage(n.app, "vaults_new", dir)
	page := newVaultPage.Create()
	n.addPage("new-vault", page)
	n.setCurrentPage("new-vault", replace)
	n.UpdateStatus()
}

package navigation

import (
	"slv.sh/slv/internal/tui/pages/environments"
	"slv.sh/slv/internal/tui/pages/help"
	"slv.sh/slv/internal/tui/pages/mainpage"
	"slv.sh/slv/internal/tui/pages/profiles"
	"slv.sh/slv/internal/tui/pages/vault_browse"
	"slv.sh/slv/internal/tui/pages/vault_new"
	"slv.sh/slv/internal/tui/pages/vault_view"
)

// ShowMainMenu displays the main menu
func (n *Navigation) ShowMainMenu() {
	// Create MainPage using the mainpage package
	mainPage := mainpage.NewMainPage(n.app)
	menu := mainPage.Create()
	n.addPage("main", menu)
	n.setCurrentPage("main")
}

// ShowVaults displays the vaults page
func (n *Navigation) ShowVaults() {
	// Create VaultPage using the pages package
	vaultPage := n.app.GetRouter().GetRegisteredPage("vaults_browse").(*vault_browse.VaultBrowsePage)
	vaults := vaultPage.Create()
	n.addPage("vaults", vaults)
	n.setCurrentPage("vaults")
	n.UpdateStatus()
}

// ShowVaultsReplace replaces the current vault page (for directory navigation)
func (n *Navigation) ShowVaultsReplace() {
	// Create VaultPage using the pages package
	vaultPage := n.app.GetRouter().GetRegisteredPage("vaults_browse").(*vault_browse.VaultBrowsePage)
	vaults := vaultPage.Create()
	n.addPage("vaults", vaults)

	// Replace current page without adding to stack
	if n.app.GetRouter().GetCurrentPage() != "" {
		n.app.GetRouter().SetCurrentPage("vaults")
		n.app.GetComponents().GetMainContentPages().SwitchToPage("vaults")
		n.UpdateStatus()
	} else {
		n.setCurrentPage("vaults")
	}
}

// ShowProfiles displays the profiles page
func (n *Navigation) ShowProfiles() {
	// Create ProfilesPage using the profiles package
	profilesPage := profiles.NewProfilesPage(n.app)
	page := profilesPage.Create()
	n.addPage("profiles", page)
	n.setCurrentPage("profiles")
	n.UpdateStatus()
}

// ShowEnvironments displays the environments page
func (n *Navigation) ShowEnvironments() {
	// Create EnvironmentsPage using the environments package
	environmentsPage := environments.NewEnvironmentsPage(n.app)
	page := environmentsPage.Create()
	n.addPage("environments", page)
	n.setCurrentPage("environments")
	n.UpdateStatus()
}

// ShowHelp displays the help page
func (n *Navigation) ShowHelp() {
	// Create HelpPage using the help package
	helpPage := help.NewHelpPage(n.app)
	page := helpPage.Create()
	n.addPage("help", page)
	n.setCurrentPage("help")
	n.UpdateStatus()
}

// ShowVaultDetails shows a vault details page
func (n *Navigation) ShowVaultDetails() {
	vaultDetailsPage := n.app.GetRouter().GetRegisteredPage("vaults_view").(*vault_view.VaultViewPage).Create()
	n.addPage("vault-details", vaultDetailsPage)
	n.setCurrentPage("vault-details")
	n.UpdateStatus()
}

// ShowNewVault shows the new vault creation page
func (n *Navigation) ShowNewVault() {
	// Create NewVaultPage using the pages package
	newVaultPage := n.app.GetRouter().GetRegisteredPage("vaults_new").(*vault_new.VaultNewPage)
	page := newVaultPage.Create()
	n.addPage("new-vault", page)
	n.setCurrentPage("new-vault")
	n.UpdateStatus()
}

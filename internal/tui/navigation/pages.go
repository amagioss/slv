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
func (n *Navigation) ShowMainMenu(replace bool) {
	// Create MainPage using the mainpage package
	mainPage := mainpage.NewMainPage(n.app)
	menu := mainPage.Create()
	n.addPage("main", menu)
	n.setCurrentPage("main", replace)
}

// ShowVaults displays the vaults page
func (n *Navigation) ShowVaults(replace bool) {
	// Create VaultPage using the pages package
	vaultPage := n.app.GetRouter().GetRegisteredPage("vaults_browse").(*vault_browse.VaultBrowsePage)
	vaults := vaultPage.Create()
	n.addPage("vaults", vaults)
	n.setCurrentPage("vaults", replace)
	n.UpdateStatus()
}

// ShowProfiles displays the profiles page
func (n *Navigation) ShowProfiles(replace bool) {
	// Create ProfilesPage using the profiles package
	profilesPage := profiles.NewProfilesPage(n.app)
	page := profilesPage.Create()
	n.addPage("profiles", page)
	n.setCurrentPage("profiles", replace)
	n.UpdateStatus()
}

// ShowEnvironments displays the environments page
func (n *Navigation) ShowEnvironments(replace bool) {
	// Create EnvironmentsPage using the environments package
	environmentsPage := environments.NewEnvironmentsPage(n.app)
	page := environmentsPage.Create()
	n.addPage("environments", page)
	n.setCurrentPage("environments", replace)
	n.UpdateStatus()
}

// ShowHelp displays the help page
func (n *Navigation) ShowHelp(replace bool) {
	// Create HelpPage using the help package
	helpPage := help.NewHelpPage(n.app)
	page := helpPage.Create()
	n.addPage("help", page)
	n.setCurrentPage("help", replace)
	n.UpdateStatus()
}

// ShowVaultDetails shows a vault details page
func (n *Navigation) ShowVaultDetails(replace bool) {
	vaultDetailsPage := n.app.GetRouter().GetRegisteredPage("vaults_view").(*vault_view.VaultViewPage).Create()
	n.addPage("vault-details", vaultDetailsPage)
	n.setCurrentPage("vault-details", replace)
	n.UpdateStatus()
}

// ShowNewVault shows the new vault creation page
func (n *Navigation) ShowNewVault(replace bool) {
	// Create NewVaultPage using the pages package
	newVaultPage := n.app.GetRouter().GetRegisteredPage("vaults_new").(*vault_new.VaultNewPage)
	page := newVaultPage.Create()
	n.addPage("new-vault", page)
	n.setCurrentPage("new-vault", replace)
	n.UpdateStatus()
}

package vault_edit

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/tui/interfaces"
	"slv.sh/slv/internal/tui/pages"
)

// VaultEditPage handles the vault editing functionality
type VaultEditPage struct {
	pages.BasePage
	vault        *vaults.Vault
	filePath     string
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

	// K8s environment
	k8sEnv *environments.Environment

	// Vault access state
	isVaultUnlocked bool // Whether the current session has access to the vault

	navigation  *FormNavigation
	currentPage tview.Primitive // Store reference to current page for modal navigation
}

// NewVaultEditPage creates a new VaultEditPage instance
func NewVaultEditPage(tui interfaces.TUIInterface, vault *vaults.Vault, filePath string) *VaultEditPage {

	vep := &VaultEditPage{
		BasePage:        *pages.NewBasePage(tui, "Edit Vault"),
		vault:           vault,
		filePath:        filePath,
		currentDir:      "", // Will be set from filePath
		publicKeys:      []string{},
		grantedEnvs:     []*environments.Environment{},
		searchEnvMap:    make(map[string]*environments.Environment),
		isVaultUnlocked: false, // Will be determined by checking access
	}

	accessors, err := vault.ListAccessors()
	if err != nil {
		vep.showError(fmt.Sprintf("Error listing accessors: %v", err))
		return nil
	}
	profile, err := profiles.GetActiveProfile()
	if err != nil {
		vep.showError(fmt.Sprintf("Error getting active profile: %v", err))
		return nil
	}
	for _, accessor := range accessors {
		accessorPubKey, err := accessor.String()
		if err != nil {
			continue
		}
		accessorEnv, err := profile.GetEnv(accessorPubKey)
		if err != nil {
			continue
		}
		if accessorEnv == nil {
			vep.grantedEnvs = append(vep.grantedEnvs,
				&environments.Environment{PublicKey: accessorPubKey,
					Name:    "Unknown",
					Email:   "Unknown",
					EnvType: environments.SERVICE,
					Tags:    []string{},
				})
		} else {
			vep.grantedEnvs = append(vep.grantedEnvs, accessorEnv)
		}
	}

	// Extract directory from filePath
	if filePath != "" {
		// Get directory from file path
		vep.currentDir = filePath[:len(filePath)-len(vault.Name+".slv.yaml")]
	}

	// Check if current session has access to the vault
	vep.checkVaultAccess()

	vep.currentPage = vep.createMainSection()
	vep.updateGrantedAccessList()
	vep.navigation = (&FormNavigation{}).NewFormNavigation(vep)
	vep.navigation.SetupNavigation()
	return vep
}

// checkVaultAccess determines if the current session has access to the vault
func (vep *VaultEditPage) checkVaultAccess() {
	// Get the current session environment
	currentEnv := environments.GetSelf()
	if currentEnv == nil {
		vep.isVaultUnlocked = false
		return
	}

	// Check if the current environment's public key is in the vault's accessors
	currentPubKey := currentEnv.PublicKey

	// Check if current environment has access to the vault
	accessors, err := vep.vault.ListAccessors()
	if err != nil {
		vep.isVaultUnlocked = false
		return
	}

	for _, accessor := range accessors {
		accessorPubKey, err := accessor.String()
		if err != nil {
			continue
		}
		if accessorPubKey == currentPubKey {
			vep.isVaultUnlocked = true
			return
		}
	}

	vep.isVaultUnlocked = false
}

// IsVaultUnlocked returns whether the vault is unlocked (current session has access)
func (vep *VaultEditPage) IsVaultUnlocked() bool {
	return vep.isVaultUnlocked
}

// Create implements the Page interface
func (vep *VaultEditPage) Create() tview.Primitive {
	flex := vep.currentPage
	// Update status bar with help text based on vault access state
	if vep.IsVaultUnlocked() {
		vep.GetTUI().UpdateStatusBar("Tab: Navigate forms | Esc: Back | Ctrl+C: Quit | Vault Unlocked - Full Access")
	} else {
		vep.GetTUI().UpdateStatusBar("Tab: Navigate forms | Esc: Back | Ctrl+C: Quit | Vault Locked - Limited Access")
	}

	vep.SetTitle("Edit Vault at " + vep.filePath)
	return vep.CreateLayout(flex)
}

// Refresh implements the Page interface
func (vep *VaultEditPage) Refresh() {
	// Recreate page using navigation system
	vep.GetTUI().GetNavigation().ShowVaultEditWithVault(vep.vault, vep.filePath, true)

	// Update help text for the current focus
	if vep.navigation != nil {
		vep.navigation.updateHelpText()
	}
}

// HandleInput implements the Page interface
func (vep *VaultEditPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// TODO: Implement vault edit page input handling
	return event
}

// GetTitle implements the Page interface
func (vep *VaultEditPage) GetTitle() string {
	return vep.BasePage.GetTitle()
}

// GetVault returns the vault
func (vep *VaultEditPage) GetVault() *vaults.Vault {
	return vep.vault
}

// GetFilePath returns the file path
func (vep *VaultEditPage) GetFilePath() string {
	return vep.filePath
}

// GetCurrentDir returns the current directory
func (vep *VaultEditPage) GetCurrentDir() string {
	return vep.currentDir
}

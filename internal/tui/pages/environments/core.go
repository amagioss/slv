package environments

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rivo/tview"

	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/tui/theme"
)

func (ep *EnvironmentsPage) getSessionEnvironment() *environments.Environment {
	currentSession, err := session.GetSession()
	if err != nil {
		return nil
	}
	env, err := currentSession.Env()
	if err != nil {
		return nil
	}
	return env
}

func (ep *EnvironmentsPage) getSelfEnvironment() *environments.Environment {
	return environments.GetSelf()
}

// loadAllEnvironments loads and displays all environments from the active profile
func (ep *EnvironmentsPage) loadAllEnvironments() {
	// Clear EDS table if it was showing
	ep.browseEnvsEDSTable = nil

	// Show list view
	ep.browseEnvsSection.Clear()
	ep.browseEnvsSection.SetTitle("Browse Environments").SetTitleAlign(tview.AlignLeft)
	ep.browseEnvsSection.AddItem(ep.browseEnvsSearch, 3, 0, true) // Fixed size: 3 rows
	ep.browseEnvsSection.AddItem(ep.browseEnvsList, 0, 1, true)   // Flexible: takes remaining space

	ep.browseEnvsList.Clear()
	ep.searchEnvMap = make(map[string]*environments.Environment) // Clear previous search results

	profile, err := profiles.GetActiveProfile()
	if err != nil {
		ep.ShowError(fmt.Sprintf("Error getting active profile: %v", err))
		ep.browseEnvsList.AddItem("❌ Error", "Failed to get active profile", 0, nil)
		return
	}

	// Get all environments
	envs, err := profile.ListEnvs()
	if err != nil {
		ep.ShowError(fmt.Sprintf("Error listing environments: %v", err))
		ep.browseEnvsList.AddItem("❌ Error", "Failed to list environments", 0, nil)
		return
	}

	// Sort environments by name for consistent display order
	sort.Slice(envs, func(i, j int) bool {
		return envs[i].Name < envs[j].Name
	})

	// Display sorted environments
	for _, env := range envs {
		// Store environment in map for modal display
		ep.searchEnvMap[env.Name] = env
		// Format with proper spacing
		mainText := fmt.Sprintf("🔍 %s", env.Name)
		secondaryText := fmt.Sprintf("Type: %s | Email: %s", string(env.EnvType), env.Email)
		ep.browseEnvsList.AddItem(mainText, secondaryText, 0, nil)
	}

	// If no environments found
	if len(envs) == 0 {
		ep.browseEnvsList.AddItem("No environments found", "No environments in active profile", 0, nil)
	}

	// Update navigation to use list view (only if navigation is initialized)
	if ep.navigation != nil {
		ep.navigation.updateFocusGroupForSearch()
	}
}

// searchEnvironments searches for environments based on query string
func (ep *EnvironmentsPage) searchEnvironments(query string) {
	ep.currentQuery = query // Store the current query for refreshing

	// Check if query is an EDS (Environment Definition String)
	// Only process if it starts with SLV_EDS_ and has at least one character after the prefix
	upperQuery := strings.ToUpper(query)
	if strings.HasPrefix(upperQuery, "SLV_EDS_") && len(query) > len("SLV_EDS_") {
		ep.handleEDSSearch(query)
		return
	}

	// Clear EDS table if it was showing
	ep.browseEnvsEDSTable = nil

	// Show list view
	ep.browseEnvsSection.Clear()
	ep.browseEnvsSection.SetTitle("Browse Environments").SetTitleAlign(tview.AlignLeft)
	ep.browseEnvsSection.AddItem(ep.browseEnvsSearch, 3, 0, true) // Fixed size: 3 rows
	ep.browseEnvsSection.AddItem(ep.browseEnvsList, 0, 1, true)   // Flexible: takes remaining space

	ep.browseEnvsList.Clear()
	ep.searchEnvMap = make(map[string]*environments.Environment) // Clear previous search results

	profile, err := profiles.GetActiveProfile()
	if err != nil {
		ep.ShowError(fmt.Sprintf("Error getting active profile: %v", err))
		ep.browseEnvsList.AddItem("❌ Error", "Failed to get active profile", 0, nil)
		return
	}

	var matchingEnvs []*environments.Environment

	if strings.HasPrefix(strings.ToUpper(query), "SLV_EPK") {
		// Search by public key
		envs, err := profile.ListEnvs()
		if err != nil {
			ep.ShowError(fmt.Sprintf("Error listing environments: %v", err))
			ep.browseEnvsList.AddItem("❌ Error", "Failed to list environments", 0, nil)
			return
		}
		for _, env := range envs {
			if env.PublicKey == query {
				matchingEnvs = append(matchingEnvs, env)
			}
		}
		if len(matchingEnvs) == 0 {
			ep.browseEnvsList.AddItem("❌ Environment not found", fmt.Sprintf("No environment found with public key: %s", query), 0, nil)
			return
		}
	} else {
		// Normal string search
		envs, err := profile.SearchEnvs([]string{query})
		if err != nil {
			ep.ShowError(fmt.Sprintf("Error searching environments: %v", err))
			ep.browseEnvsList.AddItem("❌ Error", "Failed to search environments", 0, nil)
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
		// Store environment in map for modal display
		ep.searchEnvMap[env.Name] = env
		// Format with proper spacing
		mainText := fmt.Sprintf("🔍 %s", env.Name)
		secondaryText := fmt.Sprintf("Type: %s | Email: %s", string(env.EnvType), env.Email)
		ep.browseEnvsList.AddItem(mainText, secondaryText, 0, nil)
	}

	// If no results found
	if len(matchingEnvs) == 0 {
		ep.browseEnvsList.AddItem("No environments found", fmt.Sprintf("No environments match: %s", query), 0, nil)
	}

	// Update navigation to use list view
	ep.navigation.updateFocusGroupForSearch()
}

// handleEDSSearch parses and displays an Environment Definition String
func (ep *EnvironmentsPage) handleEDSSearch(edsQuery string) {
	// Validate EDS format before parsing to avoid panics
	// EDS should be in format: SLV_EDS_<data>
	sliced := strings.Split(edsQuery, "_")
	if len(sliced) < 3 {
		// Incomplete EDS - show error in list
		ep.browseEnvsSection.Clear()
		ep.browseEnvsSection.SetTitle("Browse Environments").SetTitleAlign(tview.AlignLeft)
		ep.browseEnvsSection.AddItem(ep.browseEnvsSearch, 3, 0, true)
		ep.browseEnvsSection.AddItem(ep.browseEnvsList, 0, 1, true)
		ep.browseEnvsList.Clear()
		ep.browseEnvsList.AddItem("❌ Incomplete EDS", "EDS must be in format: SLV_EDS_<data>", 0, nil)
		ep.navigation.updateFocusGroupForSearch()
		return
	}

	// Check that the data part (sliced[2]) is not empty
	// This prevents panics in Decompress when processing empty strings
	if len(sliced) < 3 || sliced[2] == "" {
		// Empty data part - show error in list
		ep.browseEnvsSection.Clear()
		ep.browseEnvsSection.SetTitle("Browse Environments").SetTitleAlign(tview.AlignLeft)
		ep.browseEnvsSection.AddItem(ep.browseEnvsSearch, 3, 0, true)
		ep.browseEnvsSection.AddItem(ep.browseEnvsList, 0, 1, true)
		ep.browseEnvsList.Clear()
		ep.browseEnvsList.AddItem("❌ Incomplete EDS", "EDS data part is empty. EDS must be in format: SLV_EDS_<data>", 0, nil)
		ep.navigation.updateFocusGroupForSearch()
		return
	}

	// Parse the EDS with panic recovery
	var env *environments.Environment
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic during EDS parsing: %v", r)
			}
		}()
		env, err = environments.FromDefStr(edsQuery)
	}()

	if err != nil {
		// Invalid EDS - show error in list
		ep.browseEnvsSection.Clear()
		ep.browseEnvsSection.SetTitle("Browse Environments").SetTitleAlign(tview.AlignLeft)
		ep.browseEnvsSection.AddItem(ep.browseEnvsSearch, 3, 0, true)
		ep.browseEnvsSection.AddItem(ep.browseEnvsList, 0, 1, true)
		ep.browseEnvsList.Clear()
		ep.browseEnvsList.AddItem("❌ Invalid EDS", fmt.Sprintf("Error parsing EDS: %v", err), 0, nil)
		ep.navigation.updateFocusGroupForSearch()
		return
	}

	// Valid EDS - show table
	colors := theme.GetCurrentPalette()
	ep.browseEnvsEDSTable = ep.createEnvironmentTableForDetails(env, colors)

	// Update section to show table instead of list
	ep.browseEnvsSection.Clear()
	ep.browseEnvsSection.SetTitle("Environment from EDS").SetTitleAlign(tview.AlignLeft)
	ep.browseEnvsSection.AddItem(ep.browseEnvsSearch, 3, 0, true)   // Fixed size: 3 rows
	ep.browseEnvsSection.AddItem(ep.browseEnvsEDSTable, 0, 1, true) // Flexible: takes remaining space

	// Update navigation to include EDS table
	ep.navigation.updateFocusGroupForEDS()
}

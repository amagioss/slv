package environments

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/tui/theme"
)

// createMainSection creates the main content section with the specified layout using Flex
func (ep *EnvironmentsPage) createMainSection() *tview.Flex {
	colors := theme.GetCurrentPalette()

	// Create Session Environment table
	sessionEnv := ep.getSessionEnvironment()
	ep.sessionEnvTable = ep.createEnvironmentTable("Session Environment", sessionEnv, colors)

	// Create Self Environment table
	selfEnv := ep.getSelfEnvironment()
	ep.selfEnvTable = ep.createEnvironmentTable("Self Environment", selfEnv, colors)

	// Create Browse Environments section with search and results
	ep.browseEnvsSearch = ep.createSearchBox(colors)
	ep.browseEnvsList = ep.createResultsList(colors)

	// Create the browse section container (starts with search view)
	ep.browseEnvsSection = tview.NewFlex().
		SetDirection(tview.FlexRow)
	ep.browseEnvsSection.SetBorder(true)
	ep.browseEnvsSection.SetTitle("Browse Environments").SetTitleAlign(tview.AlignLeft)

	// Initially show search view
	ep.showSearchView()

	// Create first row - horizontal flex with 2 equal columns
	// Make tables focusable so they can receive input for navigation and copying
	firstRow := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(ep.sessionEnvTable, 0, 1, sessionEnv != nil). // Focusable if environment exists
		AddItem(ep.selfEnvTable, 0, 1, selfEnv != nil)        // Focusable if environment exists

	// Create main content - vertical flex
	// First row: 30% of space (2 columns)
	// Second row: 70% of space (1 column)
	ep.mainContent = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(firstRow, 8, 0, true).            // 30% of space (weight 3), focusable (tables inside) - increased to 8 for EDS field
		AddItem(ep.browseEnvsSection, 0, 1, true) // 70% of space (weight 7), focusable
	return ep.mainContent
}

// createEnvironmentTable creates a borderless table to display environment information
func (ep *EnvironmentsPage) createEnvironmentTable(title string, env *environments.Environment, colors theme.ColorPalette) *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetTitle(title).SetTitleAlign(tview.AlignLeft)
	table.SetBorderColor(colors.Border)

	// Start with table not selectable - will be enabled when focused
	table.SetSelectable(false, false)
	table.SetSelectedStyle(tcell.StyleDefault.Background(colors.Selection).Foreground(colors.TextPrimary))

	if env == nil {
		// Show "No environment configured" message
		table.SetCell(0, 0, tview.NewTableCell("No environment configured").
			SetTextColor(colors.TextMuted).
			SetAlign(tview.AlignCenter))
		table.SetCell(1, 0, tview.NewTableCell("This environment is not available.").
			SetTextColor(colors.TextMuted).
			SetAlign(tview.AlignCenter))
		return table
	}

	row := 0

	// Name field
	table.SetCell(row, 0, tview.NewTableCell("Name:").
		SetTextColor(colors.TableLabel).
		SetAlign(tview.AlignLeft).
		SetMaxWidth(12))
	nameValue := "Unnamed"
	if env.Name != "" {
		nameValue = env.Name
	}
	table.SetCell(row, 1, tview.NewTableCell(nameValue).
		SetTextColor(colors.TableValue).
		SetAlign(tview.AlignLeft).
		SetExpansion(1))
	row++

	// Type field
	if env.EnvType != "" {
		table.SetCell(row, 0, tview.NewTableCell("Type:").
			SetTextColor(colors.TableLabel).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(12))
		table.SetCell(row, 1, tview.NewTableCell(string(env.EnvType)).
			SetTextColor(colors.TableValue).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		row++
	}

	// Email field
	if env.Email != "" {
		table.SetCell(row, 0, tview.NewTableCell("Email:").
			SetTextColor(colors.TableLabel).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(12))
		table.SetCell(row, 1, tview.NewTableCell(env.Email).
			SetTextColor(colors.TableValue).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		row++
	}

	// Tags field
	if len(env.Tags) > 0 {
		table.SetCell(row, 0, tview.NewTableCell("Tags:").
			SetTextColor(colors.TableLabel).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(12))
		table.SetCell(row, 1, tview.NewTableCell(strings.Join(env.Tags, ", ")).
			SetTextColor(colors.TableValue).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		row++
	}

	// Public Key field
	if env.PublicKey != "" {
		table.SetCell(row, 0, tview.NewTableCell("Public Key:").
			SetTextColor(colors.TableLabel).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(12))
		table.SetCell(row, 1, tview.NewTableCell(env.PublicKey).
			SetTextColor(colors.TableValue).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		row++
	}

	// EDS (Environment Definition String) field
	if defStr, err := env.ToDefStr(false); err == nil && defStr != "" {
		table.SetCell(row, 0, tview.NewTableCell("EDS:").
			SetTextColor(colors.TableLabel).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(12))
		table.SetCell(row, 1, tview.NewTableCell(defStr).
			SetTextColor(colors.TableValue).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		row++
	}

	return table
}

// createSearchBox creates a search input field for browsing environments
func (ep *EnvironmentsPage) createSearchBox(colors theme.ColorPalette) *tview.InputField {
	searchBox := tview.NewInputField()
	searchBox.SetLabel("Search: ").
		SetLabelColor(colors.FormLabel).
		SetFieldTextColor(colors.FormFieldText).
		SetBorder(true).
		SetBorderColor(colors.Border)

	// Set placeholder text
	searchBox.SetPlaceholder("Type to search environments... (Or enter EDS to parse)").
		SetPlaceholderTextColor(colors.TextMuted)

	// Set up callback to trigger search on text change
	searchBox.SetChangedFunc(func(text string) {
		if text != "" {
			ep.searchEnvironments(text)
		} else {
			// When input is empty, show all environments
			ep.loadAllEnvironments()
		}
	})

	return searchBox
}

// createResultsList creates a list to display search results
func (ep *EnvironmentsPage) createResultsList(colors theme.ColorPalette) *tview.List {
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Results").SetTitleAlign(tview.AlignLeft)
	list.SetBorderColor(colors.Border)
	list.SetSelectedTextColor(colors.ListSelectedText).
		SetSelectedBackgroundColor(colors.ListSelectedBg).
		SetSecondaryTextColor(colors.ListSecondaryText).
		SetMainTextColor(colors.ListMainText).
		SetWrapAround(false) // Disable looping behavior

	// Add placeholder item
	list.AddItem("No environments found", "Use the search box above to find environments", 0, nil)

	return list
}

// showSearchView shows the search box and results list
func (ep *EnvironmentsPage) showSearchView() {
	// Hide edit form if it's open
	ep.editForm = nil
	ep.editingField = ""

	ep.browseEnvsSection.Clear()
	ep.browseEnvsSection.SetTitle("Browse Environments").SetTitleAlign(tview.AlignLeft)
	ep.browseEnvsSection.AddItem(ep.browseEnvsSearch, 3, 0, true) // Fixed size: 3 rows
	ep.browseEnvsSection.AddItem(ep.browseEnvsList, 0, 1, true)   // Flexible: takes remaining space

	// Load all environments by default
	ep.loadAllEnvironments()
}

// showDetailsView shows the environment details table
func (ep *EnvironmentsPage) showDetailsView(env *environments.Environment) {
	colors := theme.GetCurrentPalette()

	// Create or update the details table
	ep.browseEnvsDetails = ep.createEnvironmentTableForDetails(env, colors)

	ep.browseEnvsSection.Clear()
	envName := "Unnamed"
	if env.Name != "" {
		envName = env.Name
	}
	ep.browseEnvsSection.SetTitle(fmt.Sprintf("Environment: %s", envName)).SetTitleAlign(tview.AlignLeft)
	ep.browseEnvsSection.AddItem(ep.browseEnvsDetails, 0, 1, true) // Flexible: takes all space

	// Store current environment for details view
	ep.currentDetailsEnv = env
}

// createEnvironmentTableForDetails creates a table for the environment details view with copy functionality
func (ep *EnvironmentsPage) createEnvironmentTableForDetails(env *environments.Environment, colors theme.ColorPalette) *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetBorderColor(colors.Border)

	// Start with table not selectable - will be enabled when focused
	table.SetSelectable(false, false)
	table.SetSelectedStyle(tcell.StyleDefault.Background(colors.Selection).Foreground(colors.TextPrimary))

	if env == nil {
		table.SetCell(0, 0, tview.NewTableCell("No environment data").
			SetTextColor(colors.TextMuted).
			SetAlign(tview.AlignCenter))
		return table
	}

	row := 0

	// Name field
	table.SetCell(row, 0, tview.NewTableCell("Name:").
		SetTextColor(colors.TableLabel).
		SetAlign(tview.AlignLeft).
		SetMaxWidth(12))
	nameValue := "Unnamed"
	if env.Name != "" {
		nameValue = env.Name
	}
	table.SetCell(row, 1, tview.NewTableCell(nameValue).
		SetTextColor(colors.TableValue).
		SetAlign(tview.AlignLeft).
		SetExpansion(1))
	row++

	// Type field
	if env.EnvType != "" {
		table.SetCell(row, 0, tview.NewTableCell("Type:").
			SetTextColor(colors.TableLabel).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(12))
		table.SetCell(row, 1, tview.NewTableCell(string(env.EnvType)).
			SetTextColor(colors.TableValue).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		row++
	}

	// Email field
	if env.Email != "" {
		table.SetCell(row, 0, tview.NewTableCell("Email:").
			SetTextColor(colors.TableLabel).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(12))
		table.SetCell(row, 1, tview.NewTableCell(env.Email).
			SetTextColor(colors.TableValue).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		row++
	}

	// Tags field
	if len(env.Tags) > 0 {
		table.SetCell(row, 0, tview.NewTableCell("Tags:").
			SetTextColor(colors.TableLabel).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(12))
		table.SetCell(row, 1, tview.NewTableCell(strings.Join(env.Tags, ", ")).
			SetTextColor(colors.TableValue).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		row++
	}

	// Public Key field
	if env.PublicKey != "" {
		table.SetCell(row, 0, tview.NewTableCell("Public Key:").
			SetTextColor(colors.TableLabel).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(12))
		table.SetCell(row, 1, tview.NewTableCell(env.PublicKey).
			SetTextColor(colors.TableValue).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		row++
	}

	// EDS (Environment Definition String) field
	if defStr, err := env.ToDefStr(false); err == nil && defStr != "" {
		table.SetCell(row, 0, tview.NewTableCell("EDS:").
			SetTextColor(colors.TableLabel).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(12))
		table.SetCell(row, 1, tview.NewTableCell(defStr).
			SetTextColor(colors.TableValue).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		row++
	}

	return table
}

// refreshDetailsView refreshes the details view with updated environment data
func (ep *EnvironmentsPage) refreshDetailsView() {
	if ep.currentDetailsEnv == nil {
		return
	}

	colors := theme.GetCurrentPalette()
	ep.browseEnvsDetails = ep.createEnvironmentTableForDetails(ep.currentDetailsEnv, colors)

	// Re-setup input capture
	ep.navigation.setupDetailsTableInputCapture()

	// Update focus group to reference the new table
	ep.navigation.updateFocusGroupForDetails()

	// Update the section if not editing
	if ep.editForm == nil {
		ep.browseEnvsSection.Clear()
		envName := "Unnamed"
		if ep.currentDetailsEnv.Name != "" {
			envName = ep.currentDetailsEnv.Name
		}
		ep.browseEnvsSection.SetTitle(fmt.Sprintf("Environment: %s", envName)).SetTitleAlign(tview.AlignLeft)
		ep.browseEnvsSection.AddItem(ep.browseEnvsDetails, 0, 1, true)
	} else {
		// If editing, update the table in the section
		ep.browseEnvsSection.Clear()
		envName := "Unnamed"
		if ep.currentDetailsEnv.Name != "" {
			envName = ep.currentDetailsEnv.Name
		}
		ep.browseEnvsSection.SetTitle(fmt.Sprintf("Environment: %s", envName)).SetTitleAlign(tview.AlignLeft)
		ep.browseEnvsSection.AddItem(ep.browseEnvsDetails, 0, 1, false)
		ep.browseEnvsSection.AddItem(ep.editForm, 8, 0, true)
	}
}

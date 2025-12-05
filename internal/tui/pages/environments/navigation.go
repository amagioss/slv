package environments

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
	"slv.sh/slv/internal/core/environments"
)

type FormNavigation struct {
	ep           *EnvironmentsPage
	currentFocus int
	focusGroup   []tview.Primitive
	helpTexts    map[tview.Primitive]string
}

func (fn *FormNavigation) NewFormNavigation(ep *EnvironmentsPage) *FormNavigation {
	// Build focus group - only include tables if environment exists
	focusGroup := []tview.Primitive{}

	// Add Session Environment table if environment exists
	if ep.getSessionEnvironment() != nil {
		focusGroup = append(focusGroup, ep.sessionEnvTable)
	}

	// Add Self Environment table if environment exists
	if ep.getSelfEnvironment() != nil {
		focusGroup = append(focusGroup, ep.selfEnvTable)
	}

	// Always add browse environments search and list
	focusGroup = append(focusGroup, ep.browseEnvsSearch)
	focusGroup = append(focusGroup, ep.browseEnvsList)

	initialFocus := 0
	if len(focusGroup) == 0 {
		initialFocus = 0
	}

	return &FormNavigation{
		ep:           ep,
		currentFocus: initialFocus,
		focusGroup:   focusGroup,
		helpTexts:    make(map[tview.Primitive]string),
	}
}

func (fn *FormNavigation) resetSelectable() {
	// Reset all tables to not selectable - just check if tables exist, don't call environment getters
	if fn.ep.sessionEnvTable != nil {
		fn.ep.sessionEnvTable.SetSelectable(false, false)
	}
	if fn.ep.selfEnvTable != nil {
		fn.ep.selfEnvTable.SetSelectable(false, false)
	}
	// Reset details table to not selectable
	if fn.ep.browseEnvsDetails != nil {
		fn.ep.browseEnvsDetails.SetSelectable(false, false)
	}
	// Reset EDS table to not selectable
	if fn.ep.browseEnvsEDSTable != nil {
		fn.ep.browseEnvsEDSTable.SetSelectable(false, false)
	}
}

func (fn *FormNavigation) SetupNavigation() {
	// Set up help texts for each component
	fn.setupHelpTexts()

	// Set up input capture on mainContent to handle Tab navigation
	fn.ep.mainContent.SetInputCapture(fn.handleInputCapture)

	// Set up input capture for environment tables (c for copy field, d for copy EDS)
	if sessionEnv := fn.ep.getSessionEnvironment(); sessionEnv != nil {
		fn.setupEnvironmentTableInputCapture(fn.ep.sessionEnvTable, sessionEnv)
	}
	if selfEnv := fn.ep.getSelfEnvironment(); selfEnv != nil {
		fn.setupEnvironmentTableInputCapture(fn.ep.selfEnvTable, selfEnv)
	}

	// Set up input capture for browse environments list (Enter to show details)
	fn.setupBrowseEnvsListInputCapture()

	// Reset all tables to not selectable, then make focused one selectable
	fn.resetSelectable()
	if len(fn.focusGroup) > 0 {
		// Make the focused component selectable if it's a table
		if table, ok := fn.focusGroup[fn.currentFocus].(*tview.Table); ok {
			table.SetSelectable(true, false)
		}
		fn.ep.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
	}

	// Set initial help text
	fn.updateHelpText()
}

func (fn *FormNavigation) setupHelpTexts() {
	// Help text for Session Environment
	if fn.ep.getSessionEnvironment() != nil {
		fn.helpTexts[fn.ep.sessionEnvTable] = "Session Environment: ↑↓ Navigate | c: Copy field | Tab: Switch | Ctrl+N: New Environment"
	}

	// Help text for Self Environment
	if fn.ep.getSelfEnvironment() != nil {
		fn.helpTexts[fn.ep.selfEnvTable] = "Self Environment: ↑↓ Navigate | c: Copy field | Tab: Switch | Ctrl+N: New Environment"
	}

	// Help text for Browse Environments search
	fn.helpTexts[fn.ep.browseEnvsSearch] = "Search Environments: Type to search | Tab: Switch | Ctrl+N: New Environment"

	// Help text for Browse Environments list
	fn.helpTexts[fn.ep.browseEnvsList] = "Environment Results: ↑↓ Navigate | Enter: View details | Tab: Switch | Ctrl+N: New Environment"

	// Help text for Environment Details table (will be set when details view is shown)
	if fn.ep.browseEnvsDetails != nil {
		fn.helpTexts[fn.ep.browseEnvsDetails] = "Environment Details: ↑↓ Navigate | Enter: Edit (Name/Email/Tags) | c: Copy field | Backspace: Back | Tab: Switch | Ctrl+N: New Environment"
	}
}

func (fn *FormNavigation) updateHelpText() {
	if fn.currentFocus >= 0 && fn.currentFocus < len(fn.focusGroup) {
		currentComponent := fn.focusGroup[fn.currentFocus]
		if helpText, exists := fn.helpTexts[currentComponent]; exists {
			fn.ep.GetTUI().UpdateStatusBar(helpText)
		}
	}
}

// handleInputCapture handles Tab navigation for the main content
func (fn *FormNavigation) handleInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	// If form is active, let it handle Tab/Backtab/Backspace/Esc
	if fn.ep.editForm != nil {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyBacktab, tcell.KeyBackspace2, tcell.KeyBackspace, tcell.KeyEsc:
			// Let form handle these keys - don't interfere
			return event
		}
	}

	switch event.Key() {
	case tcell.KeyCtrlN:
		// Navigate to new environment page
		fn.ep.GetTUI().GetNavigation().ShowNewEnvironment(false)
		return nil
	case tcell.KeyTab:
		// Switch focus forward
		fn.ShiftFocusForward()
		return nil
	case tcell.KeyBacktab:
		// Switch focus backward
		fn.ShiftFocusBackward()
		return nil
	case tcell.KeyBackspace2, tcell.KeyBackspace:
		// If we're in details view, go back to search view instead of going to previous page
		if fn.ep.browseEnvsDetails != nil && fn.ep.currentDetailsEnv != nil {
			// Check if details table is currently in focus group (we're in details view)
			for _, component := range fn.focusGroup {
				if component == fn.ep.browseEnvsDetails {
					// We're in details view, go back to search view
					fn.ep.showSearchView()
					fn.updateFocusGroupForSearch()
					fn.ep.GetTUI().GetApplication().SetFocus(fn.ep.browseEnvsList)
					fn.updateHelpText()
					return nil // Consume the event
				}
			}
		}
		// Not in details view, let Backspace pass through for normal page navigation
		return event
	}

	return event
}

func (fn *FormNavigation) SetFocusAtIndex(index int) {
	if len(fn.focusGroup) == 0 || index < 0 || index >= len(fn.focusGroup) {
		return
	}
	// Reset all tables to not selectable
	fn.resetSelectable()

	fn.currentFocus = index

	// Make the focused component selectable if it's a table
	if table, ok := fn.focusGroup[fn.currentFocus].(*tview.Table); ok {
		table.SetSelectable(true, false)
	}

	fn.ep.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
	fn.updateHelpText()
}

// ShiftFocusForward moves focus to the next component
func (fn *FormNavigation) ShiftFocusForward() {
	if len(fn.focusGroup) == 0 {
		return
	}
	// Reset all tables to not selectable
	fn.resetSelectable()

	fn.currentFocus = (fn.currentFocus + 1) % len(fn.focusGroup)

	// Make the focused component selectable if it's a table
	if table, ok := fn.focusGroup[fn.currentFocus].(*tview.Table); ok {
		table.SetSelectable(true, false)
	}

	fn.ep.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
	fn.updateHelpText()
}

// ShiftFocusBackward moves focus to the previous component
func (fn *FormNavigation) ShiftFocusBackward() {
	if len(fn.focusGroup) == 0 {
		return
	}
	// Reset all tables to not selectable
	fn.resetSelectable()

	fn.currentFocus = (fn.currentFocus - 1 + len(fn.focusGroup)) % len(fn.focusGroup)

	// Make the focused component selectable if it's a table
	if table, ok := fn.focusGroup[fn.currentFocus].(*tview.Table); ok {
		table.SetSelectable(true, false)
	}

	fn.ep.GetTUI().GetApplication().SetFocus(fn.focusGroup[fn.currentFocus])
	fn.updateHelpText()
}

// setupEnvironmentTableInputCapture sets up input capture for copying field values and EDS
func (fn *FormNavigation) setupEnvironmentTableInputCapture(table *tview.Table, env *environments.Environment) {
	if env == nil {
		return // No environment, nothing to copy
	}

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyBacktab:
			// Allow Tab to pass through to parent for navigation between tables
			return event

		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight:
			// Allow arrow keys to pass through for table navigation
			return event

		case tcell.KeyRune:
			switch event.Rune() {
			case 'c', 'C':
				// Copy the selected field value to clipboard
				row, _ := table.GetSelection()
				if row >= 0 {
					// Always copy from column 1 (the value column)
					valueCell := table.GetCell(row, 1)
					if valueCell != nil {
						value := valueCell.Text
						if value != "" {
							clipboard.Write(clipboard.FmtText, []byte(value))
							// Get the field name from column 0
							labelCell := table.GetCell(row, 0)
							fieldName := "field"
							if labelCell != nil {
								fieldName = strings.TrimSuffix(labelCell.Text, ":")
							}
							fn.ep.UpdateStatus(fmt.Sprintf("Copied %s to clipboard", fieldName))
						} else {
							fn.ep.ShowError("No value to copy")
						}
					}
				}
				return nil

			}
		}

		return event
	})
}

// setupBrowseEnvsListInputCapture sets up input capture for the browse environments list
func (fn *FormNavigation) setupBrowseEnvsListInputCapture() {
	fn.ep.browseEnvsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event == nil {
			return event
		}

		switch event.Key() {
		case tcell.KeyEnter:
			// Show environment details view
			selected := fn.ep.browseEnvsList.GetCurrentItem()
			if selected >= 0 && selected < fn.ep.browseEnvsList.GetItemCount() {
				mainText, _ := fn.ep.browseEnvsList.GetItemText(selected)
				// Extract environment name from the formatted text
				if strings.HasPrefix(mainText, "🔍 ") {
					envName := strings.TrimPrefix(mainText, "🔍 ")
					// Find the environment in the search results map
					if env, exists := fn.ep.searchEnvMap[envName]; exists {
						fn.ep.showDetailsView(env)
						// Set up input capture for details table
						fn.setupDetailsTableInputCapture()
						// Set help text for details table
						fn.helpTexts[fn.ep.browseEnvsDetails] = "Environment Details: ↑↓ Navigate | c: Copy field | Backspace: Back | Tab: Switch | Ctrl+N: New Environment"
						// Update focus group to include details table
						fn.updateFocusGroupForDetails()
						// Set focus to details table
						fn.ep.GetTUI().GetApplication().SetFocus(fn.ep.browseEnvsDetails)
						fn.updateHelpText()
					} else {
						fn.ep.ShowError("Environment not found")
					}
				} else {
					// Not an environment item (error message, placeholder, etc.)
					fn.ep.ShowInfo("Please select an environment from the search results")
				}
			}
			return nil
		case tcell.KeyTab, tcell.KeyBacktab:
			// Allow Tab to pass through to parent for navigation
			return event
		}

		return event
	})
}

// setupDetailsTableInputCapture sets up input capture for the environment details table
func (fn *FormNavigation) setupDetailsTableInputCapture() {
	if fn.ep.browseEnvsDetails == nil || fn.ep.currentDetailsEnv == nil {
		return
	}

	fn.ep.browseEnvsDetails.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event == nil {
			return event
		}

		// If form is active, don't handle any keys - let form handle them
		if fn.ep.editForm != nil {
			return event
		}

		switch event.Key() {
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight:
			// Allow arrow keys to pass through for table navigation
			return event

		case tcell.KeyEnter:
			// Show edit form for editable fields (Name, Email, Tags)
			row, _ := fn.ep.browseEnvsDetails.GetSelection()
			if row >= 0 {
				labelCell := fn.ep.browseEnvsDetails.GetCell(row, 0)
				if labelCell != nil {
					fieldName := strings.TrimSuffix(labelCell.Text, ":")
					// Only allow editing Name, Email, Tags
					if fieldName == "Name" || fieldName == "Email" || fieldName == "Tags" {
						fn.ep.showEditForm(fieldName, row)
						return nil
					}
				}
			}
			return event

		case tcell.KeyEsc:
			// Go back to search view
			fn.ep.showSearchView()
			// Update focus group back to search view
			fn.updateFocusGroupForSearch()
			// Set focus back to results list
			fn.ep.GetTUI().GetApplication().SetFocus(fn.ep.browseEnvsList)
			fn.updateHelpText()
			return nil

		case tcell.KeyTab, tcell.KeyBacktab:
			// Allow Tab to pass through to parent for navigation
			return event

		case tcell.KeyRune:
			switch event.Rune() {
			case 'c', 'C':
				// Copy the selected field value to clipboard
				row, _ := fn.ep.browseEnvsDetails.GetSelection()
				if row >= 0 {
					// Always copy from column 1 (the value column)
					valueCell := fn.ep.browseEnvsDetails.GetCell(row, 1)
					if valueCell != nil {
						value := valueCell.Text
						if value != "" {
							clipboard.Write(clipboard.FmtText, []byte(value))
							// Get the field name from column 0
							labelCell := fn.ep.browseEnvsDetails.GetCell(row, 0)
							fieldName := "field"
							if labelCell != nil {
								fieldName = strings.TrimSuffix(labelCell.Text, ":")
							}
							fn.ep.UpdateStatus(fmt.Sprintf("Copied %s to clipboard", fieldName))
						} else {
							fn.ep.ShowError("No value to copy")
						}
					}
				}
				return nil
			}
		}

		return event
	})
}

// updateFocusGroupForDetails updates the focus group to include the details table
func (fn *FormNavigation) updateFocusGroupForDetails() {
	// Reset all tables to not selectable first
	fn.resetSelectable()

	// Rebuild focus group with details table instead of search/list
	focusGroup := []tview.Primitive{}

	// Add Session Environment table if environment exists
	if fn.ep.getSessionEnvironment() != nil {
		focusGroup = append(focusGroup, fn.ep.sessionEnvTable)
	}

	// Add Self Environment table if environment exists
	if fn.ep.getSelfEnvironment() != nil {
		focusGroup = append(focusGroup, fn.ep.selfEnvTable)
	}

	// Add details table
	if fn.ep.browseEnvsDetails != nil {
		focusGroup = append(focusGroup, fn.ep.browseEnvsDetails)
	}

	fn.focusGroup = focusGroup
	fn.currentFocus = len(fn.focusGroup) - 1 // Focus on details table

	// Make the details table selectable since it's now focused
	if fn.ep.browseEnvsDetails != nil {
		fn.ep.browseEnvsDetails.SetSelectable(true, false)
	}
}

// updateFocusGroupForSearch updates the focus group back to search view
func (fn *FormNavigation) updateFocusGroupForSearch() {
	// Reset all tables to not selectable first
	fn.resetSelectable()

	// Rebuild focus group with search/list
	focusGroup := []tview.Primitive{}

	// Add Session Environment table if environment exists
	if fn.ep.getSessionEnvironment() != nil {
		focusGroup = append(focusGroup, fn.ep.sessionEnvTable)
	}

	// Add Self Environment table if environment exists
	if fn.ep.getSelfEnvironment() != nil {
		focusGroup = append(focusGroup, fn.ep.selfEnvTable)
	}

	// Add browse environments search and list
	focusGroup = append(focusGroup, fn.ep.browseEnvsSearch)
	focusGroup = append(focusGroup, fn.ep.browseEnvsList)

	fn.focusGroup = focusGroup
	fn.currentFocus = len(fn.focusGroup) - 1 // Focus on results list
}

// updateFocusGroupForEDS updates the focus group to include the EDS table
func (fn *FormNavigation) updateFocusGroupForEDS() {
	// Reset all tables to not selectable first
	fn.resetSelectable()

	// Rebuild focus group with EDS table instead of search/list
	focusGroup := []tview.Primitive{}

	// Add Session Environment table if environment exists
	if fn.ep.getSessionEnvironment() != nil {
		focusGroup = append(focusGroup, fn.ep.sessionEnvTable)
	}

	// Add Self Environment table if environment exists
	if fn.ep.getSelfEnvironment() != nil {
		focusGroup = append(focusGroup, fn.ep.selfEnvTable)
	}

	// Add search box
	focusGroup = append(focusGroup, fn.ep.browseEnvsSearch)

	// Add EDS table
	if fn.ep.browseEnvsEDSTable != nil {
		focusGroup = append(focusGroup, fn.ep.browseEnvsEDSTable)
	}

	fn.focusGroup = focusGroup
	fn.currentFocus = len(fn.focusGroup) - 1 // Focus on EDS table

	// Make the EDS table selectable since it's now focused
	if fn.ep.browseEnvsEDSTable != nil {
		fn.ep.browseEnvsEDSTable.SetSelectable(true, false)
		// Setup input capture for EDS table
		fn.setupEDSTableInputCapture()
		// Set help text for EDS table
		fn.helpTexts[fn.ep.browseEnvsEDSTable] = "Environment from EDS: ↑↓ Navigate | c: Copy field | Tab: Switch | Ctrl+N: New Environment"
		// Set focus to EDS table
		fn.ep.GetTUI().GetApplication().SetFocus(fn.ep.browseEnvsEDSTable)
		// Update help text
		fn.updateHelpText()
	}
}

// setupEDSTableInputCapture sets up input capture for the EDS table
func (fn *FormNavigation) setupEDSTableInputCapture() {
	if fn.ep.browseEnvsEDSTable == nil {
		return
	}

	fn.ep.browseEnvsEDSTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event == nil {
			return event
		}

		switch event.Key() {
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight:
			// Allow arrow keys to pass through for table navigation
			return event

		case tcell.KeyTab, tcell.KeyBacktab:
			// Allow Tab to pass through to parent for navigation
			return event

		case tcell.KeyRune:
			switch event.Rune() {
			case 'c', 'C':
				// Copy the selected field value to clipboard
				row, _ := fn.ep.browseEnvsEDSTable.GetSelection()
				if row >= 0 {
					// Always copy from column 1 (the value column)
					valueCell := fn.ep.browseEnvsEDSTable.GetCell(row, 1)
					if valueCell != nil {
						value := valueCell.Text
						if value != "" {
							clipboard.Write(clipboard.FmtText, []byte(value))
							// Get the field name from column 0
							labelCell := fn.ep.browseEnvsEDSTable.GetCell(row, 0)
							fieldName := "field"
							if labelCell != nil {
								fieldName = strings.TrimSuffix(labelCell.Text, ":")
							}
							fn.ep.UpdateStatus(fmt.Sprintf("Copied %s to clipboard", fieldName))
						} else {
							fn.ep.ShowError("No value to copy")
						}
					}
				}
				return nil
			}
		}

		return event
	})
}

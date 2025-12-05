package environments

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/tui/theme"
	"slv.sh/slv/internal/tui/utils"
)

// showEditForm shows a form below the details table for editing a field
func (ep *EnvironmentsPage) showEditForm(fieldName string, row int) {
	if ep.currentDetailsEnv == nil {
		return
	}

	colors := theme.GetCurrentPalette()
	ep.editingField = fieldName

	// Get current value
	valueCell := ep.browseEnvsDetails.GetCell(row, 1)
	currentValue := ""
	if valueCell != nil {
		currentValue = valueCell.Text
	}

	// For Tags, get the actual tags from the environment and format as comma-separated
	if fieldName == "Tags" && ep.currentDetailsEnv != nil {
		if len(ep.currentDetailsEnv.Tags) > 0 {
			currentValue = strings.Join(ep.currentDetailsEnv.Tags, ", ")
		} else {
			currentValue = ""
		}
	}

	// Create form
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(fmt.Sprintf("Edit %s", fieldName))
	form.SetTitleAlign(tview.AlignLeft)
	form.SetBorderColor(colors.Border)

	// Add input field
	form.AddInputField(fieldName+":", currentValue, 50, nil, nil)

	// Attach paste handler
	if inputField, ok := form.GetFormItem(0).(*tview.InputField); ok {
		utils.AttachPasteHandler(inputField)
	}

	// Add buttons
	form.AddButton("Cancel", func() {
		ep.hideEditForm(false)
	})

	form.AddButton("Push to Profile", func() {
		formItem := form.GetFormItem(0)
		if inputField, ok := formItem.(*tview.InputField); ok {
			newValue := strings.TrimSpace(inputField.GetText())
			ep.pushChangesToProfile(fieldName, newValue)
		}
	})

	form.AddButton("Save Locally", func() {
		formItem := form.GetFormItem(0)
		if inputField, ok := formItem.(*tview.InputField); ok {
			newValue := strings.TrimSpace(inputField.GetText())
			ep.applyTemporaryChanges(fieldName, newValue)
		}
	})

	form.SetButtonsAlign(tview.AlignCenter)

	// Set input capture to prevent Esc and Backspace from bubbling up
	// Tab and Backtab will be handled by the form itself for navigation
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event == nil {
			return event
		}

		switch event.Key() {
		case tcell.KeyEsc:
			// Consume Esc - don't let it bubble up
			return nil
		case tcell.KeyBackspace2, tcell.KeyBackspace:
			// Consume Backspace - don't let it bubble up (form handles it for text editing)
			return event // Let form handle backspace for text editing
		case tcell.KeyTab, tcell.KeyBacktab:
			// Let form handle Tab/Backtab for navigation between field and buttons
			return event
		}

		return event
	})

	// Store form reference
	ep.editForm = form

	// Update browseEnvsSection to show table and form
	ep.browseEnvsSection.Clear()
	envName := "Unnamed"
	if ep.currentDetailsEnv.Name != "" {
		envName = ep.currentDetailsEnv.Name
	}
	ep.browseEnvsSection.SetTitle(fmt.Sprintf("Environment: %s", envName)).SetTitleAlign(tview.AlignLeft)
	ep.browseEnvsSection.AddItem(ep.browseEnvsDetails, 0, 1, false) // Table takes flexible space
	ep.browseEnvsSection.AddItem(form, 8, 0, true)                  // Form takes fixed 8 rows, focusable

	// Set focus to the form
	ep.GetTUI().GetApplication().SetFocus(form)

	// Update help text
	ep.UpdateStatus(fmt.Sprintf("Editing %s | Tab: Navigate | Cancel: Cancel | Push to Profile: Save permanently | Temporary Local: Save locally", fieldName))
}

// hideEditForm hides the edit form and shows only the details table
func (ep *EnvironmentsPage) hideEditForm(withError bool) {
	ep.editForm = nil
	ep.editingField = ""

	// Show only the details table
	ep.browseEnvsSection.Clear()
	envName := "Unnamed"
	if ep.currentDetailsEnv != nil && ep.currentDetailsEnv.Name != "" {
		envName = ep.currentDetailsEnv.Name
	}
	ep.browseEnvsSection.SetTitle(fmt.Sprintf("Environment: %s", envName)).SetTitleAlign(tview.AlignLeft)
	ep.browseEnvsSection.AddItem(ep.browseEnvsDetails, 0, 1, true) // Flexible: takes all space

	// Set focus back to details table
	if withError {
		ep.navigation.SetFocusAtIndex(0)
	} else {
		ep.navigation.SetFocusAtIndex(ep.navigation.currentFocus)
	}

	// Update help text
	ep.navigation.updateHelpText()
}

// pushChangesToProfile pushes changes to the profile and refreshes the display
func (ep *EnvironmentsPage) pushChangesToProfile(fieldName string, newValue string) {
	if ep.currentDetailsEnv == nil {
		ep.hideEditForm(true)
		ep.ShowError("No environment selected")
		return
	}

	// Apply changes to environment
	ep.applyFieldChange(ep.currentDetailsEnv, fieldName, newValue)

	// Get active profile and push changes
	profile, err := profiles.GetActiveProfile()
	if err != nil {
		ep.hideEditForm(true)
		ep.ShowError(fmt.Sprintf("Error getting active profile: %v", err))
		return
	}

	// Push to profile
	if err := profile.PutEnv(ep.currentDetailsEnv); err != nil {
		ep.hideEditForm(true)
		ep.ShowError(fmt.Sprintf("Error pushing changes to profile: %v", err))
		return
	}

	ep.UpdateStatus(fmt.Sprintf("Successfully pushed %s changes to profile", fieldName))
	ep.refreshDetailsView()
	ep.hideEditForm(false)
	// ep.navigation.SetFocusAtIndex(ep.navigation.currentFocus)
}

// applyTemporaryChanges applies changes locally without pushing to profile
func (ep *EnvironmentsPage) applyTemporaryChanges(fieldName string, newValue string) {
	if ep.currentDetailsEnv == nil {
		ep.hideEditForm(true)
		ep.ShowError("No environment selected")
		return
	}

	// Apply changes to environment
	ep.applyFieldChange(ep.currentDetailsEnv, fieldName, newValue)

	ep.UpdateStatus(fmt.Sprintf("Applied temporary local change to %s", fieldName))
	ep.refreshDetailsView()
	ep.hideEditForm(false)
	// Focus is already set correctly by hideEditForm() which uses currentFocus
	// which is set to the details table by updateFocusGroupForDetails()
}

// applyFieldChange applies a field change to an environment
func (ep *EnvironmentsPage) applyFieldChange(env *environments.Environment, fieldName string, newValue string) {
	switch fieldName {
	case "Name":
		env.Name = newValue
	case "Email":
		env.SetEmail(newValue)
	case "Tags":
		// Parse comma-separated tags
		tags := strings.Split(newValue, ",")
		env.Tags = []string{}
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				env.AddTags(tag)
			}
		}
	}
}

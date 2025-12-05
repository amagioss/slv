package environments_new

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/environments/envproviders"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/tui/theme"
	"slv.sh/slv/internal/tui/utils"
)

// showMetadataForm shows the metadata form
func (nep *NewEnvironmentPage) showMetadataForm() {
	nep.currentStep = StepMetadata
	nep.updateStepIndicator()

	colors := theme.GetCurrentPalette()

	// Create metadata form
	nep.metadataForm = tview.NewForm().
		AddInputField("Name *", nep.envName, 40, nil, func(text string) {
			nep.envName = text
		}).
		AddInputField("Email", nep.envEmail, 40, nil, func(text string) {
			nep.envEmail = text
		}).
		AddInputField("Tags (comma-separated)", strings.Join(nep.envTags, ", "), 40, nil, func(text string) {
			if text != "" {
				nep.envTags = strings.Split(text, ",")
				for i := range nep.envTags {
					nep.envTags[i] = strings.TrimSpace(nep.envTags[i])
				}
			} else {
				nep.envTags = nil
			}
		}).
		AddCheckbox("Quantum Safe", nep.quantumSafe, func(checked bool) {
			nep.quantumSafe = checked
		}).
		AddButton("Next", func() {
			if nep.envName == "" {
				nep.ShowError("Name is required")
				return
			}
			// Check if provider needs configuration
			if nep.selectedProvider == "direct" {
				// Direct doesn't need config, go to confirmation
				nep.showConfirmation()
			} else {
				// Show provider config (including password)
				nep.showProviderConfig()
			}
		}).
		AddButton("Back", func() {
			nep.showProviderSelection()
		}).
		AddButton("Cancel", func() {
			nep.GetTUI().GetNavigation().GoBack()
		})

	// Attach paste handler to input fields
	for i := 0; i < nep.metadataForm.GetFormItemCount(); i++ {
		if inputField, ok := nep.metadataForm.GetFormItem(i).(*tview.InputField); ok {
			utils.AttachPasteHandler(inputField)
		}
	}

	nep.metadataForm.SetButtonsAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle(" Environment Metadata ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(colors.Border).
		SetBackgroundColor(colors.Background)

	// Clear and add to content
	nep.contentFlex.Clear()
	nep.contentFlex.AddItem(nil, 0, 1, false).
		AddItem(nep.metadataForm, 0, 3, true).
		AddItem(nil, 0, 1, false)

	// Set focus
	nep.GetTUI().GetApplication().SetFocus(nep.metadataForm)

	// Update help text
	if nep.navigation != nil {
		nep.navigation.updateHelpText()
	}
}

// showProviderSelection shows the provider selection
func (nep *NewEnvironmentPage) showProviderSelection() {
	nep.currentStep = StepProviderSelection
	nep.updateStepIndicator()

	colors := theme.GetCurrentPalette()

	// Create provider list
	nep.providerList = tview.NewList()

	// Add Password provider (USER type)
	nep.providerList.AddItem("Password (Self/User)", "Password-protected user environment", 'p', func() {
		nep.selectedProvider = envproviders.PasswordProviderId
		nep.selectedType = environments.USER
		nep.showMetadataForm()
	})

	// Add Direct service provider
	nep.providerList.AddItem("Direct (Service)", "Self-managed service environment (returns plain secret key)", 'd', func() {
		nep.selectedProvider = "direct"
		nep.selectedType = environments.SERVICE
		nep.showMetadataForm()
	})

	// Dynamically add cloud providers (all are SERVICE type)
	for _, providerId := range envproviders.ListIds() {
		if providerId == envproviders.PasswordProviderId {
			continue // Already added above
		}

		providerName := envproviders.GetName(providerId)
		providerDesc := envproviders.GetDesc(providerId)
		pid := providerId // Capture for closure

		nep.providerList.AddItem(fmt.Sprintf("%s (Service)", providerName), providerDesc, rune(providerName[0]), func() {
			nep.selectedProvider = pid
			nep.selectedType = environments.SERVICE
			nep.showMetadataForm()
		})
	}

	nep.providerList.SetBorder(true).
		SetTitle(" Select Provider ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(colors.Border).
		SetBackgroundColor(colors.Background)
	nep.providerList.SetWrapAround(false) // Disable looping behavior

	// Create a flex with cancel button
	nep.providerCancelForm = tview.NewForm().
		AddButton("Cancel", func() {
			nep.GetTUI().GetNavigation().GoBack()
		})

	nep.providerCancelForm.SetButtonsAlign(tview.AlignCenter).
		SetBackgroundColor(colors.Background)

	// Set up input capture for Tab navigation
	nep.providerList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyBacktab:
			// Switch focus to cancel button
			nep.GetTUI().GetApplication().SetFocus(nep.providerCancelForm)
			return nil
		}
		return event
	})

	nep.providerCancelForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyBacktab:
			// Switch focus back to provider list
			nep.GetTUI().GetApplication().SetFocus(nep.providerList)
			return nil
		}
		return event
	})

	providerFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nep.providerList, 0, 1, true).
		AddItem(nep.providerCancelForm, 3, 0, false)

	// Clear and add to content
	nep.contentFlex.Clear()
	nep.contentFlex.AddItem(nil, 0, 1, false).
		AddItem(providerFlex, 0, 3, true).
		AddItem(nil, 0, 1, false)

	// Set focus
	nep.GetTUI().GetApplication().SetFocus(nep.providerList)

	// Update help text
	if nep.navigation != nil {
		nep.navigation.updateHelpText()
	}
}

// showProviderConfig shows the provider configuration form
func (nep *NewEnvironmentPage) showProviderConfig() {
	nep.currentStep = StepProviderConfig
	nep.updateStepIndicator()

	colors := theme.GetCurrentPalette()

	// Get provider arguments
	providerArgs := envproviders.GetArgs(nep.selectedProvider)

	// Create provider config form
	nep.providerForm = tview.NewForm()

	for _, arg := range providerArgs {
		argId := arg.Id()
		argName := arg.Name()
		if argName == "" {
			argName = argId
		}
		argDesc := arg.Description()
		required := arg.Required()

		label := argName
		if required {
			label += " *"
		}

		// Get existing value if any
		existingValue := nep.providerInputs[argId]

		// For password fields, use password input
		if argId == "password" {
			nep.providerForm.AddPasswordField(label, existingValue, 40, '*', func(text string) {
				nep.providerInputs[argId] = text
			})
		} else {
			// Add input field and get reference to set placeholder
			nep.providerForm.AddInputField(label, existingValue, 40, nil, func(text string) {
				nep.providerInputs[argId] = text
			})

			// Attach paste handler
			if formItemCount := nep.providerForm.GetFormItemCount(); formItemCount > 0 {
				formItem := nep.providerForm.GetFormItem(formItemCount - 1)
				if inputField, ok := formItem.(*tview.InputField); ok {
					utils.AttachPasteHandler(inputField)
				}
			}

			// Set description as placeholder if available
			if argDesc != "" {
				// Get the last added form item (the input field we just added)
				formItemCount := nep.providerForm.GetFormItemCount()
				if formItemCount > 0 {
					formItem := nep.providerForm.GetFormItem(formItemCount - 1)
					if inputField, ok := formItem.(*tview.InputField); ok {
						// Truncate long descriptions for placeholder (max 70 chars)
						placeholderText := argDesc
						if len(placeholderText) > 70 {
							placeholderText = placeholderText[:67] + "..."
						}
						inputField.SetPlaceholder(placeholderText).
							SetPlaceholderTextColor(colors.TextMuted)
					}
				}
			}
		}
	}

	nep.providerForm.AddButton("Next", func() {
		// Validate required fields
		for _, arg := range providerArgs {
			if arg.Required() && nep.providerInputs[arg.Id()] == "" {
				nep.ShowError(fmt.Sprintf("%s is required", arg.Name()))
				return
			}
		}
		nep.showConfirmation()
	}).
		AddButton("Back", func() {
			nep.showMetadataForm()
		}).
		AddButton("Cancel", func() {
			nep.GetTUI().GetNavigation().GoBack()
		})

	nep.providerForm.SetButtonsAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle(fmt.Sprintf(" Configure %s Provider ", envproviders.GetName(nep.selectedProvider))).
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(colors.Border).
		SetBackgroundColor(colors.Background)

	// Clear and add to content
	nep.contentFlex.Clear()
	nep.contentFlex.AddItem(nil, 0, 1, false).
		AddItem(nep.providerForm, 0, 3, true).
		AddItem(nil, 0, 1, false)

	// Set focus
	nep.GetTUI().GetApplication().SetFocus(nep.providerForm)

	// Update help text
	if nep.navigation != nil {
		nep.navigation.updateHelpText()
	}
}

// showConfirmation shows the confirmation screen
func (nep *NewEnvironmentPage) showConfirmation() {
	nep.currentStep = StepConfirmation
	nep.updateStepIndicator()

	colors := theme.GetCurrentPalette()

	// Build confirmation text
	var confirmLines []string
	confirmLines = append(confirmLines, "Please review your environment settings:\n")
	confirmLines = append(confirmLines, fmt.Sprintf("[%s]Type:[-] %s", colors.Accent.String(), nep.selectedType))
	confirmLines = append(confirmLines, fmt.Sprintf("[%s]Name:[-] %s", colors.Accent.String(), nep.envName))
	if nep.envEmail != "" {
		confirmLines = append(confirmLines, fmt.Sprintf("[%s]Email:[-] %s", colors.Accent.String(), nep.envEmail))
	}
	if len(nep.envTags) > 0 {
		confirmLines = append(confirmLines, fmt.Sprintf("[%s]Tags:[-] %s", colors.Accent.String(), strings.Join(nep.envTags, ", ")))
	}
	confirmLines = append(confirmLines, fmt.Sprintf("[%s]Quantum Safe:[-] %v", colors.Accent.String(), nep.quantumSafe))

	providerName := "Direct (Self-managed)"
	if nep.selectedProvider != "direct" && nep.selectedProvider != "" {
		providerName = envproviders.GetName(nep.selectedProvider)
	}
	confirmLines = append(confirmLines, fmt.Sprintf("[%s]Provider:[-] %s", colors.Accent.String(), providerName))

	nep.confirmText = tview.NewTextView().
		SetDynamicColors(true).
		SetText(strings.Join(confirmLines, "\n")).
		SetTextColor(colors.TextPrimary)
	nep.confirmText.SetBackgroundColor(colors.Background)

	nep.confirmForm = tview.NewForm().
		AddCheckbox("Add to active profile", nep.addToProfile, func(checked bool) {
			nep.addToProfile = checked
		}).
		AddButton("Create", func() {
			nep.createEnvironment()
		}).
		AddButton("Back", func() {
			if nep.selectedProvider == "direct" {
				// Direct doesn't have config, go back to metadata
				nep.showMetadataForm()
			} else {
				// Go back to provider config
				nep.showProviderConfig()
			}
		}).
		AddButton("Cancel", func() {
			nep.GetTUI().GetNavigation().GoBack()
		})

	nep.confirmForm.SetButtonsAlign(tview.AlignCenter).
		SetBackgroundColor(colors.Background)

	confirmFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nep.confirmText, 0, 1, false).
		AddItem(nep.confirmForm, 5, 0, true)

	confirmFlex.SetBorder(true).
		SetTitle(" Review Environment ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(colors.Border).
		SetBackgroundColor(colors.Background)

	// Clear and add to content
	nep.contentFlex.Clear()
	nep.contentFlex.AddItem(nil, 0, 1, false).
		AddItem(confirmFlex, 0, 3, true).
		AddItem(nil, 0, 1, false)

	// Set focus
	nep.GetTUI().GetApplication().SetFocus(nep.confirmForm)

	// Update help text
	if nep.navigation != nil {
		nep.navigation.updateHelpText()
	}
}

// showResult shows the result screen after environment creation
func (nep *NewEnvironmentPage) showResult() {
	nep.currentStep = StepResult
	nep.updateStepIndicator()

	colors := theme.GetCurrentPalette()

	// Create success message
	successText := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf("[%s]✓ Environment created successfully![-]", colors.Success.String())).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(colors.Success)

	// Create result table with environment details
	nep.resultTable = nep.createResultTable()

	// Auto-copy EDS to clipboard
	eds, err := nep.createdEnv.ToDefStr(false)
	edsCopied := false
	if err == nil && eds != "" {
		clipboard.Write(clipboard.FmtText, []byte(eds))
		edsCopied = true
	}

	// Create info text for important messages
	var infoLines []string

	if edsCopied {
		infoLines = append(infoLines, fmt.Sprintf("[%s]📋 EDS has been copied to clipboard[-]", colors.Success.String()))
		infoLines = append(infoLines, "") // Empty line for spacing
	}

	// Show secret key warning if available (for direct service creation)
	if nep.secretKey != "" {
		infoLines = append(infoLines, fmt.Sprintf("[%s]⚠ IMPORTANT: This secret key is shown only once![-]", colors.Warning.String()))
		infoLines = append(infoLines, fmt.Sprintf("[%s]Save it securely. You cannot retrieve it later.[-]", colors.Warning.String()))
		infoLines = append(infoLines, "") // Empty line for spacing
	}

	// Show profile addition status
	if nep.addToProfile {
		profile, _ := profiles.GetActiveProfile()
		if profile != nil {
			infoLines = append(infoLines, fmt.Sprintf("[%s]✓ Added to profile: %s[-]", colors.Success.String(), profile.Name()))
		}
	}

	infoText := tview.NewTextView().
		SetDynamicColors(true).
		SetText(strings.Join(infoLines, "\n")).
		SetTextColor(colors.TextPrimary).
		SetTextAlign(tview.AlignCenter).
		SetScrollable(false)
	infoText.SetBackgroundColor(colors.Background)

	nep.resultForm = tview.NewForm().
		AddButton("Create Another", func() {
			nep.Refresh()
		}).
		AddButton("View Environments", func() {
			nep.GetTUI().GetNavigation().ShowEnvironments(true)
		}).
		AddButton("Done", func() {
			nep.GetTUI().GetNavigation().GoBack()
		})

	nep.resultForm.SetButtonsAlign(tview.AlignCenter).
		SetBackgroundColor(colors.Background)

	// Add input capture to allow Shift+Tab to switch back to table
	nep.resultForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyBacktab {
			nep.GetTUI().GetApplication().SetFocus(nep.resultTable)
			return nil
		}
		return event
	})

	// Calculate info text height based on number of lines
	infoHeight := len(infoLines)
	if infoHeight == 0 {
		infoHeight = 1 // Minimum height
	}

	resultFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(successText, 2, 0, false).
		AddItem(nep.resultTable, 9, 1, true).
		AddItem(infoText, infoHeight, 0, false).
		AddItem(nep.resultForm, 3, 0, false)

	resultFlex.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(colors.Border).
		SetBackgroundColor(colors.Background)

	// Clear and add to content
	nep.contentFlex.Clear()
	nep.contentFlex.AddItem(nil, 0, 2, false).
		AddItem(resultFlex, 25, 1, true).
		AddItem(nil, 0, 2, false)

	// Set focus to table for navigation and copying
	nep.GetTUI().GetApplication().SetFocus(nep.resultTable)

	// Update help text
	if nep.navigation != nil {
		nep.navigation.updateHelpText()
	}
}

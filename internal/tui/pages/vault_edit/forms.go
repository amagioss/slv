package vault_edit

import (
	"strings"

	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/tui/utils"
)

func (vep *VaultEditPage) createVaultConfigForm() *tview.Form {
	vaultConfigForm := tview.NewForm()

	// Vault Metadata Section
	vaultConfigForm.AddInputField("Vault Name", "", 30, nil, nil).
		AddInputField("File Name", "", 40, nil, nil).
		AddInputField("K8s Namespace (optional)", "", 30, nil, nil)

	// Attach paste handler to all input fields
	for i := 0; i < vaultConfigForm.GetFormItemCount(); i++ {
		if inputField, ok := vaultConfigForm.GetFormItem(i).(*tview.InputField); ok {
			utils.AttachPasteHandler(inputField)
		}
	}

	vaultConfigForm.SetBorder(true).
		SetTitle("Vault Configuration").
		SetTitleAlign(tview.AlignLeft)

	vaultConfigForm.GetFormItem(0).(*tview.InputField).SetText(vep.vault.Name)
	vaultConfigForm.GetFormItem(1).(*tview.InputField).SetText(strings.Split(vep.filePath, "/")[len(strings.Split(vep.filePath, "/"))-1])
	vaultConfigForm.GetFormItem(2).(*tview.InputField).SetText(vep.vault.Namespace)

	vep.vaultConfigForm = vaultConfigForm
	return vaultConfigForm
}

func (vep *VaultEditPage) createVaultOptionsForm() *tview.Form {

	optionsForm := tview.NewForm()

	// Vault Options
	optionsForm.AddCheckbox("Enable Hashing", false, nil).
		AddCheckbox("Quantum Safe", false, nil)

	optionsForm.SetBorder(true).
		SetTitle("Options (Read-Only)").
		SetTitleAlign(tview.AlignLeft)

	// Set checkbox character to tick mark
	for i := 0; i < optionsForm.GetFormItemCount(); i++ {
		if checkbox, ok := optionsForm.GetFormItem(i).(*tview.Checkbox); ok {
			checkbox.SetCheckedString("✓")
		}
	}

	vep.applyDisabledStyling(optionsForm)
	vep.optionsForm = optionsForm
	return optionsForm
}

func (vep *VaultEditPage) createVaultGrantAccessForm() *tview.Form {
	grantAccessForm := tview.NewForm()
	grantAccessForm.AddInputField("Public Key / Search String", "", 80, nil, func(text string) {
		// Only allow search if vault is unlocked
		if vep.IsVaultUnlocked() {
			// Handle search string based on prefix
			if text != "" {
				vep.searchEnvironments(text)
			} else {
				// Clear results when input is empty
				vep.searchResults.Clear()
				vep.searchResults.AddItem("", "", 0, nil)
			}
		}
	})

	// Attach paste handler to the input field
	if inputField, ok := grantAccessForm.GetFormItem(0).(*tview.InputField); ok {
		utils.AttachPasteHandler(inputField)
	}

	// Set title based on vault access state
	title := "Grant Access"
	if !vep.IsVaultUnlocked() {
		title = "Grant Access (Locked - No Access)"
		// Apply disabled styling when vault is locked
		vep.applyDisabledStyling(grantAccessForm)
		// Disable the input field
		grantAccessForm.GetFormItem(0).(*tview.InputField).SetDisabled(true)
	}

	grantAccessForm.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	vep.grantAccessForm = grantAccessForm
	return grantAccessForm
}

func (vep *VaultEditPage) createVaultShareWithSelfForm() *tview.Form {

	shareWithSelfForm := tview.NewForm()
	shareWithSelfForm.AddCheckbox("Share with Self", false, func(checked bool) {
		// Only allow changes if vault is unlocked
		if vep.IsVaultUnlocked() {
			vep.handleShareWithSelfChange(checked)
		}
	})
	shareWithSelfForm.SetBorder(false)

	// Set checkbox character to tick mark and store reference
	if checkbox, ok := shareWithSelfForm.GetFormItem(0).(*tview.Checkbox); ok {
		checkbox.SetCheckedString("✓")
		// Disable checkbox if vault is locked
		if !vep.IsVaultUnlocked() {
			checkbox.SetDisabled(true)
			// Apply disabled styling to the form
			vep.applyDisabledStyling(shareWithSelfForm)
		}
	}

	selfEnv := environments.GetSelf()
	for _, env := range vep.grantedEnvs {
		if env.PublicKey == selfEnv.PublicKey {
			shareWithSelfForm.GetFormItem(0).(*tview.Checkbox).SetChecked(true)
			break
		}
	}
	vep.shareWithSelfForm = shareWithSelfForm
	return shareWithSelfForm
}

func (vep *VaultEditPage) createVaultShareWithK8sForm() *tview.Form {
	shareWithK8sForm := tview.NewForm()
	shareWithK8sForm.AddCheckbox("Share with K8s Context", false, func(checked bool) {
		// Only allow changes if vault is unlocked
		if vep.IsVaultUnlocked() {
			vep.handleShareWithK8sChange(checked)
		}
	})
	shareWithK8sForm.SetBorder(false)

	// Set checkbox character to tick mark
	if checkbox, ok := shareWithK8sForm.GetFormItem(0).(*tview.Checkbox); ok {
		checkbox.SetCheckedString("✓")
		// Disable checkbox if vault is locked
		if !vep.IsVaultUnlocked() {
			checkbox.SetDisabled(true)
			// Apply disabled styling to the form
			vep.applyDisabledStyling(shareWithK8sForm)
		}
	}
	vep.shareWithK8sForm = shareWithK8sForm
	return shareWithK8sForm
}

func (vep *VaultEditPage) createSubmitButton() *tview.Button {

	// Create the submit button
	submitButton := tview.NewButton("Edit Vault").
		SetSelectedFunc(func() {
			vep.editVaultFromForm()
		})

	// Style the submit button
	submitButton.SetBorder(true).
		// SetTitle("Actions").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(vep.GetTheme().GetAccent()).
		SetBackgroundColor(vep.GetTheme().GetBackground())

	// Set button text and background colors for better visibility
	submitButton.SetLabelColor(vep.GetTheme().GetTextPrimary()).
		SetLabelColorActivated(vep.GetTheme().GetBackground()).
		SetBackgroundColorActivated(vep.GetTheme().GetAccent())

	vep.submitButton = submitButton
	return submitButton
}

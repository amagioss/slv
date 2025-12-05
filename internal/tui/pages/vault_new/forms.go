package vault_new

import (
	"fmt"

	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/utils"
)

func (vnp *VaultNewPage) createVaultConfigForm() *tview.Form {
	vaultConfigForm := tview.NewForm()

	// Vault Metadata Section
	vaultConfigForm.AddInputField("Vault Name", "", 40, nil, func(text string) {
		// Auto-update file name based on vault name
		if text != "" {
			fileName := text + ".slv.yaml"
			vaultConfigForm.GetFormItem(1).(*tview.InputField).SetText(fileName)
		} else {
			vaultConfigForm.GetFormItem(1).(*tview.InputField).SetText("")
		}
		vnp.SetTitle(fmt.Sprintf("New Vault at %s/%s.slv.yaml", vnp.currentDir, text))
	}).
		AddInputField("File Name", "", 40, nil, nil).
		AddInputField("K8s Namespace (optional)", "", 40, nil, nil)

	// Attach paste handler to all input fields
	for i := 0; i < vaultConfigForm.GetFormItemCount(); i++ {
		if inputField, ok := vaultConfigForm.GetFormItem(i).(*tview.InputField); ok {
			utils.AttachPasteHandler(inputField)
		}
	}

	vaultConfigForm.SetBorder(true).
		SetTitle("Vault Configuration").
		SetTitleAlign(tview.AlignLeft)

	vnp.vaultConfigForm = vaultConfigForm
	return vaultConfigForm
}

func (vnp *VaultNewPage) createVaultOptionsForm() *tview.Form {

	optionsForm := tview.NewForm()

	// Vault Options
	optionsForm.AddCheckbox("Enable Hashing", false, nil).
		AddCheckbox("Quantum Safe", false, nil)

	optionsForm.SetBorder(true).
		SetTitle("Options").
		SetTitleAlign(tview.AlignLeft)

	// Set checkbox character to tick mark
	for i := 0; i < optionsForm.GetFormItemCount(); i++ {
		if checkbox, ok := optionsForm.GetFormItem(i).(*tview.Checkbox); ok {
			checkbox.SetCheckedString("✓")
		}
	}
	vnp.optionsForm = optionsForm
	return optionsForm
}

func (vnp *VaultNewPage) createVaultGrantAccessForm() *tview.Form {
	grantAccessForm := tview.NewForm()
	grantAccessForm.AddInputField("Public Key / Search String", "", 80, nil, func(text string) {
		// Handle search string based on prefix
		if text != "" {
			vnp.searchEnvironments(text)
		} else {
			// Clear results when input is empty
			vnp.searchResults.Clear()
			vnp.searchResults.AddItem("", "", 0, nil)
		}
	})

	// Attach paste handler to the input field
	if inputField, ok := grantAccessForm.GetFormItem(0).(*tview.InputField); ok {
		utils.AttachPasteHandler(inputField)
	}

	grantAccessForm.SetBorder(true).SetTitle("Grant Access").SetTitleAlign(tview.AlignLeft)
	vnp.grantAccessForm = grantAccessForm
	return grantAccessForm
}

func (vnp *VaultNewPage) createVaultShareWithSelfForm() *tview.Form {

	shareWithSelfForm := tview.NewForm()
	shareWithSelfForm.AddCheckbox("Share with Self", false, func(checked bool) {
		vnp.handleShareWithSelfChange(checked)
	})
	shareWithSelfForm.SetBorder(false)

	// Set checkbox character to tick mark and store reference
	if checkbox, ok := shareWithSelfForm.GetFormItem(0).(*tview.Checkbox); ok {
		checkbox.SetCheckedString("✓")
	}
	vnp.shareWithSelfForm = shareWithSelfForm
	return shareWithSelfForm
}

func (vnp *VaultNewPage) createVaultShareWithK8sForm() *tview.Form {
	shareWithK8sForm := tview.NewForm()
	shareWithK8sForm.AddCheckbox("Share with K8s Context", false, func(checked bool) {
		vnp.handleShareWithK8sChange(checked)
	})
	shareWithK8sForm.SetBorder(false)

	// Set checkbox character to tick mark
	if checkbox, ok := shareWithK8sForm.GetFormItem(0).(*tview.Checkbox); ok {
		checkbox.SetCheckedString("✓")
	}
	vnp.shareWithK8sForm = shareWithK8sForm
	return shareWithK8sForm
}

func (vnp *VaultNewPage) createSubmitButton() *tview.Button {

	// Create the submit button
	submitButton := tview.NewButton("Create Vault").
		SetSelectedFunc(func() {
			vnp.createVaultFromForm()
		})

	// Style the submit button
	submitButton.SetBorder(true).
		// SetTitle("Actions").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(vnp.GetTheme().GetAccent()).
		SetBackgroundColor(vnp.GetTheme().GetBackground())

	// Set button text and background colors for better visibility
	submitButton.SetLabelColor(vnp.GetTheme().GetTextPrimary()).
		SetLabelColorActivated(vnp.GetTheme().GetBackground()).
		SetBackgroundColorActivated(vnp.GetTheme().GetAccent())

	vnp.submitButton = submitButton
	return submitButton
}

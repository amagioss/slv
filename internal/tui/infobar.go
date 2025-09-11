package tui // createInfoBar creates the shared info bar

import (
	"github.com/rivo/tview"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
)

func (t *TUI) createInfoBar() {
	profileName := getProfileName()
	selfEnvironment := environments.GetSelf()

	var leftContent string

	if selfEnvironment != nil {
		// Environment exists - show all details
		envName := selfEnvironment.Name
		envEmail := selfEnvironment.Email
		envType := string(selfEnvironment.EnvType)

		var publicKey string
		if pubKey, err := selfEnvironment.GetPublicKey(); err == nil && pubKey != nil {
			if keyStr, err := pubKey.String(); err == nil {
				publicKey = keyStr
			} else {
				publicKey = "Error getting key"
			}
		} else {
			publicKey = "No key available"
		}

		leftContent = "[cyan]Profile: [white]" + profileName + "\n" +
			"[cyan]Environment: [white]" + envName + "\n" +
			"[cyan]Environment Email: [white]" + envEmail + "\n" +
			"[cyan]Environment Type: [white]" + envType + "\n" +
			"[cyan]Environment Public Key: [white]" + publicKey
	} else {
		// No environment - show minimal info
		leftContent = "[cyan]Profile: [white]" + profileName + "\n" +
			"[yellow]No self environment is set"
	}

	// Create left column for info
	leftColumn := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetWrap(true).
		SetText(leftContent)

		// Create right column for logo
	// 	rightContent := `[cyan]╔══════════════╗
	// ║   [yellow]SLV[cyan]  [white]║
	// ║  [yellow]SECURE[cyan][white]║
	// ║ [yellow]LOCAL[cyan]  [white]║
	// ║ [yellow]VAULT[cyan]  [white]║
	// ╚══════════════╝`

	rightContent := `  _________.____ ____   ____
 /   _____/|    |\   \ /   /
 \_____  \ |    | \   Y   / 
 /        \|    |__\     /  
/_______  /|_______ \___/   
        \/         \/       `

	rightColumn := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetWrap(false).
		SetText(rightContent)

	// Create flex container to hold both columns and give it a border/title
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).   // Set direction to column (horizontal)
		AddItem(leftColumn, 0, 1, false). // Left column takes remaining space
		AddItem(rightColumn, 30, 0, false)

	flex.SetBorder(true).
		SetBorderColor(t.theme.Accent).
		SetTitle("Secure Local Vault").
		SetTitleAlign(tview.AlignCenter)

	t.infoBar = flex
}

func getProfileName() string {
	profile, err := profiles.GetActiveProfile()
	if err != nil || profile == nil {
		return "No Profile"
	}
	return profile.Name()
}

package cli

import (
	"os"

	"savesecrets.org/slv/cli/internal/commands"
	"savesecrets.org/slv/core/environments/providers"
)

func RunCLI() {
	providers.LoadDefaults()
	if err := commands.SlvCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

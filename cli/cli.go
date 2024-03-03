package cli

import (
	"os"

	"savesecrets.org/slv/cli/internal/commands"
)

func RunCLI() {
	if err := commands.SlvCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

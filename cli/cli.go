package cli

import (
	"os"

	"github.com/amagimedia/slv/cli/internal/commands"
	"github.com/amagimedia/slv/core/environments/providers"
)

func RunCLI() {
	providers.LoadDefaults()
	if err := commands.SlvCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

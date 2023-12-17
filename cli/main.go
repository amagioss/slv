package main

import (
	"os"

	"github.com/shibme/slv/cli/commands"
	"github.com/shibme/slv/core/environments/providers"
)

func main() {
	providers.RegisterDefaultProviders()
	if err := commands.SlvCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

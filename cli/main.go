package main

import (
	"os"

	"github.com/shibme/slv/cli/commands"
	"github.com/shibme/slv/core/secretkeystore/providers/kms"
)

func main() {
	kms.RegisterKMSProviders()
	if err := commands.SlvCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

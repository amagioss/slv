package main

import (
	"os"

	"github.com/shibme/slv/cli/commands"
)

func main() {
	if err := commands.SlvCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

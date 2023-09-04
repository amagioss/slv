package main

import (
	"github.com/shibme/slv/cli/commands"
)

func main() {
	if err := commands.SlvCommand().Execute(); err != nil {
		commands.PrintErrorAndExit(err)
	}
}

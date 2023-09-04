package main

import (
	"fmt"
	"os"

	"github.com/shibme/slv/cli/commands"
	"github.com/spf13/cobra"
)

var slvCmd = &cobra.Command{
	Use:   "slv",
	Short: "SLV is a tool to encrypt secrets locally",
	Long:  `SLV is a tool for storing and managing secrets in an encrypted format.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func CLI() {
	slvCmd.AddCommand(commands.EnvCommand())
	slvCmd.AddCommand(commands.ConfigCommand())
	slvCmd.AddCommand(commands.VaultCommand())
	slvCmd.AddCommand(commands.SecretsCommand())
	if err := slvCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	CLI()
}

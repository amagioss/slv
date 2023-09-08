package commands

import (
	"github.com/spf13/cobra"
)

func SlvCommand() *cobra.Command {
	if slvCmd != nil {
		return slvCmd
	}
	slvCmd = &cobra.Command{
		Use:   "slv",
		Short: "SLV is a tool to encrypt secrets locally",
		Long:  `SLV is a tool for storing and managing secrets in an encrypted format.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	slvCmd.AddCommand(envCommand())
	slvCmd.AddCommand(profileCommand())
	slvCmd.AddCommand(vaultCommand())
	slvCmd.AddCommand(secretCommand())
	return slvCmd
}

package cmdsystem

import (
	"github.com/spf13/cobra"
)

func SystemCommand() *cobra.Command {
	if systemCmd != nil {
		return systemCmd
	}
	systemCmd = &cobra.Command{
		Use:     "system",
		Aliases: []string{"systems"},
		Short:   "System level commands",
		Long:    `System level operations can be carried out using this command`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	systemCmd.AddCommand(systemResetCommand())
	return systemCmd
}

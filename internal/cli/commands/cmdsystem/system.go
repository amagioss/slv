package cmdsystem

import (
	"github.com/spf13/cobra"
)

func SystemCommand() *cobra.Command {
	if systemCmd == nil {
		systemCmd = &cobra.Command{
			Use:     "system",
			Aliases: []string{"systems"},
			Short:   "Manage SLV system-level settings",
			Long:    `System-level operations such as resetting the local SLV installation.`,
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		systemCmd.AddCommand(systemResetCommand())
	}
	return systemCmd
}

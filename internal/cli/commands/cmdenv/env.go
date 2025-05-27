package cmdenv

import (
	"github.com/spf13/cobra"
)

func EnvCommand() *cobra.Command {
	if envCmd == nil {
		envCmd = &cobra.Command{
			Use:     "env",
			Aliases: []string{"envs", "environment", "environments"},
			Short:   "Manage SLV environments",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		envCmd.AddCommand(envNewCommand())
		envCmd.AddCommand(envListCommand())
		envCmd.AddCommand(envDeleteCommand())
		envCmd.AddCommand(envSetSelfCommand())
		envCmd.AddCommand(envShowCommand())
		envCmd.AddCommand(envAddCommand())
	}
	return envCmd
}

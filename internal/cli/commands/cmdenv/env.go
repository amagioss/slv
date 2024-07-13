package cmdenv

import (
	"github.com/spf13/cobra"
)

func EnvCommand() *cobra.Command {
	if envCmd != nil {
		return envCmd
	}
	envCmd = &cobra.Command{
		Use:     "env",
		Aliases: []string{"envs", "environment", "environments"},
		Short:   "Environment operations",
		Long:    `Environment operations in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	envCmd.AddCommand(envNewCommand())
	envCmd.AddCommand(envAddCommand())
	envCmd.AddCommand(envListSearchCommand())
	envCmd.AddCommand(envDeleteCommand())
	envCmd.AddCommand(envSetCommand())
	envCmd.AddCommand(envShowCommand())
	envCmd.AddCommand(envShowK8sCommand())
	return envCmd
}

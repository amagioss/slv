package cmdenv

import (
	"github.com/spf13/cobra"
	"savesecrets.org/slv/cli/internal/commands/utils"
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
	envCmd.PersistentFlags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage)
	envCmd.AddCommand(envNewCommand())
	envCmd.AddCommand(envAddCommand())
	envCmd.AddCommand(envListSearchCommand())
	envCmd.AddCommand(envDeleteCommand())
	envCmd.AddCommand(envSelfCommand())
	return envCmd
}

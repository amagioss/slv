package cmdenv

import (
	"fmt"

	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/input"
	"oss.amagi.com/slv/internal/core/profiles"
)

func envDeleteCommand() *cobra.Command {
	if envDeleteCmd != nil {
		return envDeleteCmd
	}
	envDeleteCmd = &cobra.Command{
		Use:     "del",
		Aliases: []string{"delete", "rm", "remove"},
		Short:   "Deletes an environment from current profile",
		Run: func(cmd *cobra.Command, args []string) {
			profile, err := profiles.GetDefaultProfile()
			if err != nil {
				utils.ExitOnError(err)
			}
			queries, err := cmd.Flags().GetStringSlice(EnvSearchFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			envs, err := profile.SearchEnvs(queries)
			if err != nil {
				utils.ExitOnError(err)
			}
			if envs != nil {
				for _, env := range envs {
					ShowEnv(*env, false, false)
					fmt.Println()
				}
				confirm, err := input.GetConfirmation("Are you sure you wish to delete the above environment(s) [yes/no]: ", "yes")
				if err != nil {
					utils.ExitOnError(err)
				}
				if confirm {
					for _, env := range envs {
						if profile.DeleteEnv(env.PublicKey); err != nil {
							utils.ExitOnError(err)
						}
						fmt.Printf("Environment %s deleted successfully\n", env.Name)
					}
				}
			}
			utils.SafeExit()
		},
	}
	envDeleteCmd.Flags().StringSliceP(EnvSearchFlag.Name, EnvSearchFlag.Shorthand, []string{}, EnvSearchFlag.Usage)
	envDeleteCmd.Flags().StringSliceP(EnvPublicKeysFlag.Name, EnvPublicKeysFlag.Shorthand, []string{}, EnvPublicKeysFlag.Usage)
	envDeleteCmd.Flags().BoolP(EnvSelfFlag.Name, EnvSelfFlag.Shorthand, false, EnvSelfFlag.Usage)
	envDeleteCmd.Flags().BoolP(EnvK8sFlag.Name, EnvK8sFlag.Shorthand, false, EnvK8sFlag.Usage)
	envDeleteCmd.Flags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage)
	return envDeleteCmd
}

package cmdenv

import (
	"fmt"

	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/input"
	"slv.sh/slv/internal/core/profiles"
)

func envDeleteCommand() *cobra.Command {
	if envDeleteCmd == nil {
		envDeleteCmd = &cobra.Command{
			Use:     "rm",
			Aliases: []string{"remove", "delete", "del"},
			Short:   "Removes an environment from active profile",
			Run: func(cmd *cobra.Command, args []string) {
				profile, err := profiles.GetActiveProfile()
				if err != nil {
					utils.ExitOnError(err)
				}
				if !profile.IsPushSupported() {
					utils.ExitOnError(fmt.Errorf("profile (%s) does not support deleting environments", profile.Name()))
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
		if err := envDeleteCmd.RegisterFlagCompletionFunc(EnvSearchFlag.Name, EnvSearchCompletion); err != nil {
			utils.ExitOnError(err)
		}
		envDeleteCmd.Flags().StringSliceP(EnvPublicKeysFlag.Name, EnvPublicKeysFlag.Shorthand, []string{}, EnvPublicKeysFlag.Usage)
		envDeleteCmd.Flags().BoolP(EnvSelfFlag.Name, EnvSelfFlag.Shorthand, false, EnvSelfFlag.Usage)
		envDeleteCmd.Flags().BoolP(EnvK8sFlag.Name, EnvK8sFlag.Shorthand, false, EnvK8sFlag.Usage)
	}
	return envDeleteCmd
}

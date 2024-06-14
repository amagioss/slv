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
					utils.ShowEnv(*env, false, false)
					fmt.Println()
				}
				confirm, err := input.GetConfirmation("Are you sure you want to delete the above environment(s) [yes/no]: ", "yes")
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
	envDeleteCmd.MarkFlagRequired(EnvSearchFlag.Name)
	return envDeleteCmd
}

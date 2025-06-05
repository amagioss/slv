package cmdenv

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
)

func envAddCommand() *cobra.Command {
	if envAddCmd == nil {
		envAddCmd = &cobra.Command{
			Use:     "add",
			Aliases: []string{"put"},
			Short:   "Adds/updates an environment into the active profile",
			Run: func(cmd *cobra.Command, args []string) {
				profile, err := profiles.GetActiveProfile()
				if err != nil {
					utils.ExitOnError(err)
				}
				if !profile.IsPushSupported() {
					utils.ExitOnError(fmt.Errorf("profile (%s) does not support adding environments", profile.Name()))
				}
				envdef, _ := cmd.Flags().GetString(envDefFlag.Name)
				setAsRoot, _ := cmd.Flags().GetBool(envSetRootFlag.Name)
				if setAsRoot && len(envdef) > 1 {
					utils.ExitOnError(fmt.Errorf("cannot set more than one environment as root"))
				}
				var successMessage string
				var env *environments.Environment
				if env, err = environments.FromDefStr(envdef); err == nil && env != nil {
					if setAsRoot {
						err = profile.SetRoot(env)
						successMessage = fmt.Sprintf("Successfully set %s as root environment for profile %s", color.GreenString(env.Name), color.GreenString(profile.Name()))
					} else {
						err = profile.PutEnv(env)
					}
				}
				if err != nil {
					utils.ExitOnError(err)
				}
				if successMessage == "" {
					successMessage = fmt.Sprintf("Successfully added environment to profile (%s)", color.GreenString(profile.Name()))
				}
				fmt.Println(successMessage)
				utils.SafeExit()
			},
		}
		envAddCmd.Flags().StringP(envDefFlag.Name, envDefFlag.Shorthand, "", envDefFlag.Usage)
		envAddCmd.Flags().Bool(envSetRootFlag.Name, false, envSetRootFlag.Usage)
		envAddCmd.MarkFlagRequired(envDefFlag.Name)
	}
	return envAddCmd
}

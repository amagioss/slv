package cmdenv

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/profiles"
)

func envAddCommand() *cobra.Command {
	if envAddCmd != nil {
		return envAddCmd
	}
	envAddCmd = &cobra.Command{
		Use:     "add",
		Aliases: []string{"set", "put", "store", "a"},
		Short:   "Adds an environment to the current profile",
		Run: func(cmd *cobra.Command, args []string) {
			envdefs, err := cmd.Flags().GetStringSlice(envDefFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			profile, err := profiles.GetDefaultProfile()
			if err != nil {
				utils.ExitOnError(err)
			}
			setAsRoot, _ := cmd.Flags().GetBool(envSetRootFlag.Name)
			if setAsRoot && len(envdefs) > 1 {
				utils.ExitOnError(fmt.Errorf("cannot set more than one environment as root"))
			}
			var successMessage string
			for _, envdef := range envdefs {
				var env *environments.Environment
				if env, err = environments.FromEnvDef(envdef); err == nil && env != nil {
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
			}
			if successMessage == "" {
				successMessage = fmt.Sprintf("Successfully added %d environments to profile %s", len(envdefs), color.GreenString(profile.Name()))
			}
			fmt.Println(successMessage)
			utils.SafeExit()
		},
	}
	envAddCmd.Flags().StringSliceP(envDefFlag.Name, envDefFlag.Shorthand, []string{}, envDefFlag.Usage)
	envAddCmd.Flags().Bool(envSetRootFlag.Name, false, envSetRootFlag.Usage)
	envAddCmd.MarkFlagRequired(envDefFlag.Name)
	return envAddCmd
}

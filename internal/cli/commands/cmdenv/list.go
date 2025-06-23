package cmdenv

import (
	"fmt"

	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
)

func envListCommand() *cobra.Command {
	if envListCmd == nil {
		envListCmd = &cobra.Command{
			Use:     "ls",
			Aliases: []string{"list", "search", "find", "get"},
			Short:   "List/Search environments from the active profile",
			Run: func(cmd *cobra.Command, args []string) {
				profile, err := profiles.GetActiveProfile()
				if err != nil {
					utils.ExitOnError(err)
				}
				queries, err := cmd.Flags().GetStringSlice(EnvSearchFlag.Name)
				if err != nil {
					utils.ExitOnError(err)
				}
				showEnvDef, _ := cmd.Flags().GetBool(showEnvDefFlag.Name)
				var envs []*environments.Environment
				if len(queries) == 0 {
					envs, err = profile.ListEnvs()
				} else {
					envs, err = profile.SearchEnvs(queries)
				}
				if err != nil {
					utils.ExitOnError(err)
				}
				for i, env := range envs {
					ShowEnv(*env, showEnvDef, false)
					if i < len(envs)-1 {
						fmt.Println()
					}
				}
				utils.SafeExit()
			},
		}
		envListCmd.Flags().StringSliceP(EnvSearchFlag.Name, EnvSearchFlag.Shorthand, []string{}, EnvSearchFlag.Usage)
		envListCmd.Flags().BoolP(showEnvDefFlag.Name, showEnvDefFlag.Shorthand, false, showEnvDefFlag.Usage)
		if err := envListCmd.RegisterFlagCompletionFunc(EnvSearchFlag.Name, EnvSearchCompletion); err != nil {
			utils.ExitOnError(err)
		}
	}
	return envListCmd
}

func EnvSearchCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	profile, err := profiles.GetActiveProfile()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var envs []*environments.Environment
	if toComplete == "" {
		envs, err = profile.ListEnvs()
	} else {
		envs, err = profile.SearchEnvs([]string{toComplete})
	}
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var envNames []string
	for _, env := range envs {
		if env.Name != "" {
			envNames = append(envNames, env.Name)
		}
	}
	return envNames, cobra.ShellCompDirectiveNoFileComp
}

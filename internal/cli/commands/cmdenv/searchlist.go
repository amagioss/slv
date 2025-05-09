package cmdenv

import (
	"fmt"

	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
)

func envListSearchCommand() *cobra.Command {
	if envListSearchCmd == nil {
		envListSearchCmd = &cobra.Command{
			Use:     "get",
			Aliases: []string{"list", "ls", "search", "find"},
			Short:   "List/Search environments from profile",
			Run: func(cmd *cobra.Command, args []string) {
				profile, err := profiles.GetCurrentProfile()
				if err != nil {
					utils.ExitOnError(err)
				}
				queries, err := cmd.Flags().GetStringSlice(EnvSearchFlag.Name)
				if err != nil {
					utils.ExitOnError(err)
				}
				var envs []*environments.Environment
				if len(queries) == 0 {
					envs, err = profile.ListEnvs()
				} else {
					envs, err = profile.SearchEnvs(queries)
				}
				if err != nil {
					utils.ExitOnError(err)
				}
				for _, env := range envs {
					ShowEnv(*env, false, false)
					fmt.Println()
				}
				utils.SafeExit()
			},
		}
		envListSearchCmd.Flags().StringSliceP(EnvSearchFlag.Name, EnvSearchFlag.Shorthand, []string{}, EnvSearchFlag.Usage)
		if err := envListSearchCmd.RegisterFlagCompletionFunc(EnvSearchFlag.Name, EnvSearchCompletion); err != nil {
			utils.ExitOnError(err)
		}
	}
	return envListSearchCmd
}

func EnvSearchCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	profile, err := profiles.GetCurrentProfile()
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

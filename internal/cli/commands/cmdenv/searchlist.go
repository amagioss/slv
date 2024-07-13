package cmdenv

import (
	"fmt"

	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/profiles"
)

func envListSearchCommand() *cobra.Command {
	if envListSearchCmd != nil {
		return envListSearchCmd
	}
	envListSearchCmd = &cobra.Command{
		Use:     "get",
		Aliases: []string{"list", "ls", "search", "find"},
		Short:   "List/Search environments from profile",
		Run: func(cmd *cobra.Command, args []string) {
			profile, err := profiles.GetDefaultProfile()
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
	return envListSearchCmd
}

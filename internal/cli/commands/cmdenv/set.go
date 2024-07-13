package cmdenv

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/input"
)

func envSetCommand() *cobra.Command {
	if envSetSCmd != nil {
		return envSetSCmd
	}
	envSetSCmd = &cobra.Command{
		Use:     "set",
		Aliases: []string{"put", "update"},
		Short:   "Set/update an environments",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	envSetSCmd.AddCommand(envSetSelfCommand())
	return envSetSCmd
}

func envSetSelfCommand() *cobra.Command {
	if envSetSelfSCmd != nil {
		return envSetSelfSCmd
	}
	envSetSelfSCmd = &cobra.Command{
		Use:     "set",
		Aliases: []string{"save", "put", "store", "s"},
		Short:   "Sets a given environment as self",
		Run: func(cmd *cobra.Command, args []string) {
			envDef := cmd.Flag(envDefFlag.Name).Value.String()
			env, err := environments.FromEnvDef(envDef)
			if err != nil {
				utils.ExitOnError(err)
			}
			if env.EnvType != environments.USER {
				utils.ExitOnError(fmt.Errorf("only user environments can be registered as self"))
			}
			if env.SecretBinding == "" {
				secretBinding, err := input.GetVisibleInput("Enter the secret binding: ")
				if err != nil {
					utils.ExitOnError(err)
				}
				env.SecretBinding = secretBinding
			}
			if err = env.MarkAsSelf(); err != nil {
				utils.ExitOnError(err)
			}
			ShowEnv(*env, true, true)
			fmt.Println(color.GreenString("Successfully registered self environment"))
		},
	}
	envSetSelfSCmd.Flags().StringP(envDefFlag.Name, envDefFlag.Shorthand, "", envDefFlag.Usage)
	envSetSelfSCmd.MarkFlagRequired(envDefFlag.Name)
	return envSetSelfSCmd
}

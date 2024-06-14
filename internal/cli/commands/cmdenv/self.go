package cmdenv

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/input"
)

func envSelfCommand() *cobra.Command {
	if envSelfCmd != nil {
		return envSelfCmd
	}
	envSelfCmd = &cobra.Command{
		Use:     "self",
		Aliases: []string{"user", "me", "my", "current"},
		Short:   "Shows the current user environment if registered",
		Run: func(cmd *cobra.Command, args []string) {
			env := environments.GetSelf()
			if env == nil {
				fmt.Println("No environment registered as self.")
			} else {
				utils.ShowEnv(*env, true, true)
			}
			utils.SafeExit()
		},
	}
	envSelfCmd.AddCommand(envSelfSetCommand())
	return envSelfCmd
}

func envSelfSetCommand() *cobra.Command {
	if envSelfSetCmd != nil {
		return envSelfSetCmd
	}
	envSelfSetCmd = &cobra.Command{
		Use:     "set",
		Aliases: []string{"save", "put", "store", "s"},
		Short:   "Shows the current environment if registered",
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
			utils.ShowEnv(*env, true, true)
			fmt.Println(color.GreenString("Successfully registered self environment"))
		},
	}
	envSelfSetCmd.Flags().StringP(envDefFlag.Name, envDefFlag.Shorthand, "", envDefFlag.Usage)
	envSelfSetCmd.MarkFlagRequired(envDefFlag.Name)
	return envSelfSetCmd
}

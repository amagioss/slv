package cmdenv

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/input"
)

func envSetSelfCommand() *cobra.Command {
	if envSetSelfSCmd == nil {
		envSetSelfSCmd = &cobra.Command{
			Use:     "set-self",
			Aliases: []string{"self-set", "register-self", "self-register", "register"},
			Short:   "Registers an environment in the current host as self",
			Run: func(cmd *cobra.Command, args []string) {
				envDef := cmd.Flag(envDefFlag.Name).Value.String()
				env, err := environments.FromDefStr(envDef)
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
				if err = env.SetAsSelf(); err != nil {
					utils.ExitOnError(err)
				}
				ShowEnv(*env, true, true)
				fmt.Println(color.GreenString("Successfully registered self environment"))
			},
		}
		envSetSelfSCmd.Flags().StringP(envDefFlag.Name, envDefFlag.Shorthand, "", envDefFlag.Usage)
		envSetSelfSCmd.MarkFlagRequired(envDefFlag.Name)
	}
	return envSetSelfSCmd
}

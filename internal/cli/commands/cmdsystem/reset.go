package cmdsystem

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/cmdenv"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/input"
)

func systemResetCommand() *cobra.Command {
	if systemResetCmd == nil {
		systemResetCmd = &cobra.Command{
			Use:     "reset",
			Aliases: []string{"purge", "prune", "clean", "clear"},
			Short:   "Reset the local SLV installation",
			Long:    `Removes all profiles, environments, and locally stored SLV data. Remote profile data is not affected.`,
			Run: func(cmd *cobra.Command, args []string) {
				selfEnv := environments.GetSelf()
				confirm, _ := cmd.Flags().GetBool(yesFlag.Name)
				if !confirm || selfEnv != nil {
					if selfEnv != nil {
						fmt.Println(color.YellowString("You have a registered environment. Consider backing it up before continuing:"))
						cmdenv.ShowEnv(*selfEnv, true, true)
					}
					var err error
					if confirm, err = input.GetConfirmation("Are you sure you wish to proceed? (yes/no): ", "yes"); err != nil {
						utils.ExitOnError(err)
					}
				}
				if confirm {
					err := config.ResetAppDataDir()
					if err == nil {
						fmt.Println(color.GreenString("System reset successful"))
					} else {
						utils.ExitOnError(err)
					}
				} else {
					fmt.Println(color.YellowString("System reset aborted"))
				}
				utils.SafeExit()
			},
		}
		systemResetCmd.Flags().BoolP(yesFlag.Name, yesFlag.Shorthand, false, yesFlag.Usage)
	}
	return systemResetCmd
}

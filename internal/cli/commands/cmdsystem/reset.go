package cmdsystem

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/cmdenv"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/config"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/input"
)

func systemResetCommand() *cobra.Command {
	if systemResetCmd == nil {
		systemResetCmd = &cobra.Command{
			Use:     "reset",
			Aliases: []string{"reset", "pruge", "prune", "clean", "clear"},
			Short:   "Reset the system",
			Long:    `Cleans all existing profiles and any other data`,
			Run: func(cmd *cobra.Command, args []string) {
				selfEnv := environments.GetSelf()
				confirm, _ := cmd.Flags().GetBool(yesFlag.Name)
				if !confirm || selfEnv != nil {
					if selfEnv != nil {
						fmt.Println(color.YellowString("You have a configured environment which you might have to consider backing up:"))
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

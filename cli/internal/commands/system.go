package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"savesecrets.org/slv/core/config"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/input"
)

func systemCommand() *cobra.Command {
	if systemCmd != nil {
		return systemCmd
	}
	systemCmd = &cobra.Command{
		Use:     "system",
		Aliases: []string{"systems"},
		Short:   "System level commands",
		Long:    `System level operations can be carried out using this command`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	systemCmd.AddCommand(systemResetCommand())
	return systemCmd
}

func systemResetCommand() *cobra.Command {
	if systemResetCmd != nil {
		return systemResetCmd
	}
	systemResetCmd = &cobra.Command{
		Use:     "reset",
		Aliases: []string{"reset", "pruge", "prune", "clean", "clear"},
		Short:   "Reset the system",
		Long:    `Cleans all existing profiles and any other data`,
		Run: func(cmd *cobra.Command, args []string) {
			selfEnv := environments.GetSelf()
			confirm, _ := cmd.Flags().GetBool(yesFlag.name)
			if !confirm || selfEnv != nil {
				if selfEnv != nil {
					fmt.Println(color.YellowString("You have a configured environment which you might have to consider backing up:"))
					showEnv(*selfEnv, true, true)
				}
				var err error
				if confirm, err = input.GetConfirmation("Are you sure you want to proceed? (y/n): ", "y"); err != nil {
					exitOnError(err)
				}
			}
			if confirm {
				err := config.ResetAppDataDir()
				if err == nil {
					fmt.Println(color.GreenString("System reset successful"))
				} else {
					exitOnError(err)
				}
			} else {
				fmt.Println(color.YellowString("System reset aborted"))
			}
			safeExit()
		},
	}
	systemResetCmd.Flags().BoolP(yesFlag.name, yesFlag.shorthand, false, yesFlag.usage)
	return systemResetCmd
}

package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"savesecrets.org/slv/core/config"
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
			confirm, _ := cmd.Flags().GetBool(yesFlag.name)
			if !confirm {
				fmt.Print("Are you sure you want to proceed? (y/n): ")
				var confirmation string
				fmt.Scanln(&confirmation)
				confirm = confirmation == "y"
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

package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
)

func vaultDeleteCommand() *cobra.Command {
	if vaultDeleteCmd != nil {
		return vaultDeleteCmd
	}
	vaultDeleteCmd = &cobra.Command{
		Use:     "rm",
		Aliases: []string{"del", "remove", "delete", "destroy", "erase"},
		Short:   "Removes secret from the vault",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				utils.ExitOnError(err)
			}
			secretNames, err := cmd.Flags().GetStringSlice(itemNameFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			if len(secretNames) == 0 {
				if err = vault.Delete(); err != nil {
					utils.ExitOnError(err)
				}
				fmt.Printf(color.GreenString("Successfully deleted the vault: %s\n"), vaultFile)
			} else {
				if err = vault.DeleteItems(secretNames); err != nil {
					utils.ExitOnError(err)
				}
				fmt.Printf(color.GreenString("Successfully deleted the secrets: %v from the vault: %s\n"), secretNames, vaultFile)
			}
		},
	}
	vaultDeleteCmd.Flags().StringSliceP(itemNameFlag.Name, itemNameFlag.Shorthand, []string{}, itemNameFlag.Usage)
	return vaultDeleteCmd
}

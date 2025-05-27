package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/vaults"
)

func vaultDeleteCommand() *cobra.Command {
	if vaultDeleteCmd == nil {
		vaultDeleteCmd = &cobra.Command{
			Use:     "rm",
			Aliases: []string{"del", "remove", "delete", "destroy", "erase"},
			Short:   "Removes an item from the vault",
			Run: func(cmd *cobra.Command, args []string) {
				vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
				vault, err := vaults.Get(vaultFile)
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
		if err := vaultDeleteCmd.RegisterFlagCompletionFunc(itemNameFlag.Name, vaultItemNameCompletion); err != nil {
			utils.ExitOnError(err)
		}
	}
	return vaultDeleteCmd
}

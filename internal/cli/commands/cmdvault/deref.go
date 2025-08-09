package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
)

func vaultDerefCommand() *cobra.Command {
	if vaultDerefCmd == nil {
		vaultDerefCmd = &cobra.Command{
			Use:   "deref",
			Short: "Dereferences and updates values from a vault to a given file with vault references",
			Run: func(cmd *cobra.Command, args []string) {
				envSecretKey, err := session.GetSecretKey()
				if err != nil {
					utils.ExitOnError(err)
				}
				vaultFile, err := cmd.Flags().GetString(vaultFileFlag.Name)
				if err != nil {
					utils.ExitOnError(err)
				}
				file, err := cmd.Flags().GetString(vaultRefFileFlag.Name)
				if err != nil {
					utils.ExitOnError(err)
				}
				previewOnlyMode, _ := cmd.Flags().GetBool(secretSubstitutionPreviewOnlyFlag.Name)
				vault, err := vaults.Get(vaultFile)
				if err != nil {
					utils.ExitOnError(err)
				}
				err = vault.Unlock(envSecretKey)
				if err != nil {
					utils.ExitOnError(err)
				}
				result, err := vault.DeRef(file, previewOnlyMode)
				if err != nil {
					utils.ExitOnError(err)
				}
				if previewOnlyMode {
					if len(result) > 0 && result[len(result)-1] == '\n' {
						fmt.Print(result)
					} else {
						fmt.Println(result)
					}
				} else {
					fmt.Println("Dereferenced", color.GreenString(file), "with the vault", color.GreenString(vaultFile))
				}
				utils.SafeExit()
			},
		}
		vaultDerefCmd.Flags().StringP(vaultFileFlag.Name, vaultFileFlag.Shorthand, "", vaultFileFlag.Usage)
		vaultDerefCmd.Flags().StringP(vaultRefFileFlag.Name, vaultRefFileFlag.Shorthand, "", vaultRefFileFlag.Usage)
		vaultDerefCmd.Flags().BoolP(secretSubstitutionPreviewOnlyFlag.Name, secretSubstitutionPreviewOnlyFlag.Shorthand, false, secretSubstitutionPreviewOnlyFlag.Usage)
		vaultDerefCmd.MarkFlagRequired(vaultRefFileFlag.Name)
	}
	return vaultDerefCmd
}

package cmdvault

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/input"
)

func vaultPutCommand() *cobra.Command {
	if vaultPutCmd != nil {
		return vaultPutCmd
	}
	vaultPutCmd = &cobra.Command{
		Use:     "put",
		Aliases: []string{"add", "set", "create", "load", "import"},
		Short:   "Adds secret to the vault",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			secretName := cmd.Flag(secretNameFlag.Name).Value.String()
			secretValue := cmd.Flag(secretValueFlag.Name).Value.String()
			importFile := cmd.Flag(vaultImportFileFlag.Name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				utils.ExitOnError(err)
			}
			forceUpdate, _ := cmd.Flags().GetBool(secretForceUpdateFlag.Name)
			if secretName != "" {
				if !forceUpdate && vault.Exists(secretName) {
					confirmation, err := input.GetVisibleInput("Secret already exists. Do you wish to overwrite it? (y/n): ")
					if err != nil {
						utils.ExitOnError(err)
					}
					if confirmation != "y" {
						fmt.Println(color.YellowString("Operation aborted"))
						utils.SafeExit()
					}
				}
				var secret []byte
				if secretValue == "" {
					secret, err = input.GetMultiLineHiddenInput("Enter the secret value for " + secretName + ": ")
					if err != nil {
						utils.ExitOnError(err)
					}
				} else {
					secret = []byte(secretValue)
				}
				err = vault.Put(secretName, secret, true)
				if err != nil {
					utils.ExitOnError(err)
				}
				fmt.Println("Updated secret: ", color.GreenString(secretName), " to vault: ", color.GreenString(vaultFile))
			}
			if importFile != "" || secretName == "" {
				var importData []byte
				if importFile == "" {
					importData, err = input.GetMultiLineHiddenInput("Enter the YAML/JSON data to be imported: ")
				} else {
					importData, err = os.ReadFile(importFile)
				}
				if err != nil {
					utils.ExitOnError(err)
				}
				if err = vault.Import(importData, forceUpdate, true); err != nil {
					utils.ExitOnError(err)
				}
				fmt.Printf("Successfully imported secrets from %s into the vault %s\n", color.GreenString(importFile), color.GreenString(vaultFile))
			}
			utils.SafeExit()
		},
	}
	vaultPutCmd.Flags().StringP(secretNameFlag.Name, secretNameFlag.Shorthand, "", secretNameFlag.Usage)
	vaultPutCmd.Flags().StringP(secretValueFlag.Name, secretValueFlag.Shorthand, "", secretValueFlag.Usage)
	vaultPutCmd.Flags().StringP(vaultImportFileFlag.Name, vaultImportFileFlag.Shorthand, "", vaultImportFileFlag.Usage)
	vaultPutCmd.Flags().Bool(secretForceUpdateFlag.Name, false, secretForceUpdateFlag.Usage)
	return vaultPutCmd
}

package cmdvault

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/input"
	"slv.sh/slv/internal/core/vaults"
)

func vaultPutCommand() *cobra.Command {
	if vaultPutCmd == nil {
		vaultPutCmd = &cobra.Command{
			Use:     "put",
			Aliases: []string{"add", "set", "create", "load", "import"},
			Short:   "Adds, updates or imports secrets to the vault",
			Run: func(cmd *cobra.Command, args []string) {
				vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
				itemName := cmd.Flag(itemNameFlag.Name).Value.String()
				itemValue := cmd.Flag(itemValueFlag.Name).Value.String()
				if itemValue == "" {
					itemValue = cmd.Flag(deprecatedSecretFlag.Name).Value.String()
				}
				importFile := cmd.Flag(vaultImportFileFlag.Name).Value.String()
				plaintextValue, _ := cmd.Flags().GetBool(plaintextValueFlag.Name)
				vault, err := vaults.Get(vaultFile)
				if err != nil {
					utils.ExitOnError(err)
				}
				forceUpdate, _ := cmd.Flags().GetBool(secretForceUpdateFlag.Name)
				if itemName != "" {
					if !forceUpdate && vault.ItemExists(itemName) {
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
					switch itemValue {
					case "":
						if secret, err = input.GetMultiLineHiddenInput("Enter the secret value for " + itemName + ": "); err != nil {
							utils.ExitOnError(err)
						}
					case "-":
						if secret, err = input.ReadBufferFromStdin(""); err != nil {
							utils.ExitOnError(err)
						}
					default:
						secret = []byte(itemValue)
					}
					if err = vault.Put(itemName, secret, !plaintextValue); err != nil {
						utils.ExitOnError(err)
					}
					fmt.Printf("Successfully added/updated secret %s into the vault %s\n", color.GreenString(itemName), color.GreenString(vaultFile))
				}
				if importFile != "" || itemName == "" {
					var importData []byte
					if importFile == "" {
						importData, err = input.GetMultiLineHiddenInput("Enter the YAML/JSON/ENV format data to be imported: ")
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
		vaultPutCmd.Flags().StringP(itemNameFlag.Name, itemNameFlag.Shorthand, "", itemNameFlag.Usage)
		if err := vaultPutCmd.RegisterFlagCompletionFunc(itemNameFlag.Name, vaultItemNameCompletion); err != nil {
			utils.ExitOnError(err)
		}
		vaultPutCmd.Flags().StringP(itemValueFlag.Name, itemValueFlag.Shorthand, "", itemValueFlag.Usage)
		vaultPutCmd.Flags().StringP(deprecatedSecretFlag.Name, deprecatedSecretFlag.Shorthand, "", deprecatedSecretFlag.Usage)
		vaultPutCmd.Flags().MarkDeprecated(deprecatedSecretFlag.Name, "use --"+itemValueFlag.Name+" instead")
		vaultPutCmd.Flags().StringP(vaultImportFileFlag.Name, vaultImportFileFlag.Shorthand, "", vaultImportFileFlag.Usage)
		vaultPutCmd.Flags().Bool(plaintextValueFlag.Name, false, plaintextValueFlag.Usage)
		vaultPutCmd.Flags().Bool(secretForceUpdateFlag.Name, false, secretForceUpdateFlag.Usage)
	}
	return vaultPutCmd
}

package cmdvault

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/vaults"
)

func refineType(data any, foundType string) string {
	switch data.(type) {
	case string:
		return ""
	default:
		return foundType
	}
}

func detectType(refFile string) (string, error) {
	var data any
	content, err := os.ReadFile(refFile)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(content, &data); err == nil {
		return refineType(data, "json"), nil
	}
	if err = yaml.Unmarshal(content, &data); err == nil {
		return refineType(data, "yaml"), nil
	}
	return "", nil
}

func vaultRefCommand() *cobra.Command {
	if vaultRefCmd == nil {
		vaultRefCmd = &cobra.Command{
			Use:     "ref",
			Aliases: []string{"reference"},
			Short:   "References and updates secrets to a vault from a given yaml, json or the whole file content",
			Run: func(cmd *cobra.Command, args []string) {
				vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
				vault, err := vaults.Get(vaultFile)
				if err != nil {
					utils.ExitOnError(err)
				}
				refFile := cmd.Flag(vaultRefFileFlag.Name).Value.String()
				secretNamePrefix := cmd.Flag(itemNameFlag.Name).Value.String()
				refType, err := detectType(refFile)
				if err != nil {
					utils.ExitOnError(err)
				}
				previewMode, _ := cmd.Flags().GetBool(secretSubstitutionPreviewOnlyFlag.Name)
				forceUpdate, _ := cmd.Flags().GetBool(secretForceUpdateFlag.Name)
				if secretNamePrefix == "" && refType == "" {
					utils.ExitOnErrorWithMessage("please provide --" + itemNameFlag.Name + " since the file is neither json nor yaml")
				}
				result, conflicting, err := vault.Ref(refType, refFile, secretNamePrefix, forceUpdate, true, previewMode)
				if conflicting {
					utils.ExitOnErrorWithMessage("conflict found. please use the --" + itemNameFlag.Name + " flag to set a different name or --" + secretForceUpdateFlag.Name + " flag to overwrite them.")
				} else if err != nil {
					utils.ExitOnError(err)
				}
				if previewMode {
					fmt.Println(result)
				} else {
					if refType == "" {
						fmt.Println("Auto referenced", color.GreenString(refFile), "with vault", color.GreenString(vaultFile))
					} else {
						fmt.Println("Auto referenced", color.GreenString(refFile), "("+strings.ToUpper(refType)+")", "with vault", color.GreenString(vaultFile))
					}
				}
				utils.SafeExit()
			},
		}
		vaultRefCmd.Flags().StringP(vaultRefFileFlag.Name, vaultRefFileFlag.Shorthand, "", vaultRefFileFlag.Usage)
		vaultRefCmd.Flags().StringP(itemNameFlag.Name, itemNameFlag.Shorthand, "", itemNameFlag.Usage)
		vaultRefCmd.Flags().BoolP(secretSubstitutionPreviewOnlyFlag.Name, secretSubstitutionPreviewOnlyFlag.Shorthand, false, secretSubstitutionPreviewOnlyFlag.Usage)
		vaultRefCmd.Flags().BoolP(secretForceUpdateFlag.Name, secretForceUpdateFlag.Shorthand, false, secretForceUpdateFlag.Usage)
		vaultRefCmd.MarkFlagRequired(vaultRefFileFlag.Name)
	}
	return vaultRefCmd
}

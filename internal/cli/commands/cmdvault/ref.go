package cmdvault

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
)

func vaultRefCommand() *cobra.Command {
	if vaultRefCmd != nil {
		return vaultRefCmd
	}
	vaultRefCmd = &cobra.Command{
		Use:     "ref",
		Aliases: []string{"reference"},
		Short:   "References and updates secrets to a vault from a given yaml or json file",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				utils.ExitOnError(err)
			}
			refFile := cmd.Flag(vaultRefFileFlag.Name).Value.String()
			secretNamePrefix := cmd.Flag(secretNameFlag.Name).Value.String()
			refType := strings.ToLower(cmd.Flag(vaultRefTypeFlag.Name).Value.String())
			previewOnly, _ := cmd.Flags().GetBool(secretRefPreviewOnlyFlag.Name)
			forceUpdate, _ := cmd.Flags().GetBool(secretForceUpdateFlag.Name)
			if secretNamePrefix == "" && refType == "" {
				utils.ExitOnErrorWithMessage("please provide at least one of --" + secretNameFlag.Name + " or --" + vaultRefTypeFlag.Name + " flag")
			}
			if refType != "" && refType != "yaml" {
				utils.ExitOnErrorWithMessage("only yaml auto reference is supported at the moment")
			}
			result, conflicting, err := vault.RefSecrets(refType, refFile, secretNamePrefix, forceUpdate, previewOnly)
			if conflicting {
				utils.ExitOnErrorWithMessage("conflict found. please use the --" + secretNameFlag.Name + " flag to set a different name or --" + secretForceUpdateFlag.Name + " flag to overwrite them.")
			} else if err != nil {
				utils.ExitOnError(err)
			}
			if previewOnly {
				fmt.Println(result)
			} else {
				fmt.Println("Auto referenced", color.GreenString(refFile), "with vault", color.GreenString(vaultFile))
			}
			utils.SafeExit()
		},
	}
	vaultRefCmd.Flags().StringP(vaultRefFileFlag.Name, vaultRefFileFlag.Shorthand, "", vaultRefFileFlag.Usage)
	vaultRefCmd.Flags().StringP(secretNameFlag.Name, secretNameFlag.Shorthand, "", secretNameFlag.Usage)
	vaultRefCmd.Flags().StringP(vaultRefTypeFlag.Name, vaultRefTypeFlag.Shorthand, "", vaultRefTypeFlag.Usage)
	vaultRefCmd.Flags().BoolP(secretRefPreviewOnlyFlag.Name, secretRefPreviewOnlyFlag.Shorthand, false, secretRefPreviewOnlyFlag.Usage)
	vaultRefCmd.Flags().BoolP(secretForceUpdateFlag.Name, secretForceUpdateFlag.Shorthand, false, secretForceUpdateFlag.Usage)
	vaultRefCmd.MarkFlagRequired(vaultRefFileFlag.Name)
	return vaultRefCmd
}

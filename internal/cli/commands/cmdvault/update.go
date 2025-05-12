package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/vaults"
)

func vaultUpdateCommand() *cobra.Command {
	if vaultUpdateCmd == nil {
		vaultUpdateCmd = &cobra.Command{
			Use:   "update",
			Short: "Update attributes of an existing SLV vault",
			Run: func(cmd *cobra.Command, args []string) {
				vaultFilePath := cmd.Flag(vaultFileFlag.Name).Value.String()
				name := cmd.Flag(vaultNameFlag.Name).Value.String()
				namespace := cmd.Flag(vaultK8sNamespaceFlag.Name).Value.String()
				vault, err := vaults.Get(vaultFilePath)
				if err != nil {
					utils.ExitOnError(err)
				}
				if err = vault.Update(name, namespace, nil); err != nil {
					utils.ExitOnError(err)
				}
				fmt.Printf("Vault %s transformed to K8s resource %s\n", color.GreenString(vaultFilePath), color.GreenString(name))
			},
		}
		vaultUpdateCmd.Flags().StringP(vaultNameFlag.Name, vaultNameFlag.Shorthand, "", vaultNameFlag.Usage)
		vaultUpdateCmd.Flags().StringP(vaultK8sNamespaceFlag.Name, vaultK8sNamespaceFlag.Shorthand, "", vaultK8sNamespaceFlag.Usage)
	}
	return vaultUpdateCmd
}

package cmdvault

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/input"
	"slv.sh/slv/internal/core/vaults"
)

func vaultUpdateCommand() *cobra.Command {
	if vaultUpdateCmd == nil {
		vaultUpdateCmd = &cobra.Command{
			Use:   "update",
			Short: "Update attributes of a vault",
			Run: func(cmd *cobra.Command, args []string) {
				vaultFilePath := cmd.Flag(vaultFileFlag.Name).Value.String()
				name := cmd.Flag(vaultNameFlag.Name).Value.String()
				namespace := cmd.Flag(vaultK8sNamespaceFlag.Name).Value.String()
				vault, err := vaults.Get(vaultFilePath)
				if err != nil {
					utils.ExitOnError(err)
				}
				k8sSecret := cmd.Flag(vaultK8sSecretFlag.Name).Value.String()
				var data []byte
				if k8sSecret != "" {
					if strings.HasSuffix(k8sSecret, ".yaml") || strings.HasSuffix(k8sSecret, ".yml") ||
						strings.HasSuffix(k8sSecret, ".json") {
						data, err = os.ReadFile(k8sSecret)
					} else if k8sSecret == "-" {
						data, err = input.ReadBufferFromStdin("Input the k8s secret object as yaml/json: ")
					} else {
						utils.ExitOnErrorWithMessage("invalid k8s secret resource file")
					}
					if err != nil {
						utils.ExitOnError(err)
					}
				}
				secretType := cmd.Flag(vaultK8sSecretTypeFlag.Name).Value.String()
				if err = vault.Update(name, namespace, secretType, data); err != nil {
					utils.ExitOnError(err)
				}
				fmt.Printf("Vault %s transformed to K8s resource %s\n", color.GreenString(vaultFilePath), color.GreenString(name))
			},
		}
		vaultUpdateCmd.Flags().StringP(vaultNameFlag.Name, vaultNameFlag.Shorthand, "", vaultNameFlag.Usage)
		vaultUpdateCmd.Flags().StringP(vaultK8sNamespaceFlag.Name, vaultK8sNamespaceFlag.Shorthand, "", vaultK8sNamespaceFlag.Usage)
		vaultUpdateCmd.Flags().StringP(vaultK8sSecretFlag.Name, vaultK8sSecretFlag.Shorthand, "", vaultK8sSecretFlag.Usage)
		vaultUpdateCmd.Flags().StringP(vaultK8sSecretTypeFlag.Name, vaultK8sSecretTypeFlag.Shorthand, "", vaultK8sSecretTypeFlag.Usage)
	}
	return vaultUpdateCmd
}

package cmdvault

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/crypto"
	"oss.amagi.com/slv/internal/core/input"
	"oss.amagi.com/slv/internal/core/vaults"
)

func newK8sVault(filePath, k8sName, k8sNamespace, k8sSecret string, hash, pq bool, publicKeys ...*crypto.PublicKey) (*vaults.Vault, error) {
	var data []byte
	if k8sSecret != "" {
		var err error
		if strings.HasSuffix(k8sSecret, ".yaml") || strings.HasSuffix(k8sSecret, ".yml") ||
			strings.HasSuffix(k8sSecret, ".json") {
			data, err = os.ReadFile(k8sSecret)
		} else if k8sSecret == "-" {
			data, err = input.ReadBufferFromStdin("Input the k8s secret object as yaml/json: ")
		} else {
			return nil, fmt.Errorf("invalid k8s secret resource file")
		}
		if err != nil {
			return nil, err
		}
	}
	return vaults.New(filePath, k8sName, k8sNamespace, data, hash, pq, publicKeys...)
}

func vaultToK8sCommand() *cobra.Command {
	if vaultToK8sCmd != nil {
		return vaultToK8sCmd
	}
	vaultToK8sCmd = &cobra.Command{
		Use:     "tok8s",
		Aliases: []string{"k8s", "tok8slv"},
		Short:   "Transform an existing SLV vault file to a K8s compatible one",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFilePath := cmd.Flag(vaultFileFlag.Name).Value.String()
			name := cmd.Flag(vaultK8sNameFlag.Name).Value.String()
			namespace := cmd.Flag(vaultK8sNamespaceFlag.Name).Value.String()
			vault, err := getVault(vaultFilePath)
			if err != nil {
				utils.ExitOnError(err)
			}
			if err = vault.ToK8s(name, namespace, nil); err != nil {
				utils.ExitOnError(err)
			}
			fmt.Printf("Vault %s transformed to K8s resource %s\n", color.GreenString(vaultFilePath), color.GreenString(name))
		},
	}
	vaultToK8sCmd.Flags().StringP(vaultK8sNameFlag.Name, vaultK8sNameFlag.Shorthand, "", vaultK8sNameFlag.Usage)
	vaultToK8sCmd.Flags().StringP(vaultK8sNamespaceFlag.Name, vaultK8sNamespaceFlag.Shorthand, "", vaultK8sNamespaceFlag.Usage)
	vaultToK8sCmd.MarkFlagRequired(vaultK8sNameFlag.Name)
	return vaultToK8sCmd
}

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

func newK8sVault(filePath, k8sNameOrSecretFile string, hashLength uint8, pq bool, rootPublicKey *crypto.PublicKey, publicKeys ...*crypto.PublicKey) (*vaults.Vault, error) {
	if strings.HasSuffix(k8sNameOrSecretFile, ".yaml") || strings.HasSuffix(k8sNameOrSecretFile, ".yml") ||
		strings.HasSuffix(k8sNameOrSecretFile, ".json") || k8sNameOrSecretFile == "-" {
		var data []byte
		var err error
		if k8sNameOrSecretFile == "-" {
			data, err = input.ReadBufferFromStdin("Input the k8s secret object as yaml/json: ")
		} else {
			data, err = os.ReadFile(k8sNameOrSecretFile)
		}
		if err != nil {
			return nil, err
		}
		return vaults.New(filePath, "", data, hashLength, pq, rootPublicKey, publicKeys...)
	} else {
		return vaults.New(filePath, k8sNameOrSecretFile, nil, hashLength, pq, rootPublicKey, publicKeys...)
	}
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
			k8sResourceName := cmd.Flag(vaultK8sNameFlag.Name).Value.String()
			vault, err := getVault(vaultFilePath)
			if err != nil {
				utils.ExitOnError(err)
			}
			if err = vault.ToK8s(k8sResourceName, nil); err != nil {
				utils.ExitOnError(err)
			}
			fmt.Printf("Vault %s transformed to K8s resource %s\n", color.GreenString(vaultFilePath), color.GreenString(k8sResourceName))
		},
	}
	vaultToK8sCmd.Flags().StringP(vaultK8sNameFlag.Name, vaultK8sNameFlag.Shorthand, "", vaultK8sNameFlag.Usage)
	vaultToK8sCmd.MarkFlagRequired(vaultK8sNameFlag.Name)
	return vaultToK8sCmd
}

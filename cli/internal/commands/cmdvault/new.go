package cmdvault

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/cli/internal/commands/cmdenv"
	"oss.amagi.com/slv/cli/internal/commands/utils"
	"oss.amagi.com/slv/core/commons"
	"oss.amagi.com/slv/core/crypto"
	"oss.amagi.com/slv/core/input"
	"oss.amagi.com/slv/core/vaults"
)

func newK8sVault(filePath, k8sValue string, hashLength uint8, pq bool, rootPublicKey *crypto.PublicKey, publicKeys ...*crypto.PublicKey) (*vaults.Vault, error) {
	k8slvName := k8sValue
	var secretDataMap map[string]string
	if strings.HasSuffix(k8sValue, ".yaml") || strings.HasSuffix(k8sValue, ".yml") || strings.HasSuffix(k8sValue, ".json") || k8sValue == "-" {
		var data []byte
		var err error
		if k8sValue == "-" {
			data, err = input.ReadBufferFromStdin("Input the k8s secret object: ")
		} else {
			data, err = os.ReadFile(k8sValue)
		}
		if err != nil {
			return nil, err
		}
		secret, err := k8sSecretFromData(data)
		if err != nil {
			return nil, err
		}
		k8slvName = secret.Metadata.Name
		secretDataMap = secret.Data
	}
	vault, err := vaults.New(filePath, k8sVaultField, hashLength, pq, rootPublicKey, publicKeys...)
	if err != nil {
		return nil, err
	}
	if len(secretDataMap) > 0 {
		for key, value := range secretDataMap {
			decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(value))
			secretValue, err := io.ReadAll(decoder)
			if err != nil {
				return nil, err
			}
			if err = vault.PutSecret(key, secretValue); err != nil {
				return nil, err
			}
		}
	}
	var obj map[string]interface{}
	if err := commons.ReadFromYAML(filePath, &obj); err != nil {
		return nil, err
	}
	obj["apiVersion"] = k8sApiVersion
	obj["kind"] = k8sKind
	obj["metadata"] = map[string]interface{}{
		"name": k8slvName,
	}
	return vault, commons.WriteToYAML(filePath, "", obj)
}

func vaultNewCommand() *cobra.Command {
	if vaultNewCmd != nil {
		return vaultNewCmd
	}
	vaultNewCmd = &cobra.Command{
		Use:   "new",
		Short: "Creates a new vault",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			publicKeyStrings, err := cmd.Flags().GetStringSlice(vaultAccessPublicKeysFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			queries, err := cmd.Flags().GetStringSlice(cmdenv.EnvSearchFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			shareWithSelf, _ := cmd.Flags().GetBool(cmdenv.EnvSelfFlag.Name)
			publicKeys, rootPublicKey, err := getPublicKeys(publicKeyStrings, queries, shareWithSelf)
			if err != nil {
				utils.ExitOnError(err)
			}
			enableHash, _ := cmd.Flags().GetBool(vaultEnableHashingFlag.Name)
			var hashLength uint8 = 0
			if enableHash {
				hashLength = 4
			}
			pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
			k8sName := cmd.Flag(vaultK8sFlag.Name).Value.String()
			if k8sName == "" {
				_, err = vaults.New(vaultFile, "", hashLength, pq, rootPublicKey, publicKeys...)
			} else {
				_, err = newK8sVault(vaultFile, k8sName, hashLength, pq, rootPublicKey, publicKeys...)
			}
			if err != nil {
				utils.ExitOnError(err)
			}
			fmt.Println("Created vault:", color.GreenString(vaultFile))
			utils.SafeExit()
		},
	}
	vaultNewCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.Name, vaultAccessPublicKeysFlag.Shorthand, []string{}, vaultAccessPublicKeysFlag.Usage)
	vaultNewCmd.Flags().StringSliceP(cmdenv.EnvSearchFlag.Name, cmdenv.EnvSearchFlag.Shorthand, []string{}, cmdenv.EnvSearchFlag.Usage)
	vaultNewCmd.Flags().BoolP(cmdenv.EnvSelfFlag.Name, cmdenv.EnvSelfFlag.Shorthand, false, cmdenv.EnvSelfFlag.Usage)
	vaultNewCmd.Flags().StringP(vaultK8sFlag.Name, vaultK8sFlag.Shorthand, "", vaultK8sFlag.Usage)
	vaultNewCmd.Flags().BoolP(vaultEnableHashingFlag.Name, vaultEnableHashingFlag.Shorthand, false, vaultEnableHashingFlag.Usage)
	vaultNewCmd.Flags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage)
	return vaultNewCmd
}

package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"savesecrets.org/slv/cli/internal/commands/cmdenv"
	"savesecrets.org/slv/cli/internal/commands/utils"
	"savesecrets.org/slv/core/commons"
	"savesecrets.org/slv/core/crypto"
	"savesecrets.org/slv/core/vaults"
)

func newK8sVault(filePath, name string, hashLength uint8, rootPublicKey *crypto.PublicKey, publicKeys ...*crypto.PublicKey) (*vaults.Vault, error) {
	vault, err := vaults.New(filePath, k8sVaultField, hashLength, rootPublicKey, publicKeys...)
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	if err := commons.ReadFromYAML(filePath, &obj); err != nil {
		return nil, err
	}
	obj["apiVersion"] = k8sApiVersion
	obj["kind"] = k8sKind
	obj["metadata"] = map[string]interface{}{
		"name": name,
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
			query := cmd.Flag(cmdenv.EnvSearchFlag.Name).Value.String()
			shareWithSelf, _ := cmd.Flags().GetBool(cmdenv.EnvSelfFlag.Name)
			publicKeys, rootPublicKey, err := getPublicKeys(publicKeyStrings, query, shareWithSelf)
			if err != nil {
				utils.ExitOnError(err)
			}
			enableHash, _ := cmd.Flags().GetBool(vaultEnableHashingFlag.Name)
			var hashLength uint8 = 0
			if enableHash {
				hashLength = 4
			}
			k8sName := cmd.Flag(vaultK8sFlag.Name).Value.String()
			if k8sName == "" {
				_, err = vaults.New(vaultFile, "", hashLength, rootPublicKey, publicKeys...)
			} else {
				_, err = newK8sVault(vaultFile, k8sName, hashLength, rootPublicKey, publicKeys...)
			}
			if err != nil {
				utils.ExitOnError(err)
			}
			fmt.Println("Created vault:", color.GreenString(vaultFile))
			utils.SafeExit()
		},
	}
	vaultNewCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.Name, vaultAccessPublicKeysFlag.Shorthand, []string{}, vaultAccessPublicKeysFlag.Usage)
	vaultNewCmd.Flags().StringP(cmdenv.EnvSearchFlag.Name, cmdenv.EnvSearchFlag.Shorthand, "", cmdenv.EnvSearchFlag.Usage)
	vaultNewCmd.Flags().BoolP(cmdenv.EnvSelfFlag.Name, cmdenv.EnvSelfFlag.Shorthand, false, cmdenv.EnvSelfFlag.Usage)
	vaultNewCmd.Flags().StringP(vaultK8sFlag.Name, vaultK8sFlag.Shorthand, "", vaultK8sFlag.Usage)
	vaultNewCmd.Flags().BoolP(vaultEnableHashingFlag.Name, vaultEnableHashingFlag.Shorthand, false, vaultEnableHashingFlag.Usage)
	return vaultNewCmd
}

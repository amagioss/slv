package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/cmdenv"
	"oss.amagi.com/slv/internal/cli/commands/utils"
)

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
			shareWithK8s, _ := cmd.Flags().GetBool(vaultAccessK8sFlag.Name)
			pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
			publicKeys, rootPublicKey, err := getPublicKeys(publicKeyStrings, queries, shareWithSelf, shareWithK8s, pq)
			if err != nil {
				utils.ExitOnError(err)
			}
			enableHash, _ := cmd.Flags().GetBool(vaultEnableHashingFlag.Name)
			k8sName := cmd.Flag(vaultK8sFlag.Name).Value.String()
			if _, err = newK8sVault(vaultFile, k8sName, enableHash, pq, rootPublicKey, publicKeys...); err != nil {
				utils.ExitOnError(err)
			}
			fmt.Println("Created vault:", color.GreenString(vaultFile))
			utils.SafeExit()
		},
	}
	vaultNewCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.Name, vaultAccessPublicKeysFlag.Shorthand, []string{}, vaultAccessPublicKeysFlag.Usage)
	vaultNewCmd.Flags().StringSliceP(cmdenv.EnvSearchFlag.Name, cmdenv.EnvSearchFlag.Shorthand, []string{}, cmdenv.EnvSearchFlag.Usage)
	vaultNewCmd.Flags().BoolP(cmdenv.EnvSelfFlag.Name, cmdenv.EnvSelfFlag.Shorthand, false, cmdenv.EnvSelfFlag.Usage)
	vaultNewCmd.Flags().BoolP(vaultAccessK8sFlag.Name, vaultAccessK8sFlag.Shorthand, false, vaultAccessK8sFlag.Usage)
	vaultNewCmd.Flags().StringP(vaultK8sFlag.Name, vaultK8sFlag.Shorthand, "", vaultK8sFlag.Usage)
	vaultNewCmd.Flags().BoolP(vaultEnableHashingFlag.Name, vaultEnableHashingFlag.Shorthand, false, vaultEnableHashingFlag.Usage)
	vaultNewCmd.Flags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage)
	return vaultNewCmd
}

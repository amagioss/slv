package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/cmdenv"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/vaults"
)

func vaultNewCommand() *cobra.Command {
	if vaultNewCmd == nil {
		vaultNewCmd = &cobra.Command{
			Use:     "new",
			Aliases: []string{"create"},
			Short:   "Creates a new vault",
			Run: func(cmd *cobra.Command, args []string) {
				vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
				pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
				publicKeys, err := cmdenv.GetPublicKeys(cmd, true, pq)
				if err != nil {
					utils.ExitOnError(err)
				}
				enableHash, _ := cmd.Flags().GetBool(vaultEnableHashingFlag.Name)
				name := cmd.Flag(vaultNameFlag.Name).Value.String()
				k8sNamespace := cmd.Flag(vaultK8sNamespaceFlag.Name).Value.String()
				if _, err = vaults.New(vaultFile, name, k8sNamespace, enableHash, pq, publicKeys...); err != nil {
					utils.ExitOnError(err)
				}
				fmt.Println("Created vault:", color.GreenString(vaultFile))
				utils.SafeExit()
			},
		}
		vaultNewCmd.Flags().StringSliceP(cmdenv.EnvPublicKeysFlag.Name, cmdenv.EnvPublicKeysFlag.Shorthand, []string{}, cmdenv.EnvPublicKeysFlag.Usage)
		vaultNewCmd.Flags().StringSliceP(cmdenv.EnvSearchFlag.Name, cmdenv.EnvSearchFlag.Shorthand, []string{}, cmdenv.EnvSearchFlag.Usage)
		if err := vaultNewCmd.RegisterFlagCompletionFunc(cmdenv.EnvSearchFlag.Name, cmdenv.EnvSearchCompletion); err != nil {
			utils.ExitOnError(err)
		}
		vaultNewCmd.Flags().BoolP(cmdenv.EnvSelfFlag.Name, cmdenv.EnvSelfFlag.Shorthand, false, cmdenv.EnvSelfFlag.Usage)
		vaultNewCmd.Flags().BoolP(cmdenv.EnvK8sFlag.Name, cmdenv.EnvK8sFlag.Shorthand, false, cmdenv.EnvK8sFlag.Usage)
		vaultNewCmd.Flags().StringP(vaultNameFlag.Name, vaultNameFlag.Shorthand, "", vaultNameFlag.Usage)
		vaultNewCmd.Flags().StringP(vaultK8sNamespaceFlag.Name, vaultK8sNamespaceFlag.Shorthand, "", vaultK8sNamespaceFlag.Usage)
		vaultNewCmd.Flags().BoolP(vaultEnableHashingFlag.Name, vaultEnableHashingFlag.Shorthand, false, vaultEnableHashingFlag.Usage)
		vaultNewCmd.Flags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage)
	}
	return vaultNewCmd
}

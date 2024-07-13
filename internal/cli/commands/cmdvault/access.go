package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/cmdenv"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/crypto"
	"oss.amagi.com/slv/internal/core/secretkey"
)

func vaultAccessCommand() *cobra.Command {
	if vaultAccessCmd != nil {
		return vaultAccessCmd
	}
	vaultAccessCmd = &cobra.Command{
		Use:     "access",
		Aliases: []string{"rights", "privilege", "permission", "permissions"},
		Short:   "Managing access to a vault",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	vaultAccessCmd.PersistentFlags().StringSliceP(vaultAccessPublicKeysFlag.Name, vaultAccessPublicKeysFlag.Shorthand, []string{}, vaultAccessPublicKeysFlag.Usage)
	vaultAccessCmd.PersistentFlags().StringSliceP(cmdenv.EnvSearchFlag.Name, cmdenv.EnvSearchFlag.Shorthand, []string{}, cmdenv.EnvSearchFlag.Usage)
	vaultAccessCmd.PersistentFlags().BoolP(cmdenv.EnvSelfFlag.Name, cmdenv.EnvSelfFlag.Shorthand, false, cmdenv.EnvSelfFlag.Usage)
	vaultAccessCmd.PersistentFlags().BoolP(vaultAccessK8sFlag.Name, vaultAccessK8sFlag.Shorthand, false, vaultAccessK8sFlag.Usage)
	vaultAccessCmd.PersistentFlags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage+" (used with k8s environment)")
	vaultAccessCmd.AddCommand(vaultAccessAddCommand())
	vaultAccessCmd.AddCommand(vaultAccessRemoveCommand())
	return vaultAccessCmd
}

func vaultAccessAddCommand() *cobra.Command {
	if vaultAccessAddCmd != nil {
		return vaultAccessAddCmd
	}
	vaultAccessAddCmd = &cobra.Command{
		Use:     "add",
		Aliases: []string{"allow", "grant", "share"},
		Short:   "Adds access to a vault for the given environments/public keys",
		Run: func(cmd *cobra.Command, args []string) {
			envSecretKey, err := secretkey.Get()
			if err != nil {
				utils.ExitOnError(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			publicKeyStrings, err := cmd.Flags().GetStringSlice(vaultAccessPublicKeysFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			queries, err := cmd.Flags().GetStringSlice(cmdenv.EnvSearchFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			selfEnv, _ := cmd.Flags().GetBool(cmdenv.EnvSelfFlag.Name)
			k8sEnv, _ := cmd.Flags().GetBool(vaultAccessK8sFlag.Name)
			k8sPQ, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
			publicKeys, _, err := getPublicKeys(publicKeyStrings, queries, selfEnv, k8sEnv, k8sPQ)
			if err != nil {
				utils.ExitOnError(err)
			}
			vault, err := getVault(vaultFile)
			if err == nil {
				err = vault.Unlock(*envSecretKey)
				if err == nil {
					for _, publicKey := range publicKeys {
						if _, err = vault.Share(publicKey); err != nil {
							break
						}
					}
					if err == nil {
						fmt.Println("Shared vault:", color.GreenString(vaultFile))
						utils.SafeExit()
					}
				}
			}
			utils.ExitOnError(err)
		},
	}
	return vaultAccessAddCmd
}

func vaultAccessRemoveCommand() *cobra.Command {
	if vaultAccessRemoveCmd != nil {
		return vaultAccessRemoveCmd
	}
	vaultAccessRemoveCmd = &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "deny", "revoke", "restrict", "delete", "del"},
		Short:   "Remove access to a vault for the given environments/public keys",
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
			selfEnv, _ := cmd.Flags().GetBool(cmdenv.EnvSelfFlag.Name)
			k8sEnv, _ := cmd.Flags().GetBool(vaultAccessK8sFlag.Name)
			k8sPQ, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
			publicKeys, _, err := getPublicKeys(publicKeyStrings, queries, selfEnv, k8sEnv, k8sPQ)
			if err != nil {
				utils.ExitOnError(err)
			}
			vault, err := getVault(vaultFile)
			pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
			if err == nil {
				var envSecretKey *crypto.SecretKey
				if envSecretKey, err = secretkey.Get(); err == nil {
					err = vault.Unlock(*envSecretKey)
				}
				if err == nil {
					if err = vault.Revoke(publicKeys, pq); err == nil {
						fmt.Println("Shared vault:", color.GreenString(vaultFile))
						utils.SafeExit()
					}
				}
			}
			utils.ExitOnError(err)
		},
	}
	return vaultAccessRemoveCmd
}

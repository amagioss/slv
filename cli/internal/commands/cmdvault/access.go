package cmdvault

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"savesecrets.org/slv"
	"savesecrets.org/slv/cli/internal/commands/cmdenv"
	"savesecrets.org/slv/cli/internal/commands/utils"
	"savesecrets.org/slv/core/crypto"
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
			envSecretKey, err := slv.GetSecretKey()
			if err != nil {
				utils.ExitOnError(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			publicKeyStrings, err := cmd.Flags().GetStringSlice(vaultAccessPublicKeysFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			query := cmd.Flag(cmdenv.EnvSearchFlag.Name).Value.String()
			selfEnv, _ := cmd.Flags().GetBool(cmdenv.EnvSelfFlag.Name)
			publicKeys, _, err := getPublicKeys(publicKeyStrings, query, selfEnv)
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
	vaultAccessAddCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.Name, vaultAccessPublicKeysFlag.Shorthand, []string{}, vaultAccessPublicKeysFlag.Usage)
	vaultAccessAddCmd.Flags().StringP(cmdenv.EnvSearchFlag.Name, cmdenv.EnvSearchFlag.Shorthand, "", cmdenv.EnvSearchFlag.Usage)
	vaultAccessAddCmd.Flags().BoolP(cmdenv.EnvSelfFlag.Name, cmdenv.EnvSelfFlag.Shorthand, false, cmdenv.EnvSelfFlag.Usage)
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
			query := cmd.Flag(cmdenv.EnvSearchFlag.Name).Value.String()
			selfEnv, _ := cmd.Flags().GetBool(cmdenv.EnvSelfFlag.Name)
			publicKeys, _, err := getPublicKeys(publicKeyStrings, query, selfEnv)
			if err != nil {
				utils.ExitOnError(err)
			}
			vault, err := getVault(vaultFile)
			if err == nil {
				var envSecretKey *crypto.SecretKey
				if envSecretKey, err = slv.GetSecretKey(); err == nil {
					err = vault.Unlock(*envSecretKey)
				}
				if err == nil {
					if err = vault.Revoke(publicKeys); err == nil {
						fmt.Println("Shared vault:", color.GreenString(vaultFile))
						utils.SafeExit()
					}
				}
			}
			utils.ExitOnError(err)
		},
	}
	vaultAccessRemoveCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.Name, vaultAccessPublicKeysFlag.Shorthand, []string{}, vaultAccessPublicKeysFlag.Usage)
	vaultAccessRemoveCmd.Flags().StringP(cmdenv.EnvSearchFlag.Name, cmdenv.EnvSearchFlag.Shorthand, "", cmdenv.EnvSearchFlag.Usage)
	vaultAccessRemoveCmd.Flags().BoolP(cmdenv.EnvSelfFlag.Name, cmdenv.EnvSelfFlag.Shorthand, false, cmdenv.EnvSelfFlag.Usage)
	return vaultAccessRemoveCmd
}

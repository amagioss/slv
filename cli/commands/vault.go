package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/keystore"
	"github.com/shibme/slv/core/vaults"
	"github.com/spf13/cobra"
)

func vaultCommand() *cobra.Command {
	if vaultCmd != nil {
		return vaultCmd
	}
	vaultCmd = &cobra.Command{
		Use:   "vault",
		Short: "Vault operations",
		Long:  `Vault operations in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	vaultCmd.AddCommand(vaultNewCommand())
	vaultCmd.AddCommand(vaultShareCommand())
	return vaultCmd
}

func vaultNewCommand() *cobra.Command {
	if vaultNewCmd != nil {
		return vaultNewCmd
	}
	vaultNewCmd = &cobra.Command{
		Use:   "new",
		Short: "Creates a new vault",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			publicKeyStrings, err := cmd.Flags().GetStringSlice(vaultAccessPublicKeysFlag.name)
			if err != nil {
				PrintErrorAndExit(err)
			}
			var publicKeys []crypto.PublicKey
			for _, publicKeyString := range publicKeyStrings {
				publicKey, err := crypto.PublicKeyFromString(publicKeyString)
				if err != nil {
					PrintErrorAndExit(err)
				}
				publicKeys = append(publicKeys, *publicKey)
			}
			enableHash, _ := cmd.Flags().GetBool(vaultEnableHashingFlag.name)
			var hashLength uint32 = 0
			if enableHash {
				hashLength = 4
			}
			_, err = vaults.New(vaultFile, hashLength, publicKeys...)
			if err != nil {
				PrintErrorAndExit(err)
			}
			fmt.Println("Created vault:", color.GreenString(vaultFile))
			os.Exit(0)
		},
	}
	vaultNewCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	vaultNewCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.name, vaultAccessPublicKeysFlag.shorthand, []string{}, vaultAccessPublicKeysFlag.usage)
	vaultNewCmd.Flags().BoolP(vaultEnableHashingFlag.name, vaultEnableHashingFlag.shorthand, false, vaultEnableHashingFlag.usage)
	vaultNewCmd.MarkFlagRequired(vaultFileFlag.name)
	vaultNewCmd.MarkFlagRequired(vaultAccessPublicKeysFlag.name)
	return vaultNewCmd
}

func vaultShareCommand() *cobra.Command {
	if vaultShareCmd != nil {
		return vaultShareCmd
	}
	vaultShareCmd = &cobra.Command{
		Use:   "share",
		Short: "Shares a vault with another environment or group",
		Run: func(cmd *cobra.Command, args []string) {
			var envSecretKey *crypto.SecretKey
			envSecretKeyString, err := keystore.GetSecretKeyFromEnvar()
			if err == nil {
				envSecretKey, err = crypto.SecretKeyFromString(envSecretKeyString)
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			publicKeyStrings, err := cmd.Flags().GetStringSlice(vaultAccessPublicKeysFlag.name)
			if err != nil {
				PrintErrorAndExit(err)
			}
			var publicKeys []crypto.PublicKey
			for _, publicKeyString := range publicKeyStrings {
				publicKey, err := crypto.PublicKeyFromString(publicKeyString)
				if err != nil {
					PrintErrorAndExit(err)
				}
				publicKeys = append(publicKeys, *publicKey)
			}
			vault, err := vaults.Get(vaultFile)
			if err == nil {
				err = vault.Unlock(*envSecretKey)
				if err == nil {
					for _, pubKey := range publicKeys {
						if err = vault.ShareAccessToKey(pubKey); err != nil {
							break
						}
					}
					if err == nil {
						fmt.Println("Shared vault:", color.GreenString(vaultFile))
						os.Exit(0)
					}
				}
			}
			PrintErrorAndExit(err)
		},
	}
	vaultShareCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	vaultShareCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.name, vaultAccessPublicKeysFlag.shorthand, []string{}, vaultAccessPublicKeysFlag.usage)
	vaultShareCmd.MarkFlagRequired(vaultFileFlag.name)
	vaultShareCmd.MarkFlagRequired(vaultAccessPublicKeysFlag.name)
	return vaultShareCmd
}

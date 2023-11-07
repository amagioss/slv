package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/profiles"
	"github.com/shibme/slv/core/secretkeystore"
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
				exitOnError(err)
			}
			query := cmd.Flag(envSearchFlag.name).Value.String()
			if len(publicKeyStrings) == 0 && query == "" {
				exitOnError(fmt.Errorf("either --" + envSearchFlag.name +
					" or --" + vaultAccessPublicKeysFlag.name + " must be specified"))
			}
			var publicKeys []crypto.PublicKey
			for _, publicKeyString := range publicKeyStrings {
				publicKey, err := crypto.PublicKeyFromString(publicKeyString)
				if err != nil {
					exitOnError(err)
				}
				publicKeys = append(publicKeys, *publicKey)
			}
			if query != "" {
				prof, err := profiles.GetDefaultProfile()
				if err != nil {
					exitOnError(err)
				}
				envManifest, err := prof.GetEnvManifest()
				if err != nil {
					exitOnError(err)
				}
				for _, env := range envManifest.SearchEnv(query) {
					publicKeys = append(publicKeys, env.PublicKey)
				}
				if len(publicKeys) == 0 {
					exitOnError(fmt.Errorf("no matching environments found for search query: " + query))
				}
			}
			enableHash, _ := cmd.Flags().GetBool(vaultEnableHashingFlag.name)
			var hashLength uint32 = 0
			if enableHash {
				hashLength = 4
			}
			_, err = vaults.New(vaultFile, hashLength, publicKeys...)
			if err != nil {
				exitOnError(err)
			}
			fmt.Println("Created vault:", color.GreenString(vaultFile))
			safeExit()
		},
	}
	vaultNewCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	vaultNewCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.name, vaultAccessPublicKeysFlag.shorthand, []string{}, vaultAccessPublicKeysFlag.usage)
	vaultNewCmd.Flags().StringP(envSearchFlag.name, envSearchFlag.shorthand, "", envSearchFlag.usage)
	vaultNewCmd.Flags().BoolP(vaultEnableHashingFlag.name, vaultEnableHashingFlag.shorthand, false, vaultEnableHashingFlag.usage)
	vaultNewCmd.MarkFlagRequired(vaultFileFlag.name)
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
			envSecretKey, err := secretkeystore.GetSecretKey()
			if err != nil {
				exitOnError(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			publicKeyStrings, err := cmd.Flags().GetStringSlice(vaultAccessPublicKeysFlag.name)
			if err != nil {
				exitOnError(err)
			}
			query := cmd.Flag(envSearchFlag.name).Value.String()
			if len(publicKeyStrings) == 0 && query == "" {
				exitOnError(fmt.Errorf("either --" + envSearchFlag.name +
					" or --" + vaultAccessPublicKeysFlag.name + " must be specified"))
			}
			var publicKeys []crypto.PublicKey
			for _, publicKeyString := range publicKeyStrings {
				publicKey, err := crypto.PublicKeyFromString(publicKeyString)
				if err != nil {
					exitOnError(err)
				}
				publicKeys = append(publicKeys, *publicKey)
			}
			if query != "" {
				prof, err := profiles.GetDefaultProfile()
				if err != nil {
					exitOnError(err)
				}
				envManifest, err := prof.GetEnvManifest()
				if err != nil {
					exitOnError(err)
				}
				for _, env := range envManifest.SearchEnv(query) {
					publicKeys = append(publicKeys, env.PublicKey)
				}
				if len(publicKeys) == 0 {
					exitOnError(fmt.Errorf("no matching environments found for search query: " + query))
				}
			}
			vault, err := vaults.Get(vaultFile)
			if err == nil {
				err = vault.Unlock(*envSecretKey)
				if err == nil {
					for _, pubKey := range publicKeys {
						if _, err = vault.ShareAccessToKey(pubKey); err != nil {
							break
						}
					}
					if err == nil {
						fmt.Println("Shared vault:", color.GreenString(vaultFile))
						safeExit()
					}
				}
			}
			exitOnError(err)
		},
	}
	vaultShareCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	vaultShareCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.name, vaultAccessPublicKeysFlag.shorthand, []string{}, vaultAccessPublicKeysFlag.usage)
	vaultShareCmd.Flags().StringP(envSearchFlag.name, envSearchFlag.shorthand, "", envSearchFlag.usage)
	vaultShareCmd.MarkFlagRequired(vaultFileFlag.name)
	return vaultShareCmd
}

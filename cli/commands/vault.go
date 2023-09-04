package commands

import (
	"fmt"
	"os"

	"github.com/shibme/slv/crypto"
	"github.com/shibme/slv/vaults"
	"github.com/spf13/cobra"
)

func VaultCommand() *cobra.Command {
	env := &cobra.Command{
		Use:   "vault",
		Short: "Vault operations",
		Long:  `Vault operations in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	env.AddCommand(newVaultCommand())
	env.AddCommand(SecretsCommand())
	return env
}

func newVaultCommand() *cobra.Command {
	newVaultCmd := &cobra.Command{
		Use:   "new",
		Short: "Creates a new vault",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag("file").Value.String()
			publicKeyStrings, err := cmd.Flags().GetStringSlice("public-keys")
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
			_, err = vaults.New(vaultFile, publicKeys...)
			if err != nil {
				PrintErrorAndExit(err)
			}
			fmt.Println("Created vault: ", green, vaultFile)
			os.Exit(0)
		},
	}

	// Adding the flags
	newVaultCmd.Flags().StringP("file", "f", "", "Name of the environment")
	newVaultCmd.Flags().StringSliceP("public-keys", "k", []string{}, "Public keys of environments or groups that can access the vault")

	// Marking the flags as required
	newVaultCmd.MarkFlagRequired("file")
	newVaultCmd.MarkFlagRequired("public-keys")
	return newVaultCmd
}

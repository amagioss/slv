package commands

import (
	"fmt"
	"os"

	"github.com/shibme/slv/crypto"
	"github.com/shibme/slv/vaults"
	"github.com/spf13/cobra"
)

func SecretsCommand() *cobra.Command {
	env := &cobra.Command{
		Use:     "secret",
		Aliases: []string{"secrets"},
		Short:   "Working with secrets",
		Long:    `Working with secrets in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	env.AddCommand(addDirectSecretToVault())
	env.AddCommand(getDirectSecretFromVault())
	return env
}

func addDirectSecretToVault() *cobra.Command {
	addSecretCmd := &cobra.Command{
		Use:   "add",
		Short: "Adds a secret to the vault",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag("vault-file").Value.String()
			secretName := cmd.Flag("name").Value.String()
			secretValue := cmd.Flag("value").Value.String()
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				PrintErrorAndExit(err)
			}
			err = vault.AddDirectSecret(secretName, secretValue)
			if err != nil {
				PrintErrorAndExit(err)
			}
			fmt.Println("Added secret: ", green, secretName, reset, " to vault: ", green, vaultFile)
			os.Exit(0)
		},
	}

	// Adding the flags
	addSecretCmd.Flags().StringP("vault-file", "f", "", "Path to the vault file")
	addSecretCmd.Flags().StringP("name", "n", "", "Name of the secret")
	addSecretCmd.Flags().StringP("value", "v", "", "Value of the secret")

	// Marking the flags as required
	addSecretCmd.MarkFlagRequired("vault-file")
	addSecretCmd.MarkFlagRequired("name")
	addSecretCmd.MarkFlagRequired("value")
	return addSecretCmd
}

func getDirectSecretFromVault() *cobra.Command {
	addSecretCmd := &cobra.Command{
		Use:   "get",
		Short: "Gets a secret from the vault",
		Run: func(cmd *cobra.Command, args []string) {
			envPrivateKeyString := os.Getenv("SLV_ENV_SECRET_KEY")
			if envPrivateKeyString == "" {
				PrintErrorMessageAndExit("secret key not set in SLV_ENV_SECRET_KEY")
			}
			envPrivateKey, err := crypto.PrivateKeyFromString(envPrivateKeyString)
			if err != nil {
				PrintErrorAndExit(err)
			}
			vaultFile := cmd.Flag("vault-file").Value.String()
			secretName := cmd.Flag("name").Value.String()
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				PrintErrorAndExit(err)
			}
			err = vault.Unlock(*envPrivateKey)
			if err != nil {
				PrintErrorAndExit(err)
			}
			secret, err := vault.GetDirectSecret(secretName)
			if err != nil {
				PrintErrorAndExit(err)
			}
			fmt.Println(secret)
			os.Exit(0)
		},
	}

	// Adding the flags
	addSecretCmd.Flags().StringP("vault-file", "f", "", "Path to the vault file")
	addSecretCmd.Flags().StringP("name", "n", "", "Name of the secret")

	// Marking the flags as required
	addSecretCmd.MarkFlagRequired("vault-file")
	addSecretCmd.MarkFlagRequired("name")
	return addSecretCmd
}

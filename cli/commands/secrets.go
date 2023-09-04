package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/vaults"
	"github.com/spf13/cobra"
)

func secretCommand() *cobra.Command {
	if secretCmd != nil {
		return secretCmd
	}
	secretCmd = &cobra.Command{
		Use:     "secret",
		Aliases: []string{"secrets"},
		Short:   "Working with secrets",
		Long:    `Working with secrets in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	secretCmd.AddCommand(secretAddCommand())
	secretCmd.AddCommand(secretGetCommand())
	return secretCmd
}

func secretAddCommand() *cobra.Command {
	if secretAddCmd != nil {
		return secretAddCmd
	}
	secretAddCmd = &cobra.Command{
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
			fmt.Println("Added secret: ", color.GreenString(secretName), " to vault: ", color.GreenString(vaultFile))
			os.Exit(0)
		},
	}
	secretAddCmd.Flags().StringP("vault-file", "f", "", "Path to the vault file")
	secretAddCmd.Flags().StringP("name", "n", "", "Name of the secret")
	secretAddCmd.Flags().StringP("value", "v", "", "Value of the secret")
	secretAddCmd.MarkFlagRequired("vault-file")
	secretAddCmd.MarkFlagRequired("name")
	secretAddCmd.MarkFlagRequired("value")
	return secretAddCmd
}

func secretGetCommand() *cobra.Command {
	if secretGetCmd != nil {
		return secretGetCmd
	}
	secretGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Gets a secret from the vault",
		Run: func(cmd *cobra.Command, args []string) {
			envPrivateKeyString := getEnvSecretKey()
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
	secretGetCmd.Flags().StringP("vault-file", "f", "", "Path to the vault file")
	secretGetCmd.Flags().StringP("name", "n", "", "Name of the secret")
	secretGetCmd.MarkFlagRequired("vault-file")
	secretGetCmd.MarkFlagRequired("name")
	return secretGetCmd
}

package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/keyreader"
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
	secretCmd.AddCommand(secretRefCommand())
	secretCmd.AddCommand(secretDerefCommand())
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
			name := cmd.Flag("name").Value.String()
			secret := cmd.Flag("secret").Value.String()
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				PrintErrorAndExit(err)
			}
			err = vault.AddDirectSecret(name, secret)
			if err != nil {
				PrintErrorAndExit(err)
			}
			fmt.Println("Added secret: ", color.GreenString(name), " to vault: ", color.GreenString(vaultFile))
			os.Exit(0)
		},
	}
	secretAddCmd.Flags().StringP("vault-file", "v", "", "Path to the vault file")
	secretAddCmd.Flags().StringP("name", "n", "", "Name of the secret")
	secretAddCmd.Flags().StringP("secret", "s", "", "Value of the secret")
	secretAddCmd.MarkFlagRequired("vault-file")
	secretAddCmd.MarkFlagRequired("name")
	secretAddCmd.MarkFlagRequired("secret")
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
			var envPrivateKey *crypto.PrivateKey
			envPrivateKeyString, err := keyreader.GetFromEnvar()
			if err == nil {
				envPrivateKey, err = crypto.PrivateKeyFromString(envPrivateKeyString)
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
			vaultFile := cmd.Flag("vault-file").Value.String()
			name := cmd.Flag("name").Value.String()
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				PrintErrorAndExit(err)
			}
			err = vault.Unlock(*envPrivateKey)
			if err != nil {
				PrintErrorAndExit(err)
			}
			secret, err := vault.GetDirectSecret(name)
			if err != nil {
				PrintErrorAndExit(err)
			}
			fmt.Println(secret)
			os.Exit(0)
		},
	}
	secretGetCmd.Flags().StringP("vault-file", "v", "", "Path to the vault file")
	secretGetCmd.Flags().StringP("name", "n", "", "Name of the secret")
	secretGetCmd.MarkFlagRequired("vault-file")
	secretGetCmd.MarkFlagRequired("name")
	return secretGetCmd
}

func secretRefCommand() *cobra.Command {
	if secretRefCmd != nil {
		return secretRefCmd
	}
	secretRefCmd = &cobra.Command{
		Use:   "ref",
		Short: "References and updates secrets to a vault from a given yaml or json file",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag("vault-file").Value.String()
			file := cmd.Flag("file").Value.String()
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				PrintErrorAndExit(err)
			}
			preview, _ := cmd.Flags().GetBool("preview")
			result, err := vault.ReferenceSecrets(file, preview)
			if err != nil {
				PrintErrorAndExit(err)
			}
			if preview {
				fmt.Println(result)
			} else {
				fmt.Println("Referenced", color.GreenString(file), "from vault", color.GreenString(vaultFile))
			}
			os.Exit(0)
		},
	}
	secretRefCmd.Flags().StringP("vault-file", "v", "", "Path to the vault file")
	secretRefCmd.Flags().StringP("file", "f", "", "Path to the yaml or json file")
	secretRefCmd.Flags().BoolP("preview", "p", false, "Enable preview mode")
	secretRefCmd.MarkFlagRequired("vault-file")
	secretRefCmd.MarkFlagRequired("file")
	return secretRefCmd
}

func secretDerefCommand() *cobra.Command {
	if secretDerefCmd != nil {
		return secretDerefCmd
	}
	secretDerefCmd = &cobra.Command{
		Use:   "deref",
		Short: "Dereferences and updates secrets from a vault to a given yaml or json file",
		Run: func(cmd *cobra.Command, args []string) {
			var envPrivateKey *crypto.PrivateKey
			envPrivateKeyString, err := keyreader.GetFromEnvar()
			if err == nil {
				envPrivateKey, err = crypto.PrivateKeyFromString(envPrivateKeyString)
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
			vaultFile := cmd.Flag("vault-file").Value.String()
			file := cmd.Flag("file").Value.String()
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				PrintErrorAndExit(err)
			}
			err = vault.Unlock(*envPrivateKey)
			if err != nil {
				PrintErrorAndExit(err)
			}
			_, err = vault.DereferenceSecrets(file)
			if err != nil {
				PrintErrorAndExit(err)
			}
			fmt.Println("Dereferenced ", color.GreenString(file), "with the vault", color.GreenString(vaultFile))
			os.Exit(0)
		},
	}
	secretDerefCmd.Flags().StringP("vault-file", "v", "", "Path to the vault file")
	secretDerefCmd.Flags().StringP("file", "f", "", "Path to the yaml or json file")
	secretDerefCmd.Flags().BoolP("preview", "p", false, "Enable preview mode")
	secretDerefCmd.MarkFlagRequired("vault-file")
	secretDerefCmd.MarkFlagRequired("file")
	return secretDerefCmd
}

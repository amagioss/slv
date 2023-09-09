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
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			name := cmd.Flag(secretNameFlag.name).Value.String()
			secret := cmd.Flag(secretValueFlag.name).Value.String()
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
	secretAddCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	secretAddCmd.Flags().StringP(secretNameFlag.name, secretNameFlag.shorthand, "", secretNameFlag.usage)
	secretAddCmd.Flags().StringP(secretValueFlag.name, secretValueFlag.shorthand, "", secretValueFlag.usage)
	secretAddCmd.MarkFlagRequired(vaultFileFlag.name)
	secretAddCmd.MarkFlagRequired(secretNameFlag.name)
	secretAddCmd.MarkFlagRequired(secretValueFlag.name)
	return secretAddCmd
}

func secretGetCommand() *cobra.Command {
	if secretGetCmd != nil {
		return secretGetCmd
	}
	secretGetCmd = &cobra.Command{
		Use:     "get",
		Aliases: []string{"show", "view", "read"},
		Short:   "Gets a secret from the vault",
		Run: func(cmd *cobra.Command, args []string) {
			var envSecretKey *crypto.SecretKey
			envSecretKeyString, err := keyreader.GetFromEnvar()
			if err == nil {
				envSecretKey, err = crypto.SecretKeyFromString(envSecretKeyString)
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			name := cmd.Flag(secretNameFlag.name).Value.String()
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				PrintErrorAndExit(err)
			}
			err = vault.Unlock(*envSecretKey)
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
	secretGetCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	secretGetCmd.Flags().StringP(secretNameFlag.name, secretNameFlag.shorthand, "", secretNameFlag.usage)
	secretGetCmd.MarkFlagRequired(vaultFileFlag.name)
	secretGetCmd.MarkFlagRequired(secretNameFlag.name)
	return secretGetCmd
}

func secretRefCommand() *cobra.Command {
	if secretRefCmd != nil {
		return secretRefCmd
	}
	secretRefCmd = &cobra.Command{
		Use:     "ref",
		Aliases: []string{"reference"},
		Short:   "References and updates secrets to a vault from a given yaml or json file",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			file := cmd.Flag(secretRefFileFlag.name).Value.String()
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				PrintErrorAndExit(err)
			}
			previewOnly, _ := cmd.Flags().GetBool(secretRefPreviewOnlyFlag.name)
			result, err := vault.ReferenceSecrets(file, previewOnly)
			if err != nil {
				PrintErrorAndExit(err)
			}
			if previewOnly {
				fmt.Println(result)
			} else {
				fmt.Println("Referenced", color.GreenString(file), "from vault", color.GreenString(vaultFile))
			}
			os.Exit(0)
		},
	}
	secretRefCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	secretRefCmd.Flags().StringP(secretRefFileFlag.name, secretRefFileFlag.shorthand, "", secretRefFileFlag.usage)
	secretRefCmd.Flags().BoolP(secretRefPreviewOnlyFlag.name, secretRefPreviewOnlyFlag.shorthand, false, secretRefPreviewOnlyFlag.usage)
	secretRefCmd.MarkFlagRequired(vaultFileFlag.name)
	secretRefCmd.MarkFlagRequired(secretRefFileFlag.name)
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
			var envSecretKey *crypto.SecretKey
			envSecretKeyString, err := keyreader.GetFromEnvar()
			if err == nil {
				envSecretKey, err = crypto.SecretKeyFromString(envSecretKeyString)
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			file := cmd.Flag(secretRefFileFlag.name).Value.String()
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				PrintErrorAndExit(err)
			}
			err = vault.Unlock(*envSecretKey)
			if err != nil {
				PrintErrorAndExit(err)
			}
			previewOnly, _ := cmd.Flags().GetBool(secretRefPreviewOnlyFlag.name)
			result, err := vault.DereferenceSecrets(file, previewOnly)
			if err != nil {
				PrintErrorAndExit(err)
			}
			if previewOnly {
				fmt.Println(result)
			} else {
				fmt.Println("Dereferenced ", color.GreenString(file), "with the vault", color.GreenString(vaultFile))
			}
			os.Exit(0)
		},
	}
	secretDerefCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	secretDerefCmd.Flags().StringP(secretRefFileFlag.name, secretRefFileFlag.shorthand, "", secretRefFileFlag.usage)
	secretDerefCmd.Flags().BoolP(secretRefPreviewOnlyFlag.name, secretRefPreviewOnlyFlag.shorthand, false, secretRefPreviewOnlyFlag.usage)
	secretDerefCmd.MarkFlagRequired(vaultFileFlag.name)
	secretDerefCmd.MarkFlagRequired(secretRefFileFlag.name)
	return secretDerefCmd
}

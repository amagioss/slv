package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/shibme/slv/core/secretkeystore"
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
	secretCmd.AddCommand(secretPutCommand())
	secretCmd.AddCommand(secretGetCommand())
	secretCmd.AddCommand(secretRefCommand())
	secretCmd.AddCommand(secretDerefCommand())
	return secretCmd
}

func secretPutCommand() *cobra.Command {
	if secretPutCmd != nil {
		return secretPutCmd
	}
	secretPutCmd = &cobra.Command{
		Use:     "put",
		Aliases: []string{"add", "set", "create"},
		Short:   "Adds a secret to the vault",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			name := cmd.Flag(secretNameFlag.name).Value.String()
			secret := cmd.Flag(secretValueFlag.name).Value.String()
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				exitOnError(err)
			}
			err = vault.AddDirectSecret(name, secret)
			if err != nil {
				exitOnError(err)
			}
			fmt.Println("Added secret: ", color.GreenString(name), " to vault: ", color.GreenString(vaultFile))
			safeExit()
		},
	}
	secretPutCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	secretPutCmd.Flags().StringP(secretNameFlag.name, secretNameFlag.shorthand, "", secretNameFlag.usage)
	secretPutCmd.Flags().StringP(secretValueFlag.name, secretValueFlag.shorthand, "", secretValueFlag.usage)
	secretPutCmd.MarkFlagRequired(vaultFileFlag.name)
	secretPutCmd.MarkFlagRequired(secretNameFlag.name)
	secretPutCmd.MarkFlagRequired(secretValueFlag.name)
	return secretPutCmd
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
			envSecretKey, err := secretkeystore.GetSecretKey()
			if err != nil {
				exitOnError(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			name := cmd.Flag(secretNameFlag.name).Value.String()
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				exitOnError(err)
			}
			err = vault.Unlock(*envSecretKey)
			if err != nil {
				exitOnError(err)
			}
			secret, err := vault.GetDirectSecret(name)
			if err != nil {
				exitOnError(err)
			}
			fmt.Println(secret)
			safeExit()
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
			files, err := cmd.Flags().GetStringSlice(secretRefFileFlag.name)
			if err != nil {
				exitOnError(err)
			}
			vault, err := vaults.Get(vaultFile)
			if err != nil {
				exitOnError(err)
			}
			previewOnly := false
			if len(files) > 1 {
				previewOnly, _ = cmd.Flags().GetBool(secretRefPreviewOnlyFlag.name)
			}
			for _, file := range files {
				result, err := vault.RefSecrets(file, previewOnly)
				if err != nil {
					exitOnError(err)
				}
				if previewOnly {
					fmt.Println(result)
				} else {
					fmt.Println("Referenced", color.GreenString(file), "from vault", color.GreenString(vaultFile))
				}
			}
			safeExit()
		},
	}
	secretRefCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	secretRefCmd.Flags().StringSliceP(secretRefFileFlag.name, secretRefFileFlag.shorthand, []string{}, secretRefFileFlag.usage)
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
			envSecretKey, err := secretkeystore.GetSecretKey()
			if err != nil {
				exitOnError(err)
			}
			vaultFiles, err := cmd.Flags().GetStringSlice(vaultFileFlag.name)
			if err != nil {
				exitOnError(err)
			}
			files, err := cmd.Flags().GetStringSlice(secretRefFileFlag.name)
			if err != nil {
				exitOnError(err)
			}
			previewOnly := false
			if len(vaultFiles) > 1 || len(files) > 1 {
				previewOnly, _ = cmd.Flags().GetBool(secretRefPreviewOnlyFlag.name)
			}
			for _, vaultFile := range vaultFiles {
				vault, err := vaults.Get(vaultFile)
				if err != nil {
					exitOnError(err)
				}
				err = vault.Unlock(*envSecretKey)
				if err != nil {
					exitOnError(err)
				}
				for _, file := range files {
					result, err := vault.DeRefSecrets(file, previewOnly)
					if err != nil {
						exitOnError(err)
					}
					if previewOnly {
						fmt.Println(result)
					} else {
						fmt.Println("Dereferenced ", color.GreenString(file), "with the vault", color.GreenString(vaultFile))
					}
				}
			}
			safeExit()
		},
	}
	secretDerefCmd.Flags().StringSliceP(vaultFileFlag.name, vaultFileFlag.shorthand, []string{}, vaultFileFlag.usage)
	secretDerefCmd.Flags().StringSliceP(secretRefFileFlag.name, secretRefFileFlag.shorthand, []string{}, secretRefFileFlag.usage)
	secretDerefCmd.Flags().BoolP(secretRefPreviewOnlyFlag.name, secretRefPreviewOnlyFlag.shorthand, false, secretRefPreviewOnlyFlag.usage)
	secretDerefCmd.MarkFlagRequired(vaultFileFlag.name)
	secretDerefCmd.MarkFlagRequired(secretRefFileFlag.name)
	return secretDerefCmd
}

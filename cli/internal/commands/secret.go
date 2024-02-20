package commands

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"savesecrets.org/slv"
)

func toBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

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
	secretCmd.AddCommand(secretExportCommand())
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
			vault, err := getVault(vaultFile)
			if err != nil {
				exitOnError(err)
			}
			forceUpdate, _ := cmd.Flags().GetBool(secretForceUpdateFlag.name)
			if !forceUpdate && vault.SecretExists(name) {
				exitOnErrorWithMessage("secret already exists. please use the --" + secretForceUpdateFlag.name + " flag to overwrite it.")
			}
			err = vault.PutSecret(name, []byte(secret))
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
	secretPutCmd.Flags().Bool(secretForceUpdateFlag.name, false, secretForceUpdateFlag.usage)
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
			envSecretKey, err := slv.GetSecretKey()
			if err != nil {
				exitOnError(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			name := cmd.Flag(secretNameFlag.name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				exitOnError(err)
			}
			err = vault.Unlock(*envSecretKey)
			if err != nil {
				exitOnError(err)
			}
			secret, err := vault.GetSecret(name)
			if err != nil {
				exitOnError(err)
			}
			encodeToBase64, _ := cmd.Flags().GetBool(secretEncodeBase64Flag.name)
			if encodeToBase64 {
				fmt.Println(toBase64(secret))
			} else {
				fmt.Println(string(secret))
			}
			safeExit()
		},
	}
	secretGetCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	secretGetCmd.Flags().StringP(secretNameFlag.name, secretNameFlag.shorthand, "", secretNameFlag.usage)
	secretGetCmd.Flags().BoolP(secretEncodeBase64Flag.name, secretEncodeBase64Flag.shorthand, false, secretEncodeBase64Flag.usage)
	secretGetCmd.MarkFlagRequired(vaultFileFlag.name)
	secretGetCmd.MarkFlagRequired(secretNameFlag.name)
	return secretGetCmd
}

func secretExportCommand() *cobra.Command {
	if secretExportCmd != nil {
		return secretExportCmd
	}
	secretExportCmd = &cobra.Command{
		Use:     "export",
		Aliases: []string{"dump", "get-all", "show-all", "view-all", "read-all"},
		Short:   "Exports all secrets from the vault",
		Run: func(cmd *cobra.Command, args []string) {
			envSecretKey, err := slv.GetSecretKey()
			if err != nil {
				exitOnError(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				exitOnError(err)
			}
			err = vault.Unlock(*envSecretKey)
			if err != nil {
				exitOnError(err)
			}
			secrets, err := vault.GetAllSecrets()
			if err != nil {
				exitOnError(err)
			}
			secretOutputMap := make(map[string]string)
			encodeToBase64, _ := cmd.Flags().GetBool(secretEncodeBase64Flag.name)
			for name, secret := range secrets {
				if encodeToBase64 {
					secretOutputMap[name] = toBase64(secret)
				} else {
					secretOutputMap[name] = string(secret)
				}
			}
			listFormat := cmd.Flag(secretListFormatFlag.name).Value.String()
			if listFormat == "" {
				listFormat = "table"
			}
			switch listFormat {
			case "table":
				tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
				for key, value := range secretOutputMap {
					fmt.Fprintf(tw, "%s\t%s\n", key, value)
				}
				tw.Flush()
			case "json":
				jsonData, err := json.MarshalIndent(secretOutputMap, "", "  ")
				if err != nil {
					exitOnError(err)
				}
				fmt.Println(string(jsonData))
			case "yaml", "yml":
				yamlData, err := yaml.Marshal(secretOutputMap)
				if err != nil {
					exitOnError(err)
				}
				fmt.Println(string(yamlData))
			case "envars", "envar", "env":
				for key, value := range secretOutputMap {
					fmt.Printf("%s=%s\n", key, value)
				}
			default:
				exitOnErrorWithMessage("invalid format: " + listFormat)
			}
			safeExit()
		},
	}
	secretExportCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	secretExportCmd.Flags().StringP(secretListFormatFlag.name, secretListFormatFlag.shorthand, "", secretListFormatFlag.usage)
	secretExportCmd.Flags().BoolP(secretEncodeBase64Flag.name, secretEncodeBase64Flag.shorthand, false, secretEncodeBase64Flag.usage)
	secretExportCmd.MarkFlagRequired(vaultFileFlag.name)
	return secretExportCmd
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
			vault, err := getVault(vaultFile)
			if err != nil {
				exitOnError(err)
			}
			refFile := cmd.Flag(secretRefFileFlag.name).Value.String()
			secretNamePrefix := cmd.Flag(secretNameFlag.name).Value.String()
			refType := strings.ToLower(cmd.Flag(secretRefTypeFlag.name).Value.String())
			previewOnly, _ := cmd.Flags().GetBool(secretRefPreviewOnlyFlag.name)
			forceUpdate, _ := cmd.Flags().GetBool(secretForceUpdateFlag.name)
			if secretNamePrefix == "" && refType == "" {
				exitOnErrorWithMessage("please provide at least one of --" + secretNameFlag.name + " or --" + secretRefTypeFlag.name + " flag")
			}
			result, conflicting, err := vault.RefSecrets(refType, refFile, secretNamePrefix, forceUpdate, previewOnly)
			if conflicting {
				exitOnErrorWithMessage("conflict found. please use the --" + secretNameFlag.name + " flag to set a different name or --" + secretForceUpdateFlag.name + " flag to overwrite them.")
			} else if err != nil {
				exitOnError(err)
			}
			if previewOnly {
				fmt.Println(result)
			} else {
				fmt.Println("Auto referenced", color.GreenString(refFile), "with vault", color.GreenString(vaultFile))
			}
			safeExit()
		},
	}
	secretRefCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	secretRefCmd.Flags().StringP(secretRefFileFlag.name, secretRefFileFlag.shorthand, "", secretRefFileFlag.usage)
	secretRefCmd.Flags().StringP(secretNameFlag.name, secretNameFlag.shorthand, "", secretNameFlag.usage)
	secretRefCmd.Flags().StringP(secretRefTypeFlag.name, secretRefTypeFlag.shorthand, "", secretRefTypeFlag.usage)
	secretRefCmd.Flags().BoolP(secretRefPreviewOnlyFlag.name, secretRefPreviewOnlyFlag.shorthand, false, secretRefPreviewOnlyFlag.usage)
	secretRefCmd.Flags().BoolP(secretForceUpdateFlag.name, secretForceUpdateFlag.shorthand, false, secretForceUpdateFlag.usage)
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
			envSecretKey, err := slv.GetSecretKey()
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
				vault, err := getVault(vaultFile)
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

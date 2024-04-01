package cmdvault

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"oss.amagi.com/slv"
	"oss.amagi.com/slv/cli/internal/commands/utils"
)

func vaultGetCommand() *cobra.Command {
	if vaultGetCmd != nil {
		return vaultGetCmd
	}
	vaultGetCmd = &cobra.Command{
		Use:     "get",
		Aliases: []string{"show", "view", "read", "export", "dump"},
		Short:   "Get a secret from the vault",
		Run: func(cmd *cobra.Command, args []string) {
			envSecretKey, err := slv.GetSecretKey()
			if err != nil {
				utils.ExitOnError(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			name := cmd.Flag(secretNameFlag.Name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				utils.ExitOnError(err)
			}
			err = vault.Unlock(*envSecretKey)
			if err != nil {
				utils.ExitOnError(err)
			}
			encodeToBase64, _ := cmd.Flags().GetBool(secretEncodeBase64Flag.Name)
			exportFormat := cmd.Flag(vaultExportFormatFlag.Name).Value.String()
			secretMap := make(map[string]string)
			if name == "" {
				secrets, err := vault.GetAllSecrets()
				if err != nil {
					utils.ExitOnError(err)
				}
				for name, secret := range secrets {
					if encodeToBase64 {
						secretMap[name] = base64.StdEncoding.EncodeToString(secret)
					} else {
						secretMap[name] = string(secret)
					}
				}
			} else {
				secret, err := vault.GetSecret(name)
				if err != nil {
					utils.ExitOnError(err)
				}
				if encodeToBase64 {
					secretMap[name] = base64.StdEncoding.EncodeToString(secret)
				} else {
					secretMap[name] = string(secret)
				}
			}
			if exportFormat == "" {
				if name != "" {
					fmt.Println(secretMap[name])
					utils.SafeExit()
				}
				exportFormat = "envar"
			}
			switch exportFormat {
			case "json":
				jsonData, err := json.MarshalIndent(secretMap, "", "  ")
				if err != nil {
					utils.ExitOnError(err)
				}
				fmt.Println(string(jsonData))
			case "yaml", "yml":
				yamlData, err := yaml.Marshal(secretMap)
				if err != nil {
					utils.ExitOnError(err)
				}
				fmt.Println(string(yamlData))
			case "envars", "envar", ".env":
				for key, value := range secretMap {
					value = strings.ReplaceAll(value, "\\", "\\\\")
					value = strings.ReplaceAll(value, "\"", "\\\"")
					fmt.Printf("%s=\"%s\"\n", key, value)
				}
			default:
				utils.ExitOnErrorWithMessage("invalid format: " + exportFormat)
			}
			utils.SafeExit()
		},
	}
	vaultGetCmd.Flags().StringP(secretNameFlag.Name, secretNameFlag.Shorthand, "", secretNameFlag.Usage)
	vaultGetCmd.Flags().BoolP(secretEncodeBase64Flag.Name, secretEncodeBase64Flag.Shorthand, false, secretEncodeBase64Flag.Usage)
	vaultGetCmd.Flags().StringP(vaultExportFormatFlag.Name, vaultExportFormatFlag.Shorthand, "", vaultExportFormatFlag.Usage)
	return vaultGetCmd
}

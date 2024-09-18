package cmdvault

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/secretkey"
	"oss.amagi.com/slv/internal/core/vaults"
)

func getDecryptedDataMap(vault *vaults.Vault, itemName string, encodeToBase64, withMetadata bool) map[string]any {
	type vaultItem struct {
		Value     string `json:"value,omitempty" yaml:"value,omitempty"`
		Secret    bool   `json:"secret,omitempty" yaml:"secret,omitempty"`
		UpdatedAt string `json:"updatedAt,omitempty" yaml:"updatedAt,omitempty"`
		Hash      string `json:"hash,omitempty" yaml:"hash,omitempty"`
	}
	dataMap := make(map[string]any)
	var vaultData map[string]*vaults.VaultData
	var err error
	if itemName == "" {
		if vaultData, err = vault.List(true); err != nil {
			utils.ExitOnError(err)
		}
	} else {
		var item *vaults.VaultData
		if item, err = vault.Get(itemName); err != nil {
			utils.ExitOnError(err)
		}
		vaultData = map[string]*vaults.VaultData{itemName: item}
	}
	for name, value := range vaultData {
		var val string
		if encodeToBase64 {
			val = base64.StdEncoding.EncodeToString(value.Value())
		} else {
			val = string(value.Value())
		}
		if withMetadata {
			vi := vaultItem{
				Value:  val,
				Secret: value.IsSecret(),
			}
			if value.UpdatedAt() != nil {
				vi.UpdatedAt = value.UpdatedAt().Format(time.RFC3339)
			}
			if value.Hash() != "" {
				vi.Hash = value.Hash()
			}
			dataMap[name] = vi
		} else {
			dataMap[name] = val
		}
	}
	return dataMap
}

func vaultGetCommand() *cobra.Command {
	if vaultGetCmd != nil {
		return vaultGetCmd
	}
	vaultGetCmd = &cobra.Command{
		Use:     "get",
		Aliases: []string{"show", "view", "read", "export", "dump"},
		Short:   "Get a secret from the vault",
		Run: func(cmd *cobra.Command, args []string) {
			envSecretKey, err := secretkey.Get()
			if err != nil {
				utils.ExitOnError(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			itemName := cmd.Flag(secretNameFlag.Name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				utils.ExitOnError(err)
			}
			err = vault.Unlock(envSecretKey)
			if err != nil {
				utils.ExitOnError(err)
			}
			encodeToBase64, _ := cmd.Flags().GetBool(valueEncodeBase64Flag.Name)
			withMetadata, _ := cmd.Flags().GetBool(valueWithMetadata.Name)
			exportFormat := cmd.Flag(vaultExportFormatFlag.Name).Value.String()
			switch exportFormat {
			case "json":
				dataMap := getDecryptedDataMap(vault, itemName, encodeToBase64, withMetadata)
				jsonData, err := json.MarshalIndent(dataMap, "", "  ")
				if err != nil {
					utils.ExitOnError(err)
				}
				fmt.Println(string(jsonData))
			case "yaml", "yml":
				dataMap := getDecryptedDataMap(vault, itemName, encodeToBase64, withMetadata)
				yamlData, err := yaml.Marshal(dataMap)
				if err != nil {
					utils.ExitOnError(err)
				}
				fmt.Println(string(yamlData))
			case "envars", "envar", "env", ".env":
				dataMap := getDecryptedDataMap(vault, itemName, encodeToBase64, false)
				for key, value := range dataMap {
					strValue := value.(string)
					strValue = strings.ReplaceAll(strValue, "\\", "\\\\")
					strValue = strings.ReplaceAll(strValue, "\"", "\\\"")
					fmt.Printf("%s=\"%s\"\n", key, strValue)
				}
			default:
				if itemName == "" {
					showVaultData(vault)
				} else {
					if item, err := vault.Get(itemName); err != nil {
						utils.ExitOnError(err)
					} else {
						fmt.Println(string(item.Value()))
					}
				}
			}
			utils.SafeExit()
		},
	}
	vaultGetCmd.Flags().StringP(secretNameFlag.Name, secretNameFlag.Shorthand, "", secretNameFlag.Usage)
	vaultGetCmd.Flags().BoolP(valueEncodeBase64Flag.Name, valueEncodeBase64Flag.Shorthand, false, valueEncodeBase64Flag.Usage)
	vaultGetCmd.Flags().BoolP(valueWithMetadata.Name, valueWithMetadata.Shorthand, false, valueWithMetadata.Usage)
	vaultGetCmd.Flags().StringP(vaultExportFormatFlag.Name, vaultExportFormatFlag.Shorthand, "", vaultExportFormatFlag.Usage)
	return vaultGetCmd
}

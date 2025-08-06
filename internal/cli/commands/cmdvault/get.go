package cmdvault

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
)

func unlockVault(vault *vaults.Vault) {
	if vault.IsLocked() {
		envSecretKey, err := session.GetSecretKey()
		if err != nil {
			utils.ExitOnError(err)
		}
		if err = vault.Unlock(envSecretKey); err != nil {
			utils.ExitOnError(err)
		}
	}
}

func getVaultItemMap(vault *vaults.Vault, itemName string, encodeToBase64, withMetadata bool) map[string]any {
	type itemInfo struct {
		Value       string `json:"value,omitempty" yaml:"value,omitempty"`
		IsPlaintext bool   `json:"isPlaintext,omitempty" yaml:"isPlaintext,omitempty"`
		EncryptedAt string `json:"encryptedAt,omitempty" yaml:"encryptedAt,omitempty"`
		Hash        string `json:"hash,omitempty" yaml:"hash,omitempty"`
	}
	dataMap := make(map[string]any)
	var vaultItemMap map[string]*vaults.VaultItem
	var err error
	if itemName == "" {
		unlockVault(vault)
		if vaultItemMap, err = vault.GetAllItems(); err != nil {
			utils.ExitOnError(err)
		}
	} else {
		item, _ := vault.Get(itemName)
		if item == nil {
			utils.ExitOnError(fmt.Errorf("item %s not found", itemName))
		}
		if !item.IsPlaintext() {
			unlockVault(vault)
		}
		if item, err = vault.Get(itemName); err != nil {
			utils.ExitOnError(err)
		}
		vaultItemMap = map[string]*vaults.VaultItem{itemName: item}
	}
	for name, item := range vaultItemMap {
		itemValue, err := item.Value()
		if err != nil {
			utils.ExitOnError(fmt.Errorf("error getting value for %s: %w", name, err))
		}
		var valueStr string
		if encodeToBase64 {
			valueStr = base64.StdEncoding.EncodeToString(itemValue)
		} else {
			valueStr = string(itemValue)
		}
		if withMetadata {
			fi := itemInfo{
				Value:       valueStr,
				IsPlaintext: item.IsPlaintext(),
			}
			if item.EncryptedAt() != nil {
				fi.EncryptedAt = item.EncryptedAt().Format(time.RFC3339)
			}
			if item.Hash() != "" {
				fi.Hash = item.Hash()
			}
			dataMap[name] = fi
		} else {
			dataMap[name] = valueStr
		}
	}
	return dataMap
}

func vaultGetCommand() *cobra.Command {
	if vaultGetCmd == nil {
		vaultGetCmd = &cobra.Command{
			Use:     "get",
			Aliases: []string{"show", "view", "read", "export", "dump"},
			Short:   "Get one or more values or list the vault in desired format",
			Run: func(cmd *cobra.Command, args []string) {
				vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
				itemName := cmd.Flag(itemNameFlag.Name).Value.String()
				vault, err := vaults.Get(vaultFile)
				if err != nil {
					utils.ExitOnError(err)
				}
				encodeToBase64, _ := cmd.Flags().GetBool(valueEncodeBase64Flag.Name)
				withMetadata, _ := cmd.Flags().GetBool(valueWithMetadata.Name)
				exportFormat := cmd.Flag(vaultExportFormatFlag.Name).Value.String()
				switch exportFormat {
				case "json":
					viMap := getVaultItemMap(vault, itemName, encodeToBase64, withMetadata)
					jsonData, err := json.MarshalIndent(viMap, "", "  ")
					if err != nil {
						utils.ExitOnError(err)
					}
					fmt.Println(string(jsonData))
				case "yaml", "yml":
					dataMap := getVaultItemMap(vault, itemName, encodeToBase64, withMetadata)
					yamlData, err := yaml.Marshal(dataMap)
					if err != nil {
						utils.ExitOnError(err)
					}
					fmt.Println(string(yamlData))
				case "envars", "envar", "env", ".env":
					dataMap := getVaultItemMap(vault, itemName, encodeToBase64, false)
					for key, value := range dataMap {
						strValue := value.(string)
						strValue = strings.ReplaceAll(strValue, "\\", "\\\\")
						strValue = strings.ReplaceAll(strValue, "\"", "\\\"")
						fmt.Printf("%s=\"%s\"\n", key, strValue)
					}
				default:
					if itemName == "" {
						unlockVault(vault)
						showVault(vault)
					} else {
						if !vault.ItemExists(itemName) {
							utils.ExitOnError(fmt.Errorf("item %s not found", itemName))
						}
						item, err := vault.Get(itemName)
						if err != nil {
							utils.ExitOnError(err)
						}
						if !item.IsPlaintext() {
							unlockVault(vault)
						}
						if itemValueStr, err := item.ValueString(); err != nil {
							utils.ExitOnError(err)
						} else {
							fmt.Println(itemValueStr)
						}
					}
				}
				utils.SafeExit()
			},
		}
		vaultGetCmd.Flags().StringP(itemNameFlag.Name, itemNameFlag.Shorthand, "", itemNameFlag.Usage)
		if err := vaultGetCmd.RegisterFlagCompletionFunc(itemNameFlag.Name, vaultItemNameCompletion); err != nil {
			utils.ExitOnError(err)
		}
		vaultGetCmd.Flags().BoolP(valueEncodeBase64Flag.Name, valueEncodeBase64Flag.Shorthand, false, valueEncodeBase64Flag.Usage)
		vaultGetCmd.Flags().BoolP(valueWithMetadata.Name, valueWithMetadata.Shorthand, false, valueWithMetadata.Usage)
		vaultGetCmd.Flags().StringP(vaultExportFormatFlag.Name, vaultExportFormatFlag.Shorthand, "", vaultExportFormatFlag.Usage)
	}
	return vaultGetCmd
}

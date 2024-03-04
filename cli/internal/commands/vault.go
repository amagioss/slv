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
	"savesecrets.org/slv/core/commons"
	"savesecrets.org/slv/core/config"
	"savesecrets.org/slv/core/crypto"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/input"
	"savesecrets.org/slv/core/profiles"
	"savesecrets.org/slv/core/vaults"
)

const (
	k8sApiVersion = config.K8SLVGroup + "/" + config.K8SLVVersion
	k8sKind       = config.K8SLVKind
	k8sVaultField = config.K8SLVVaultField
)

func getVault(filePath string) (*vaults.Vault, error) {
	vault, err := vaults.Get(filePath)
	if err != nil || vault.Config.PublicKey == "" {
		vault, err = vaults.GetFromField(filePath, k8sVaultField)
	}
	return vault, err
}

func newK8sVault(filePath, name string, hashLength uint8, rootPublicKey *crypto.PublicKey, publicKeys ...*crypto.PublicKey) (*vaults.Vault, error) {
	vault, err := vaults.New(filePath, k8sVaultField, hashLength, rootPublicKey, publicKeys...)
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	if err := commons.ReadFromYAML(filePath, &obj); err != nil {
		return nil, err
	}
	obj["apiVersion"] = k8sApiVersion
	obj["kind"] = k8sKind
	obj["metadata"] = map[string]interface{}{
		"name": name,
	}
	return vault, commons.WriteToYAML(filePath, "", obj)
}

func vaultCommand() *cobra.Command {
	if vaultCmd != nil {
		return vaultCmd
	}
	vaultCmd = &cobra.Command{
		Use:   "vault",
		Short: "Vault operations",
		Long:  `Vault operations in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	vaultCmd.PersistentFlags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	vaultCmd.MarkPersistentFlagRequired(vaultFileFlag.name)
	vaultCmd.AddCommand(vaultNewCommand())
	vaultCmd.AddCommand(vaultShareCommand())
	vaultCmd.AddCommand(vaultInfoCommand())
	vaultCmd.AddCommand(vaultPutCommand())
	vaultCmd.AddCommand(vaultImportCommand())
	vaultCmd.AddCommand(vaultGetCommand())
	vaultCmd.AddCommand(vaultExportCommand())
	vaultCmd.AddCommand(vaultRefCommand())
	vaultCmd.AddCommand(vaultDerefCommand())
	return vaultCmd
}

func vaultInfoCommand() *cobra.Command {
	if vaultInfoCmd != nil {
		return vaultInfoCmd
	}
	vaultInfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Displays information about a vault",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				exitOnError(err)
			}
			sealedSecretsMap, err := vault.ListSealedSecrets()
			if err != nil {
				exitOnError(err)
			}
			accessors, err := vault.ListAccessors()
			if err != nil {
				exitOnError(err)
			}
			profile, _ := profiles.GetDefaultProfile()
			self := environments.GetSelf()
			envMap := make(map[string]string, len(accessors))
			for _, accessor := range accessors {
				var env *environments.Environment
				envId := accessor.String()
				selfEnv := false
				if self != nil && self.PublicKey == accessor.String() {
					env = self
					selfEnv = true
				} else if profile != nil {
					env, err = profile.GetEnv(envId)
					if err != nil {
						exitOnError(err)
					}
				}
				if env != nil {
					if selfEnv {
						envMap[envId] = envId + "\t(" + color.CyanString("Self"+": "+env.Name) + ")"
					} else {
						envMap[envId] = envId + "\t(" + env.Name + ")"
					}
				} else {
					envMap[envId] = envId
				}
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
			fmt.Fprintln(w, "Vault ID\t:\t", vault.Config.PublicKey)
			fmt.Fprintln(w, "Secrets:")
			for name, sealedSecret := range sealedSecretsMap {
				hash := sealedSecret.Hash()
				if hash == "" {
					fmt.Fprintln(w, "  -", name, "\t:\t", sealedSecret.EncryptedAt().Format("Jan _2, 2006 03:04:05 PM MST"))
				} else {
					fmt.Fprintln(w, "  -", name, "\t:\t", sealedSecret.EncryptedAt().Format("Jan _2, 2006 03:04:05 PM MST"), "\t(", hash, ")")
				}
			}
			fmt.Fprintln(w, "Accessible by:")
			for _, envDesc := range envMap {
				fmt.Fprintln(w, "  -", envDesc)
			}
			w.Flush()
			safeExit()
		},
	}
	return vaultInfoCmd

}

func getPublicKeys(pubKeyStrSlice []string, query string, self bool) (publicKeys []*crypto.PublicKey,
	rootPublicKey *crypto.PublicKey, err error) {
	if len(pubKeyStrSlice) == 0 && query == "" && !self {
		return nil, nil, fmt.Errorf("Specify atleast one of the following flags:\n" +
			" --" + envSearchFlag.name + "\n" +
			" --" + vaultAccessPublicKeysFlag.name + "\n" +
			" --" + envSelfFlag.name)
	}
	for _, pubKeyStr := range pubKeyStrSlice {
		publicKey, err := crypto.PublicKeyFromString(pubKeyStr)
		if err != nil {
			return nil, nil, err
		}
		publicKeys = append(publicKeys, publicKey)
	}
	profile, err := profiles.GetDefaultProfile()
	if query != "" {
		if err != nil {
			return nil, nil, err
		}
		envs, err := profile.SearchEnvs(query)
		if err != nil {
			return nil, nil, err
		}
		for _, env := range envs {
			publicKey, err := crypto.PublicKeyFromString(env.PublicKey)
			if err != nil {
				return nil, nil, err
			}
			publicKeys = append(publicKeys, publicKey)
		}
		if len(publicKeys) == 0 {
			return nil, nil, fmt.Errorf("no matching environments found for search query: " + query)
		}
	}
	if self {
		selfEnv := environments.GetSelf()
		if selfEnv != nil {
			publicKey, err := crypto.PublicKeyFromString(selfEnv.PublicKey)
			if err != nil {
				return nil, nil, err
			}
			publicKeys = append(publicKeys, publicKey)
		}
	}
	if profile != nil {
		rootPublicKey, err = profile.RootPublicKey()
		if err != nil {
			return nil, nil, err
		}
	}
	return publicKeys, rootPublicKey, nil
}

func vaultNewCommand() *cobra.Command {
	if vaultNewCmd != nil {
		return vaultNewCmd
	}
	vaultNewCmd = &cobra.Command{
		Use:   "new",
		Short: "Creates a new vault",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			publicKeyStrings, err := cmd.Flags().GetStringSlice(vaultAccessPublicKeysFlag.name)
			if err != nil {
				exitOnError(err)
			}
			query := cmd.Flag(envSearchFlag.name).Value.String()
			envSelfFlag, _ := cmd.Flags().GetBool(envSelfFlag.name)
			publicKeys, rootPublicKey, err := getPublicKeys(publicKeyStrings, query, envSelfFlag)
			if err != nil {
				exitOnError(err)
			}
			enableHash, _ := cmd.Flags().GetBool(vaultEnableHashingFlag.name)
			var hashLength uint8 = 0
			if enableHash {
				hashLength = 4
			}
			k8sName := cmd.Flag(vaultK8sFlag.name).Value.String()
			if k8sName == "" {
				_, err = vaults.New(vaultFile, "", hashLength, rootPublicKey, publicKeys...)
			} else {
				_, err = newK8sVault(vaultFile, k8sName, hashLength, rootPublicKey, publicKeys...)
			}
			if err != nil {
				exitOnError(err)
			}
			fmt.Println("Created vault:", color.GreenString(vaultFile))
			safeExit()
		},
	}
	vaultNewCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.name, vaultAccessPublicKeysFlag.shorthand, []string{}, vaultAccessPublicKeysFlag.usage)
	vaultNewCmd.Flags().StringP(envSearchFlag.name, envSearchFlag.shorthand, "", envSearchFlag.usage)
	vaultNewCmd.Flags().BoolP(envSelfFlag.name, envSelfFlag.shorthand, false, envSelfFlag.usage)
	vaultNewCmd.Flags().StringP(vaultK8sFlag.name, vaultK8sFlag.shorthand, "", vaultK8sFlag.usage)
	vaultNewCmd.Flags().BoolP(vaultEnableHashingFlag.name, vaultEnableHashingFlag.shorthand, false, vaultEnableHashingFlag.usage)
	return vaultNewCmd
}

func vaultShareCommand() *cobra.Command {
	if vaultShareCmd != nil {
		return vaultShareCmd
	}
	vaultShareCmd = &cobra.Command{
		Use:   "share",
		Short: "Shares a vault with another environment or group",
		Run: func(cmd *cobra.Command, args []string) {
			envSecretKey, err := slv.GetSecretKey()
			if err != nil {
				exitOnError(err)
			}
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			publicKeyStrings, err := cmd.Flags().GetStringSlice(vaultAccessPublicKeysFlag.name)
			if err != nil {
				exitOnError(err)
			}
			query := cmd.Flag(envSearchFlag.name).Value.String()
			envSelfFlag, _ := cmd.Flags().GetBool(envSelfFlag.name)
			publicKeys, _, err := getPublicKeys(publicKeyStrings, query, envSelfFlag)
			if err != nil {
				exitOnError(err)
			}
			vault, err := getVault(vaultFile)
			if err == nil {
				err = vault.Unlock(*envSecretKey)
				if err == nil {
					for _, publicKey := range publicKeys {
						if _, err = vault.Share(publicKey); err != nil {
							break
						}
					}
					if err == nil {
						fmt.Println("Shared vault:", color.GreenString(vaultFile))
						safeExit()
					}
				}
			}
			exitOnError(err)
		},
	}
	vaultShareCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.name, vaultAccessPublicKeysFlag.shorthand, []string{}, vaultAccessPublicKeysFlag.usage)
	vaultShareCmd.Flags().StringP(envSearchFlag.name, envSearchFlag.shorthand, "", envSearchFlag.usage)
	vaultShareCmd.Flags().BoolP(envSelfFlag.name, envSelfFlag.shorthand, false, envSelfFlag.usage)
	return vaultShareCmd
}

func vaultPutCommand() *cobra.Command {
	if vaultPutCmd != nil {
		return vaultPutCmd
	}
	vaultPutCmd = &cobra.Command{
		Use:     "put",
		Aliases: []string{"add", "set", "create"},
		Short:   "Adds a secret to the vault",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			name := cmd.Flag(secretNameFlag.name).Value.String()
			secretStr := cmd.Flag(secretValueFlag.name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				exitOnError(err)
			}
			forceUpdate, _ := cmd.Flags().GetBool(secretForceUpdateFlag.name)
			if !forceUpdate && vault.SecretExists(name) {
				confirmation, err := input.GetVisibleInput("Secret already exists. Do you wish to overwrite it? (y/n): ")
				if err != nil {
					exitOnError(err)
				}
				if confirmation != "y" {
					fmt.Println(color.YellowString("Operation aborted"))
					safeExit()
				}
			}
			var secret []byte
			if secretStr == "" {
				secret, err = input.GetHiddenInput("Enter the secret value for " + name + ": ")
				if err != nil {
					exitOnError(err)
				}
			} else {
				secret = []byte(secretStr)
			}
			err = vault.PutSecret(name, secret)
			if err != nil {
				exitOnError(err)
			}
			fmt.Println("Updated secret: ", color.GreenString(name), " to vault: ", color.GreenString(vaultFile))
			safeExit()
		},
	}
	vaultPutCmd.Flags().StringP(secretNameFlag.name, secretNameFlag.shorthand, "", secretNameFlag.usage)
	vaultPutCmd.Flags().StringP(secretValueFlag.name, secretValueFlag.shorthand, "", secretValueFlag.usage)
	vaultPutCmd.Flags().Bool(secretForceUpdateFlag.name, false, secretForceUpdateFlag.usage)
	vaultPutCmd.MarkFlagRequired(secretNameFlag.name)
	return vaultPutCmd
}

func vaultImportCommand() *cobra.Command {
	if vaultImportCmd != nil {
		return vaultImportCmd
	}
	vaultImportCmd = &cobra.Command{
		Use:     "import",
		Aliases: []string{"load", "put-all", "add-all", "set-all", "create-all"},
		Short:   "Imports secrets into the vault from YAML or JSON",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				exitOnError(err)
			}
			forceUpdate, _ := cmd.Flags().GetBool(secretForceUpdateFlag.name)
			importFile := cmd.Flag(vaultImportFileFlag.name).Value.String()
			var importData []byte
			if importFile == "" {
				importData, err = input.GetHiddenInput("Enter the YAML/JSON data to be imported: ")
			} else {
				importData, err = os.ReadFile(importFile)
			}
			if err != nil {
				exitOnError(err)
			}
			if err = vault.ImportSecrets(importData, forceUpdate); err != nil {
				exitOnError(err)
			}
			fmt.Printf("Successfully imported secrets from %s into the vault %s\n", color.GreenString(importFile), color.GreenString(vaultFile))
			safeExit()
		},
	}
	vaultImportCmd.Flags().StringP(vaultImportFileFlag.name, vaultImportFileFlag.shorthand, "", vaultImportFileFlag.usage)
	vaultImportCmd.Flags().Bool(secretForceUpdateFlag.name, false, secretForceUpdateFlag.usage)
	return vaultImportCmd
}

func vaultGetCommand() *cobra.Command {
	if vaultGetCmd != nil {
		return vaultGetCmd
	}
	vaultGetCmd = &cobra.Command{
		Use:     "get",
		Aliases: []string{"show", "view", "read"},
		Short:   "Get a secret from the vault",
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
				fmt.Println(base64.StdEncoding.EncodeToString(secret))
			} else {
				fmt.Println(string(secret))
			}
			safeExit()
		},
	}
	vaultGetCmd.Flags().StringP(secretNameFlag.name, secretNameFlag.shorthand, "", secretNameFlag.usage)
	vaultGetCmd.Flags().BoolP(secretEncodeBase64Flag.name, secretEncodeBase64Flag.shorthand, false, secretEncodeBase64Flag.usage)
	vaultGetCmd.MarkFlagRequired(secretNameFlag.name)
	return vaultGetCmd
}

func vaultExportCommand() *cobra.Command {
	if vaultExportCmd != nil {
		return vaultExportCmd
	}
	vaultExportCmd = &cobra.Command{
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
					secretOutputMap[name] = base64.StdEncoding.EncodeToString(secret)
				} else {
					secretOutputMap[name] = string(secret)
				}
			}
			exportFormat := cmd.Flag(vaultExportFormatFlag.name).Value.String()
			if exportFormat == "" {
				exportFormat = "envar"
			}
			switch exportFormat {
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
			case "envars", "envar", ".env":
				for key, value := range secretOutputMap {
					value = strings.ReplaceAll(value, "\\", "\\\\")
					value = strings.ReplaceAll(value, "\"", "\\\"")
					fmt.Printf("%s=\"%s\"\n", key, value)
				}
			default:
				exitOnErrorWithMessage("invalid format: " + exportFormat)
			}
			safeExit()
		},
	}
	vaultExportCmd.Flags().StringP(vaultExportFormatFlag.name, vaultExportFormatFlag.shorthand, "", vaultExportFormatFlag.usage)
	vaultExportCmd.Flags().BoolP(secretEncodeBase64Flag.name, secretEncodeBase64Flag.shorthand, false, secretEncodeBase64Flag.usage)
	return vaultExportCmd
}

func vaultRefCommand() *cobra.Command {
	if vaultRefCmd != nil {
		return vaultRefCmd
	}
	vaultRefCmd = &cobra.Command{
		Use:     "ref",
		Aliases: []string{"reference"},
		Short:   "References and updates secrets to a vault from a given yaml or json file",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				exitOnError(err)
			}
			refFile := cmd.Flag(vaultRefFileFlag.name).Value.String()
			secretNamePrefix := cmd.Flag(secretNameFlag.name).Value.String()
			refType := strings.ToLower(cmd.Flag(vaultRefTypeFlag.name).Value.String())
			previewOnly, _ := cmd.Flags().GetBool(secretRefPreviewOnlyFlag.name)
			forceUpdate, _ := cmd.Flags().GetBool(secretForceUpdateFlag.name)
			if secretNamePrefix == "" && refType == "" {
				exitOnErrorWithMessage("please provide at least one of --" + secretNameFlag.name + " or --" + vaultRefTypeFlag.name + " flag")
			}
			if refType != "" && refType != "yaml" {
				exitOnErrorWithMessage("only yaml auto reference is supported at the moment")
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
	vaultRefCmd.Flags().StringP(vaultRefFileFlag.name, vaultRefFileFlag.shorthand, "", vaultRefFileFlag.usage)
	vaultRefCmd.Flags().StringP(secretNameFlag.name, secretNameFlag.shorthand, "", secretNameFlag.usage)
	vaultRefCmd.Flags().StringP(vaultRefTypeFlag.name, vaultRefTypeFlag.shorthand, "", vaultRefTypeFlag.usage)
	vaultRefCmd.Flags().BoolP(secretRefPreviewOnlyFlag.name, secretRefPreviewOnlyFlag.shorthand, false, secretRefPreviewOnlyFlag.usage)
	vaultRefCmd.Flags().BoolP(secretForceUpdateFlag.name, secretForceUpdateFlag.shorthand, false, secretForceUpdateFlag.usage)
	vaultRefCmd.MarkFlagRequired(vaultRefFileFlag.name)
	return vaultRefCmd
}

func vaultDerefCommand() *cobra.Command {
	if vaultDerefCmd != nil {
		return vaultDerefCmd
	}
	vaultDerefCmd = &cobra.Command{
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
			paths, err := cmd.Flags().GetStringSlice(vaultDerefPathFlag.name)
			if err != nil {
				exitOnError(err)
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
				for _, path := range paths {
					if err = vault.DeRefSecrets(path); err != nil {
						exitOnError(err)
					}
					fmt.Println("Dereferenced", color.GreenString(path), "with the vault", color.GreenString(vaultFile))
				}
			}
			safeExit()
		},
	}
	vaultDerefCmd.Flags().StringSliceP(vaultFileFlag.name, vaultFileFlag.shorthand, []string{}, vaultFileFlag.usage)
	vaultDerefCmd.Flags().StringSliceP(vaultDerefPathFlag.name, vaultDerefPathFlag.shorthand, []string{}, vaultDerefPathFlag.usage)
	vaultDerefCmd.MarkFlagRequired(vaultDerefPathFlag.name)
	return vaultDerefCmd
}

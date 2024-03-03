package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"savesecrets.org/slv"
	"savesecrets.org/slv/core/commons"
	"savesecrets.org/slv/core/crypto"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/profiles"
	"savesecrets.org/slv/core/vaults"
)

const (
	k8sApiVersion = "slv.savesecrets.org/v1"
	k8sKind       = "SLV"
	k8sVaultField = "spec"
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
	vaultCmd.AddCommand(vaultInfoCommand())
	vaultCmd.AddCommand(vaultNewCommand())
	vaultCmd.AddCommand(vaultShareCommand())
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
			envMap := make(map[string]string, len(accessors))
			for _, accessor := range accessors {
				var env *environments.Environment
				envId := accessor.String()
				if profile != nil {
					env, err = profile.GetEnv(envId)
					if err != nil {
						exitOnError(err)
					}
				}
				if env != nil {
					envMap[envId] = envId + "\t(" + env.Name + ")"
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
	vaultInfoCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	vaultInfoCmd.MarkFlagRequired(vaultFileFlag.name)
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
	return
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
	vaultNewCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	vaultNewCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.name, vaultAccessPublicKeysFlag.shorthand, []string{}, vaultAccessPublicKeysFlag.usage)
	vaultNewCmd.Flags().StringP(envSearchFlag.name, envSearchFlag.shorthand, "", envSearchFlag.usage)
	vaultNewCmd.Flags().BoolP(envSelfFlag.name, envSelfFlag.shorthand, false, envSelfFlag.usage)
	vaultNewCmd.Flags().StringP(vaultK8sFlag.name, vaultK8sFlag.shorthand, "", vaultK8sFlag.usage)
	vaultNewCmd.Flags().BoolP(vaultEnableHashingFlag.name, vaultEnableHashingFlag.shorthand, false, vaultEnableHashingFlag.usage)
	vaultNewCmd.MarkFlagRequired(vaultFileFlag.name)
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
	vaultShareCmd.Flags().StringP(vaultFileFlag.name, vaultFileFlag.shorthand, "", vaultFileFlag.usage)
	vaultShareCmd.Flags().StringSliceP(vaultAccessPublicKeysFlag.name, vaultAccessPublicKeysFlag.shorthand, []string{}, vaultAccessPublicKeysFlag.usage)
	vaultShareCmd.Flags().StringP(envSearchFlag.name, envSearchFlag.shorthand, "", envSearchFlag.usage)
	vaultShareCmd.Flags().BoolP(envSelfFlag.name, envSelfFlag.shorthand, false, envSelfFlag.usage)
	vaultShareCmd.MarkFlagRequired(vaultFileFlag.name)
	return vaultShareCmd
}

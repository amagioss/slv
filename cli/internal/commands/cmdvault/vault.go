package cmdvault

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/cli/internal/commands/utils"
	"oss.amagi.com/slv/core/config"
	"oss.amagi.com/slv/core/environments"
	"oss.amagi.com/slv/core/profiles"
	"oss.amagi.com/slv/core/vaults"
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

func VaultCommand() *cobra.Command {
	if vaultCmd != nil {
		return vaultCmd
	}
	vaultCmd = &cobra.Command{
		Use:   "vault",
		Short: "Manage vaults and secrets in them",
		Long:  `Handle vault operations in SLV. SLV Vaults are files that store secrets.`,
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				utils.ExitOnError(err)
			}
			sealedSecretsMap, err := vault.ListSealedSecrets()
			if err != nil {
				utils.ExitOnError(err)
			}
			accessors, err := vault.ListAccessors()
			if err != nil {
				utils.ExitOnError(err)
			}
			profile, _ := profiles.GetDefaultProfile()
			self := environments.GetSelf()
			envMap := make(map[string]string, len(accessors))
			for _, accessor := range accessors {
				var env *environments.Environment
				accessorStr, err := accessor.String()
				if err != nil {
					utils.ExitOnError(err)
				}
				selfEnv := false
				rootEnv := false
				if self != nil && self.PublicKey == accessorStr {
					env = self
					selfEnv = true
				} else if profile != nil {
					env, err = profile.GetEnv(accessorStr)
					if err != nil {
						utils.ExitOnError(err)
					}
					if env == nil {
						root, err := profile.GetRoot()
						if err != nil {
							utils.ExitOnError(err)
						}
						if root != nil && root.PublicKey == accessorStr {
							rootEnv = true
							env = root
						}
					}
				}
				if env != nil {
					if selfEnv {
						envMap[accessorStr] = accessorStr + "\t(" + color.CyanString("Self"+": "+env.Name) + ")"
					} else if rootEnv {
						envMap[accessorStr] = accessorStr + "\t(" + color.CyanString("Root"+": "+env.Name) + ")"
					} else {
						envMap[accessorStr] = accessorStr + "\t(" + env.Name + ")"
					}
				} else {
					envMap[accessorStr] = accessorStr
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
			utils.SafeExit()
		},
	}
	vaultCmd.PersistentFlags().StringP(vaultFileFlag.Name, vaultFileFlag.Shorthand, "", vaultFileFlag.Usage)
	vaultCmd.MarkPersistentFlagRequired(vaultFileFlag.Name)
	vaultCmd.AddCommand(vaultNewCommand())
	vaultCmd.AddCommand(vaultSecretsCommand())
	vaultCmd.AddCommand(vaultPutCommand())
	vaultCmd.AddCommand(vaultGetCommand())
	vaultCmd.AddCommand(vaultDeleteCommand())
	vaultCmd.AddCommand(vaultRefCommand())
	vaultCmd.AddCommand(vaultDerefCommand())
	vaultCmd.AddCommand(vaultAccessCommand())
	return vaultCmd
}

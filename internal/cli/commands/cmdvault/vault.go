package cmdvault

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/profiles"
	"oss.amagi.com/slv/internal/core/vaults"
)

func getVault(filePath string) (*vaults.Vault, error) {
	return vaults.Get(filePath)
}

func VaultCommand() *cobra.Command {
	if vaultCmd != nil {
		return vaultCmd
	}
	vaultCmd = &cobra.Command{
		Use:     "vault",
		Aliases: []string{"v", "vaults", "secret", "secrets"},
		Short:   "Manage vaults/secrets with SLV",
		Long:    `Manage vaults/secrets using SLV. SLV Vaults are files that store secrets in a key-value format.`,
		Run: func(cmd *cobra.Command, args []string) {
			vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
			vault, err := getVault(vaultFile)
			if err != nil {
				utils.ExitOnError(err)
			}
			accessors, err := vault.ListAccessors()
			if err != nil {
				utils.ExitOnError(err)
			}
			profile, _ := profiles.GetDefaultProfile()
			self := environments.GetSelf()

			accessorTable := tablewriter.NewWriter(os.Stdout)
			accessorTable.SetHeader([]string{"Public Key", "Type", "Name"})
			accessorTable.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor})
			for _, accessor := range accessors {
				var env *environments.Environment
				accessorPK, err := accessor.String()
				if err != nil {
					utils.ExitOnError(err)
				}
				row := []string{accessorPK}
				selfEnv := false
				rootEnv := false
				if self != nil && self.PublicKey == accessorPK {
					env = self
					selfEnv = true
				} else if profile != nil {
					env, err = profile.GetEnv(accessorPK)
					if err != nil {
						utils.ExitOnError(err)
					}
					if env == nil {
						root, err := profile.GetRoot()
						if err != nil {
							utils.ExitOnError(err)
						}
						if root != nil && root.PublicKey == accessorPK {
							rootEnv = true
							env = root
						}
					}
				}
				if env != nil {
					if selfEnv {
						row = append(row, "Self")
					} else if rootEnv {
						row = append(row, "Root")
					} else {
						if env.EnvType == environments.USER {
							row = append(row, "User")
						} else {
							row = append(row, "Service")
						}
					}
					row = append(row, env.Name)
				} else {
					row = append(row, "Unknown", "")
				}
				accessorTable.Append(row)
			}
			fmt.Println("Vault ID: ", vault.Config.PublicKey)
			fmt.Println("Vault Data:")
			dataTable := tablewriter.NewWriter(os.Stdout)
			tableHeaderColors := []tablewriter.Colors{{tablewriter.Bold, tablewriter.FgHiWhiteColor},
				{tablewriter.Bold, tablewriter.FgHiWhiteColor},
				{tablewriter.Bold, tablewriter.FgHiWhiteColor}}
			hashAdded := false
			rows := [][]string{}
			secretsMap, ptMap := vault.List()
			for name, sealedSecret := range secretsMap {
				row := []string{name, "(Encrypted)", sealedSecret.EncryptedAt().Format("02-Jan-2006 15:04:05")}
				if hash := sealedSecret.Hash(); hash != "" {
					row = append(row, hash)
					hashAdded = true
				}
				rows = append(rows, row)
			}
			for name, value := range ptMap {
				rows = append(rows, []string{name, value, "N/A"})
			}
			header := []string{"Name", "Value", "Created At"}
			if hashAdded {
				header = append(header, "Hash")
				for i := range rows {
					if len(rows[i]) < 4 {
						rows[i] = append(rows[i], "")
					}
				}
				tableHeaderColors = append(tableHeaderColors, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor})
			}
			dataTable.SetHeader(header)
			dataTable.SetHeaderColor(tableHeaderColors...)
			dataTable.AppendBulk(rows)
			dataTable.Render()
			fmt.Println("Accessible by:")
			accessorTable.Render()
			utils.SafeExit()
		},
	}
	vaultCmd.PersistentFlags().StringP(vaultFileFlag.Name, vaultFileFlag.Shorthand, "", vaultFileFlag.Usage)
	vaultCmd.MarkPersistentFlagRequired(vaultFileFlag.Name)
	vaultCmd.AddCommand(vaultNewCommand())
	vaultCmd.AddCommand(vaultToK8sCommand())
	vaultCmd.AddCommand(vaultPutCommand())
	vaultCmd.AddCommand(vaultGetCommand())
	vaultCmd.AddCommand(vaultShellCommand())
	vaultCmd.AddCommand(vaultDeleteCommand())
	vaultCmd.AddCommand(vaultRefCommand())
	vaultCmd.AddCommand(vaultDerefCommand())
	vaultCmd.AddCommand(vaultAccessCommand())
	return vaultCmd
}

package cmdvault

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/helpers"
)

func showVault(vault *vaults.Vault) {
	vaultItemMap, err := vault.GetAllItems()
	if err != nil {
		utils.ExitOnError(err)
	}
	hashFound := false
	for _, data := range vaultItemMap {
		if hashFound = data.Hash() != ""; hashFound {
			break
		}
	}
	dataTable := table.NewWriter()
	dataTable.SetOutputMirror(os.Stdout)
	dataTableHeader := table.Row{
		text.Colors{text.Bold}.Sprint("Name"),
		text.Colors{text.Bold}.Sprint("Value"),
		text.Colors{text.Bold}.Sprint("Type"),
		text.Colors{text.Bold}.Sprint("Encrypted At"),
	}
	if hashFound {
		dataTableHeader = append(dataTableHeader, text.Colors{text.Bold}.Sprint("Hash"))
	}
	dataTable.AppendHeader(dataTableHeader)
	dataTableRows := make([]table.Row, 0, len(vaultItemMap))
	for name, item := range vaultItemMap {
		row := table.Row{name}
		if vault.IsLocked() {
			row = append(row, "(Locked)")
		} else {
			itemValueStr, err := item.ValueString()
			if err != nil {
				utils.ExitOnError(fmt.Errorf("error getting value for %s: %w", name, err))
			}
			row = append(row, itemValueStr)
		}
		if item.IsPlaintext() {
			row = append(row, "Plain Text")
		} else {
			row = append(row, "Secret")
		}
		if item.EncryptedAt() != nil {
			row = append(row, item.EncryptedAt().Format("02-Jan-2006 15:04:05"))
		} else {
			row = append(row, "N/A")
		}
		if hashFound {
			row = append(row, item.Hash())
		}
		dataTableRows = append(dataTableRows, row)
	}
	dataTable.AppendRows(dataTableRows)

	accessors, err := vault.ListAccessors()
	if err != nil {
		utils.ExitOnError(err)
	}
	profile, _ := profiles.GetActiveProfile()
	var root *environments.Environment
	if profile != nil {
		root, _ = profile.GetRoot()
	}
	self := environments.GetSelf()
	accessTable := table.NewWriter()
	accessTable.SetOutputMirror(os.Stdout)
	accessTable.AppendHeader(table.Row{
		text.Colors{text.Bold}.Sprint("Public Key"),
		text.Colors{text.Bold}.Sprint("Type"),
		text.Colors{text.Bold}.Sprint("Name"),
	})
	accessTableRows := make([]table.Row, 0, len(accessors))
	for _, accessor := range accessors {
		accessorPubKey, err := accessor.String()
		if err != nil {
			utils.ExitOnError(err)
		}
		row := table.Row{accessorPubKey}
		if self != nil && self.PublicKey == accessorPubKey {
			row = append(row, "Self", self.Name)
		} else if root != nil && root.PublicKey == accessorPubKey {
			row = append(row, "Root", root.Name)
		} else if profile != nil {
			if env, _ := profile.GetEnv(accessorPubKey); env != nil {
				if env.EnvType == environments.USER {
					row = append(row, "User", env.Name)
				} else {
					row = append(row, "Service", env.Name)
				}
			}
		}
		if len(row) < 3 {
			row = append(row, "Unknown", "")
		}
		accessTableRows = append(accessTableRows, row)
	}
	accessTable.AppendRows(accessTableRows)

	fmt.Printf("Vault Name: %s\n", vault.Name)
	fmt.Println("Vault Data:")
	dataTable.SetStyle(table.StyleLight)
	dataTable.Render()
	fmt.Println("Accessible by:")
	accessTable.SetStyle(table.StyleLight)
	accessTable.Render()
}

func VaultCommand() *cobra.Command {
	if vaultCmd == nil {
		vaultCmd = &cobra.Command{
			Use:     "vault",
			Aliases: []string{"v", "vaults", "secret", "secrets"},
			Short:   "Manage SLV vaults",
			Run: func(cmd *cobra.Command, args []string) {
				vaultFile := cmd.Flag(vaultFileFlag.Name).Value.String()
				vault, err := vaults.Get(vaultFile)
				if err != nil {
					utils.ExitOnError(err)
				}
				envSecretKey, _ := session.GetSecretKey()
				if envSecretKey != nil {
					vault.Unlock(envSecretKey)
				}
				showVault(vault)
			},
		}
		vaultCmd.PersistentFlags().StringP(vaultFileFlag.Name, vaultFileFlag.Shorthand, "", vaultFileFlag.Usage)
		vaultCmd.MarkPersistentFlagRequired(vaultFileFlag.Name)
		if err := vaultCmd.RegisterFlagCompletionFunc(vaultFileFlag.Name, vaultFilePathCompletion); err != nil {
			utils.ExitOnError(err)
		}
		vaultCmd.AddCommand(vaultListCommand())
		vaultCmd.AddCommand(vaultNewCommand())
		vaultCmd.AddCommand(vaultUpdateCommand())
		vaultCmd.AddCommand(vaultPutCommand())
		vaultCmd.AddCommand(vaultGetCommand())
		vaultCmd.AddCommand(vaultRunCommand())
		vaultCmd.AddCommand(vaultDeleteCommand())
		vaultCmd.AddCommand(vaultRefCommand())
		vaultCmd.AddCommand(vaultDerefCommand())
		vaultCmd.AddCommand(vaultAccessCommand())
	}
	return vaultCmd
}

func vaultFilePathCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if vaultFiles, err := helpers.ListVaultFiles("", true); err != nil {
		return nil, cobra.ShellCompDirectiveError
	} else {
		return vaultFiles, cobra.ShellCompDirectiveDefault
	}
}

func vaultItemNameCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if vault, err := vaults.Get(cmd.Flag(vaultFileFlag.Name).Value.String()); err == nil {
		return vault.GetItemNames(), cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveError
}

package cmdvault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/helpers"
)

var (
	vaultListCmd *cobra.Command

	listDirFlag = utils.FlagDef{
		Name:      "dir",
		Shorthand: "d",
		Usage:     "Directory to search for vaults (default: current directory)",
	}

	listRecursiveFlag = utils.FlagDef{
		Name:      "recursive",
		Shorthand: "r",
		Usage:     "Search for vaults recursively in subdirectories",
	}

	listShowDetailsFlag = utils.FlagDef{
		Name:      "details",
		Shorthand: "l",
		Usage:     "Show detailed information about each vault",
	}
)

func vaultListCommand() *cobra.Command {
	if vaultListCmd == nil {
		vaultListCmd = &cobra.Command{
			Use:     "list",
			Aliases: []string{"ls", "find", "search"},
			Short:   "List all vaults in a directory",
			Long: `List all vault files in the specified directory.
By default, lists vaults in the current directory.
Use --recursive to search in subdirectories.
Use --details to show vault metadata.`,
			PreRun: func(cmd *cobra.Command, args []string) {
				// The list command doesn't need the vault flag
				cmd.Parent().PersistentFlags().Lookup(vaultFileFlag.Name).Changed = true
			},
			Run: func(cmd *cobra.Command, args []string) {
				dir := cmd.Flag(listDirFlag.Name).Value.String()
				recursive, _ := cmd.Flags().GetBool(listRecursiveFlag.Name)
				showDetails, _ := cmd.Flags().GetBool(listShowDetailsFlag.Name)

				if dir == "" {
					var err error
					dir, err = os.Getwd()
					if err != nil {
						utils.ExitOnError(err)
					}
				}

				vaultFiles, err := helpers.ListVaultFiles(dir, recursive)
				if err != nil {
					utils.ExitOnError(err)
				}

				if len(vaultFiles) == 0 {
					fmt.Println("No vaults found in the specified directory.")
					utils.SafeExit()
				}

				if showDetails {
					displayVaultsWithDetails(dir, vaultFiles)
				} else {
					displayVaultsList(dir, vaultFiles)
				}
			},
		}

		vaultListCmd.Flags().StringP(listDirFlag.Name, listDirFlag.Shorthand, "", listDirFlag.Usage)
		vaultListCmd.Flags().BoolP(listRecursiveFlag.Name, listRecursiveFlag.Shorthand, false, listRecursiveFlag.Usage)
		vaultListCmd.Flags().BoolP(listShowDetailsFlag.Name, listShowDetailsFlag.Shorthand, false, listShowDetailsFlag.Usage)
	}
	return vaultListCmd
}

func displayVaultsList(baseDir string, vaultFiles []string) {
	fmt.Printf("Found %d vault(s):\n\n", len(vaultFiles))
	for _, vaultFile := range vaultFiles {
		fmt.Printf("  • %s\n", vaultFile)
	}
	fmt.Println()
	fmt.Println("Use 'slv vault list --details' to see more information.")
}

func displayVaultsWithDetails(baseDir string, vaultFiles []string) {
	fmt.Printf("Found %d vault(s):\n\n", len(vaultFiles))
	
	vaultTable := table.NewWriter()
	vaultTable.SetOutputMirror(os.Stdout)
	vaultTable.AppendHeader(table.Row{
		text.Colors{text.Bold}.Sprint("Vault File"),
		text.Colors{text.Bold}.Sprint("Name"),
		text.Colors{text.Bold}.Sprint("Namespace"),
		text.Colors{text.Bold}.Sprint("Secrets"),
		text.Colors{text.Bold}.Sprint("Accessors"),
	})

	for _, vaultFile := range vaultFiles {
		fullPath := filepath.Join(baseDir, vaultFile)
		vault, err := vaults.Get(fullPath)
		if err != nil {
			// If we can't load the vault, show basic info with error
			vaultTable.AppendRow(table.Row{
				vaultFile,
				text.Colors{text.FgRed}.Sprint("(Error loading)"),
				"-",
				"-",
				"-",
			})
			continue
		}

		// Get vault metadata
		vaultName := vault.Name
		if vaultName == "" {
			vaultName = "-"
		}

		namespace := vault.Namespace
		if namespace == "" {
			namespace = "-"
		}

		// Count secrets
		itemNames := vault.GetItemNames()
		secretCount := len(itemNames)

		// Count accessors
		accessors, err := vault.ListAccessors()
		accessorCount := 0
		if err == nil {
			accessorCount = len(accessors)
		}

		vaultTable.AppendRow(table.Row{
			vaultFile,
			vaultName,
			namespace,
			fmt.Sprintf("%d", secretCount),
			fmt.Sprintf("%d", accessorCount),
		})
	}
	
	// Table styling with proper wrapping for terminal readability
	vaultTable.SetStyle(table.StyleLight)
	vaultTable.Style().Options.SeparateRows = false
	vaultTable.Style().Options.SeparateColumns = true
	vaultTable.Style().Options.DrawBorder = true
	
	// Set column configurations with wrapping enabled inside cells
	vaultTable.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: false, WidthMax: 35, WidthMin: 20, Align: text.AlignLeft},      // Vault File
		{Number: 2, AutoMerge: false, WidthMax: 20, WidthMin: 10, Align: text.AlignLeft},      // Name
		{Number: 3, AutoMerge: false, WidthMax: 15, WidthMin: 10, Align: text.AlignLeft},      // Namespace
		{Number: 4, Align: text.AlignRight, WidthMax: 8, WidthMin: 8},                         // Secrets
		{Number: 5, Align: text.AlignRight, WidthMax: 10, WidthMin: 10},                       // Accessors
	})
	
	vaultTable.Render()
	fmt.Println()
}

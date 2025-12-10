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
	
	for i, vaultFile := range vaultFiles {
		fullPath := filepath.Join(baseDir, vaultFile)
		vault, err := vaults.Get(fullPath)
		
		// Display vault number and path
		fmt.Printf("%d. %s\n", i+1, text.Colors{text.Bold, text.FgCyan}.Sprint(vaultFile))
		
		if err != nil {
			fmt.Printf("   %s\n\n", text.Colors{text.FgRed}.Sprint("⚠ Error loading vault"))
			continue
		}
		
		// Get vault metadata
		vaultName := vault.Name
		if vaultName == "" {
			vaultName = text.Colors{text.Faint}.Sprint("(unnamed)")
		} else {
			vaultName = text.Colors{text.FgGreen}.Sprint(vaultName)
		}
		
		namespace := vault.Namespace
		if namespace == "" {
			namespace = text.Colors{text.Faint}.Sprint("(default)")
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
		
		// Display details in a readable format
		fmt.Printf("   Name:       %s\n", vaultName)
		fmt.Printf("   Namespace:  %s\n", namespace)
		fmt.Printf("   Secrets:    %s\n", text.Colors{text.Bold}.Sprintf("%d", secretCount))
		fmt.Printf("   Accessors:  %s\n", text.Colors{text.Bold}.Sprintf("%d", accessorCount))
		
		// Show accessor details if available
		if accessorCount > 0 && len(accessors) > 0 {
			fmt.Printf("   Accessible by: ")
			for idx, accessor := range accessors {
				if idx > 0 {
					fmt.Printf(", ")
				}
				accStr, err := accessor.String()
				if err == nil {
					// Show truncated accessor key
					if len(accStr) > 16 {
						fmt.Printf("%s", text.Colors{text.FgYellow}.Sprintf("%s...", accStr[:16]))
					} else {
						fmt.Printf("%s", text.Colors{text.FgYellow}.Sprint(accStr))
					}
				}
				if idx >= 2 { // Show max 3 accessors
					remaining := len(accessors) - 3
					if remaining > 0 {
						fmt.Printf(" %s", text.Colors{text.Faint}.Sprintf("+%d more", remaining))
					}
					break
				}
			}
			fmt.Println()
		}
		
		fmt.Println()
	}
	
	fmt.Println(text.Colors{text.Faint}.Sprint("💡 Tip: Use 'slv vault list --details' for table view"))
}

func displayVaultsWithDetails(baseDir string, vaultFiles []string) {
	fmt.Printf("Found %d vault(s):\n\n", len(vaultFiles))
	
	vaultTable := table.NewWriter()
	vaultTable.SetOutputMirror(os.Stdout)
	vaultTable.AppendHeader(table.Row{
		text.Colors{text.Bold}.Sprint("#"),
		text.Colors{text.Bold}.Sprint("Vault File"),
		text.Colors{text.Bold}.Sprint("Name"),
		text.Colors{text.Bold}.Sprint("Namespace"),
		text.Colors{text.Bold}.Sprint("Secrets"),
		text.Colors{text.Bold}.Sprint("Accessors"),
		text.Colors{text.Bold}.Sprint("Access List"),
	})

	for idx, vaultFile := range vaultFiles {
		fullPath := filepath.Join(baseDir, vaultFile)
		vault, err := vaults.Get(fullPath)
		if err != nil {
			// If we can't load the vault, show basic info with error
			vaultTable.AppendRow(table.Row{
				fmt.Sprintf("%d", idx+1),
				vaultFile,
				text.Colors{text.FgRed}.Sprint("(Error loading)"),
				"-",
				"-",
				"-",
				"-",
			})
			continue
		}

		// Get vault metadata
		vaultName := vault.Name
		if vaultName == "" {
			vaultName = text.Colors{text.Faint}.Sprint("(unnamed)")
		}

		namespace := vault.Namespace
		if namespace == "" {
			namespace = text.Colors{text.Faint}.Sprint("(default)")
		}

		// Count secrets
		itemNames := vault.GetItemNames()
		secretCount := len(itemNames)

		// Count accessors
		accessors, err := vault.ListAccessors()
		accessorCount := 0
		accessorList := "-"
		if err == nil {
			accessorCount = len(accessors)
			if accessorCount > 0 {
				if accessorCount <= 2 {
					accessorList = ""
					for i, acc := range accessors {
						if i > 0 {
							accessorList += ", "
						}
						accStr, err := acc.String()
						if err == nil {
							// Show first 12 characters of the accessor
							if len(accStr) > 12 {
								accessorList += accStr[:12] + "..."
							} else {
								accessorList += accStr
							}
						}
					}
				} else {
					// Show first accessor and count
					accStr, err := accessors[0].String()
					if err == nil {
						if len(accStr) > 12 {
							accessorList = fmt.Sprintf("%s... +%d more", accStr[:12], accessorCount-1)
						} else {
							accessorList = fmt.Sprintf("%s +%d more", accStr, accessorCount-1)
						}
					} else {
						accessorList = fmt.Sprintf("%d accessors", accessorCount)
					}
				}
			}
		}

		vaultTable.AppendRow(table.Row{
			fmt.Sprintf("%d", idx+1),
			vaultFile,
			vaultName,
			namespace,
			fmt.Sprintf("%d", secretCount),
			fmt.Sprintf("%d", accessorCount),
			accessorList,
		})
	}
	
	// Table styling with proper wrapping for terminal readability
	vaultTable.SetStyle(table.StyleLight)
	vaultTable.Style().Options.SeparateRows = false
	vaultTable.Style().Options.SeparateColumns = true
	vaultTable.Style().Options.DrawBorder = true
	
	// Set column configurations with wrapping enabled inside cells
	vaultTable.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignCenter, WidthMax: 4, WidthMin: 3},                        // #
		{Number: 2, AutoMerge: false, WidthMax: 30, WidthMin: 20, Align: text.AlignLeft},      // Vault File
		{Number: 3, AutoMerge: false, WidthMax: 18, WidthMin: 10, Align: text.AlignLeft},      // Name
		{Number: 4, AutoMerge: false, WidthMax: 12, WidthMin: 10, Align: text.AlignLeft},      // Namespace
		{Number: 5, Align: text.AlignRight, WidthMax: 8, WidthMin: 8},                         // Secrets
		{Number: 6, Align: text.AlignRight, WidthMax: 10, WidthMin: 10},                       // Accessors
		{Number: 7, AutoMerge: false, WidthMax: 25, WidthMin: 12, Align: text.AlignLeft},      // Access List
	})
	
	vaultTable.Render()
	
	// Add summary footer
	fmt.Println()
	totalSecrets := 0
	totalAccessors := 0
	for _, vaultFile := range vaultFiles {
		fullPath := filepath.Join(baseDir, vaultFile)
		vault, err := vaults.Get(fullPath)
		if err != nil {
			continue
		}
		totalSecrets += len(vault.GetItemNames())
		accessors, err := vault.ListAccessors()
		if err == nil {
			totalAccessors += len(accessors)
		}
	}
	
	fmt.Printf("📊 Summary: %s vaults | %s total secrets | %s unique accessors\n",
		text.Colors{text.Bold}.Sprintf("%d", len(vaultFiles)),
		text.Colors{text.Bold, text.FgGreen}.Sprintf("%d", totalSecrets),
		text.Colors{text.Bold, text.FgCyan}.Sprintf("%d", totalAccessors),
	)
	fmt.Println()
}

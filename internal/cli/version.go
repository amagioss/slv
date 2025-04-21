package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"slv.sh/slv/internal/core/config"
)

func versionCommand() *cobra.Command {
	if versionCmd == nil {
		versionCmd = &cobra.Command{
			Use:   "version",
			Short: "Show version information",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(config.VersionInfo())
			},
		}
	}
	return versionCmd
}

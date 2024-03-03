package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"savesecrets.org/slv"
)

func versionCommand() *cobra.Command {
	if versionCmd != nil {
		return versionCmd
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(slv.VersionInfo())
		},
	}
	return versionCmd
}

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"savesecrets.org/slv"
	"savesecrets.org/slv/core/config"
)

func showVersionInfo() {
	fmt.Println(config.AppNameUpperCase, slv.Version)
}

func versionCommand() *cobra.Command {
	if versionCmd != nil {
		return versionCmd
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			showVersionInfo()
		},
	}
	return versionCmd
}

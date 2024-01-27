package commands

import (
	"fmt"

	"github.com/amagimedia/slv"
	"github.com/spf13/cobra"
)

func showVersionInfo() {
	fmt.Println(slv.AppName, slv.Version)
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

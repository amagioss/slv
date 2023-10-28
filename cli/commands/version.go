package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"
var AppName = "app"

func showVersionInfo() {
	fmt.Println(AppName, Version)
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

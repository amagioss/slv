package commands

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"savesecrets.org/slv"
	"savesecrets.org/slv/core/config"
)

func showVersionInfo() {
	buildAt := "unknown"
	if builtAtTime, err := time.Parse(time.RFC3339, slv.BuildDate); err == nil {
		builtAtLocalTime := builtAtTime.Local()
		buildAt = builtAtLocalTime.Format("02 Jan 2006 03:04:05 PM MST")
	}
	fmt.Println(config.Art)
	fmt.Println(config.AppNameUpperCase + ": " + config.Description)
	fmt.Println("-------------------------------------------------")
	fmt.Printf("SLV Version : %s\n", slv.Version)
	fmt.Printf("Built At    : %s\n", buildAt)
	fmt.Printf("Git Commit  : %s\n", slv.Commit)
	fmt.Printf("Web         : %s\n", config.Website)
	fmt.Printf("Platform    : %s\n", runtime.GOOS+"/"+runtime.GOARCH)
	fmt.Printf("Go Version  : %s\n", runtime.Version())
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

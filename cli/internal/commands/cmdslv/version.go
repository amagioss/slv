package cmdslv

import (
	"fmt"

	"github.com/spf13/cobra"
	"oss.amagi.com/slv"
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

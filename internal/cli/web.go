package cli

import (
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/api"
)

func webCommand() *cobra.Command {
	if webCmd == nil {
		webCmd = &cobra.Command{
			Use:   "web",
			Short: "Starts the web server",
			Run: func(cmd *cobra.Command, args []string) {
				api.Run()
			},
		}
	}
	return webCmd
}

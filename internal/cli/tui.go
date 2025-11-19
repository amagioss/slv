package cli

import (
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/tui"
)

func tuiCommand() *cobra.Command {
	if tuiCmd == nil {
		tuiCmd = &cobra.Command{
			Use:   "tui",
			Short: "Starts the SLV TUI",
			Run: func(cmd *cobra.Command, args []string) {
				tui.RunTUIWithErrorHandling()
			},
		}
	}
	return tuiCmd
}

package cli

import (
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/api"
	"slv.sh/slv/internal/cli/commands/utils"
)

var (
	portFlag = utils.FlagDef{
		Name:      "port",
		Shorthand: "p",
		Usage:     "Port to serve the SLV Web Interface",
	}
)

func webCommand() *cobra.Command {
	if webCmd == nil {
		webCmd = &cobra.Command{
			Use:    "web",
			Short:  "Starts the SLV Web Interface",
			Hidden: true,
			Run: func(cmd *cobra.Command, args []string) {
				port, _ := cmd.Flags().GetUint16(portFlag.Name)
				api.Run(port)
			},
		}
		webCmd.Flags().Uint16P(portFlag.Name, portFlag.Shorthand, 0, portFlag.Usage)
	}
	return webCmd
}

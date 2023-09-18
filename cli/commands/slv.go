package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func SlvCommand() *cobra.Command {
	if slvCmd != nil {
		return slvCmd
	}
	slvCmd = &cobra.Command{
		Use:   "slv",
		Short: "SLV is a tool to encrypt secrets locally",
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := cmd.Flags().GetBool("version")
			if version {
				showVersionInfo()
			} else {
				cmd.Help()
			}
		},
	}
	slvCmd.Flags().BoolP("version", "v", false, "Shows version")
	slvCmd.AddCommand(systemCommand())
	slvCmd.AddCommand(envCommand())
	slvCmd.AddCommand(profileCommand())
	slvCmd.AddCommand(vaultCommand())
	slvCmd.AddCommand(secretCommand())
	return slvCmd
}

var Version = "dev"

func showVersionInfo() {
	fmt.Println("SLV", Version)
}

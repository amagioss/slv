package cmdslv

import (
	"fmt"

	"github.com/spf13/cobra"
	"oss.amagi.com/slv"
	"oss.amagi.com/slv/cli/internal/commands/cmdenv"
	"oss.amagi.com/slv/cli/internal/commands/cmdprofile"
	"oss.amagi.com/slv/cli/internal/commands/cmdsystem"
	"oss.amagi.com/slv/cli/internal/commands/cmdvault"
	"oss.amagi.com/slv/cli/internal/commands/utils"
)

var (
	slvCmd     *cobra.Command
	versionCmd *cobra.Command

	versionFlag = utils.FlagDef{
		Name:      "version",
		Shorthand: "v",
		Usage:     "Shows version info",
	}
)

func slvCommand() *cobra.Command {
	if slvCmd != nil {
		return slvCmd
	}
	slvCmd = &cobra.Command{
		Use:   "slv",
		Short: "SLV is a tool to encrypt secrets locally",
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := cmd.Flags().GetBool(versionFlag.Name)
			if version {
				fmt.Println(slv.VersionInfo())
			} else {
				cmd.Help()
			}
		},
	}
	slvCmd.Flags().BoolP(versionFlag.Name, versionFlag.Shorthand, false, versionFlag.Usage)
	slvCmd.AddCommand(versionCommand())
	slvCmd.AddCommand(cmdsystem.SystemCommand())
	slvCmd.AddCommand(cmdenv.EnvCommand())
	slvCmd.AddCommand(cmdprofile.ProfileCommand())
	slvCmd.AddCommand(cmdvault.VaultCommand())
	return slvCmd
}

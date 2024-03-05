package cmdslv

import (
	"fmt"

	"github.com/spf13/cobra"
	"savesecrets.org/slv"
	"savesecrets.org/slv/cli/internal/commands/cmdenv"
	"savesecrets.org/slv/cli/internal/commands/cmdprofile"
	"savesecrets.org/slv/cli/internal/commands/cmdsystem"
	"savesecrets.org/slv/cli/internal/commands/cmdvault"
	"savesecrets.org/slv/cli/internal/commands/utils"
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
	// slvCmd.AddCommand(secretCommand())
	return slvCmd
}

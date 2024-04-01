package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/cli/internal/commands/utils"
	"oss.amagi.com/slv/core/profiles"
)

func profileDeleteCommand() *cobra.Command {
	if profileDelCmd != nil {
		return profileDelCmd
	}
	profileDelCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "rm", "remove"},
		Short:   "Deletes a profile",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString(profileNameFlag.Name)
			if err := profiles.Delete(name); err != nil {
				utils.ExitOnError(err)
			} else {
				fmt.Println("Deleted profile: ", color.GreenString(name))
				utils.SafeExit()
			}
		},
	}
	profileDelCmd.Flags().StringP(profileNameFlag.Name, profileNameFlag.Shorthand, "", profileNameFlag.Usage)
	profileDelCmd.MarkFlagRequired(profileNameFlag.Name)
	return profileDelCmd
}

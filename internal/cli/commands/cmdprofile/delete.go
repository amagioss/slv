package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/profiles"
)

func profileDeleteCommand() *cobra.Command {
	if profileDelCmd == nil {
		profileDelCmd = &cobra.Command{
			Use:     "rm",
			Aliases: []string{"remove", "del", "delete"},
			Short:   "Removes a profile",
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
		if err := profileDelCmd.RegisterFlagCompletionFunc(profileNameFlag.Name, profileNameCompletion); err != nil {
			utils.ExitOnError(err)
		}
	}
	return profileDelCmd
}

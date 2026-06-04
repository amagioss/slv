package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/profiles"
)

func profileSyncCommand() *cobra.Command {
	if profileSyncCmd == nil {
		profileSyncCmd = &cobra.Command{
			Use:     "sync",
			Aliases: []string{"pull"},
			Short:   "Sync the active profile with its remote source",
			Run: func(cmd *cobra.Command, args []string) {
				profile, err := profiles.GetActiveProfile()
				if err != nil {
					utils.ExitOnError(err)
				}
				if err = profile.Pull(); err != nil {
					utils.ExitOnError(err)
				}
				fmt.Printf("Profile %s is updated from remote successfully\n", color.GreenString(profile.Name()))
			},
		}
	}
	return profileSyncCmd
}

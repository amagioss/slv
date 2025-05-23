package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/profiles"
)

func profileListCommand() *cobra.Command {
	if profileListCmd == nil {
		profileListCmd = &cobra.Command{
			Use:     "list",
			Aliases: []string{"ls"},
			Short:   "Lists all profiles",
			Run: func(cmd *cobra.Command, args []string) {
				profileNames, err := profiles.List()
				if err != nil {
					utils.ExitOnError(err)
				} else {
					activeProfileName, _ := profiles.GetActiveProfileName()
					for _, profileName := range profileNames {
						if profileName == activeProfileName {
							fmt.Println(color.GreenString(profileName))
						} else {
							fmt.Println(profileName)
						}
					}
				}
			},
		}
	}
	return profileListCmd
}

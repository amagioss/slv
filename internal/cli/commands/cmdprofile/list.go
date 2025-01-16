package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/profiles"
)

func profileListCommand() *cobra.Command {
	if profileListCmd == nil {
		profileListCmd = &cobra.Command{
			Use:   "list",
			Short: "Lists all profiles",
			Run: func(cmd *cobra.Command, args []string) {
				profileNames, err := profiles.List()
				if err != nil {
					utils.ExitOnError(err)
				} else {
					defaultProfileName, _ := profiles.GetDefaultProfileName()
					for _, profileName := range profileNames {
						if profileName == defaultProfileName {
							fmt.Println("*", color.GreenString(profileName))
						} else {
							fmt.Println(" ", profileName)
						}
					}
				}
			},
		}
	}
	return profileListCmd
}

package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"savesecrets.org/slv/cli/internal/commands/utils"
	"savesecrets.org/slv/core/profiles"
)

func profileListCommand() *cobra.Command {
	if profileListCmd != nil {
		return profileListCmd
	}
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
	return profileListCmd
}

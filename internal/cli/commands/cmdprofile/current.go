package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/profiles"
)

func profileSetCurrentCommand() *cobra.Command {
	if profileSetCurrentCmd == nil {
		profileSetCurrentCmd = &cobra.Command{
			Use:     "set",
			Aliases: []string{"set-current", "current", "set-default", "default"},
			Short:   "Set a profile as the active current profile",
			Run: func(cmd *cobra.Command, args []string) {
				profileNames, err := profiles.List()
				if err != nil {
					utils.ExitOnError(err)
				}
				name, _ := cmd.Flags().GetString(profileNameFlag.Name)
				for _, profileName := range profileNames {
					if profileName == name {
						profiles.SetCurrentProfile(name)
						fmt.Printf("Successfully set %s as current profile\n", color.GreenString(name))
						utils.SafeExit()
					}
				}
				utils.ExitOnError(fmt.Errorf("profile %s not found", name))
			},
		}
		profileSetCurrentCmd.Flags().StringP(profileNameFlag.Name, profileNameFlag.Shorthand, "", profileNameFlag.Usage)
		profileSetCurrentCmd.MarkFlagRequired(profileNameFlag.Name)
	}
	return profileSetCurrentCmd
}

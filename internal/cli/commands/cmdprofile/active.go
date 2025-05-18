package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/profiles"
)

func profileSetActiveCommand() *cobra.Command {
	if profileSetActiveCmd == nil {
		profileSetActiveCmd = &cobra.Command{
			Use:     "activate",
			Aliases: []string{"active", "set-active", "current", "default", "set"},
			Short:   "Set a profile as the active profile",
			Run: func(cmd *cobra.Command, args []string) {
				profileNames, err := profiles.List()
				if err != nil {
					utils.ExitOnError(err)
				}
				name, _ := cmd.Flags().GetString(profileNameFlag.Name)
				for _, profileName := range profileNames {
					if profileName == name {
						profiles.SetActiveProfile(name)
						fmt.Printf("Successfully set %s as the active profile\n", color.GreenString(name))
						utils.SafeExit()
					}
				}
				utils.ExitOnError(fmt.Errorf("profile %s not found", name))
			},
		}
		profileSetActiveCmd.Flags().StringP(profileNameFlag.Name, profileNameFlag.Shorthand, "", profileNameFlag.Usage)
		profileSetActiveCmd.MarkFlagRequired(profileNameFlag.Name)
		if err := profileSetActiveCmd.RegisterFlagCompletionFunc(profileNameFlag.Name, profileNameCompletion); err != nil {
			utils.ExitOnError(err)
		}
	}
	return profileSetActiveCmd
}

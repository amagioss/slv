package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/cli/internal/commands/utils"
	"oss.amagi.com/slv/core/profiles"
)

func profileDefaultCommand() *cobra.Command {
	if profileSetCmd != nil {
		return profileSetCmd
	}
	profileSetCmd = &cobra.Command{
		Use:     "default",
		Aliases: []string{"set-default"},
		Short:   "Set a profile as default",
		Run: func(cmd *cobra.Command, args []string) {
			profileNames, err := profiles.List()
			if err != nil {
				utils.ExitOnError(err)
			}
			name, _ := cmd.Flags().GetString(profileNameFlag.Name)
			for _, profileName := range profileNames {
				if profileName == name {
					profiles.SetDefault(name)
					fmt.Printf("Successfully set %s as default profile\n", color.GreenString(name))
					utils.SafeExit()
				}
			}
			utils.ExitOnError(fmt.Errorf("profile %s not found", name))
		},
	}
	profileSetCmd.Flags().StringP(profileNameFlag.Name, profileNameFlag.Shorthand, "", profileNameFlag.Usage)
	profileSetCmd.MarkFlagRequired(profileNameFlag.Name)
	return profileSetCmd
}

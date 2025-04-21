package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/profiles"
)

func profileNewCommand() *cobra.Command {
	if profileNewCmd == nil {
		profileNewCmd = &cobra.Command{
			Use:   "new",
			Short: "Creates a new profile",
			Run: func(cmd *cobra.Command, args []string) {
				name, _ := cmd.Flags().GetString(profileNameFlag.Name)
				gitURI, _ := cmd.Flags().GetString(profileGitURI.Name)
				gitBranch, _ := cmd.Flags().GetString(profileGitBranch.Name)
				err := profiles.New(name, gitURI, gitBranch)
				if err == nil {
					fmt.Println("Created profile: ", color.GreenString(name))
					utils.SafeExit()
				} else {
					utils.ExitOnError(err)
				}
			},
		}
		profileNewCmd.Flags().StringP(profileNameFlag.Name, profileNameFlag.Shorthand, "", profileNameFlag.Usage)
		profileNewCmd.Flags().StringP(profileGitURI.Name, profileGitURI.Shorthand, "", profileGitURI.Usage)
		profileNewCmd.Flags().StringP(profileGitBranch.Name, profileGitBranch.Shorthand, "", profileGitBranch.Usage)
		profileNewCmd.MarkFlagRequired(profileNameFlag.Name)
	}
	return profileNewCmd
}

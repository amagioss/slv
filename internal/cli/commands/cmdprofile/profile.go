package cmdprofile

import (
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/core/profiles"
)

func ProfileCommand() *cobra.Command {
	if profileCmd == nil {
		profileCmd = &cobra.Command{
			Use:     "profile",
			Aliases: []string{"profiles"},
			Short:   "Manage profiles and components within them",
			Long:    `Profile management along with environments and preferences within profiles are handled in this command`,
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		profileCmd.AddCommand(profileNewCommand())
		profileCmd.AddCommand(profileSetCurrentCommand())
		profileCmd.AddCommand(profileListCommand())
		profileCmd.AddCommand(profileDeleteCommand())
		profileCmd.AddCommand(profilePullCommand())
	}
	return profileCmd
}

func profileNameCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	profileNames, err := profiles.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	return profileNames, cobra.ShellCompDirectiveNoFileComp
}

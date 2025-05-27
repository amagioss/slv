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
			Short:   "Manage SLV profiles",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		profileCmd.AddCommand(profileNewCommand())
		profileCmd.AddCommand(profileSetActiveCommand())
		profileCmd.AddCommand(profileListCommand())
		profileCmd.AddCommand(profileDeleteCommand())
		profileCmd.AddCommand(profileSyncCommand())
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

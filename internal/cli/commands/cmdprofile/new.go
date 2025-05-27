package cmdprofile

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/profiles"
)

func profileNewCommand() *cobra.Command {
	if profileNewCmd == nil {
		profileNewCmd = &cobra.Command{
			Use:     "new",
			Aliases: []string{"setup"},
			Short:   "Sets up a new profile from a given remote",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
	}
	for _, remoteType := range profiles.ListRemoteNames() {
		profileNewCmd.AddCommand(getRemoteProfileCommand(remoteType))
	}
	return profileNewCmd
}

func getRemoteProfileCommand(remoteType string) *cobra.Command {
	remoteArgs := profiles.GetRemoteTypeArgs(remoteType)
	remoteProfileCommand := &cobra.Command{
		Use:   remoteType,
		Short: "Sets up a profile based on a remote profile (" + remoteType + ")",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString(profileNameFlag.Name)
			updateInterval, err := cmd.Flags().GetDuration(profileSyncInterval.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			remoteConfig := make(map[string]string)
			for _, arg := range remoteArgs {
				if value, _ := cmd.Flags().GetString(arg.Name()); value != "" {
					remoteConfig[arg.Name()] = value
				}
			}
			readOnly, _ := cmd.Flags().GetBool(profileReadOnlyFlag.Name)
			if err = profiles.New(name, remoteType, readOnly, updateInterval, remoteConfig); err != nil {
				utils.ExitOnError(err)
			}
			fmt.Printf("Created profile %s from remote (%s)\n", color.GreenString(name), color.GreenString(remoteType))
		},
	}
	remoteProfileCommand.Flags().StringP(profileNameFlag.Name, profileNameFlag.Shorthand, "", profileNameFlag.Usage)
	remoteProfileCommand.MarkFlagRequired(profileNameFlag.Name)
	remoteProfileCommand.Flags().BoolP(profileReadOnlyFlag.Name, "", false, profileReadOnlyFlag.Usage)
	remoteProfileCommand.Flags().DurationP(profileSyncInterval.Name, "", time.Hour, profileSyncInterval.Usage)
	for _, arg := range remoteArgs {
		remoteProfileCommand.Flags().StringP(arg.Name(), "", "", arg.Description())
		if arg.Required() {
			remoteProfileCommand.MarkFlagRequired(arg.Name())
		}
	}
	return remoteProfileCommand
}

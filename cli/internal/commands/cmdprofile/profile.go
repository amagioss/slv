package cmdprofile

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"savesecrets.org/slv/cli/internal/commands/utils"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/profiles"
)

func ProfileCommand() *cobra.Command {
	if profileCmd != nil {
		return profileCmd
	}
	profileCmd = &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles",
		Long:  `Manage profiles in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	profileCmd.AddCommand(profileNewCommand())
	profileCmd.AddCommand(profileDefaultCommand())
	profileCmd.AddCommand(profileListCommand())
	profileCmd.AddCommand(profileAddEnvCommand())
	profileCmd.AddCommand(profileDeleteCommand())
	profileCmd.AddCommand(profilePullCommand())
	profileCmd.AddCommand(profilePushCommand())
	return profileCmd
}

func profileNewCommand() *cobra.Command {
	if profileNewCmd != nil {
		return profileNewCmd
	}
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
	return profileNewCmd
}

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

func profileAddEnvCommand() *cobra.Command {
	if envAddCmd != nil {
		return envAddCmd
	}
	envAddCmd = &cobra.Command{
		Use:     "addenv",
		Aliases: []string{"add-env", "add-environment", "add-environments", "addenvs", "envadd", "env-add", "env-adds", "env-adds"},
		Short:   "Adds an environment to a profile",
		Run: func(cmd *cobra.Command, args []string) {
			envdefs, err := cmd.Flags().GetStringSlice(profileEnvDefFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			profileName := cmd.Flag(profileNameFlag.Name).Value.String()
			var profile *profiles.Profile
			if profileName != "" {
				profile, err = profiles.Get(profileName)
			} else {
				profile, err = profiles.GetDefaultProfile()
			}
			if err != nil {
				utils.ExitOnError(err)
			}
			setAsRoot, _ := cmd.Flags().GetBool(profileSetRootEnvFlag.Name)
			if setAsRoot && len(envdefs) > 1 {
				utils.ExitOnError(fmt.Errorf("cannot set more than one environment as root"))
			}
			var successMessage string
			for _, envdef := range envdefs {
				var env *environments.Environment
				if env, err = environments.FromEnvDef(envdef); err == nil && env != nil {
					if setAsRoot {
						err = profile.SetRoot(env)
						successMessage = fmt.Sprintf("Successfully set %s as root environment for profile %s", color.GreenString(env.Name), color.GreenString(profile.Name()))
					} else {
						err = profile.PutEnv(env)
					}
				}
				if err != nil {
					utils.ExitOnError(err)
				}
			}
			if successMessage == "" {
				successMessage = fmt.Sprintf("Successfully added %d environments to profile %s", len(envdefs), color.GreenString(profile.Name()))
			}
			fmt.Println(successMessage)
			utils.SafeExit()
		},
	}
	envAddCmd.Flags().StringP(profileNameFlag.Name, profileNameFlag.Shorthand, "", profileNameFlag.Usage)
	envAddCmd.Flags().StringSliceP(profileEnvDefFlag.Name, profileEnvDefFlag.Shorthand, []string{}, profileEnvDefFlag.Usage)
	envAddCmd.Flags().Bool(profileSetRootEnvFlag.Name, false, profileSetRootEnvFlag.Usage)
	envAddCmd.MarkFlagRequired(profileEnvDefFlag.Name)
	return envAddCmd
}

func profileDeleteCommand() *cobra.Command {
	if profileDelCmd != nil {
		return profileDelCmd
	}
	profileDelCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "rm", "remove"},
		Short:   "Deletes a profile",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString(profileNameFlag.Name)
			if err := profiles.Delete(name); err != nil {
				utils.ExitOnError(err)
			} else {
				fmt.Println("Deleted profile: ", color.GreenString(name))
				utils.SafeExit()
			}
		},
	}
	profileDelCmd.Flags().StringP(profileNameFlag.Name, profileNameFlag.Shorthand, "", profileNameFlag.Usage)
	profileDelCmd.MarkFlagRequired(profileNameFlag.Name)
	return profileDelCmd
}

func profilePullCommand() *cobra.Command {
	if profilePullCmd != nil {
		return profilePullCmd
	}
	profilePullCmd = &cobra.Command{
		Use:     "pull",
		Aliases: []string{"sync"},
		Short:   "Pulls the latest changes for the current profile from remote repository",
		Run: func(cmd *cobra.Command, args []string) {
			profile, err := profiles.GetDefaultProfile()
			if err != nil {
				utils.ExitOnError(err)
			}
			if err = profile.Pull(); err != nil {
				utils.ExitOnError(err)
			}
			fmt.Printf("Successfully pulled changes into profile: %s\n", color.GreenString(profile.Name()))
		},
	}
	return profilePullCmd
}

func profilePushCommand() *cobra.Command {
	if profilePushCmd != nil {
		return profilePushCmd
	}
	profilePushCmd = &cobra.Command{
		Use:   "push",
		Short: "Pushes the changes in the current profile to the pre-configured remote repository",
		Run: func(cmd *cobra.Command, args []string) {
			profile, err := profiles.GetDefaultProfile()
			if err != nil {
				utils.ExitOnError(err)
			}
			if err = profile.Push(); err != nil {
				utils.ExitOnError(err)
			}
			fmt.Printf("Successfully pushed changes from profile: %s\n", color.GreenString(profile.Name()))
		},
	}
	return profilePushCmd
}

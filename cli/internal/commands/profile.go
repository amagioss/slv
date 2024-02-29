package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/profiles"
)

func profileCommand() *cobra.Command {
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
			name, _ := cmd.Flags().GetString(profileNameFlag.name)
			gitURI, _ := cmd.Flags().GetString(profileGitURI.name)
			gitBranch, _ := cmd.Flags().GetString(profileGitBranch.name)
			err := profiles.New(name, gitURI, gitBranch)
			if err == nil {
				fmt.Println("Created profile: ", color.GreenString(name))
				safeExit()
			} else {
				exitOnError(err)
			}
		},
	}
	profileNewCmd.Flags().StringP(profileNameFlag.name, profileNameFlag.shorthand, "", profileNameFlag.usage)
	profileNewCmd.Flags().StringP(profileGitURI.name, profileGitURI.shorthand, "", profileGitURI.usage)
	profileNewCmd.Flags().StringP(profileGitBranch.name, profileGitBranch.shorthand, "", profileGitBranch.usage)
	profileNewCmd.MarkFlagRequired(profileNameFlag.name)
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
				exitOnError(err)
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
				exitOnError(err)
			}
			name, _ := cmd.Flags().GetString(profileNameFlag.name)
			for _, profileName := range profileNames {
				if profileName == name {
					profiles.SetDefault(name)
					fmt.Printf("Successfully set %s as default profile\n", color.GreenString(name))
					safeExit()
				}
			}
			exitOnError(fmt.Errorf("profile %s not found", name))
		},
	}
	profileSetCmd.Flags().StringP(profileNameFlag.name, profileNameFlag.shorthand, "", profileNameFlag.usage)
	profileSetCmd.MarkFlagRequired(profileNameFlag.name)
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
			envdefs, err := cmd.Flags().GetStringSlice(envDefFlag.name)
			if err != nil {
				exitOnError(err)
			}
			profileName := cmd.Flag(profileNameFlag.name).Value.String()
			var profile *profiles.Profile
			if profileName != "" {
				profile, err = profiles.Get(profileName)
			} else {
				profile, err = profiles.GetDefaultProfile()
			}
			if err != nil {
				exitOnError(err)
			}
			setAsRoot, _ := cmd.Flags().GetBool(profileSetRootEnvFlag.name)
			if setAsRoot && len(envdefs) > 1 {
				exitOnError(fmt.Errorf("cannot set more than one environment as root"))
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
					exitOnError(err)
				}
			}
			if successMessage == "" {
				successMessage = fmt.Sprintf("Successfully added %d environments to profile %s", len(envdefs), color.GreenString(profile.Name()))
			}
			fmt.Println(successMessage)
			safeExit()
		},
	}
	envAddCmd.Flags().StringP(profileNameFlag.name, profileNameFlag.shorthand, "", profileNameFlag.usage)
	envAddCmd.Flags().StringSliceP(envDefFlag.name, envDefFlag.shorthand, []string{}, envDefFlag.usage)
	envAddCmd.Flags().Bool(profileSetRootEnvFlag.name, false, profileSetRootEnvFlag.usage)
	envAddCmd.MarkFlagRequired(envDefFlag.name)
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
			name, _ := cmd.Flags().GetString(profileNameFlag.name)
			if err := profiles.Delete(name); err != nil {
				exitOnError(err)
			} else {
				fmt.Println("Deleted profile: ", color.GreenString(name))
				safeExit()
			}
		},
	}
	profileDelCmd.Flags().StringP(profileNameFlag.name, profileNameFlag.shorthand, "", profileNameFlag.usage)
	profileDelCmd.MarkFlagRequired(profileNameFlag.name)
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
				exitOnError(err)
			}
			if err = profile.Pull(); err != nil {
				exitOnError(err)
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
				exitOnError(err)
			}
			if err = profile.Push(); err != nil {
				exitOnError(err)
			}
			fmt.Printf("Successfully pushed changes from profile: %s\n", color.GreenString(profile.Name()))
		},
	}
	return profilePushCmd
}

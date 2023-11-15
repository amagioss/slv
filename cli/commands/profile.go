package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/shibme/slv/core/environments"
	"github.com/shibme/slv/core/profiles"
	"github.com/spf13/cobra"
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
			err := profiles.New(name)
			if err == nil {
				fmt.Println("Created profile: ", color.GreenString(name))
				safeExit()
			} else {
				exitOnError(err)
			}
		},
	}
	profileNewCmd.Flags().StringP(profileNameFlag.name, profileNameFlag.shorthand, "", profileNameFlag.usage)
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
			envdefs, err := cmd.Flags().GetStringSlice(profileEnvDefFlag.name)
			if err != nil {
				exitOnError(err)
			}
			profileName := cmd.Flag(profileNameFlag.name).Value.String()
			var prof *profiles.Profile
			if profileName != "" {
				prof, err = profiles.GetProfile(profileName)
			} else {
				prof, err = profiles.GetDefaultProfile()
			}
			if err != nil {
				exitOnError(err)
			}
			setAsRoot, _ := cmd.Flags().GetBool(profileSetRootEnvFlag.name)
			if setAsRoot && len(envdefs) > 1 {
				exitOnError(fmt.Errorf("cannot set more than one environment as root"))
			}
			for _, envdef := range envdefs {
				var env *environments.Environment
				if env, err = environments.FromEnvDef(envdef); err == nil && env != nil {
					if setAsRoot {
						err = prof.SetRoot(env)
					} else {
						err = prof.AddEnv(env)
					}
				}
				if err != nil {
					exitOnError(err)
				}
			}
		},
	}
	envAddCmd.Flags().StringP(profileNameFlag.name, profileNameFlag.shorthand, "", profileNameFlag.usage)
	envAddCmd.Flags().StringSliceP(profileEnvDefFlag.name, profileEnvDefFlag.shorthand, []string{}, profileEnvDefFlag.usage)
	envAddCmd.Flags().Bool(profileSetRootEnvFlag.name, false, profileSetRootEnvFlag.usage)
	envAddCmd.MarkFlagRequired(profileEnvDefFlag.name)
	return envAddCmd
}

package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
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
	profileCmd.AddCommand(profileInitRootCommand())
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
				os.Exit(0)
			} else {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
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
				fmt.Println(err)
				os.Exit(1)
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
				PrintErrorAndExit(err)
			}
			name, _ := cmd.Flags().GetString(profileNameFlag.name)
			for _, profileName := range profileNames {
				if profileName == name {
					profiles.SetDefault(name)
					fmt.Printf("Successfully set %s as default profile\n", color.GreenString(name))
					os.Exit(0)
				}
			}
			PrintErrorAndExit(fmt.Errorf("profile %s not found", name))
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
				PrintErrorAndExit(err)
			}
			profileName := cmd.Flag(profileNameFlag.name).Value.String()
			var prof *profiles.Profile
			if profileName != "" {
				prof, err = profiles.GetProfile(profileName)
			} else {
				prof, err = profiles.GetDefaultProfile()
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
			for _, envdef := range envdefs {
				err = prof.AddEnvDef(envdef)
				if err != nil {
					PrintErrorAndExit(err)
				}
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
		},
	}
	envAddCmd.Flags().StringP(profileNameFlag.name, profileNameFlag.shorthand, "", profileNameFlag.usage)
	envAddCmd.Flags().StringSliceP(profileEnvDefFlag.name, profileEnvDefFlag.shorthand, []string{}, profileEnvDefFlag.usage)
	envAddCmd.MarkFlagRequired(profileEnvDefFlag.name)
	return envAddCmd
}

func profileInitRootCommand() *cobra.Command {
	if profileInitRootCmd != nil {
		return profileInitRootCmd
	}
	profileInitRootCmd = &cobra.Command{
		Use:     "initroot",
		Aliases: []string{"rootinit", "init-root", "root-init"},
		Short:   "Initializes the root environment in a profile",
		Run: func(cmd *cobra.Command, args []string) {
			profileName := cmd.Flag(profileNameFlag.name).Value.String()
			var prof *profiles.Profile
			var err error
			if profileName != "" {
				prof, err = profiles.GetProfile(profileName)
			} else {
				prof, err = profiles.GetDefaultProfile()
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
			secretKey, err := prof.InitRoot()
			if err != nil {
				PrintErrorAndExit(err)
			}
			fmt.Println("Root environment initialized with secret key:", secretKey)
		},
	}
	profileInitRootCmd.Flags().StringP(profileNameFlag.name, profileNameFlag.shorthand, "", profileNameFlag.usage)
	return profileInitRootCmd
}

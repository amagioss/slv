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
	profileCmd.AddCommand(profileSetCommand())
	profileCmd.AddCommand(profileListCommand())
	return profileCmd
}

func profileNewCommand() *cobra.Command {
	if profileNewCmd != nil {
		return profileNewCmd
	}
	profileNewCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new profile",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
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
	profileNewCmd.Flags().StringP("name", "n", "", "Name for the profile")
	profileNewCmd.MarkFlagRequired("name")
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

func profileSetCommand() *cobra.Command {
	if profileSetCmd != nil {
		return profileSetCmd
	}
	profileSetCmd = &cobra.Command{
		Use:     "default",
		Aliases: []string{"set-default"},
		Short:   "Set a profile as default profile",
		Run: func(cmd *cobra.Command, args []string) {
			profileNames, err := profiles.List()
			if err != nil {
				PrintErrorAndExit(err)
			}
			name, _ := cmd.Flags().GetString("name")
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
	profileSetCmd.Flags().StringP("name", "n", "", "Name of the profile to be set as default")
	profileSetCmd.MarkFlagRequired("name")
	return profileSetCmd
}

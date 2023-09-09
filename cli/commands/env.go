package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shibme/slv/core/environments"
	"github.com/shibme/slv/core/profiles"
	"github.com/spf13/cobra"
)

func envCommand() *cobra.Command {
	if envCmd != nil {
		return envCmd
	}
	envCmd = &cobra.Command{
		Use:   "env",
		Short: "Environment operations",
		Long:  `Environment operations in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	envCmd.AddCommand(envNewCommand())
	envCmd.AddCommand(envListCommand())
	envCmd.AddCommand(envUserRegisterCommand())
	return envCmd
}

func envNewCommand() *cobra.Command {
	if envNewCmd != nil {
		return envNewCmd
	}
	envNewCmd = &cobra.Command{
		Use:   "new",
		Short: "Creates a service environment",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString(envNameFlag.name)
			email, _ := cmd.Flags().GetString(envEmailFlag.name)
			tags, err := cmd.Flags().GetStringSlice(envTagsFlag.name)
			if err != nil {
				PrintErrorAndExit(err)
				os.Exit(1)
			}
			env, secretKey, _ := environments.New(name, email, environments.SERVICE)
			env.AddTags(tags...)
			envDef, err := env.ToEnvDef()
			if err != nil {
				PrintErrorAndExit(err)
				os.Exit(1)
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
			fmt.Fprintln(w, "Secret Key:\t", secretKey)
			fmt.Fprintln(w)
			fmt.Fprintln(w, "Public Key:\t", env.PublicKey)
			fmt.Fprintln(w, "Name:\t", env.Name)
			fmt.Fprintln(w, "Email:\t", env.Email)
			fmt.Fprintln(w, "Tags:\t", env.Tags)
			fmt.Fprintln(w)
			fmt.Fprintln(w, "Environment Definition:\t", envDef)
			w.Flush()

			// Adding env to a specified profile
			addToProfileFlag, _ := cmd.Flags().GetBool(envAddFlag.name)
			var cfg *profiles.Profile
			if addToProfileFlag {
				cfg, err = profiles.GetDefaultProfile()
				if err != nil {
					PrintErrorAndExit(err)
				}
				err = cfg.AddEnv(envDef)
				if err != nil {
					PrintErrorAndExit(err)
				}
			}
			os.Exit(0)
		},
	}
	envNewCmd.Flags().StringP(envNameFlag.name, envNameFlag.shorthand, "", envNameFlag.usage)
	envNewCmd.Flags().StringP(envEmailFlag.name, envEmailFlag.shorthand, "", envEmailFlag.usage)
	envNewCmd.Flags().StringSliceP(envTagsFlag.name, envTagsFlag.shorthand, []string{}, envTagsFlag.usage)
	envNewCmd.Flags().BoolP(envAddFlag.name, envAddFlag.shorthand, false, envAddFlag.usage)
	envNewCmd.MarkFlagRequired(envNameFlag.name)
	return envNewCmd
}

func envListCommand() *cobra.Command {
	if envListCmd != nil {
		return envListCmd
	}
	envListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists environments from profile",
		Run: func(cmd *cobra.Command, args []string) {
			profileName := cmd.Flag(profileNameFlag.name).Value.String()
			var cfg *profiles.Profile
			var err error
			if profileName != "" {
				cfg, err = profiles.GetProfile(profileName)
			} else {
				cfg, err = profiles.GetDefaultProfile()
			}
			if err != nil {
				PrintErrorAndExit(err)
			}
			envManifest, err := cfg.GetEnvManifest()
			if err != nil {
				PrintErrorAndExit(err)
			}
			query := cmd.Flag(envSearchFlag.name).Value.String()
			var envs []*environments.Environment
			if query != "" {
				envs = envManifest.SearchEnv(query)
			} else {
				envs = envManifest.ListEnv()
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
			for _, env := range envs {
				fmt.Fprintln(w, env.Id()+":")
				fmt.Fprintln(w, "Public Key:\t", env.PublicKey)
				fmt.Fprintln(w, "Name:\t", env.Name)
				fmt.Fprintln(w, "Email:\t", env.Email)
				fmt.Fprintln(w, "Tags:\t", env.Tags)
				fmt.Fprintln(w)
			}
			w.Flush()
			os.Exit(0)

		},
	}
	envListCmd.Flags().StringP(profileNameFlag.name, profileNameFlag.shorthand, "", profileNameFlag.usage)
	envListCmd.Flags().StringP(envSearchFlag.name, envSearchFlag.shorthand, "", envSearchFlag.usage)
	return envListCmd
}

func envUserRegisterCommand() *cobra.Command {
	if envUserRegisterCmd != nil {
		return envUserRegisterCmd
	}
	envUserRegisterCmd = &cobra.Command{
		Use:   "new",
		Short: "Creates an environment for user and registers it locally",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString(envNameFlag.name)
			email, _ := cmd.Flags().GetString(envEmailFlag.name)
			tags, err := cmd.Flags().GetStringSlice(envTagsFlag.name)
			if err != nil {
				PrintErrorAndExit(err)
				os.Exit(1)
			}
			env, secretKey, _ := environments.New(name, email, environments.USER)
			env.AddTags(tags...)
			envDef, err := env.ToEnvDef()
			if err != nil {
				PrintErrorAndExit(err)
				os.Exit(1)
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
			fmt.Fprintln(w, "Secret Key:\t", secretKey)
			fmt.Fprintln(w)
			fmt.Fprintln(w, "Public Key:\t", env.PublicKey)
			fmt.Fprintln(w, "Name:\t", env.Name)
			fmt.Fprintln(w, "Email:\t", env.Email)
			fmt.Fprintln(w, "Tags:\t", env.Tags)
			fmt.Fprintln(w)
			fmt.Fprintln(w, "Environment Definition:\t", envDef)
			w.Flush()

			// Adding env to a specified profile
			addToProfileFlag, _ := cmd.Flags().GetBool(envAddFlag.name)
			var cfg *profiles.Profile
			if addToProfileFlag {
				cfg, err = profiles.GetDefaultProfile()
				if err != nil {
					PrintErrorAndExit(err)
				}
				err = cfg.AddEnv(envDef)
				if err != nil {
					PrintErrorAndExit(err)
				}
			}
			os.Exit(0)
		},
	}
	envUserRegisterCmd.Flags().StringP(envNameFlag.name, envNameFlag.shorthand, "", envNameFlag.usage)
	envUserRegisterCmd.Flags().StringP(envEmailFlag.name, envEmailFlag.shorthand, "", envEmailFlag.usage)
	envUserRegisterCmd.Flags().StringSliceP(envTagsFlag.name, envTagsFlag.shorthand, []string{}, envTagsFlag.usage)
	envUserRegisterCmd.Flags().BoolP(envAddFlag.name, envAddFlag.shorthand, false, envAddFlag.usage)
	envUserRegisterCmd.MarkFlagRequired(envNameFlag.name)
	return envUserRegisterCmd
}

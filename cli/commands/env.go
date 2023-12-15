package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/shibme/slv/core/crypto"
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
	return envCmd
}

func showEnv(env environments.Environment, includeEDS bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "ID (Public Key):\t", env.PublicKey)
	fmt.Fprintln(w, "Name:\t", env.Name)
	fmt.Fprintln(w, "Email:\t", env.Email)
	fmt.Fprintln(w, "Tags:\t", env.Tags)
	if env.ProviderBinding != "" {
		fmt.Fprintln(w, "Provider Binding:\t", env.ProviderBinding)
	}
	if includeEDS {
		if envDef, err := env.ToEnvData(); err == nil {
			fmt.Fprintln(w, "\nEnv Data:\t", color.CyanString(envDef))
		}
	}
	w.Flush()
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
				exitOnError(err)
			}
			userEnv, _ := cmd.Flags().GetBool(envSelfFlag.name)
			envType := environments.SERVICE

			var env *environments.Environment
			var secretKey *crypto.SecretKey
			if userEnv {
				envType = environments.USER
			}
			env, secretKey, err = environments.NewEnvironment(name, envType)
			if err != nil {
				exitOnError(err)
			}
			env.SetEmail(email)
			env.AddTags(tags...)
			showEnv(*env, true)
			if secretKey != nil {
				fmt.Println("\nSecret Key:\t", color.HiBlackString(secretKey.String()))
			}
			addToProfileFlag, _ := cmd.Flags().GetBool(envAddFlag.name)
			var prof *profiles.Profile
			if addToProfileFlag {
				prof, err = profiles.GetDefaultProfile()
				if err != nil {
					exitOnError(err)
				}
				err = prof.AddEnv(env)
				if err != nil {
					exitOnError(err)
				}
			}
			safeExit()
		},
	}

	envNewCmd.Flags().StringP(envNameFlag.name, envNameFlag.shorthand, "", envNameFlag.usage)
	envNewCmd.Flags().StringP(envEmailFlag.name, envEmailFlag.shorthand, "", envEmailFlag.usage)
	envNewCmd.Flags().StringSliceP(envTagsFlag.name, envTagsFlag.shorthand, []string{}, envTagsFlag.usage)
	envNewCmd.Flags().BoolP(envAddFlag.name, envAddFlag.shorthand, false, envAddFlag.usage)
	envNewCmd.Flags().BoolP(envSelfFlag.name, envSelfFlag.shorthand, false, envSelfFlag.usage)
	envNewCmd.MarkFlagRequired(envNameFlag.name)

	envNewCmd.AddCommand(newKMSEnvCommand("aws", "Create environment for AWS KMS using RSA-4096 public key", kmsAWSARNFlag, kmsRSAPublicKey))

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
			var prof *profiles.Profile
			var err error
			if profileName != "" {
				prof, err = profiles.GetProfile(profileName)
			} else {
				prof, err = profiles.GetDefaultProfile()
			}
			if err != nil {
				exitOnError(err)
			}
			envManifest, err := prof.GetEnvManifest()
			if err != nil {
				exitOnError(err)
			}
			query := cmd.Flag(envSearchFlag.name).Value.String()
			var envs []*environments.Environment
			if query != "" {
				envs = envManifest.SearchEnv(query)
			} else {
				envs = envManifest.ListEnv()
			}
			for _, env := range envs {
				showEnv(*env, false)
				fmt.Println()
			}
			safeExit()
		},
	}
	envListCmd.Flags().StringP(profileNameFlag.name, profileNameFlag.shorthand, "", profileNameFlag.usage)
	envListCmd.Flags().StringP(envSearchFlag.name, envSearchFlag.shorthand, "", envSearchFlag.usage)
	return envListCmd
}

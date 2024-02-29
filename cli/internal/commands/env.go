package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"savesecrets.org/slv/core/crypto"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/input"
	"savesecrets.org/slv/core/profiles"
)

func envCommand() *cobra.Command {
	if envCmd != nil {
		return envCmd
	}
	envCmd = &cobra.Command{
		Use:     "env",
		Aliases: []string{"envs", "environment", "environments"},
		Short:   "Environment operations",
		Long:    `Environment operations in SLV`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	envCmd.AddCommand(envNewCommand())
	envCmd.AddCommand(envListCommand())
	envCmd.AddCommand(envSelfCommand())
	return envCmd
}

func showEnv(env environments.Environment, includeEDS bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "ID (Public Key):\t", env.PublicKey)
	fmt.Fprintln(w, "Name:\t", env.Name)
	fmt.Fprintln(w, "Email:\t", env.Email)
	fmt.Fprintln(w, "Tags:\t", env.Tags)
	if env.SecretBinding != "" {
		fmt.Fprintln(w, "Secret Binding:\t", env.SecretBinding)
	}
	if includeEDS {
		if envDef, err := env.ToEnvDef(); err == nil {
			fmt.Fprintln(w, "\nEnv Definition:\t", color.CyanString(envDef))
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
		Short: "Create a new environment",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	envNewCmd.AddCommand(envNewServiceCommand())
	envNewCmd.AddCommand(envNewUserCommand())
	envNewCmd.AddCommand(newKMSEnvCommand("aws", "Create an environment that works with AWS KMS", awsARNFlag))
	return envNewCmd
}

func envNewServiceCommand() *cobra.Command {
	if envNewServiceCmd != nil {
		return envNewServiceCmd
	}
	envNewServiceCmd = &cobra.Command{
		Use:     "service",
		Aliases: []string{"serv", "svc", "s"},
		Short:   "Creates a new service environment",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString(envNameFlag.name)
			email, _ := cmd.Flags().GetString(envEmailFlag.name)
			tags, err := cmd.Flags().GetStringSlice(envTagsFlag.name)
			if err != nil {
				exitOnError(err)
			}
			var env *environments.Environment
			var secretKey *crypto.SecretKey
			env, secretKey, err = environments.NewEnvironment(name, environments.SERVICE)
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
			if addToProfileFlag {
				profile, err := profiles.GetDefaultProfile()
				if err != nil {
					exitOnError(err)
				}
				err = profile.PutEnv(env)
				if err != nil {
					exitOnError(err)
				}
			}
			safeExit()
		},
	}
	envNewServiceCmd.Flags().StringP(envNameFlag.name, envNameFlag.shorthand, "", envNameFlag.usage)
	envNewServiceCmd.Flags().StringP(envEmailFlag.name, envEmailFlag.shorthand, "", envEmailFlag.usage)
	envNewServiceCmd.Flags().StringSliceP(envTagsFlag.name, envTagsFlag.shorthand, []string{}, envTagsFlag.usage)
	envNewServiceCmd.Flags().BoolP(envAddFlag.name, envAddFlag.shorthand, false, envAddFlag.usage)
	envNewServiceCmd.MarkFlagRequired(envNameFlag.name)
	return envNewServiceCmd
}

func envNewUserCommand() *cobra.Command {
	if envNewUserCmd != nil {
		return envNewUserCmd
	}
	envNewUserCmd = &cobra.Command{
		Use:     "user",
		Aliases: []string{"usr", "u"},
		Short:   "Register as a new user environment",
		Run: func(cmd *cobra.Command, args []string) {
			envName, _ := cmd.Flags().GetString(envNameFlag.name)
			envEmail, _ := cmd.Flags().GetString(envEmailFlag.name)
			envTags, err := cmd.Flags().GetStringSlice(envTagsFlag.name)
			if err != nil {
				exitOnError(err)
			}
			inputs := make(map[string][]byte)
			password, err := input.GetPasswordFromUser(true, input.GetDefaultPasswordPolicy())
			if err != nil {
				exitOnError(err)
			}
			inputs["password"] = password
			var env *environments.Environment
			env, err = environments.NewEnvForProvider("password", envName, environments.USER, inputs)
			if err != nil {
				exitOnError(err)
			}
			env.SetEmail(envEmail)
			env.AddTags(envTags...)
			if err = env.MarkAsSelf(); err != nil {
				exitOnError(err)
			}
			env.SecretBinding = ""
			showEnv(*env, true)
			addToProfileFlag, _ := cmd.Flags().GetBool(envAddFlag.name)
			if addToProfileFlag {
				profile, err := profiles.GetDefaultProfile()
				if err != nil {
					exitOnError(err)
				}
				err = profile.PutEnv(env)
				if err != nil {
					exitOnError(err)
				}
			}
			safeExit()
		},
	}
	envNewUserCmd.Flags().StringP(envNameFlag.name, envNameFlag.shorthand, "", envNameFlag.usage)
	envNewUserCmd.Flags().StringP(envEmailFlag.name, envEmailFlag.shorthand, "", envEmailFlag.usage)
	envNewUserCmd.Flags().StringSliceP(envTagsFlag.name, envTagsFlag.shorthand, []string{}, envTagsFlag.usage)
	envNewUserCmd.Flags().BoolP(envAddFlag.name, envAddFlag.shorthand, false, envAddFlag.usage)
	envNewUserCmd.MarkFlagRequired(envNameFlag.name)
	return envNewUserCmd
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
			var profile *profiles.Profile
			var err error
			if profileName != "" {
				profile, err = profiles.Get(profileName)
			} else {
				profile, err = profiles.GetDefaultProfile()
			}
			if err != nil {
				exitOnError(err)
			}
			query := cmd.Flag(envSearchFlag.name).Value.String()
			var envs []*environments.Environment
			if query != "" {
				envs, err = profile.SearchEnvs(query)
			} else {
				envs, err = profile.ListEnvs()
			}
			if err != nil {
				exitOnError(err)
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

func envSelfCommand() *cobra.Command {
	if envSelfCmd != nil {
		return envSelfCmd
	}
	envSelfCmd = &cobra.Command{
		Use:     "self",
		Aliases: []string{"me", "my", "current"},
		Short:   "Shows the current environment if registered",
		Run: func(cmd *cobra.Command, args []string) {
			env := environments.GetSelf()
			if env == nil {
				fmt.Println("No environment registered as self.")
			} else {
				showBinding, _ := cmd.Flags().GetBool(envShowBindingFlag.name)
				if !showBinding {
					env.SecretBinding = ""
				}
				showDef, _ := cmd.Flags().GetBool(envShowDefFlag.name)
				showEnv(*env, showDef)
			}
			safeExit()
		},
	}
	envSelfCmd.Flags().BoolP(envShowBindingFlag.name, envShowBindingFlag.shorthand, false, envShowBindingFlag.usage)
	envSelfCmd.Flags().BoolP(envShowDefFlag.name, envShowDefFlag.shorthand, false, envShowDefFlag.usage)
	envSelfCmd.AddCommand(envSelfSetCommand())
	return envSelfCmd
}

func envSelfSetCommand() *cobra.Command {
	if envSelfSetCmd != nil {
		return envSelfSetCmd
	}
	envSelfSetCmd = &cobra.Command{
		Use:     "set",
		Aliases: []string{"save", "put", "store", "s"},
		Short:   "Shows the current environment if registered",
		Run: func(cmd *cobra.Command, args []string) {
			envDef := cmd.Flag(envDefFlag.name).Value.String()
			env, err := environments.FromEnvDef(envDef)
			if err != nil {
				exitOnError(err)
			}
			if err = env.MarkAsSelf(); err != nil {
				exitOnError(err)
			}
			env.SecretBinding = ""
			showEnv(*env, true)
			fmt.Println(color.GreenString("Successfully registered as self environment"))
		},
	}
	envSelfSetCmd.Flags().StringP(envDefFlag.name, envDefFlag.shorthand, "", envDefFlag.usage)
	envSelfSetCmd.MarkFlagRequired(envDefFlag.name)
	return envSelfSetCmd
}

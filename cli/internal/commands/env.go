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
	envCmd.AddCommand(envListSearchCommand())
	envCmd.AddCommand(envSelfCommand())
	return envCmd
}

func showEnv(env environments.Environment, includeEDS, excludeBindingFromEds bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "ID (Public Key):\t", env.PublicKey)
	fmt.Fprintln(w, "Name:\t", env.Name)
	fmt.Fprintln(w, "Email:\t", env.Email)
	fmt.Fprintln(w, "Tags:\t", env.Tags)
	if env.SecretBinding != "" {
		fmt.Fprintln(w, "Secret Binding:\t", env.SecretBinding)
	}
	if includeEDS {
		secretBinding := env.SecretBinding
		if excludeBindingFromEds {
			env.SecretBinding = ""
		}
		if envDef, err := env.ToEnvDef(); err == nil {
			fmt.Fprintln(w, "------------------------------------------------------------")
			fmt.Fprintln(w, "Env Definition:\t", color.CyanString(envDef))
		}
		env.SecretBinding = secretBinding
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
	envNewCmd.AddCommand(newKMSEnvCommand("gcp", "Create an environment that works with GCP KMS", gcpKmsResNameFlag))
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
			showEnv(*env, true, false)
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
		Use:     "self",
		Aliases: []string{"user", "usr", "u"},
		Short:   "Register as a new user environment",
		Run: func(cmd *cobra.Command, args []string) {
			selfEnv := environments.GetSelf()
			if selfEnv != nil {
				showEnv(*selfEnv, true, true)
				confirmed, err := input.GetConfirmation("You are already registered as an environment, "+
					"this will replace the existing environment. Proceed? (yes/no): ", "yes")
				if err != nil {
					exitOnError(err)
				}
				if !confirmed {
					fmt.Println(color.YellowString("Operation aborted"))
					safeExit()
				}
			}
			envName, _ := cmd.Flags().GetString(envNameFlag.name)
			envEmail, _ := cmd.Flags().GetString(envEmailFlag.name)
			envTags, err := cmd.Flags().GetStringSlice(envTagsFlag.name)
			if err != nil {
				exitOnError(err)
			}
			inputs := make(map[string][]byte)
			password, err := input.NewPasswordFromUser(input.DefaultPasswordPolicy())
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
			secretBinding := env.SecretBinding
			showEnv(*env, true, true)
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
			fmt.Println(color.GreenString("Successfully registered as self environment"))
			if secretBinding != "" {
				fmt.Println(color.YellowString("Please note down the \"Secret Binding\" somewhere safe so that you don't lose it.\n" +
					"It is required to access your registered environment."))
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

func envListSearchCommand() *cobra.Command {
	if envListSearchCmd != nil {
		return envListSearchCmd
	}
	envListSearchCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "search", "find"},
		Short:   "List/Search environments from profile",
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
				showEnv(*env, false, false)
				fmt.Println()
			}
			safeExit()
		},
	}
	envListSearchCmd.Flags().StringP(profileNameFlag.name, profileNameFlag.shorthand, "", profileNameFlag.usage)
	envListSearchCmd.Flags().StringP(envSearchFlag.name, envSearchFlag.shorthand, "", envSearchFlag.usage)
	return envListSearchCmd
}

func envSelfCommand() *cobra.Command {
	if envSelfCmd != nil {
		return envSelfCmd
	}
	envSelfCmd = &cobra.Command{
		Use:     "self",
		Aliases: []string{"user", "me", "my", "current"},
		Short:   "Shows the current user environment if registered",
		Run: func(cmd *cobra.Command, args []string) {
			env := environments.GetSelf()
			if env == nil {
				fmt.Println("No environment registered as self.")
			} else {
				showEnv(*env, true, true)
			}
			safeExit()
		},
	}
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
			if env.EnvType != environments.USER {
				exitOnError(fmt.Errorf("only user environments can be registered as self"))
			}
			if env.SecretBinding == "" {
				secretBinding, err := input.GetVisibleInput("Enter the secret binding: ")
				if err != nil {
					exitOnError(err)
				}
				env.SecretBinding = secretBinding
			}
			if err = env.MarkAsSelf(); err != nil {
				exitOnError(err)
			}
			showEnv(*env, true, true)
			fmt.Println(color.GreenString("Successfully registered self environment"))
		},
	}
	envSelfSetCmd.Flags().StringP(envDefFlag.name, envDefFlag.shorthand, "", envDefFlag.usage)
	envSelfSetCmd.MarkFlagRequired(envDefFlag.name)
	return envSelfSetCmd
}

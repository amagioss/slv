package cmdenv

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"savesecrets.org/slv/cli/internal/commands/utils"
	"savesecrets.org/slv/core/crypto"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/environments/providers"
	"savesecrets.org/slv/core/input"
	"savesecrets.org/slv/core/profiles"
)

func EnvCommand() *cobra.Command {
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
			name, _ := cmd.Flags().GetString(envNameFlag.Name)
			email, _ := cmd.Flags().GetString(envEmailFlag.Name)
			tags, err := cmd.Flags().GetStringSlice(envTagsFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			var env *environments.Environment
			var secretKey *crypto.SecretKey
			env, secretKey, err = environments.NewEnvironment(name, environments.SERVICE)
			if err != nil {
				utils.ExitOnError(err)
			}
			env.SetEmail(email)
			env.AddTags(tags...)
			utils.ShowEnv(*env, true, false)
			if secretKey != nil {
				fmt.Println("\nSecret Key:\t", color.HiBlackString(secretKey.String()))
			}
			addToProfileFlag, _ := cmd.Flags().GetBool(envAddFlag.Name)
			if addToProfileFlag {
				profile, err := profiles.GetDefaultProfile()
				if err != nil {
					utils.ExitOnError(err)
				}
				err = profile.PutEnv(env)
				if err != nil {
					utils.ExitOnError(err)
				}
			}
			utils.SafeExit()
		},
	}
	envNewServiceCmd.Flags().StringP(envNameFlag.Name, envNameFlag.Shorthand, "", envNameFlag.Usage)
	envNewServiceCmd.Flags().StringP(envEmailFlag.Name, envEmailFlag.Shorthand, "", envEmailFlag.Usage)
	envNewServiceCmd.Flags().StringSliceP(envTagsFlag.Name, envTagsFlag.Shorthand, []string{}, envTagsFlag.Usage)
	envNewServiceCmd.Flags().BoolP(envAddFlag.Name, envAddFlag.Shorthand, false, envAddFlag.Usage)
	envNewServiceCmd.MarkFlagRequired(envNameFlag.Name)
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
				utils.ShowEnv(*selfEnv, true, true)
				confirmed, err := input.GetConfirmation("You are already registered as an environment, "+
					"this will replace the existing environment. Proceed? (yes/no): ", "yes")
				if err != nil {
					utils.ExitOnError(err)
				}
				if !confirmed {
					fmt.Println(color.YellowString("Operation aborted"))
					utils.SafeExit()
				}
			}
			envName, _ := cmd.Flags().GetString(envNameFlag.Name)
			envEmail, _ := cmd.Flags().GetString(envEmailFlag.Name)
			envTags, err := cmd.Flags().GetStringSlice(envTagsFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			inputs := make(map[string][]byte)
			password, err := input.NewPasswordFromUser(input.DefaultPasswordPolicy())
			if err != nil {
				utils.ExitOnError(err)
			}
			inputs["password"] = password
			var env *environments.Environment
			env, err = providers.NewEnvForProvider("password", envName, environments.USER, inputs)
			if err != nil {
				utils.ExitOnError(err)
			}
			env.SetEmail(envEmail)
			env.AddTags(envTags...)
			if err = env.MarkAsSelf(); err != nil {
				utils.ExitOnError(err)
			}
			secretBinding := env.SecretBinding
			utils.ShowEnv(*env, true, true)
			addToProfileFlag, _ := cmd.Flags().GetBool(envAddFlag.Name)
			if addToProfileFlag {
				profile, err := profiles.GetDefaultProfile()
				if err != nil {
					utils.ExitOnError(err)
				}
				err = profile.PutEnv(env)
				if err != nil {
					utils.ExitOnError(err)
				}
			}
			fmt.Println(color.GreenString("Successfully registered as self environment"))
			if secretBinding != "" {
				fmt.Println(color.YellowString("Please note down the \"Secret Binding\" somewhere safe so that you don't lose it.\n" +
					"It is required to access your registered environment."))
			}
			utils.SafeExit()
		},
	}
	envNewUserCmd.Flags().StringP(envNameFlag.Name, envNameFlag.Shorthand, "", envNameFlag.Usage)
	envNewUserCmd.Flags().StringP(envEmailFlag.Name, envEmailFlag.Shorthand, "", envEmailFlag.Usage)
	envNewUserCmd.Flags().StringSliceP(envTagsFlag.Name, envTagsFlag.Shorthand, []string{}, envTagsFlag.Usage)
	envNewUserCmd.Flags().BoolP(envAddFlag.Name, envAddFlag.Shorthand, false, envAddFlag.Usage)
	envNewUserCmd.MarkFlagRequired(envNameFlag.Name)
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
			profile, err := profiles.GetDefaultProfile()
			if err != nil {
				utils.ExitOnError(err)
			}
			queries, err := cmd.Flags().GetStringSlice(EnvSearchFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}
			var envs []*environments.Environment
			for _, query := range queries {
				result, err := profile.SearchEnvs(query)
				if err != nil {
					utils.ExitOnError(err)
				}
				envs = append(envs, result...)
			}
			for _, env := range envs {
				utils.ShowEnv(*env, false, false)
				fmt.Println()
			}
			utils.SafeExit()
		},
	}
	envListSearchCmd.Flags().StringSliceP(EnvSearchFlag.Name, EnvSearchFlag.Shorthand, []string{}, EnvSearchFlag.Usage)
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
				utils.ShowEnv(*env, true, true)
			}
			utils.SafeExit()
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
			envDef := cmd.Flag(envDefFlag.Name).Value.String()
			env, err := environments.FromEnvDef(envDef)
			if err != nil {
				utils.ExitOnError(err)
			}
			if env.EnvType != environments.USER {
				utils.ExitOnError(fmt.Errorf("only user environments can be registered as self"))
			}
			if env.SecretBinding == "" {
				secretBinding, err := input.GetVisibleInput("Enter the secret binding: ")
				if err != nil {
					utils.ExitOnError(err)
				}
				env.SecretBinding = secretBinding
			}
			if err = env.MarkAsSelf(); err != nil {
				utils.ExitOnError(err)
			}
			utils.ShowEnv(*env, true, true)
			fmt.Println(color.GreenString("Successfully registered self environment"))
		},
	}
	envSelfSetCmd.Flags().StringP(envDefFlag.Name, envDefFlag.Shorthand, "", envDefFlag.Usage)
	envSelfSetCmd.MarkFlagRequired(envDefFlag.Name)
	return envSelfSetCmd
}

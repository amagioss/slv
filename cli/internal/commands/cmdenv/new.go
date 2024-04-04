package cmdenv

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"oss.amagi.com/slv/cli/internal/commands/utils"
	"oss.amagi.com/slv/core/crypto"
	"oss.amagi.com/slv/core/environments"
	"oss.amagi.com/slv/core/environments/providers"
	"oss.amagi.com/slv/core/input"
	"oss.amagi.com/slv/core/profiles"
)

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
	envNewCmd.PersistentFlags().BoolP(utils.QuantumSafeFlag.Name, utils.QuantumSafeFlag.Shorthand, false, utils.QuantumSafeFlag.Usage)
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
			pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
			env, secretKey, err = environments.NewEnvironment(name, environments.SERVICE, pq)
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
			pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
			env, err = providers.NewEnvForProvider("password", envName, environments.USER, inputs, pq)
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

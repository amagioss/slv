package cmdenv

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/environments/providers"
	"slv.sh/slv/internal/core/profiles"
)

func newKMSEnvCommand(kmsName, kmsProviderDesc string, keyIdFlag utils.FlagDef) *cobra.Command {
	newKMSEnvCmd := &cobra.Command{
		Use:   kmsName,
		Short: kmsProviderDesc,
		Long: kmsProviderDesc +
			" - Uses RSA 4096 key with SHA-256 hashing in case of asymmetric binding. Create a KMS key accordingly.",
		Run: func(cmd *cobra.Command, args []string) {
			envName, _ := cmd.Flags().GetString(envNameFlag.Name)
			envEmail, _ := cmd.Flags().GetString(envEmailFlag.Name)
			envTags, err := cmd.Flags().GetStringSlice(envTagsFlag.Name)
			if err != nil {
				utils.ExitOnError(err)
			}

			addToProfileFlag, _ := cmd.Flags().GetBool(envAddFlag.Name)
			var profile *profiles.Profile
			if addToProfileFlag {
				if profile, err = profiles.GetActiveProfile(); err != nil {
					utils.ExitOnError(err)
				}
				if !profile.IsPushSupported() {
					utils.ExitOnError(fmt.Errorf("profile (%s) does not support adding environments", profile.Name()))
				}
			}
			inputs := make(map[string][]byte)
			keyIdFlagValue := cmd.Flag(keyIdFlag.Name).Value.String()
			inputs[keyIdFlag.Name] = []byte(keyIdFlagValue)
			pubKeyFilePath := cmd.Flag(kmsRSAPublicKey.Name).Value.String()
			var rsaPublicKey []byte
			var env *environments.Environment
			if pubKeyFilePath != "" {
				if rsaPublicKey, err = os.ReadFile(pubKeyFilePath); err == nil {
					inputs[kmsRSAPublicKey.Name] = rsaPublicKey
				}
			}
			pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
			env, err = providers.NewEnv(kmsName, envName, environments.SERVICE, inputs, pq)
			if err != nil {
				utils.ExitOnError(err)
			}
			env.SetEmail(envEmail)
			env.AddTags(envTags...)
			ShowEnv(*env, true, false)
			if addToProfileFlag {
				if err = profile.PutEnv(env); err != nil {
					utils.ExitOnError(fmt.Errorf("failed to add environment to profile (%s): %w", profile.Name(), err))
				}
				fmt.Printf("Successfully added the environment to profile (%s)\n", color.GreenString(profile.Name()))
			}
			utils.SafeExit()
		},
	}
	newKMSEnvCmd.Flags().StringP(envNameFlag.Name, envNameFlag.Shorthand, "", envNameFlag.Usage)
	newKMSEnvCmd.Flags().StringP(envEmailFlag.Name, envEmailFlag.Shorthand, "", envEmailFlag.Usage)
	newKMSEnvCmd.Flags().StringSliceP(envTagsFlag.Name, envTagsFlag.Shorthand, []string{}, envTagsFlag.Usage)
	newKMSEnvCmd.Flags().StringP(keyIdFlag.Name, keyIdFlag.Shorthand, "", keyIdFlag.Usage)
	newKMSEnvCmd.Flags().StringP(kmsRSAPublicKey.Name, kmsRSAPublicKey.Shorthand, "", kmsRSAPublicKey.Usage)
	newKMSEnvCmd.Flags().BoolP(envAddFlag.Name, envAddFlag.Shorthand, false, envAddFlag.Usage)
	newKMSEnvCmd.MarkFlagRequired(envNameFlag.Name)
	newKMSEnvCmd.MarkFlagRequired(keyIdFlag.Name)
	return newKMSEnvCmd
}

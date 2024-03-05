package cmdenv

import (
	"os"

	"github.com/spf13/cobra"
	"savesecrets.org/slv/cli/internal/commands/utils"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/environments/providers"
	"savesecrets.org/slv/core/profiles"
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
			env, err = providers.NewEnvForProvider(kmsName, envName, environments.SERVICE, inputs)
			if err != nil {
				utils.ExitOnError(err)
			}
			env.SetEmail(envEmail)
			env.AddTags(envTags...)
			utils.ShowEnv(*env, true, false)
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

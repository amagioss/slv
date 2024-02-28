package commands

import (
	"os"

	"github.com/spf13/cobra"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/profiles"
)

func newKMSEnvCommand(kmsName, kmsProviderDesc string, keyIdFlag FlagDef) *cobra.Command {
	newKMSEnvCmd := &cobra.Command{
		Use:   kmsName,
		Short: kmsProviderDesc,
		Run: func(cmd *cobra.Command, args []string) {
			envName, _ := cmd.Flags().GetString(envNameFlag.name)
			envEmail, _ := cmd.Flags().GetString(envEmailFlag.name)
			envTags, err := cmd.Flags().GetStringSlice(envTagsFlag.name)
			if err != nil {
				exitOnError(err)
			}
			inputs := make(map[string][]byte)
			keyIdFlagValue := cmd.Flag(keyIdFlag.name).Value.String()
			inputs[keyIdFlag.name] = []byte(keyIdFlagValue)
			pubKeyFilePath := cmd.Flag(kmsRSAPublicKey.name).Value.String()
			var rsaPublicKey []byte
			var env *environments.Environment
			if pubKeyFilePath != "" {
				if rsaPublicKey, err = os.ReadFile(pubKeyFilePath); err == nil {
					inputs[kmsRSAPublicKey.name] = rsaPublicKey
				}
			}
			env, err = environments.NewEnvForProvider(kmsName, envName, environments.SERVICE, inputs)
			if err != nil {
				exitOnError(err)
			}
			env.SetEmail(envEmail)
			env.AddTags(envTags...)
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

	newKMSEnvCmd.Flags().StringP(envNameFlag.name, envNameFlag.shorthand, "", envNameFlag.usage)
	newKMSEnvCmd.Flags().StringP(envEmailFlag.name, envEmailFlag.shorthand, "", envEmailFlag.usage)
	newKMSEnvCmd.Flags().StringSliceP(envTagsFlag.name, envTagsFlag.shorthand, []string{}, envTagsFlag.usage)
	newKMSEnvCmd.Flags().StringP(keyIdFlag.name, keyIdFlag.shorthand, "", keyIdFlag.usage)
	newKMSEnvCmd.Flags().StringP(kmsRSAPublicKey.name, kmsRSAPublicKey.shorthand, "", kmsRSAPublicKey.usage)
	newKMSEnvCmd.Flags().BoolP(envAddFlag.name, envAddFlag.shorthand, false, envAddFlag.usage)
	newKMSEnvCmd.MarkFlagRequired(envNameFlag.name)
	newKMSEnvCmd.MarkFlagRequired(keyIdFlag.name)
	return newKMSEnvCmd
}

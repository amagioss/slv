package commands

import (
	"os"

	"github.com/shibme/slv/core/environments"
	"github.com/shibme/slv/core/profiles"
	"github.com/spf13/cobra"
)

func newKMSEnvCommand(kmsName, kmsProviderDesc string, keyIdFlag FlagDef, pubkeyFileFlag FlagDef) *cobra.Command {
	newKMSEnvCmd := &cobra.Command{
		Use:   kmsName,
		Short: kmsProviderDesc,
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString(envNameFlag.name)
			email, _ := cmd.Flags().GetString(envEmailFlag.name)
			tags, err := cmd.Flags().GetStringSlice(envTagsFlag.name)
			if err != nil {
				exitOnError(err)
			}
			inputs := make(map[string][]byte)
			keyIdFlagValue := cmd.Flag(keyIdFlag.name).Value.String()
			providerName := "kms-" + kmsName
			inputs[keyIdFlag.name] = []byte(keyIdFlagValue)
			pubKeyFilePath := cmd.Flag(pubkeyFileFlag.name).Value.String()
			var rsaPublicKey []byte
			var env *environments.Environment
			if pubKeyFilePath != "" {
				if rsaPublicKey, err = os.ReadFile(pubKeyFilePath); err == nil {
					inputs[pubkeyFileFlag.name] = rsaPublicKey
				}
			}
			env, err = environments.NewEnvForProvider(providerName, name, environments.SERVICE, inputs)
			if err != nil {
				exitOnError(err)
			}
			env.SetEmail(email)
			env.AddTags(tags...)
			showEnv(*env, true)
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

	newKMSEnvCmd.Flags().StringP(envNameFlag.name, envNameFlag.shorthand, "", envNameFlag.usage)
	newKMSEnvCmd.Flags().StringP(envEmailFlag.name, envEmailFlag.shorthand, "", envEmailFlag.usage)
	newKMSEnvCmd.Flags().StringSliceP(envTagsFlag.name, envTagsFlag.shorthand, []string{}, envTagsFlag.usage)
	newKMSEnvCmd.Flags().StringP(keyIdFlag.name, keyIdFlag.shorthand, "", keyIdFlag.usage)
	newKMSEnvCmd.Flags().StringP(pubkeyFileFlag.name, pubkeyFileFlag.shorthand, "", pubkeyFileFlag.usage)
	newKMSEnvCmd.Flags().BoolP(envAddFlag.name, envAddFlag.shorthand, false, envAddFlag.usage)
	newKMSEnvCmd.MarkFlagRequired(envNameFlag.name)
	newKMSEnvCmd.MarkFlagRequired(keyIdFlag.name)
	// newKMSEnvCmd.MarkFlagRequired(pubkeyFileFlag.name)
	return newKMSEnvCmd
}

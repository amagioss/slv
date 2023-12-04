package commands

import (
	"os"

	"github.com/shibme/slv/core/environments"
	"github.com/shibme/slv/core/profiles"
	"github.com/shibme/slv/core/secretkeystore/providers"
	"github.com/spf13/cobra"
)

func newKMSEnvCommand(kmsProviderName, kmsProviderDesc string, keyIdFlag FlagDef, pubkeyFileFlag FlagDef) *cobra.Command {
	newKMSEnvCmd := &cobra.Command{
		Use:   kmsProviderName,
		Short: kmsProviderDesc,
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString(envNameFlag.name)
			email, _ := cmd.Flags().GetString(envEmailFlag.name)
			tags, err := cmd.Flags().GetStringSlice(envTagsFlag.name)
			if err != nil {
				exitOnError(err)
			}
			inputs := make(map[string]any)
			keyIdFlagValue := cmd.Flag(keyIdFlag.name).Value.String()
			inputs[keyIdFlag.name] = keyIdFlagValue
			pubKeyFilePath := cmd.Flag(pubkeyFileFlag.name).Value.String()
			var rsaPublicKey []byte
			var env *environments.Environment
			if rsaPublicKey, err = os.ReadFile(pubKeyFilePath); err == nil {
				inputs[pubkeyFileFlag.name] = rsaPublicKey
				env, err = providers.NewEnvForProvider(name, email, environments.SERVICE, kmsProviderName, inputs)
			}
			if err != nil {
				exitOnError(err)
			}
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
	newKMSEnvCmd.MarkFlagRequired(pubkeyFileFlag.name)
	return newKMSEnvCmd
}

package cmdenv

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/environments/envproviders"
	"slv.sh/slv/internal/core/profiles"
)

func getEnvProviderCommand(providerId string) *cobra.Command {
	providerArgs := envproviders.GetArgs(providerId)
	envProviderCmd := &cobra.Command{
		Use:   providerId,
		Short: "Creates a new service environment using " + envproviders.GetName(providerId),
		Long:  "Creates a new service environment using " + envproviders.GetDesc(providerId),
		Run: func(cmd *cobra.Command, args []string) {
			envName, _ := cmd.Flags().GetString(envNameFlag.Name)
			envEmail, _ := cmd.Flags().GetString(envEmailFlag.Name)
			envTags, err := cmd.Flags().GetStringSlice(envTagsFlag.Name)
			pq, _ := cmd.Flags().GetBool(utils.QuantumSafeFlag.Name)
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
			inputs := make(map[string]string)
			for _, arg := range providerArgs {
				if value, _ := cmd.Flags().GetString(arg.Id()); value != "" {
					inputs[arg.Id()] = value
				}
			}
			var env *environments.Environment
			if env, err = envproviders.NewEnv(providerId, envName, environments.SERVICE, inputs, pq); err != nil {
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
		},
	}
	for _, arg := range providerArgs {
		envProviderCmd.Flags().StringP(arg.Id(), "", "", arg.Description())
		if arg.Required() {
			envProviderCmd.MarkFlagRequired(arg.Id())
		}
	}
	return envProviderCmd
}

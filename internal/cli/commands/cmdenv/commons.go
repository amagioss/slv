package cmdenv

import (
	"fmt"

	"github.com/spf13/cobra"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/session"
)

func GetPublicKeys(cmd *cobra.Command, root, pq bool) (publicKeys []*crypto.PublicKey, err error) {
	publicKeyStrings, err := cmd.Flags().GetStringSlice(EnvPublicKeysFlag.Name)
	if err != nil {
		utils.ExitOnError(err)
	}
	queries, err := cmd.Flags().GetStringSlice(EnvSearchFlag.Name)
	if err != nil {
		utils.ExitOnError(err)
	}
	shareWithSelf, _ := cmd.Flags().GetBool(EnvSelfFlag.Name)
	shareWithK8s, _ := cmd.Flags().GetBool(EnvK8sFlag.Name)
	if len(publicKeyStrings) == 0 && len(queries) == 0 && !shareWithSelf && !shareWithK8s {
		return nil, fmt.Errorf("specify at least one of the following flags:\n --%s [search keyword]\n --%s [env public key]\n --%s\n --%s",
			EnvSearchFlag.Name, EnvPublicKeysFlag.Name, EnvSelfFlag.Name, EnvK8sFlag.Name)
	}
	for _, pubKeyStr := range publicKeyStrings {
		publicKey, err := crypto.PublicKeyFromString(pubKeyStr)
		if err != nil {
			return nil, err
		}
		publicKeys = append(publicKeys, publicKey)
	}
	profile, err := profiles.GetActiveProfile()
	if err != nil && len(queries) > 0 {
		return nil, err
	}
	if len(queries) > 0 {
		envs, err := profile.SearchEnvs(queries)
		if err != nil {
			return nil, err
		}
		for _, env := range envs {
			publicKey, err := crypto.PublicKeyFromString(env.PublicKey)
			if err != nil {
				return nil, err
			}
			publicKeys = append(publicKeys, publicKey)
		}
	}
	if shareWithSelf {
		selfEnv := environments.GetSelf()
		if selfEnv != nil {
			publicKey, err := crypto.PublicKeyFromString(selfEnv.PublicKey)
			if err != nil {
				return nil, err
			}
			publicKeys = append(publicKeys, publicKey)
		}
	}
	if shareWithK8s {
		pk, err := session.GetPublicKeyFromK8s(config.AppNameLowerCase, pq)
		if err != nil {
			return nil, err
		}
		publicKey, err := crypto.PublicKeyFromString(pk)
		if err != nil {
			return nil, err
		}
		publicKeys = append(publicKeys, publicKey)
	}
	if len(publicKeys) == 0 {
		return nil, fmt.Errorf("no matching environments found")
	}
	if profile != nil && root {
		rootEnv, err := profile.GetRoot()
		if err != nil {
			return nil, err
		}
		if rootPublicKey, err := rootEnv.GetPublicKey(); err != nil {
			return nil, err
		} else if rootPublicKey != nil {
			publicKeys = append(publicKeys, rootPublicKey)
		}
	}
	return publicKeys, nil
}

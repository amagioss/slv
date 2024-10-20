package cmdenv

import (
	"fmt"

	"github.com/spf13/cobra"
	"oss.amagi.com/slv/internal/cli/commands/utils"
	"oss.amagi.com/slv/internal/core/config"
	"oss.amagi.com/slv/internal/core/crypto"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/profiles"
	k8sutils "oss.amagi.com/slv/internal/k8s/utils"
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
		return nil, fmt.Errorf("Specify atleast one of the following flags:\n" +
			" --" + EnvSearchFlag.Name + " [search keyword]\n" +
			" --" + EnvPublicKeysFlag.Name + " [env public key]\n" +
			" --" + EnvSelfFlag.Name + "\n" +
			" --" + EnvK8sFlag.Name)
	}
	for _, pubKeyStr := range publicKeyStrings {
		publicKey, err := crypto.PublicKeyFromString(pubKeyStr)
		if err != nil {
			return nil, err
		}
		publicKeys = append(publicKeys, publicKey)
	}
	profile, err := profiles.GetDefaultProfile()
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
		pk, err := k8sutils.GetPublicKeyFromK8s(config.AppNameLowerCase, pq)
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
	if profile != nil {
		if rootPublicKey, err := profile.RootPublicKey(); err != nil {
			return nil, err
		} else if rootPublicKey != nil {
			publicKeys = append(publicKeys, rootPublicKey)
		}
	}
	return publicKeys, nil
}

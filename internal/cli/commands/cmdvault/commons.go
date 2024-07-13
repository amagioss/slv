package cmdvault

import (
	"fmt"

	"oss.amagi.com/slv/internal/cli/commands/cmdenv"
	"oss.amagi.com/slv/internal/core/config"
	"oss.amagi.com/slv/internal/core/crypto"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/profiles"
	k8sutils "oss.amagi.com/slv/internal/k8s/utils"
)

func getPublicKeys(pubKeyStrSlice, queries []string, self, k8sCluster, k8sPQ bool) (publicKeys []*crypto.PublicKey,
	rootPublicKey *crypto.PublicKey, err error) {
	if len(pubKeyStrSlice) == 0 && len(queries) == 0 && !self && !k8sCluster {
		return nil, nil, fmt.Errorf("Specify atleast one of the following flags:\n" +
			" --" + cmdenv.EnvSearchFlag.Name + " [search keyword]\n" +
			" --" + vaultAccessPublicKeysFlag.Name + " [env public key]\n" +
			" --" + cmdenv.EnvSelfFlag.Name + "\n" +
			" --" + vaultAccessK8sFlag.Name)
	}
	for _, pubKeyStr := range pubKeyStrSlice {
		publicKey, err := crypto.PublicKeyFromString(pubKeyStr)
		if err != nil {
			return nil, nil, err
		}
		publicKeys = append(publicKeys, publicKey)
	}
	profile, err := profiles.GetDefaultProfile()
	if len(queries) > 0 {
		if err != nil {
			return nil, nil, err
		}
		envs, err := profile.SearchEnvs(queries)
		if err != nil {
			return nil, nil, err
		}
		for _, env := range envs {
			publicKey, err := crypto.PublicKeyFromString(env.PublicKey)
			if err != nil {
				return nil, nil, err
			}
			publicKeys = append(publicKeys, publicKey)
		}
		if len(publicKeys) == 0 {
			return nil, nil, fmt.Errorf("no matching environments found for the given search queries")
		}
	}
	if self {
		selfEnv := environments.GetSelf()
		if selfEnv != nil {
			publicKey, err := crypto.PublicKeyFromString(selfEnv.PublicKey)
			if err != nil {
				return nil, nil, err
			}
			publicKeys = append(publicKeys, publicKey)
		}
	}
	if k8sCluster {
		pk, err := k8sutils.GetPublicKeyFromK8s(config.AppNameLowerCase, k8sPQ)
		if err != nil {
			return nil, nil, err
		}
		publicKey, err := crypto.PublicKeyFromString(pk)
		if err != nil {
			return nil, nil, err
		}
		publicKeys = append(publicKeys, publicKey)
	}
	if profile != nil {
		rootPublicKey, err = profile.RootPublicKey()
		if err != nil {
			return nil, nil, err
		}
	}
	return publicKeys, rootPublicKey, nil
}

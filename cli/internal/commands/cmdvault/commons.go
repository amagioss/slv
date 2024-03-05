package cmdvault

import (
	"fmt"

	"savesecrets.org/slv/cli/internal/commands/cmdenv"
	"savesecrets.org/slv/core/crypto"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/profiles"
	"savesecrets.org/slv/core/vaults"
)

func getVault(filePath string) (*vaults.Vault, error) {
	vault, err := vaults.Get(filePath)
	if err != nil || vault.Config.PublicKey == "" {
		vault, err = vaults.GetFromField(filePath, k8sVaultField)
	}
	return vault, err
}

func getPublicKeys(pubKeyStrSlice []string, query string, self bool) (publicKeys []*crypto.PublicKey,
	rootPublicKey *crypto.PublicKey, err error) {
	if len(pubKeyStrSlice) == 0 && query == "" && !self {
		return nil, nil, fmt.Errorf("Specify atleast one of the following flags:\n" +
			" --" + cmdenv.EnvSearchFlag.Name + "\n" +
			" --" + vaultAccessPublicKeysFlag.Name + "\n" +
			" --" + cmdenv.EnvSelfFlag.Name)
	}
	for _, pubKeyStr := range pubKeyStrSlice {
		publicKey, err := crypto.PublicKeyFromString(pubKeyStr)
		if err != nil {
			return nil, nil, err
		}
		publicKeys = append(publicKeys, publicKey)
	}
	profile, err := profiles.GetDefaultProfile()
	if query != "" {
		if err != nil {
			return nil, nil, err
		}
		envs, err := profile.SearchEnvs(query)
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
			return nil, nil, fmt.Errorf("no matching environments found for search query: " + query)
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
	if profile != nil {
		rootPublicKey, err = profile.RootPublicKey()
		if err != nil {
			return nil, nil, err
		}
	}
	return publicKeys, rootPublicKey, nil
}
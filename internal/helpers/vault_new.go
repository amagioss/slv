package helpers

import (
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/vaults"
)

func NewVault(vaultFile, name, k8sNamespace string, enableHash, pq bool, pkStrList []string) (*vaults.Vault, error) {
	var pubKeys []*crypto.PublicKey
	for _, pkStr := range pkStrList {
		pubKey, err := crypto.PublicKeyFromString(pkStr)
		if err != nil {
			return nil, err
		}
		pubKeys = append(pubKeys, pubKey)
	}
	return vaults.New(vaultFile, name, k8sNamespace, enableHash, pq, pubKeys...)
}

func UpdateVault(vaultFile, name, k8sNamespace, secretType string, k8SecretContent []byte) (*vaults.Vault, error) {
	vlt, err := vaults.Get(vaultFile)
	if err != nil {
		return nil, err
	}
	if err = vlt.Update(name, k8sNamespace, secretType, k8SecretContent); err != nil {
		return nil, err
	}
	return vlt, nil
}

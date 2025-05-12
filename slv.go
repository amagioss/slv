package slv

import (
	"slv.sh/slv/internal/core/secretkey"
	"slv.sh/slv/internal/core/vaults"
)

func unlockVault(vaultFile string) (*vaults.Vault, error) {
	secretKey, err := secretkey.Get()
	if err != nil {
		return nil, err
	}
	vault, err := vaults.Get(vaultFile)
	if err != nil {
		return nil, err
	}
	if err = vault.Unlock(secretKey); err != nil {
		return nil, err
	}
	return vault, nil
}

// GetAllVaultItems returns all secrets from the vault
func GetAllVaultItems(vaultFile string) (map[string]*vaults.VaultItem, error) {
	vault, err := unlockVault(vaultFile)
	if err != nil {
		return nil, err
	}
	return vault.GetAllItems()
}

// GetVaultItem returns a named secret from the vault
func GetVaultItem(vaultFile, name string) (vaultItem *vaults.VaultItem, err error) {
	if vault, err := unlockVault(vaultFile); err != nil {
		return nil, err
	} else {
		return vault.Get(name)
	}
}

// PutVaultItem writes a secret to the vault
func PutVaultItem(vaultFile, secretName string, secretValue []byte, encrypt bool) error {
	vault, err := vaults.Get(vaultFile)
	if err != nil {
		return err
	}
	return vault.Put(secretName, secretValue, encrypt)
}

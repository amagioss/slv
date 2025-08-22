package slv

import (
	"slv.sh/slv/internal/core/session"
	"slv.sh/slv/internal/core/vaults"
)

func unlockVault(vaultFileOrURL string) (*vaults.Vault, error) {
	secretKey, err := session.GetSecretKey()
	if err != nil {
		return nil, err
	}
	vault, err := vaults.Get(vaultFileOrURL)
	if err != nil {
		return nil, err
	}
	if err = vault.Unlock(secretKey); err != nil {
		return nil, err
	}
	return vault, nil
}

// GetAllVaultItems returns all secrets from the vault
func GetAllVaultItems(vaultFileOrURL string) (map[string]*vaults.VaultItem, error) {
	vault, err := unlockVault(vaultFileOrURL)
	if err != nil {
		return nil, err
	}
	return vault.GetAllItems()
}

// GetVaultItem returns a named secret from the vault
func GetVaultItem(vaultFileOrURL, name string) (vaultItem *vaults.VaultItem, err error) {
	if vault, err := unlockVault(vaultFileOrURL); err != nil {
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

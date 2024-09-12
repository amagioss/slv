package slv

import (
	"oss.amagi.com/slv/internal/core/secretkey"
	"oss.amagi.com/slv/internal/core/vaults"
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

// GetAllVaultData returns all secrets from the vault
func GetAllVaultData(vaultFile string) (map[string]*vaults.VaultData, error) {
	vault, err := unlockVault(vaultFile)
	if err != nil {
		return nil, err
	}
	return vault.GetAll()
}

// GetVaultData returns a named secret from the vault
func GetVaultData(vaultFile, name string) (vaultData *vaults.VaultData, err error) {
	if vault, err := unlockVault(vaultFile); err != nil {
		return nil, err
	} else {
		return vault.Get(name)
	}
}

// PutVaultData writes a secret to the vault
func PutVaultData(vaultFile, secretName string, secretValue []byte, encrypt bool) error {
	vault, err := vaults.Get(vaultFile)
	if err != nil {
		return err
	}
	return vault.Put(secretName, secretValue, encrypt)
}

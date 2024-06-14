package slv

import (
	"oss.amagi.com/slv/internal/core/secretkey"
	"oss.amagi.com/slv/internal/core/vaults"
)

func getVaultUnlocked(vaultFile string) (*vaults.Vault, error) {
	secretKey, err := secretkey.Get()
	if err != nil {
		return nil, err
	}
	vault, err := vaults.Get(vaultFile)
	if err != nil {
		return nil, err
	}
	if err = vault.Unlock(*secretKey); err != nil {
		return nil, err
	}
	return vault, nil
}

// GetAllSecrets returns all secrets from the vault
func GetAllSecrets(vaultFile string) (map[string][]byte, error) {
	vault, err := getVaultUnlocked(vaultFile)
	if err != nil {
		return nil, err
	}
	return vault.GetAllSecrets()
}

// GetSecret returns a named secret from the vault
func GetSecret(vaultFile, secretName string) ([]byte, error) {
	vault, err := getVaultUnlocked(vaultFile)
	if err != nil {
		return nil, err
	}
	return vault.GetSecret(secretName)
}

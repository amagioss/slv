package vaults

import (
	"oss.amagi.com/slv/internal/core/crypto"
)

func (vlt *Vault) Share(publicKey *crypto.PublicKey) (bool, error) {
	return vlt.share(publicKey, true)
}

func (vlt *Vault) share(publicKey *crypto.PublicKey, commit bool) (bool, error) {
	if vlt.IsLocked() {
		return false, errVaultLocked
	}
	if publicKey.Type() == VaultKey {
		return false, errVaultCannotBeSharedWithVault
	}
	for _, wrappedKeyStr := range vlt.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		if err := wrappedKey.FromString(wrappedKeyStr); err != nil {
			return false, err
		}
		if wrappedKey.IsEncryptedBy(publicKey) {
			return false, nil
		}
	}
	wrappedKey, err := publicKey.EncryptKey(*vlt.secretKey)
	if err == nil {
		vlt.Config.WrappedKeys = append(vlt.Config.WrappedKeys, wrappedKey.String())
		if commit {
			err = vlt.commit()
		}
	}
	return err == nil, err
}

func (vlt *Vault) Revoke(publicKeys []*crypto.PublicKey, quantumSafe bool) error {
	vaultDataMap, err := vlt.List(true)
	if err != nil {
		return err
	}
	accessors, err := vlt.ListAccessors()
	if err != nil {
		return err
	}
	var newAccessors []crypto.PublicKey
	for _, accessor := range accessors {
		found := false
		for _, publicKey := range publicKeys {
			publicKeyStr, err := publicKey.String()
			if err != nil {
				return err
			}
			accessorStr, err := accessor.String()
			if err != nil {
				return err
			}
			if publicKeyStr == accessorStr {
				found = true
				break
			}
		}
		if !found {
			newAccessors = append(newAccessors, accessor)
		}
	}
	if len(newAccessors) == len(accessors) {
		return nil
	}
	vaultSecretKey, err := crypto.NewSecretKey(VaultKey)
	if err != nil {
		return err
	}
	vaultPublicKey, err := vaultSecretKey.PublicKey(quantumSafe)
	if err != nil {
		return err
	}
	vlt.publicKey = vaultPublicKey
	vaultPublicKeyStr, err := vaultPublicKey.String()
	if err != nil {
		return err
	}
	vlt.Config.PublicKey = vaultPublicKeyStr
	vlt.secretKey = vaultSecretKey
	vlt.Config.WrappedKeys = []string{}
	for _, accessor := range newAccessors {
		wrappedKey, err := accessor.EncryptKey(*vlt.secretKey)
		if err == nil {
			vlt.Config.WrappedKeys = append(vlt.Config.WrappedKeys, wrappedKey.String())
		} else {
			return err
		}
	}
	for name, vaultData := range vaultDataMap {
		if err = vlt.putWithoutCommit(name, vaultData.value, vaultData.isSecret); err != nil {
			return err
		}
	}
	return vlt.commit()
}

func (vlt *Vault) ListAccessors() ([]crypto.PublicKey, error) {
	var accessors []crypto.PublicKey
	for _, wrappedKeyStr := range vlt.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		err := wrappedKey.FromString(wrappedKeyStr)
		if err != nil {
			return nil, err
		}
		encryptedBy, err := wrappedKey.EncryptedByPublicKey()
		if err != nil {
			return nil, err
		}
		accessors = append(accessors, *encryptedBy)
	}
	return accessors, nil
}

func (vlt *Vault) Unlock(secretKey *crypto.SecretKey) error {
	if !vlt.IsLocked() {
		return nil
	}
	for _, wrappedKeyStr := range vlt.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		if err := wrappedKey.FromString(wrappedKeyStr); err != nil {
			return err
		}
		decryptedKey, err := secretKey.DecryptKey(*wrappedKey)
		if err == nil {
			vlt.secretKey = decryptedKey
			return nil
		}
	}
	return errVaultNotAccessible
}

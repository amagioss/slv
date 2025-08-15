package vaults

import (
	"slv.sh/slv/internal/core/crypto"
)

func (vlt *Vault) Share(publicKey *crypto.PublicKey) (bool, error) {
	if !vlt.Spec.writable {
		return false, errVaultNotWritable
	}
	return vlt.share(publicKey, true)
}

func (vlt *Vault) share(publicKey *crypto.PublicKey, commit bool) (bool, error) {
	if vlt.IsLocked() {
		return false, errVaultLocked
	}
	if publicKey.Type() == VaultKey {
		return false, errVaultCannotBeSharedWithVault
	}
	for _, wrappedKeyStr := range vlt.Spec.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		if err := wrappedKey.FromString(wrappedKeyStr); err != nil {
			return false, err
		}
		if wrappedKey.IsEncryptedBy(publicKey) {
			return false, nil
		}
	}
	wrappedKey, err := publicKey.EncryptKey(*vlt.Spec.secretKey)
	if err == nil {
		vlt.Spec.Config.WrappedKeys = append(vlt.Spec.Config.WrappedKeys, wrappedKey.String())
		if commit {
			err = vlt.commit()
		}
	}
	return err == nil, err
}

func (vlt *Vault) Revoke(publicKeys []*crypto.PublicKey, quantumSafe bool) (err error) {
	if !vlt.Spec.writable {
		return errVaultNotWritable
	}
	var vaultItemsMap map[string]*VaultItem
	if vaultItemsMap, err = vlt.GetAllItems(); err != nil {
		return err
	}
	for _, vaultItem := range vaultItemsMap {
		if _, err = vaultItem.Value(); err != nil {
			return err
		}
	}
	var accessors []crypto.PublicKey
	if accessors, err = vlt.ListAccessors(); err != nil {
		return err
	}
	var newAccessors []crypto.PublicKey
	for _, accessor := range accessors {
		found := false
		for _, publicKey := range publicKeys {
			var publicKeyStr string
			if publicKeyStr, err = publicKey.String(); err != nil {
				return err
			}
			var accessorStr string
			if accessorStr, err = accessor.String(); err != nil {
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
	vlt.Spec.publicKey = vaultPublicKey
	vaultPublicKeyStr, err := vaultPublicKey.String()
	if err != nil {
		return err
	}
	vlt.Spec.Config.PublicKey = vaultPublicKeyStr
	vlt.Spec.secretKey = vaultSecretKey
	vlt.Spec.Config.WrappedKeys = []string{}
	for _, accessor := range newAccessors {
		wrappedKey, err := accessor.EncryptKey(*vlt.Spec.secretKey)
		if err == nil {
			vlt.Spec.Config.WrappedKeys = append(vlt.Spec.Config.WrappedKeys, wrappedKey.String())
		} else {
			return err
		}
	}
	for name, vaultItem := range vaultItemsMap {
		if err = vlt.putWithoutCommit(name, vaultItem.value, !vaultItem.IsPlaintext()); err != nil {
			return err
		}
	}
	err = vlt.commit()
	vlt.clearCache()
	return
}

func (vlt *Vault) ListAccessors() ([]crypto.PublicKey, error) {
	var accessors []crypto.PublicKey
	for _, wrappedKeyStr := range vlt.Spec.Config.WrappedKeys {
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

func (vlt *Vault) IsAccessibleBy(secretKey *crypto.SecretKey) bool {
	pubKeyEC, err := secretKey.PublicKey(false)
	if err != nil {
		return false
	}
	pubKeyPQ, err := secretKey.PublicKey(true)
	if err != nil {
		return false
	}
	for _, wrappedKeyStr := range vlt.Spec.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		if err := wrappedKey.FromString(wrappedKeyStr); err == nil {
			if wrappedKey.IsEncryptedBy(pubKeyEC) || wrappedKey.IsEncryptedBy(pubKeyPQ) {
				return true
			}
		}
	}
	return false
}

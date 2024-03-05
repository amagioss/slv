package vaults

import "savesecrets.org/slv/core/crypto"

func (vlt *Vault) Share(publicKey *crypto.PublicKey) (bool, error) {
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
		err = vlt.commit()
	}
	return err == nil, err
}

func (vlt *Vault) Revoke(publicKey *crypto.PublicKey, rotateVaultKeyPair bool) error {
	if rotateVaultKeyPair {
		if vlt.IsLocked() {
			return errVaultLocked
		}
		accessors, err := vlt.ListAccessors()
		if err != nil {
			return err
		}
		var newAccessors []crypto.PublicKey
		for _, accessor := range accessors {
			if accessor.String() != publicKey.String() {
				newAccessors = append(newAccessors, accessor)
			}
		}
		if len(newAccessors) == len(accessors) {
			return nil
		}
		secretsMap, err := vlt.GetAllSecrets()
		if err != nil {
			return err
		}
		vaultSecretKey, err := crypto.NewSecretKey(VaultKey)
		if err != nil {
			return err
		}
		vaultPublicKey, err := vaultSecretKey.PublicKey()
		if err != nil {
			return err
		}
		vlt.publicKey = vaultPublicKey
		vlt.Config.PublicKey = vaultPublicKey.String()
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
		for secretName, secretValue := range secretsMap {
			if err = vlt.putSecretWithoutCommit(secretName, secretValue); err != nil {
				return err
			}
		}
		return vlt.commit()
	} else {
		for i, wrappedKeyStr := range vlt.Config.WrappedKeys {
			wrappedKey := &crypto.WrappedKey{}
			if err := wrappedKey.FromString(wrappedKeyStr); err != nil {
				return err
			}
			if !wrappedKey.IsEncryptedBy(publicKey) {
				vlt.Config.WrappedKeys = append(vlt.Config.WrappedKeys[:i], vlt.Config.WrappedKeys[i+1:]...)
				return vlt.commit()
			}
		}
	}
	return nil
}

func (vlt *Vault) ListAccessors() ([]crypto.PublicKey, error) {
	var accessors []crypto.PublicKey
	for _, wrappedKeyStr := range vlt.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		err := wrappedKey.FromString(wrappedKeyStr)
		if err != nil {
			return nil, err
		}
		accessors = append(accessors, wrappedKey.EncryptedBy())
	}
	return accessors, nil
}

func (vlt *Vault) Unlock(secretKey crypto.SecretKey) error {
	publicKey, err := secretKey.PublicKey()
	if err != nil || (!vlt.IsLocked() && *vlt.unlockedBy == publicKey.String()) {
		return err
	}
	for _, wrappedKeyStr := range vlt.Config.WrappedKeys {
		wrappedKey := &crypto.WrappedKey{}
		if err = wrappedKey.FromString(wrappedKeyStr); err != nil {
			return err
		}
		decryptedKey, err := secretKey.DecryptKey(*wrappedKey)
		if err == nil {
			vlt.secretKey = decryptedKey
			vlt.unlockedBy = new(string)
			*vlt.unlockedBy = publicKey.String()
			return nil
		}
	}
	return errVaultNotAccessible
}

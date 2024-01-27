package vaults

func (vlt *Vault) getSecretFromCache(secretName string) (decryptedSecret []byte) {
	if vlt.decryptedSecrets != nil {
		decryptedSecret = vlt.decryptedSecrets[secretName]
	}
	return
}

func (vlt *Vault) putSecretToCache(secretName string, secretValue []byte) {
	if vlt.decryptedSecrets == nil {
		vlt.decryptedSecrets = make(map[string][]byte)
	}
	vlt.decryptedSecrets[secretName] = secretValue
}

func (vlt *Vault) deleteSecretFromCache(secretName string) {
	if vlt.decryptedSecrets != nil {
		delete(vlt.decryptedSecrets, secretName)
	}
}

func (vlt *Vault) clearSecretCache() {
	vlt.decryptedSecrets = nil
}

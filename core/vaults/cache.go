package vaults

func (vlt *Vault) getSecretFromCache(secretName string) (decryptedSecret []byte) {
	if vlt.decryptedSecretCache != nil {
		decryptedSecret = vlt.decryptedSecretCache[secretName]
	}
	return
}

func (vlt *Vault) putSecretToCache(secretName string, secretValue []byte) {
	if vlt.decryptedSecretCache == nil {
		vlt.decryptedSecretCache = make(map[string][]byte)
	}
	vlt.decryptedSecretCache[secretName] = secretValue
}

func (vlt *Vault) deleteSecretFromCache(secretName string) {
	if vlt.decryptedSecretCache != nil {
		delete(vlt.decryptedSecretCache, secretName)
	}
}

func (vlt *Vault) clearSecretCache() {
	vlt.decryptedSecretCache = nil
}

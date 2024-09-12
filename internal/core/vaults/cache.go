package vaults

func (vlt *Vault) getFromCache(name string) (data *VaultData) {
	if vlt.cache != nil {
		data = vlt.cache[name]
	}
	return
}

func (vlt *Vault) putToCache(name string, data *VaultData) {
	if vlt.cache == nil {
		vlt.cache = make(map[string]*VaultData)
	}
	vlt.cache[name] = data
}

func (vlt *Vault) deleteFromCache(name string) {
	if vlt.cache != nil {
		delete(vlt.cache, name)
	}
}

func (vlt *Vault) clearSecretCache() {
	vlt.cache = nil
}

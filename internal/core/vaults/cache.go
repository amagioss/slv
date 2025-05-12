package vaults

func (vlt *Vault) getFromCache(name string) (item *VaultItem) {
	if vlt.Spec.cache != nil {
		item = vlt.Spec.cache[name]
	}
	return
}

func (vlt *Vault) putToCache(name string, item *VaultItem) {
	if vlt.Spec.cache == nil {
		vlt.Spec.cache = make(map[string]*VaultItem)
	}
	vlt.Spec.cache[name] = item
}

func (vlt *Vault) deleteFromCache(name string) {
	if vlt.Spec.cache != nil {
		delete(vlt.Spec.cache, name)
	}
}

func (vlt *Vault) clearCache() {
	vlt.Spec.cache = nil
}

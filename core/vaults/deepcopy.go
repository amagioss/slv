package vaults

func (v *Vault) DeepCopy() *Vault {
	if v == nil {
		return nil
	}
	out := new(Vault)
	v.DeepCopyInto(out)
	return out
}

func (v *Vault) DeepCopyInto(out *Vault) {
	*out = *v
	out.Secrets = make(map[string]string, len(v.Secrets))
	for key, val := range v.Secrets {
		out.Secrets[key] = val
	}
	out.Config = vaultConfig{
		Id:          v.Config.Id,
		PublicKey:   v.Config.PublicKey,
		HashLength:  v.Config.HashLength,
		WrappedKeys: make([]string, len(v.Config.WrappedKeys)),
	}
	copy(out.Config.WrappedKeys, v.Config.WrappedKeys)
	out.vaultSecretRefRegex = v.vaultSecretRefRegex
}

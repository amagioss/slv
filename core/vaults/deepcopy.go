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
	out.Config = vaultConfig{}
	if v.Config.PublicKey != "" {
		out.Config.PublicKey = v.Config.PublicKey
	}
	if v.Config.HashLength != nil {
		var hashLen uint8 = *v.Config.HashLength
		out.Config.HashLength = &hashLen
	}
	out.Config.WrappedKeys = make([]string, len(v.Config.WrappedKeys))
	copy(out.Config.WrappedKeys, v.Config.WrappedKeys)
	out.vaultSecretRefRegex = v.vaultSecretRefRegex
}

package vaults

import (
	"encoding/json"

	"maps"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
)

func (vlt *Vault) Update(name, namespace, secretType string, k8SecretContent []byte) (err error) {
	if !vlt.Spec.writable {
		return errVaultNotWritable
	}
	if err = vlt.validateAndUpdate(); err != nil {
		return err
	}
	if k8SecretContent != nil {
		var secretResource any
		if err = yaml.Unmarshal(k8SecretContent, &secretResource); err != nil {
			return err
		}
		jsonData, err := json.Marshal(secretResource)
		if err != nil {
			return err
		}
		k8secret := &corev1.Secret{}
		if err = json.Unmarshal(jsonData, k8secret); err != nil {
			return err
		}
		metaJson, err := json.Marshal(k8secret.ObjectMeta)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(metaJson, &vlt.ObjectMeta); err != nil {
			return err
		}
		secretDataMap := make(map[string][]byte)
		if k8secret.Data != nil {
			for key, value := range k8secret.Data {
				secretDataMap[key] = value
			}
		}
		if k8secret.StringData != nil {
			for key, value := range k8secret.StringData {
				secretDataMap[key] = []byte(value)
			}
		}
		if len(secretDataMap) > 0 {
			for key, value := range secretDataMap {
				if err = vlt.putWithoutCommit(key, value, true); err != nil {
					return err
				}
			}
			vlt.Type = string(k8secret.Type)
		}
	}
	if name != "" {
		vlt.Name = name
	}
	if vlt.Name == "" {
		return errK8sNameRequired
	}
	if namespace != "" {
		vlt.Namespace = namespace
	}
	if secretType != "" {
		vlt.Type = secretType
	}
	return vlt.commit()
}

func (v *Vault) DeepCopy() *Vault {
	if v == nil {
		return nil
	}
	out := new(Vault)
	v.DeepCopyInto(out)
	return out
}

func (v *Vault) DeepCopyInto(out *Vault) {
	if v == nil || out == nil {
		return
	}
	out.ObjectMeta = v.ObjectMeta
	out.TypeMeta = v.TypeMeta
	out.Type = v.Type
	out.Spec = &VaultSpec{}
	v.Spec.DeepCopyInto(out.Spec)
}

func (v *VaultSpec) DeepCopyInto(out *VaultSpec) {
	if v == nil || out == nil {
		return
	}
	out.path = v.path
	out.Data = make(map[string]string)
	maps.Copy(out.Data, v.Data)
	if v.secretKey != nil {
		out.secretKey = v.secretKey
	}
	out.publicKey = v.publicKey
	out.vaultSecretRefRegex = v.vaultSecretRefRegex
	out.Config = vaultConfig{
		PublicKey:   v.Config.PublicKey,
		Hash:        v.Config.Hash,
		WrappedKeys: v.Config.WrappedKeys,
	}
}

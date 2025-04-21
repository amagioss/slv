package vaults

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
)

func (vlt *Vault) ToK8s(name, namespace string, k8SecretContent []byte) (err error) {
	if name == "" && k8SecretContent == nil {
		return errK8sNameRequired
	}
	vlt.validate()
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
	*out = *v
	v.init()
	// out.Data = make(map[string]string, len(v.Data))
	// for key, val := range v.Data {
	// 	out.Data[key] = val
	// }
	// out.Config = vaultConfig{
	// 	Version:     v.Config.Version,
	// 	Id:          v.Config.Id,
	// 	PublicKey:   v.Config.PublicKey,
	// 	Hash:        v.Config.Hash,
	// 	WrappedKeys: make([]string, len(v.Config.WrappedKeys)),
	// }
	// copy(out.Config.WrappedKeys, v.Config.WrappedKeys)
	// out.vaultSecretRefRegex = v.vaultSecretRefRegex
}

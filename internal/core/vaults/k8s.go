package vaults

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type k8slv struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Type              corev1.SecretType `json:"type,omitempty"`
	Spec              *Vault            `json:"spec" yaml:"spec"`
}

func (vlt *Vault) ToK8s(k8sName string, k8SecretContent []byte) (err error) {
	if k8sName == "" && k8SecretContent == nil {
		return errK8sNameRequired
	}
	if vlt.k8s == nil {
		vlt.k8s = &k8slv{
			TypeMeta: metav1.TypeMeta{
				APIVersion: k8sApiVersion,
				Kind:       k8sKind,
			},
			Spec: vlt,
		}
	}
	if k8sName != "" {
		vlt.k8s.ObjectMeta = metav1.ObjectMeta{
			Name: k8sName,
		}
	}
	if k8SecretContent != nil {
		var secretResource interface{}
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
		if k8secret.Name != "" {
			vlt.k8s.Name = k8secret.Name
		}
		if vlt.k8s.Name == "" {
			return errK8sNameRequired
		}
		vlt.k8s.ObjectMeta = k8secret.ObjectMeta
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
				if err = vlt.putSecretWithoutCommit(key, value); err != nil {
					return err
				}
			}
			vlt.k8s.Type = k8secret.Type
		}
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
	out.Secrets = make(map[string]string, len(v.Secrets))
	for key, val := range v.Secrets {
		out.Secrets[key] = val
	}
	out.Config = vaultConfig{
		Id:          v.Config.Id,
		PublicKey:   v.Config.PublicKey,
		Hash:        v.Config.Hash,
		WrappedKeys: make([]string, len(v.Config.WrappedKeys)),
	}
	copy(out.Config.WrappedKeys, v.Config.WrappedKeys)
	out.vaultSecretRefRegex = v.vaultSecretRefRegex
}

package vaults

import (
	"encoding/json"
	"reflect"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
)

type k8slv struct {
	Kind       string                 `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string                 `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Type       corev1.SecretType      `json:"type,omitempty" yaml:"type,omitempty"`
	Spec       *Vault                 `json:"spec" yaml:"spec"`
}

func structToMap(obj interface{}, toMap map[string]interface{}) {
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if toMap == nil {
		toMap = make(map[string]interface{})
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		toMap[field.Name] = value.Interface()
	}
}

func (vlt *Vault) ToK8s(name, namespace string, k8SecretContent []byte) (err error) {
	if name == "" && k8SecretContent == nil {
		return errK8sNameRequired
	}
	if vlt.k8s == nil {
		vlt.k8s = &k8slv{
			APIVersion: k8sApiVersion,
			Kind:       k8sKind,
			Metadata:   make(map[string]interface{}),
			Spec:       vlt,
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
			vlt.k8s.Metadata["name"] = k8secret.Name
		}
		if vlt.k8s.Metadata["name"] == "" {
			return errK8sNameRequired
		}
		structToMap(k8secret.ObjectMeta, vlt.k8s.Metadata)
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
			vlt.k8s.Type = k8secret.Type
		}
	}
	if name != "" {
		vlt.k8s.Metadata["name"] = name
	}
	if namespace != "" {
		vlt.k8s.Metadata["namespace"] = namespace
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
	out.Data = make(map[string]string, len(v.Data))
	for key, val := range v.Data {
		out.Data[key] = val
	}
	out.Config = vaultConfig{
		Version:     v.Config.Version,
		Id:          v.Config.Id,
		PublicKey:   v.Config.PublicKey,
		Hash:        v.Config.Hash,
		WrappedKeys: make([]string, len(v.Config.WrappedKeys)),
	}
	copy(out.Config.WrappedKeys, v.Config.WrappedKeys)
	out.vaultSecretRefRegex = v.vaultSecretRefRegex
}

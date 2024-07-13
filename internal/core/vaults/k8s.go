package vaults

import (
	"encoding/base64"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

type k8sMeta struct {
	Name string `yaml:"name"`
}

type k8slv struct {
	ApiVersion string  `yaml:"apiVersion"`
	Kind       string  `yaml:"kind"`
	Metadata   k8sMeta `yaml:"metadata"`
	Type       string  `yaml:"type,omitempty"`
	Spec       *Vault  `yaml:"spec"`
}

type k8Secret struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   k8sMeta           `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
	StringData map[string]string `yaml:"stringData"`
	Type       string            `yaml:"type"`
}

func (vlt *Vault) ToK8s(k8sName string, k8SecretContent []byte) (err error) {
	if k8sName == "" && k8SecretContent == nil {
		return errK8sNameRequired
	}
	if vlt.k8s == nil {
		vlt.k8s = &k8slv{
			ApiVersion: k8sApiVersion,
			Kind:       k8sKind,
			Spec:       vlt,
		}
	}
	if k8sName != "" {
		vlt.k8s.Metadata = k8sMeta{Name: k8sName}
	}
	if k8SecretContent != nil {
		k8secret := &k8Secret{}
		if err = yaml.Unmarshal(k8SecretContent, k8secret); err != nil {
			return err
		}
		if k8secret.Metadata.Name != "" {
			vlt.k8s.Metadata.Name = k8secret.Metadata.Name
		}
		if vlt.k8s.Metadata.Name == "" {
			return errK8sNameRequired
		}
		secretDataMap := make(map[string][]byte)
		if k8secret.Data != nil {
			for key, value := range k8secret.Data {
				decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(value))
				secretValue, err := io.ReadAll(decoder)
				if err != nil {
					return err
				}
				secretDataMap[key] = secretValue
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
		HashLength:  v.Config.HashLength,
		WrappedKeys: make([]string, len(v.Config.WrappedKeys)),
	}
	copy(out.Config.WrappedKeys, v.Config.WrappedKeys)
	out.vaultSecretRefRegex = v.vaultSecretRefRegex
}

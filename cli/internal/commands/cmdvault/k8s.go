package cmdvault

import (
	"encoding/base64"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"oss.amagi.com/slv/core/commons"
	"oss.amagi.com/slv/core/crypto"
	"oss.amagi.com/slv/core/input"
	"oss.amagi.com/slv/core/vaults"
)

type k8Secret struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
	StringData map[string]string `yaml:"stringData"`
	Type       string            `yaml:"type"`
}

func k8sSecretFromData(data []byte) (*k8Secret, error) {
	seceret := &k8Secret{}
	if err := yaml.Unmarshal(data, seceret); err != nil {
		return nil, err
	}
	return seceret, nil
}

func newK8sVault(filePath, k8sValue string, hashLength uint8, pq bool, rootPublicKey *crypto.PublicKey, publicKeys ...*crypto.PublicKey) (*vaults.Vault, error) {
	k8slvName := k8sValue
	var secretDataMap map[string][]byte
	var k8sSecretType string
	if strings.HasSuffix(k8sValue, ".yaml") || strings.HasSuffix(k8sValue, ".yml") || strings.HasSuffix(k8sValue, ".json") || k8sValue == "-" {
		var data []byte
		var err error
		if k8sValue == "-" {
			data, err = input.ReadBufferFromStdin("Input the k8s secret object as yaml/json: ")
		} else {
			data, err = os.ReadFile(k8sValue)
		}
		if err != nil {
			return nil, err
		}
		secret, err := k8sSecretFromData(data)
		if err != nil {
			return nil, err
		}
		k8slvName = secret.Metadata.Name
		secretDataMap = make(map[string][]byte)
		if secret.Data != nil {
			for key, value := range secret.Data {
				decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(value))
				secretValue, err := io.ReadAll(decoder)
				if err != nil {
					return nil, err
				}
				secretDataMap[key] = secretValue
			}
		}
		if secret.StringData != nil {
			for key, value := range secret.StringData {
				secretDataMap[key] = []byte(value)
			}
		}
		k8sSecretType = secret.Type
	}
	vault, err := vaults.New(filePath, k8sVaultField, hashLength, pq, rootPublicKey, publicKeys...)
	if err != nil {
		return nil, err
	}
	if len(secretDataMap) > 0 {
		for key, value := range secretDataMap {
			if err = vault.PutSecret(key, value); err != nil {
				return nil, err
			}
		}
	}
	var obj map[string]interface{}
	if err := commons.ReadFromYAML(filePath, &obj); err != nil {
		return nil, err
	}
	obj["apiVersion"] = k8sApiVersion
	obj["kind"] = k8sKind
	obj["metadata"] = map[string]interface{}{
		"name": k8slvName,
	}
	obj["type"] = k8sSecretType
	return vault, commons.WriteToYAML(filePath, "", obj)
}

package vaults

import (
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func (vlt *Vault) yamlTraverseAndUpdateRefSecrets(data *map[string]interface{}, path []string, forceUpdate bool) (err error) {
	for key, value := range *data {
		switch secretValue := value.(type) {
		case map[string]interface{}:
			if err = vlt.yamlTraverseAndUpdateRefSecrets(&secretValue, append(path, key), forceUpdate); err != nil {
				return err
			}
		case []interface{}:
			for i, item := range secretValue {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if err = vlt.yamlTraverseAndUpdateRefSecrets(&itemMap,
						append(path, key+"__"+strconv.Itoa(i)+""), forceUpdate); err != nil {
						return err
					}
				}
			}
		case string:
			if !secretRefRegex.MatchString(secretValue) {
				simplifiedPathName := strings.Join(append(path, key), "_")
				if !forceUpdate && vlt.SecretExists(simplifiedPathName) {
					return ErrVaultSecretExistsAlready
				}
				err = vlt.putSecretWithoutCommit(simplifiedPathName, []byte(secretValue))
				if err == nil {
					(*data)[key] = vlt.getSecretRef(simplifiedPathName)
				}
			}
		}
	}
	return
}

func (vlt *Vault) yamlRef(data []byte, prefix string, forceUpdate bool) (result string, conflicting bool, err error) {
	var yamlMap map[string]interface{}
	err = yaml.Unmarshal(data, &yamlMap)
	if err != nil {
		return
	}
	var path []string
	if prefix != "" {
		path = append(path, prefix)
	}
	err = vlt.yamlTraverseAndUpdateRefSecrets(&yamlMap, path, forceUpdate)
	conflicting = (err == ErrVaultSecretExistsAlready)
	if err == nil {
		updatedYaml, err := yaml.Marshal(yamlMap)
		if err == nil {
			return string(updatedYaml), conflicting, err
		}
	}
	return
}

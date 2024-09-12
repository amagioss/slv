package vaults

import (
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func (vlt *Vault) yamlTraverseAndUpdateRefSecrets(data *map[string]interface{}, path []string, forceUpdate, encrypt bool) (err error) {
	for key, value := range *data {
		switch secretValue := value.(type) {
		case map[string]interface{}:
			if err = vlt.yamlTraverseAndUpdateRefSecrets(&secretValue, append(path, key), forceUpdate, encrypt); err != nil {
				return err
			}
		case []interface{}:
			for i, item := range secretValue {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if err = vlt.yamlTraverseAndUpdateRefSecrets(&itemMap,
						append(path, key+"__"+strconv.Itoa(i)+""), forceUpdate, encrypt); err != nil {
						return err
					}
				}
			}
		case string:
			if !secretRefRegex.MatchString(secretValue) {
				simplifiedPathName := strings.Join(append(path, key), "__")
				if !forceUpdate && vlt.Exists(simplifiedPathName) {
					return errVaultDataExistsAlready
				}
				err = vlt.putWithoutCommit(simplifiedPathName, []byte(secretValue), encrypt)
				if err == nil {
					(*data)[key] = vlt.getSecretRef(simplifiedPathName)
				}
			}
		}
	}
	return
}

func (vlt *Vault) yamlRef(data []byte, prefix string, forceUpdate, encrypt bool) (result string, conflicting bool, err error) {
	var yamlMap map[string]interface{}
	err = yaml.Unmarshal(data, &yamlMap)
	if err != nil {
		return
	}
	var path []string
	if prefix != "" {
		path = append(path, prefix)
	}
	err = vlt.yamlTraverseAndUpdateRefSecrets(&yamlMap, path, forceUpdate, encrypt)
	conflicting = (err == errVaultDataExistsAlready)
	if err == nil {
		updatedYaml, err := yaml.Marshal(yamlMap)
		if err == nil {
			return string(updatedYaml), conflicting, err
		}
	}
	return
}

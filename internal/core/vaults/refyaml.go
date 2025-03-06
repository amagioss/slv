package vaults

import (
	"encoding/json"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func (vlt *Vault) yamlTraverseAndUpdateRefSecrets(data *map[string]any, path []string, forceUpdate, encrypt bool) (err error) {
	for key, value := range *data {
		switch secretValue := value.(type) {
		case map[string]any:
			if err = vlt.yamlTraverseAndUpdateRefSecrets(&secretValue, append(path, key), forceUpdate, encrypt); err != nil {
				return err
			}
		case []any:
			for i, item := range secretValue {
				if itemMap, ok := item.(map[string]any); ok {
					if err = vlt.yamlTraverseAndUpdateRefSecrets(&itemMap,
						append(path, key+"__"+strconv.Itoa(i)+""), forceUpdate, encrypt); err != nil {
						return err
					}
				}
			}
		case string:
			if !secretRefRegex.MatchString(secretValue) {
				simplifiedPathName := strings.Join(append(path, key), "__")
				simplifiedPathName = cleanUnsupportedNameChars(simplifiedPathName)
				if !forceUpdate {
					simplifiedPathName = vlt.getUnusedName(simplifiedPathName)
				}
				err = vlt.putWithoutCommit(simplifiedPathName, []byte(secretValue), encrypt)
				if err == nil {
					(*data)[key] = vlt.getDataRef(simplifiedPathName)
				}
			}
		}
	}
	return
}

func (vlt *Vault) yamlJsonRef(data []byte, prefix string, forceUpdate, encrypt, toJson bool) (result string, conflicting bool, err error) {
	var yamlMap map[string]any
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
		var updatedYaml []byte
		if toJson {
			updatedYaml, err = json.Marshal(yamlMap)
		} else {
			updatedYaml, err = yaml.Marshal(yamlMap)
		}
		if err == nil {
			return string(updatedYaml), conflicting, err
		}
	}
	return
}

package vaults

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
	"slv.sh/slv/internal/core/config"
)

func (vlt *Vault) getDataRef(secretName string) string {
	return fmt.Sprintf("{{%s.%s.%s}}", config.AppNameUpperCase, vlt.Name, secretName)
}

func cleanUnsupportedNameChars(name string) string {
	return unsupportedSecretNameCharRegex.ReplaceAllString(name, "_")
}

func (vlt *Vault) getUnusedName(name string) string {
	nameToUse := name
	for i := 1; vlt.ItemExists(nameToUse); i++ {
		nameToUse = name + "_" + fmt.Sprint(i)
	}
	return nameToUse
}

func (vlt *Vault) refBlob(data []byte, secretName string, forceUpdate, encrypt bool) (result string, conflicting bool, err error) {
	if !secretNameRegex.MatchString(secretName) {
		return "", false, errInvalidVaultItemName
	}
	if !forceUpdate && vlt.ItemExists(secretName) {
		return "", true, errVaultItemExistsAlready
	}
	return vlt.getDataRef(secretName), false, vlt.putWithoutCommit(secretName, data, encrypt)
}

func (vlt *Vault) yamlTraverseAndUpdateRefSecrets(data any, path []string, forceUpdate, encrypt bool) (processed any, updated bool, err error) {
	var valueUpdated bool
	switch secretValue := data.(type) {
	case map[string]any:
		for key := range secretValue {
			value := secretValue[key]
			var processedValue any
			if processedValue, valueUpdated, err = vlt.yamlTraverseAndUpdateRefSecrets(value, append(path, key), forceUpdate, encrypt); err != nil {
				return data, false, err
			}
			secretValue[key] = processedValue
			updated = updated || valueUpdated
		}
	case []any:
		for i := range secretValue {
			value := secretValue[i]
			var processedValue any
			if processedValue, valueUpdated, err = vlt.yamlTraverseAndUpdateRefSecrets(value, append(path, strconv.Itoa(i)), forceUpdate, encrypt); err != nil {
				return data, false, err
			}
			secretValue[i] = processedValue
			updated = updated || valueUpdated
		}
	case string:
		if !secretRefRegex.MatchString(secretValue) {
			simplifiedPathName := strings.Join(path, "__")
			simplifiedPathName = cleanUnsupportedNameChars(simplifiedPathName)
			if !forceUpdate {
				simplifiedPathName = vlt.getUnusedName(simplifiedPathName)
			}
			if err = vlt.putWithoutCommit(simplifiedPathName, []byte(secretValue), encrypt); err != nil {
				return data, false, err
			}
			valueUpdated = true
			return vlt.getDataRef(simplifiedPathName), true, nil
		}
	}
	processed = data
	return
}

func (vlt *Vault) yamlJsonRef(data []byte, prefixOrName string, forceUpdate, encrypt, toJson bool) (result string, conflicting, updated bool, err error) {
	var unmarshaled any
	if err = yaml.Unmarshal(data, &unmarshaled); err != nil {
		return
	}
	var path []string
	if prefixOrName != "" {
		path = append(path, prefixOrName)
	}
	var processed any
	processed, updated, err = vlt.yamlTraverseAndUpdateRefSecrets(unmarshaled, path, forceUpdate, encrypt)
	conflicting = (err == errVaultItemExistsAlready)
	if err == nil {
		var updatedYaml []byte
		if toJson {
			updatedYaml, err = json.Marshal(processed)
		} else {
			updatedYaml, err = yaml.Marshal(processed)
		}
		if err == nil {
			result = string(updatedYaml)
		}
	}
	return
}

func (vlt *Vault) Ref(refType, file, name string, forceUpdate, encrypt, dryRun bool) (result string, conflicting bool, err error) {
	if !vlt.Spec.writable {
		return "", false, errVaultNotWritable
	}
	var data []byte
	data, err = os.ReadFile(file)
	updated := true
	if err == nil {
		switch refType {
		case "yaml", "yml":
			result, conflicting, updated, err = vlt.yamlJsonRef(data, name, forceUpdate, encrypt, false)
		case "json":
			result, conflicting, updated, err = vlt.yamlJsonRef(data, name, forceUpdate, encrypt, true)
		default:
			result, conflicting, err = vlt.refBlob(data, name, forceUpdate, encrypt)
		}
		if err == nil && !conflicting && !dryRun && updated {
			if err = vlt.commit(); err == nil {
				err = os.WriteFile(file, []byte(result), 0644)
			}
		}
		if reloadErr := vlt.reload(); reloadErr != nil {
			err = reloadErr
		}
	}
	return
}

package vaults

import (
	"os"

	"gopkg.in/yaml.v3"
)

func (vlt *Vault) UpdateReferencedSecretsYAML(yamlFile string, preview bool) (string, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return "", err
	}

	var yMap map[string]interface{}
	err = yaml.Unmarshal(data, &yMap)
	if err != nil {
		panic(err)
	}
	err = vlt.updateReferencedSecrets(&yMap, preview)
	if err == nil {
		updatedYaml, err := yaml.Marshal(yMap)
		if err == nil {
			if !preview {
				err = os.WriteFile(yamlFile, updatedYaml, 0644)
			}
			return string(updatedYaml), err
		}
	}
	return "", err
}

func (vlt *Vault) updateReferencedSecrets(data *map[string]interface{}, preview bool) error {
	for key, value := range *data {
		switch v := value.(type) {
		case map[string]interface{}:
			if err := vlt.updateReferencedSecrets(&v, preview); err != nil {
				return err
			}
		case []interface{}:
			for _, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if err := vlt.updateReferencedSecrets(&itemMap, preview); err != nil {
						return err
					}
				}
			}
		case string:
			var newValue string
			if preview {
				newValue = referencedSecretPreviewVal
			} else {
				ref, err := vlt.AddReferencedSecret(v)
				if err != nil {
					return err
				}
				newValue = ref
			}
			(*data)[key] = newValue
		}
	}
	return nil
}

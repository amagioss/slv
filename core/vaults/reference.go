package vaults

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func (vlt *Vault) ReferenceSecrets(file string, previewOnly bool) (string, error) {
	if previewOnly {
		return vlt.updateReferencedSecretsYAML(file, refActionTypePreview)
	} else {
		return vlt.updateReferencedSecretsYAML(file, refActionTypeReference)
	}
}

func (vlt *Vault) DereferenceSecrets(file string) (string, error) {
	return vlt.updateReferencedSecretsYAML(file, refActionTypeDereference)
}

func (vlt *Vault) updateReferencedSecretsYAML(file string, refType refActionType) (string, error) {
	if !strings.HasSuffix(file, ".yaml") && !strings.HasSuffix(file, ".yml") {
		return "", ErrInvalidReferenceFileFormat
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	var yMap map[string]interface{}
	err = yaml.Unmarshal(data, &yMap)
	if err != nil {
		return "", err
	}
	var act action
	switch refType {
	case refActionTypeReference:
		act = vlt.refSecret
	case refActionTypeDereference:
		act = vlt.derefSecret
	case refActionTypePreview:
		act = vlt.previewStr
	default:
		return "", ErrInvalidRefActionType
	}
	err = vlt.yamlTraverser(&yMap, act)
	if err == nil {
		updatedYaml, err := yaml.Marshal(yMap)
		if err == nil {
			if refType != refActionTypePreview {
				err = os.WriteFile(file, updatedYaml, 0644)
				vlt.commit()
			}
			return string(updatedYaml), err
		}
	}
	return "", err
}

type action func(string) string

func (vlt *Vault) refSecret(targetStr string) string {
	ref, err := vlt.addReferencedSecret(targetStr)
	if err != nil {
		return ""
	}
	return ref
}

func (vlt *Vault) derefSecret(targetStr string) string {
	if strings.HasPrefix(targetStr, secretRefPrefix) {
		decrypted, err := vlt.getReferencedSecret(targetStr)
		vlt.deleteReferencedSecret(targetStr)
		if err != nil {
			return ""
		}
		return decrypted
	}
	return targetStr
}

func (vlt *Vault) previewStr(targetStr string) string {
	return referencedSecretPreviewVal
}

func (vlt *Vault) yamlTraverser(data *map[string]interface{}, act action) error {
	for key, value := range *data {
		switch v := value.(type) {
		case map[string]interface{}:
			if err := vlt.yamlTraverser(&v, act); err != nil {
				return err
			}
		case []interface{}:
			for _, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if err := vlt.yamlTraverser(&itemMap, act); err != nil {
						return err
					}
				}
			}
		case string:
			newValue := act(v)
			if newValue == "" {
				return ErrFailedToUpdateSecretReferences
			}
			(*data)[key] = newValue
		}
	}
	return nil
}

package vaults

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func (vlt *Vault) ReferenceSecrets(file string, previewOnly bool) (string, error) {
	return vlt.updateReferencedSecretsYAML(file, refActionReference, previewOnly)
}

func (vlt *Vault) DereferenceSecrets(file string, previewOnly bool) (string, error) {
	return vlt.updateReferencedSecretsYAML(file, refActionDereference, previewOnly)
}

func (vlt *Vault) updateReferencedSecretsYAML(file string, actionType secretRefrencingAction, previewOnly bool) (string, error) {
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
	var act referenceAction
	switch actionType {
	case refActionReference:
		if previewOnly {
			act = vlt.refSecretPreview
		} else {
			act = vlt.refSecret
		}
	case refActionDereference:
		act = vlt.derefSecret
	default:
		return "", ErrInvalidRefActionType
	}
	err = vlt.yamlTraverser(&yMap, act)
	if err == nil {
		updatedYaml, err := yaml.Marshal(yMap)
		if err == nil {
			if !previewOnly {
				err = os.WriteFile(file, updatedYaml, 0644)
				vlt.commit()
			}
			return string(updatedYaml), err
		}
	}
	return "", err
}

type referenceAction func(string) (string, error)

func (vlt *Vault) refSecret(targetStr string) (string, error) {
	return vlt.addReferencedSecret("", targetStr)
}

func (vlt *Vault) refSecretPreview(targetStr string) (string, error) {
	return autoReferencedPreviewValue, nil
}

func (vlt *Vault) derefSecret(targetStr string) (string, error) {
	if strings.HasPrefix(targetStr, autoReferencedPrefix+vlt.Config.PublicKey.Id()) {
		decrypted, err := vlt.getReferencedSecret(targetStr)
		if err != nil {
			return "", err
		}
		vlt.deleteReferencedSecret(targetStr)
		return decrypted, nil
	}
	return targetStr, nil
}

func (vlt *Vault) yamlTraverser(data *map[string]interface{}, act referenceAction) error {
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
			newValue, err := act(v)
			if err != nil {
				return err
			}
			(*data)[key] = newValue
		}
	}
	return nil
}

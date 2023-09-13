package vaults

import (
	"os"
	"strconv"
	"strings"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"gopkg.in/yaml.v3"
)

func (vlt *Vault) autoReferenceSecret(path, secret string) (secretRef string, err error) {
	var sealedSecret *crypto.SealedSecret
	sealedSecret, err = vlt.Config.PublicKey.EncryptSecretString(secret, vlt.Config.HashLength)
	if err == nil {
		if vlt.vault.Secrets.Referenced == nil {
			vlt.vault.Secrets.Referenced = make(map[string]*crypto.SealedSecret)
		}
		secretRef = autoReferencedPrefix + vlt.Config.PublicKey.IdStr() + "_" +
			commons.Encode([]byte(path))
		vlt.vault.Secrets.Referenced[secretRef] = sealedSecret
		err = vlt.commit()
	}
	return
}

func (vlt *Vault) yamlTraverseAndUpdateRefSecrets(data *map[string]interface{},
	path []string, previewOnly bool) (err error) {
	for key, value := range *data {
		switch v := value.(type) {
		case map[string]interface{}:
			if err = vlt.yamlTraverseAndUpdateRefSecrets(&v, append(path, key),
				previewOnly); err != nil {
				return err
			}
		case []interface{}:
			for i, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if err = vlt.yamlTraverseAndUpdateRefSecrets(&itemMap,
						append(path, key+"["+strconv.Itoa(i)+"]"), previewOnly); err != nil {
						return err
					}
				}
			}
		case string:
			if previewOnly {
				v = autoReferencedPreviewValue
			} else {
				v, err = vlt.autoReferenceSecret(strings.Join(append(path, key), "."), v)
				if err != nil {
					return err
				}
			}
			(*data)[key] = v
		}
	}
	return nil
}

func (vlt *Vault) RefSecrets(file string, previewOnly bool) (string, error) {
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
	err = vlt.yamlTraverseAndUpdateRefSecrets(&yMap, []string{}, previewOnly)
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

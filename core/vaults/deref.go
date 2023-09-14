package vaults

// func getSecretReferencesFromFile(file string) ([]string, error) {
// 	bytes, err := os.ReadFile(file)
// 	if err != nil {
// 		panic(err)
// 	}
// 	autoRefPatter := regexp.MustCompile("" + "[A-Za-z0-9]+_[A-Za-z0-9]+")
// 	return autoRefPatter.FindAllString(string(bytes), -1), nil
// }

// func (vlt *Vault) getReferencedSecretIfRequired(secretRef string) (secret string, err error) {
// 	if vlt.IsLocked() {
// 		return secret, ErrVaultLocked
// 	}
// 	encryptedData, ok := vlt.vault.Secrets.Referenced[secretRef]
// 	if !ok {
// 		return "", ErrVaultSecretNotFound
// 	}
// 	secretBytes, err := vlt.secretKey.DecryptSecret(*encryptedData)
// 	return string(secretBytes), err
// }

// func (vlt *Vault) deleteReferencedSecret(secretReference string) {
// 	delete(vlt.vault.Secrets.Referenced, secretReference)
// }

// func (vlt *Vault) derefSecret(secretRef string) (string, error) {
// 	if strings.HasPrefix(secretRef, autoReferencedPrefix+vlt.Config.PublicKey.IdStr()) {
// 		decrypted, err := vlt.getReferencedSecretIfRequired(secretRef)
// 		if err != nil {
// 			return "", err
// 		}
// 		vlt.deleteReferencedSecret(secretRef)
// 		return decrypted, nil
// 	}
// 	return secretRef, nil
// }

// func (vlt *Vault) DeRefSecrets(file string, previewOnly bool) (string, error) {
// 	data, err := os.ReadFile(file)
// 	if err != nil {
// 		return "", err
// 	}
// 	var yMap map[string]interface{}
// 	err = yaml.Unmarshal(data, &yMap)
// 	if err != nil {
// 		return "", err
// 	}
// 	err = vlt.yamlTraverseAndUpdateRefSecrets(&yMap, []string{}, previewOnly)
// 	if err == nil {
// 		updatedYaml, err := yaml.Marshal(yMap)
// 		if err == nil {
// 			if !previewOnly {
// 				err = os.WriteFile(file, updatedYaml, 0644)
// 				vlt.commit()
// 			}
// 			return string(updatedYaml), err
// 		}
// 	}
// 	return "", err
// }

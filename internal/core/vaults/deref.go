package vaults

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"slv.sh/slv/internal/core/commons"
)

func (vlt *Vault) getDataByReference(secretRef string) ([]byte, error) {
	sliced := strings.Split(secretRef, "."+vlt.Name+".")
	if len(sliced) != 2 {
		return nil, errInvalidReferenceFormat
	}
	secretName := strings.Trim(sliced[1], " }")
	if item, err := vlt.Get(secretName); err != nil {
		return nil, err
	} else if itemValue, err := item.Value(); err != nil {
		return nil, err
	} else {
		return itemValue, nil
	}
}

func (vlt *Vault) getVaultSecretRefRegex() *regexp.Regexp {
	if vlt.Spec.vaultSecretRefRegex == nil {
		vlt.Spec.vaultSecretRefRegex = regexp.MustCompile(strings.ReplaceAll(secretRefPatternBase, vaultNamePatternPlaceholder, vlt.Name))
	}
	return vlt.Spec.vaultSecretRefRegex
}

func (vlt *Vault) deRefContent(content string) ([]byte, error) {
	vaultSecretRefRegex := vlt.getVaultSecretRefRegex()
	secretRefs := vaultSecretRefRegex.FindAllString(content, -1)
	if len(secretRefs) == 1 && len(content) == len(secretRefs[0]) {
		return vlt.getDataByReference(secretRefs[0])
	} else {
		for _, secretRef := range secretRefs {
			decrypted, err := vlt.getDataByReference(secretRef)
			if err != nil {
				return nil, err
			}
			content = strings.Replace(content, secretRef, string(decrypted), -1)
		}
	}
	return []byte(content), nil
}

func (vlt *Vault) DeRef(file string, previewOnlyMode bool) (string, error) {
	if vlt.IsLocked() {
		return "", errVaultLocked
	}
	if !commons.FileExists(file) {
		return "", fmt.Errorf("file does not exist")
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	derefedBytes, err := vlt.deRefContent(string(data))
	if err != nil {
		return "", err
	}
	if previewOnlyMode {
		return string(derefedBytes), nil
	}
	return "", os.WriteFile(file, derefedBytes, 0644)
}

package vaults

import (
	"os"
	"regexp"
	"strings"
)

func (vlt *Vault) getSecretByReference(secretRef string) (secret string, err error) {
	sliced := strings.Split(secretRef, vlt.Id()+".")
	if len(sliced) != 2 {
		return "", ErrInvalidReferenceFormat
	}
	secretName := strings.Trim(sliced[1], " }")
	return vlt.GetSecret(secretName)
}

func (vlt *Vault) getVaultSecretRefRegex() *regexp.Regexp {
	if vlt.vaultSecretRefRegex == nil {
		vlt.vaultSecretRefRegex = regexp.MustCompile(strings.ReplaceAll(secretRefPatternBase, "VAULTID", vlt.Id()))
	}
	return vlt.vaultSecretRefRegex
}

func (vlt *Vault) deRefSecretsFromContent(content string) (string, error) {
	vaultSecretRefRegex := vlt.getVaultSecretRefRegex()
	secretRefs := vaultSecretRefRegex.FindAllString(content, -1)
	for _, secretRef := range secretRefs {
		decrypted, err := vlt.getSecretByReference(secretRef)
		if err != nil {
			return "", err
		}
		content = strings.Replace(content, secretRef, decrypted, -1)
	}
	return content, nil
}

func (vlt *Vault) DeRefSecrets(file string, previewOnly bool) (derefedContent string, err error) {
	if vlt.IsLocked() {
		return "", ErrVaultLocked
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	derefedContent, err = vlt.deRefSecretsFromContent(string(data))
	if err != nil {
		return "", err
	}
	if !previewOnly {
		err = os.WriteFile(file, []byte(derefedContent), 0644)
	}
	return
}

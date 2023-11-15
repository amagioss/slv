package vaults

import (
	"os"
	"regexp"
	"strings"
)

func (vlt *Vault) getSecretByReference(secretRef string) (secret []byte, err error) {
	sliced := strings.Split(secretRef, vlt.Id()+".")
	if len(sliced) != 2 {
		return nil, ErrInvalidReferenceFormat
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

func (vlt *Vault) deRefSecretsFromContent(content string) ([]byte, error) {
	vaultSecretRefRegex := vlt.getVaultSecretRefRegex()
	secretRefs := vaultSecretRefRegex.FindAllString(content, -1)
	if len(secretRefs) == 1 && len(content) == len(secretRefs[0]) {
		return vlt.getSecretByReference(secretRefs[0])
	} else {
		for _, secretRef := range secretRefs {
			decrypted, err := vlt.getSecretByReference(secretRef)
			if err != nil {
				return nil, err
			}
			content = strings.Replace(content, secretRef, string(decrypted), -1)
		}
	}
	return []byte(content), nil
}

func (vlt *Vault) DeRefSecrets(file string, previewOnly bool) (dereferncedBytes []byte, err error) {
	if vlt.IsLocked() {
		return nil, ErrVaultLocked
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	dereferncedBytes, err = vlt.deRefSecretsFromContent(string(data))
	if err != nil {
		return nil, err
	}
	if !previewOnly {
		err = os.WriteFile(file, dereferncedBytes, 0644)
	}
	return
}

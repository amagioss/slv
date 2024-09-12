package vaults

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (vlt *Vault) getSecretByReference(secretRef string) ([]byte, error) {
	sliced := strings.Split(secretRef, vlt.Id()+".")
	if len(sliced) != 2 {
		return nil, errInvalidReferenceFormat
	}
	secretName := strings.Trim(sliced[1], " }")
	vd, err := vlt.Get(secretName)
	return vd.value, err
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

func (vlt *Vault) DeRefSecrets(path string) error {
	if vlt.IsLocked() {
		return errVaultLocked
	}
	return filepath.WalkDir(path, func(currentPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := os.ReadFile(currentPath)
		if err != nil {
			return err
		}
		dereferncedBytes, err := vlt.deRefSecretsFromContent(string(data))
		if err != nil {
			return err
		}
		return os.WriteFile(currentPath, dereferncedBytes, 0644)
	})
}

package vaults

import (
	"os"
	"regexp"
	"strings"

	"github.com/shibme/slv/core/commons"
)

func (vlt *Vault) getAutoRefdSecret(secretRef string) (secret string, err error) {
	encryptedAutoRefdSecret, ok := vlt.vault.Secrets.Referenced[secretRef]
	if !ok {
		return "", ErrVaultSecretNotFound
	}
	secretBytes, err := vlt.secretKey.DecryptSecret(*encryptedAutoRefdSecret)
	return string(secretBytes), err
}

func (vlt *Vault) getDirectRefdSecret(secretRef string) (secret string, err error) {
	sliced := strings.Split(secretRef, "_")
	if len(sliced) != 4 || sliced[0] != commons.SLV || sliced[1] != directRefAbbrev ||
		sliced[2] != vlt.Id() || sliced[3][0] != '[' || sliced[3][len(sliced[3])-1] != ']' {
		return "", ErrInvalidDirectRefString
	}
	directRefName := strings.Trim(sliced[3], "[]")
	return vlt.GetDirectSecret(directRefName)
}

func (vlt *Vault) deRefSecretsFromContent(content string) (string, error) {
	autoRefPattern := regexp.MustCompile(commons.SLV + "_" + autoRefAbbrev + "_" +
		vlt.Id() + "_[A-Za-z0-9]+_[A-Za-z0-9]+")
	directRefPattern := regexp.MustCompile(commons.SLV + "_" + directRefAbbrev + "_" +
		vlt.Id() + "_\\[[A-Za-z0-9]+\\]")
	autoRefs := autoRefPattern.FindAllString(content, -1)
	for _, autoRef := range autoRefs {
		decrypted, err := vlt.getAutoRefdSecret(autoRef)
		if err != nil {
			return "", err
		}
		content = strings.Replace(content, autoRef, decrypted, -1)
	}
	directRefs := directRefPattern.FindAllString(content, -1)
	for _, directRef := range directRefs {
		decrypted, err := vlt.getDirectRefdSecret(directRef)
		if err != nil {
			return "", err
		}
		content = strings.Replace(content, directRef, decrypted, -1)
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
		vlt.commit()
	}
	return
}

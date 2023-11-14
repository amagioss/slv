package vaults

import (
	"fmt"
	"os"

	"github.com/shibme/slv/core/commons"
)

func (vlt *Vault) getSecretRef(secretName string) string {
	return fmt.Sprintf("{{%s_%s_%s.%s}}", commons.SLV, secretRefAbbrev, vlt.Id(), secretName)
}

func (vlt *Vault) refBlob(data []byte, secretName string, forceUpdate bool) (result string, conflicting bool, err error) {
	if !forceUpdate && vlt.SecretExists(secretName) {
		return "", true, ErrVaultSecretExistsAlready
	}
	return vlt.getSecretRef(secretName), false, vlt.putSecretWithoutCommit(secretName, data)
}

func (vlt *Vault) RefSecrets(refType, file, name string, forceUpdate, dryRun bool) (result string, conflicting bool, err error) {
	var data []byte
	data, err = os.ReadFile(file)
	if err == nil {
		switch refType {
		case "yaml", "yml":
			result, conflicting, err = vlt.yamlRef(data, name, forceUpdate)
		default:
			result, conflicting, err = vlt.refBlob(data, name, forceUpdate)
		}
		if err == nil && !conflicting && !dryRun {
			if err = vlt.commit(); err == nil {
				err = os.WriteFile(file, []byte(result), 0644)
			}
		}
		if resetErr := vlt.reset(); resetErr != nil {
			err = resetErr
		}
	}
	return
}

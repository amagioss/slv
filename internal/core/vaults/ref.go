package vaults

import (
	"fmt"
	"os"
)

func (vlt *Vault) getSecretRef(secretName string) string {
	return fmt.Sprintf("{{%s.%s}}", vlt.Id(), secretName)
}

func (vlt *Vault) refBlob(data []byte, secretName string, forceUpdate, encrypt bool) (result string, conflicting bool, err error) {
	if !forceUpdate && vlt.Exists(secretName) {
		return "", true, errVaultDataExistsAlready
	}
	return vlt.getSecretRef(secretName), false, vlt.putWithoutCommit(secretName, data, encrypt)
}

func (vlt *Vault) RefSecrets(refType, file, name string, forceUpdate, encrypt, dryRun bool) (result string, conflicting bool, err error) {
	var data []byte
	data, err = os.ReadFile(file)
	if err == nil {
		switch refType {
		case "yaml", "yml":
			result, conflicting, err = vlt.yamlRef(data, name, forceUpdate, encrypt)
		default:
			result, conflicting, err = vlt.refBlob(data, name, forceUpdate, encrypt)
		}
		if err == nil && !conflicting && !dryRun {
			if err = vlt.commit(); err == nil {
				err = os.WriteFile(file, []byte(result), 0644)
			}
		}
		if reseterr := vlt.reset(); reseterr != nil {
			err = reseterr
		}
	}
	return
}

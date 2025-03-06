package vaults

import (
	"fmt"
	"os"
)

func (vlt *Vault) getDataRef(secretName string) string {
	return fmt.Sprintf("{{%s.%s}}", vlt.Id(), secretName)
}

func cleanUnsupportedNameChars(name string) string {
	return unsupportedSecretNameCharRegex.ReplaceAllString(name, "_")
}

func (vlt *Vault) getUnusedName(name string) string {
	nameToUse := name
	for i := 1; vlt.Exists(nameToUse); i++ {
		nameToUse = name + "_" + fmt.Sprint(i)
	}
	return nameToUse
}

func (vlt *Vault) refBlob(data []byte, secretName string, forceUpdate, encrypt bool) (result string, conflicting bool, err error) {
	if !secretNameRegex.MatchString(secretName) {
		return "", false, errInvalidVaultDataName
	}
	if !forceUpdate && vlt.Exists(secretName) {
		return "", true, errVaultDataExistsAlready
	}
	return vlt.getDataRef(secretName), false, vlt.putWithoutCommit(secretName, data, encrypt)
}

func (vlt *Vault) Ref(refType, file, name string, forceUpdate, encrypt, dryRun bool) (result string, conflicting bool, err error) {
	var data []byte
	data, err = os.ReadFile(file)
	if err == nil {
		switch refType {
		case "yaml", "yml":
			result, conflicting, err = vlt.yamlJsonRef(data, name, forceUpdate, encrypt, false)
		case "json":
			result, conflicting, err = vlt.yamlJsonRef(data, name, forceUpdate, encrypt, true)
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

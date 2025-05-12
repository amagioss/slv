package main

import (
	"C"

	"slv.sh/slv"
)
import (
	"encoding/json"
	"unsafe"

	"slv.sh/slv/internal/core/vaults"
)

func getSecret(vaultPath, secretName *C.char, secretValue **C.char, secretLength *C.int, errMessage **C.char, errLength *C.int) {
	vaultFile := C.GoString(vaultPath)
	name := C.GoString(secretName)
	var err error
	var vaultItem *vaults.VaultItem
	if vaultItem, err = slv.GetVaultItem(vaultFile, name); err == nil {
		var valueBytes []byte
		if valueBytes, err = vaultItem.Value(); err == nil {
			*secretValue = (*C.char)(C.CBytes(valueBytes))
			*secretLength = C.int(len(valueBytes))
			*errMessage = nil
			*errLength = 0
			return
		}
	}
	*secretValue = nil
	*secretLength = 0
	*errMessage = C.CString(err.Error())
	*errLength = C.int(len(err.Error()))
}

func getAllSecrets(vaultPath *C.char, secretsJson **C.char, secretsJsonLength *C.int, errMessage **C.char, errLength *C.int) {
	vaultFile := C.GoString(vaultPath)
	secrets, err := slv.GetAllVaultItems(vaultFile)
	if err == nil {
		var jsonBytes []byte
		if jsonBytes, err = json.Marshal(secrets); err == nil {
			*secretsJson = (*C.char)(C.CBytes(jsonBytes))
			*secretsJsonLength = C.int(len(jsonBytes))
			*errMessage = nil
			*errLength = 0
			return
		}
	}
	*secretsJson = nil
	*secretsJsonLength = 0
	*errMessage = C.CString(err.Error())
	*errLength = C.int(len(err.Error()))
}

func putSecret(vaultPath, secretName, secretValue *C.char, errMessage **C.char, errLength *C.int) {
	vaultFile := C.GoString(vaultPath)
	name := C.GoString(secretName)
	value := C.GoBytes(unsafe.Pointer(secretValue), C.int(len(C.GoString(secretValue))))
	if err := slv.PutVaultItem(vaultFile, name, value, true); err != nil {
		*errMessage = C.CString(err.Error())
		*errLength = C.int(len(err.Error()))
	} else {
		*errMessage = nil
		*errLength = 0
	}
}

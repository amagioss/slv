package main

import (
	"C"

	"slv.sh/slv"
)
import (
	"encoding/json"
	"unsafe"
)

func getSecret(vaultPath, secretName *C.char, secretValue **C.char, secretLength *C.int, errMessage **C.char, errLength *C.int) {
	vaultFile := C.GoString(vaultPath)
	name := C.GoString(secretName)
	if vaultItem, err := slv.GetVaultItem(vaultFile, name); err != nil {
		*secretValue = nil
		*secretLength = 0
		*errMessage = C.CString(err.Error())
		*errLength = C.int(len(err.Error()))
	} else {
		*secretValue = (*C.char)(C.CBytes(vaultItem.Value()))
		*secretLength = C.int(len(vaultItem.Value()))
		*errMessage = nil
		*errLength = 0
	}
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

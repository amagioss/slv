package main

import (
	"C"

	"oss.amagi.com/slv"
)
import "encoding/json"

func getSecret(vaultPath *C.char, secretName *C.char, secretValue **C.char, secretLength *C.int, errMessage **C.char, errLength *C.int) {
	vaultFile := C.GoString(vaultPath)
	name := C.GoString(secretName)
	if value, err := slv.GetSecret(vaultFile, name); err != nil {
		*secretValue = nil
		*secretLength = 0
		*errMessage = C.CString(err.Error())
		*errLength = C.int(len(err.Error()))
	} else {
		*secretValue = (*C.char)(C.CBytes(value))
		*secretLength = C.int(len(value))
		*errMessage = nil
		*errLength = 0
	}
}

func getAllSecrets(vaultPath *C.char, secretsJson **C.char, secretsJsonLength *C.int, errMessage **C.char, errLength *C.int) {
	vaultFile := C.GoString(vaultPath)
	secrets, err := slv.GetAllSecrets(vaultFile)
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

package main

import (
	"C"
	"encoding/json"

	"oss.amagi.com/slv"
)

//export GetSecret
func GetSecret(vaultPath *C.char, secretName *C.char, secretValue **C.char, secretLength *C.int, errMessage **C.char, errLength *C.int) {
	vaultFile := C.GoString(vaultPath)
	name := C.GoString(secretName)
	if value, err := slv.GetSecret(vaultFile, name); err != nil {
		*errLength = C.int(len(err.Error()))
		*errMessage = C.CString(err.Error())
		*secretValue = nil
		*secretLength = 0
	} else {
		*secretLength = C.int(len(value))
		*secretValue = (*C.char)(C.CBytes(value))
		*errMessage = nil
		*errLength = 0
	}
}

//export GetAllSecrets
func GetAllSecrets(vaultPath *C.char, secretsJson **C.char, secretsJsonLength *C.int, errMessage **C.char, errLength *C.int) {
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
	*errLength = C.int(len(err.Error()))
	*errMessage = C.CString(err.Error())
	*secretsJson = nil
	*secretsJsonLength = 0
}

func main() {}

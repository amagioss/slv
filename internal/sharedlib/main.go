package main

import (
	"C"
)

// SLVGetSecret retuns the value of a single secret from a given vault for the specified secret name
//
//export SLVGetSecret
func SLVGetSecret(vaultPath *C.char, secretName *C.char, secretValue **C.char, secretLength *C.int, errMessage **C.char, errLength *C.int) {
	getSecret(vaultPath, secretName, secretValue, secretLength, errMessage, errLength)
}

// SLVGetAllSecrets returns all the secrets from a given vault as a JSON string
//
//export SLVGetAllSecrets
func SLVGetAllSecrets(vaultPath *C.char, secretsJson **C.char, secretsJsonLength *C.int, errMessage **C.char, errLength *C.int) {
	getAllSecrets(vaultPath, secretsJson, secretsJsonLength, errMessage, errLength)
}

// SLVPutSecret writes a secret to the vault
//
//export SLVPutSecret
func SLVPutSecret(vaultPath *C.char, secretName *C.char, secretValue *C.char, errMessage **C.char, errLength *C.int) {
	putSecret(vaultPath, secretName, secretValue, errMessage, errLength)
}

func main() {}

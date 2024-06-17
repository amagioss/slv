package providers

import (
	"crypto/sha256"
	"fmt"

	"dev.shib.me/xipher"
	"github.com/zalando/go-keyring"
	"oss.amagi.com/slv/internal/core/commons"
	"oss.amagi.com/slv/internal/core/input"
)

func bindWithPassword(skBytes []byte, inputs map[string][]byte) (ref map[string][]byte, err error) {
	password := inputs["password"]
	if len(password) == 0 {
		return nil, err
	}
	xipherKey, err := xipher.NewSecretKeyForPassword(password)
	if err != nil {
		return nil, err
	}
	sealedSecretKeyBytes, err := xipherKey.Encrypt(skBytes, false)
	if err != nil {
		return nil, err
	}
	ref = make(map[string][]byte)
	ref["ssk"] = sealedSecretKeyBytes
	return
}

func getFromKeyring(sealedSecretKeyBytes []byte) (string, error) {
	sha256sum := sha256.Sum256(sealedSecretKeyBytes)
	return keyring.Get(keyringServiceName, commons.Encode(sha256sum[:]))
}

func putToKeyring(sealedSecretKeyBytes []byte, password string) error {
	sha256sum := sha256.Sum256(sealedSecretKeyBytes)
	return keyring.Set(keyringServiceName, commons.Encode(sha256sum[:]), password)
}

func unBindWithPassword(ref map[string][]byte) (secretKeyBytes []byte, err error) {
	sealedSecretKeyBytes := ref["ssk"]
	if len(sealedSecretKeyBytes) == 0 {
		return nil, errSealedSecretKeyRef
	}
	var password []byte
	setPasswordToKeyring := false
	if input.IsInteractive() == nil {
		pwd, err := getFromKeyring(sealedSecretKeyBytes)
		if err == nil {
			password = []byte(pwd)
		} else {
			if err == keyring.ErrNotFound {
				setPasswordToKeyring = true
			}
			if password, err = input.GetHiddenInput("Enter Password: "); err != nil {
				return nil, err
			}
		}
	}
	if password == nil {
		return nil, errPasswordNotSet
	}
	xipherKey, err := xipher.NewSecretKeyForPassword([]byte(password))
	if err != nil {
		return nil, err
	}
	secretKeyBytes, err = xipherKey.Decrypt(sealedSecretKeyBytes)
	if err != nil {
		return nil, errInvalidPassword
	}
	if setPasswordToKeyring {
		confirm, _ := input.GetConfirmation("Do you want to save the password in keyring? (y/n): ", "y")
		if confirm {
			if err := putToKeyring(sealedSecretKeyBytes, string(password)); err != nil {
				fmt.Println("Failed to save password in keyring: ", err.Error())
			}
		}
	}
	return
}

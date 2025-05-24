package providers

import (
	"fmt"

	"slv.sh/slv/internal/core/input"
	"slv.sh/slv/internal/core/keystore"
	"xipher.org/xipher"
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
	sealedSecretKeyBytes, err := xipherKey.Encrypt(skBytes, false, false)
	if err != nil {
		return nil, err
	}
	ref = make(map[string][]byte)
	ref["ssk"] = sealedSecretKeyBytes
	return
}

func unBindWithPassword(ref map[string][]byte) (secretKeyBytes []byte, err error) {
	sealedSecretKeyBytes := ref["ssk"]
	if len(sealedSecretKeyBytes) == 0 {
		return nil, errSealedSecretKeyRef
	}
	var password []byte
	setPasswordToKeystore := false
	if input.IsInteractive() == nil {
		if password, err = keystore.Get(sealedSecretKeyBytes, false); err == keystore.ErrNotFound {
			setPasswordToKeystore = true
			if password, err = input.GetHiddenInput("Enter Password: "); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to get password from keystore: %w", err)
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
	if setPasswordToKeystore {
		confirm, _ := input.GetConfirmation("Do you want to save the password locally? (y/n): ", "y")
		if confirm {
			if err := keystore.Put(sealedSecretKeyBytes, password, false); err != nil {
				return nil, fmt.Errorf("failed to save password to keystore: %w", err)
			}
		}
	}
	return
}

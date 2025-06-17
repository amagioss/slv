package envproviders

import (
	"errors"
	"fmt"
	"os"

	"slv.sh/slv/internal/core/input"
	"slv.sh/slv/internal/core/keystore"
	"xipher.org/xipher"
)

const (
	PasswordProviderId   = "password"
	passwordProviderName = "Password"
	passwordProviderDesc = "Password to set for the environment"
	passwordRefName      = "password"

	envar_SLV_ENV_SECRET_PASSWORD = "SLV_ENV_SECRET_PASSWORD"
)

var (
	errPasswordNotSet  = errors.New("password not set: please set password through the environment variable or use the interactive terminal to enter the password")
	errInvalidPassword = errors.New("invalid password")

	pwdArgs = []arg{
		{
			id:          passwordRefName,
			required:    true,
			description: "Password to use",
		},
	}
)

func bindWithPassword(skBytes []byte, inputs map[string]string) (ref map[string][]byte, err error) {
	password := inputs[passwordRefName]
	if len(password) == 0 {
		return nil, err
	}
	xipherKey, err := xipher.NewSecretKeyForPassword([]byte(password))
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
	if passwordStr := os.Getenv(envar_SLV_ENV_SECRET_PASSWORD); passwordStr != "" {
		password = []byte(passwordStr)
	} else if input.IsInteractive() {
		if password, err = keystore.Get(sealedSecretKeyBytes, false); err != nil {
			if err == keystore.ErrNotFound {
				setPasswordToKeystore = true
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

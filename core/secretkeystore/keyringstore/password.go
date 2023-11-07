package keyringstore

import (
	"os"
	"path/filepath"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/environments"
)

func RegisterUser(name, email, password string) (env *environments.Environment, err error) {
	if _, err = getSalt(); err != nil {
		return nil, UserAlreadyRegistered
	}
	var salt []byte
	var secretKey *crypto.SecretKey
	if secretKey, salt, err = crypto.NewSecretKeyForPassword([]byte(password), environments.EnvironmentKey); err == nil {
		if env, err = environments.NewEnvironmentForSecretKey(name, email, environments.USER, secretKey); err == nil {
			if err = setSalt(salt); err == nil {
				if err = putUserSecretToKeyring([]byte(secretKey.String())); err == nil {
					return env, nil
				}
			}
		}
	}
	return
}

func setSalt(salt []byte) error {
	userDir := filepath.Join(commons.AppDataDir(), userDataDirName)
	userDirInfo, err := os.Stat(userDir)
	if err != nil {
		err = os.MkdirAll(userDir, 0755)
		if err != nil {
			return ErrCreatingUserDataDir
		}
	} else if !userDirInfo.IsDir() {
		return ErrFileExistsInPathOfUserDataDir
	}
	return commons.WriteToFile(filepath.Join(userDir, userSaltFileName), salt)
}

func getSalt() (salt []byte, err error) {
	return os.ReadFile(filepath.Join(commons.AppDataDir(), userDataDirName, userSaltFileName))
}

func getPassphraseFromEnvar() (string, error) {
	passphrase := os.Getenv(slvEnvPasswordEnvarName)
	if passphrase == "" {
		return "", ErrPassphraseNotSet
	}
	return passphrase, nil
}

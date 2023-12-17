package keyringstore

import (
	"os"
	"path/filepath"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/environments"
)

func RegisterUser(name, password string) (env *environments.Environment, err error) {
	if _, err = getSalt(); err != nil {
		return nil, errUserAlreadyRegistered
	}
	var salt []byte
	var secretKey *crypto.SecretKey
	if secretKey, salt, err = crypto.NewSecretKeyForPassword([]byte(password), environments.EnvironmentKey); err == nil {
		publicKey, err := secretKey.PublicKey()
		if err != nil {
			return nil, err
		}
		if env, err = environments.NewEnvironmentForPublicKey(name, environments.USER, publicKey); err == nil {
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
			return errCreatingUserDataDir
		}
	} else if !userDirInfo.IsDir() {
		return errFileExistsInPathOfUserDataDir
	}
	return commons.WriteToFile(filepath.Join(userDir, userSaltFileName), salt)
}

func getSalt() (salt []byte, err error) {
	return os.ReadFile(filepath.Join(commons.AppDataDir(), userDataDirName, userSaltFileName))
}

func getPassphraseFromEnvar() (string, error) {
	passphrase := os.Getenv(slvEnvPasswordEnvarName)
	if passphrase == "" {
		return "", errPassphraseNotSet
	}
	return passphrase, nil
}

package keyringstore

import "errors"

const (
	userSaltFileName        = ".salt"
	userDataDirName         = "user"
	slvEnvPasswordEnvarName = "SLV_ENV_PASSWORD"
	slvKeyringServiceName   = "slv"
	slvKeyringItemKey       = "slv_secret_key"
)

var ErrPassphraseNotSet = errors.New(slvEnvPasswordEnvarName + " environment variable is not set")
var ErrCreatingUserDataDir = errors.New("error in creating user data directory")
var ErrFileExistsInPathOfUserDataDir = errors.New("file exists in path of user data directory")
var ErrSaltNotFound = errors.New("salt not found")
var UserAlreadyRegistered = errors.New("user already registered")

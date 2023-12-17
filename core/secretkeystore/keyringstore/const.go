package keyringstore

import "errors"

const (
	userSaltFileName        = ".salt"
	userDataDirName         = "user"
	slvEnvPasswordEnvarName = "SLV_ENV_PASSWORD"
	slvKeyringServiceName   = "slv"
	slvKeyringItemKey       = "slv_secret_key"
)

var errPassphraseNotSet = errors.New(slvEnvPasswordEnvarName + " environment variable is not set")
var errCreatingUserDataDir = errors.New("error in creating user data directory")
var errFileExistsInPathOfUserDataDir = errors.New("file exists in path of user data directory")
var errSaltNotFound = errors.New("salt not found")
var errUserAlreadyRegistered = errors.New("user already registered")

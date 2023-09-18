package keystore

import "errors"

const (
	userSaltFileName       = ".salt"
	userDataDirName        = "user"
	slvSecreKeyEnvarName   = "SLV_SECRET_KEY"
	slvPassphraseEnvarName = "SLV_USER_PASSWORD"
	slvKeyringServiceName  = "slv"
	slvKeyringItemKey      = "slv_secret_key"
)

var ErrEnvSecretNotSet = errors.New(slvSecreKeyEnvarName + " environment variable is not set")
var ErrCreatingUserDataDir = errors.New("error in creating user data directory")
var ErrFileExistsInPathOfUserDataDir = errors.New("file exists in path of user data directory")
var ErrSaltNotFound = errors.New("salt not found")
var UserAlreadyRegistered = errors.New("user already registered")

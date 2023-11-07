package keyringstore

import "errors"

const (
	userSaltFileName       = ".salt"
	userDataDirName        = "user"
	slvSecreKeyEnvarName   = "SLV_SECRET_KEY"
	slvPassphraseEnvarName = "SLV_USER_PASSWORD"
	slvKeyringServiceName  = "slv"
	slvKeyringItemKey      = "slv_secret_key"
	slvSecretIdEnvarName   = "SLV_SECRET_ID"
)

var ErrPassphraseNotSet = errors.New(slvPassphraseEnvarName + " environment variable is not set")
var ErrCreatingUserDataDir = errors.New("error in creating user data directory")
var ErrFileExistsInPathOfUserDataDir = errors.New("file exists in path of user data directory")
var ErrSaltNotFound = errors.New("salt not found")
var UserAlreadyRegistered = errors.New("user already registered")

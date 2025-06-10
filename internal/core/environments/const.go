package environments

import (
	"errors"

	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
)

const (
	envDefStringAbbrev                = "EDS" // Environment Definition String
	EnvironmentKey     crypto.KeyType = 'E'
	USER               EnvType        = "user"
	SERVICE            EnvType        = "service"
	slvPrefix                         = config.AppNameUpperCase
	selfEnvFileName                   = ".self"
)

var (
	errInvalidEnvDef                 = errors.New("invalid environment definition string")
	errInvalidEnvironmentType        = errors.New("invalid environment type")
	errEnvironmentPublicKeyNotFound  = errors.New("environment public key not found")
	errManifestPathExistsAlready     = errors.New("manifest path exists already")
	errManifestNotFound              = errors.New("manifest not found")
	errWritingManifest               = errors.New("error in writing manifest")
	errRootExistsAlready             = errors.New("root environment exists already")
	errMarkingSelfEnvBindingNotFound = errors.New("error in marking environment as self - env secret binding not found")
	errMarkingSelfNonUserEnv         = errors.New("error in marking environment as self - non user environment")
	errEnvNotFound                   = errors.New("environment not found")
)

package environments

import (
	"errors"

	"savesecrets.org/slv/core/config"
	"savesecrets.org/slv/core/crypto"
)

const (
	envDefStringAbbrev                    = "EDS" // Environment Definition String
	envSecretBindingAbbrev                = "ESB" // Environment Secret Binding
	EnvironmentKey         crypto.KeyType = 'E'
	USER                   EnvType        = "user"
	SERVICE                EnvType        = "service"
	ROOT                   EnvType        = "root"
	slvPrefix                             = config.AppNameUpperCase
	selfEnvFileName                       = ".self"
)

var (
	errInvalidEnvDef                 = errors.New("invalid environment definition string")
	errInvalidEnvironmentType        = errors.New("invalid environment type")
	errEnvironmentPublicKeyNotFound  = errors.New("environment public key not found")
	errManifestPathExistsAlready     = errors.New("manifest path exists already")
	errManifestNotFound              = errors.New("manifest not found")
	errWritingManifest               = errors.New("error in writing manifest")
	errRootExistsAlready             = errors.New("root environment exists already")
	errProviderUnknown               = errors.New("unknown provider")
	errInvalidEnvSecretBindingFormat = errors.New("invalid environment secret binding format")
	errEnvSecretBindingUnspecified   = errors.New("environment secret binding unspecified")
	errProviderRegisteredAlready     = errors.New("env secret provider registered already")
	errMarkingSelfEnvBindingNotFound = errors.New("error in marking environment as self - env secret binding not found")
	errMarkingSelfNonUserEnv         = errors.New("error in marking environment as self - non user environment")
)

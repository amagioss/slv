package environments

import (
	"errors"

	"github.com/shibme/slv/core/crypto"
)

const (
	envDataStringAbbrev                        = "EDS" // Environment Data String
	providerAccessBindingAbbrev                = "PAB" // Provider Access Binding
	EnvironmentKey              crypto.KeyType = 'E'
	USER                        EnvType        = "user"
	SERVICE                     EnvType        = "service"
	ROOT                        EnvType        = "root"
)

var (
	ErrInvalidEnvData                     = errors.New("invalid environment data string")
	ErrInvalidEnvironmentType             = errors.New("invalid environment type")
	ErrEnvironmentNotFound                = errors.New("no such environment exists")
	ErrEnvironmentPublicKeyNotFound       = errors.New("environment public key not found")
	ErrManifestPathExistsAlready          = errors.New("manifest path exists already")
	ErrManifestNotFound                   = errors.New("manifest not found")
	ErrWritingManifest                    = errors.New("error in writing manifest")
	ErrRootExistsAlready                  = errors.New("root environment exists already")
	ErrProviderUnknown                    = errors.New("unknown provider")
	ErrInvalidProviderAccessBindingFormat = errors.New("invalid provider access binding format")
	ErrProviderAccessBindingUnspecified   = errors.New("provider access binding unspecified")
)

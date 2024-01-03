package environments

import (
	"errors"

	"github.com/amagimedia/slv/core/crypto"
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
	errInvalidEnvData                     = errors.New("invalid environment data string")
	errInvalidEnvironmentType             = errors.New("invalid environment type")
	errEnvironmentNotFound                = errors.New("no such environment exists")
	errEnvironmentPublicKeyNotFound       = errors.New("environment public key not found")
	errManifestPathExistsAlready          = errors.New("manifest path exists already")
	errManifestNotFound                   = errors.New("manifest not found")
	errWritingManifest                    = errors.New("error in writing manifest")
	errRootExistsAlready                  = errors.New("root environment exists already")
	errProviderUnknown                    = errors.New("unknown provider")
	errInvalidProviderAccessBindingFormat = errors.New("invalid provider access binding format")
	errProviderAccessBindingUnspecified   = errors.New("provider access binding unspecified")
	errProviderRegisteredAlready          = errors.New("provider registered already")
)

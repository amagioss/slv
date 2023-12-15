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

var ErrInvalidEnvData = errors.New("invalid environment data string")
var ErrInvalidEnvironmentType = errors.New("invalid environment type")
var ErrEnvironmentNotFound = errors.New("no such environment exists")
var ErrManifestPathExistsAlready = errors.New("manifest path exists already")
var ErrManifestNotFound = errors.New("manifest not found")
var ErrWritingManifest = errors.New("error in writing manifest")
var ErrRootExistsAlready = errors.New("root environment exists already")

var ErrProviderUnknown = errors.New("unknown provider")
var ErrInvalidProviderAccessBindingFormat = errors.New("invalid provider access binding format")
var ErrProviderAccessBindingUnspecified = errors.New("provider access binding unspecified")

package environments

import (
	"errors"

	"github.com/shibme/slv/core/crypto"
)

const (
	envDataStringAbbrev                   = "EDS" // Environment Data String
	envAccessBindingAbbrev                = "EAB" // Environment Access Binding
	EnvironmentKey         crypto.KeyType = 'E'
	USER                   EnvType        = "user"
	SERVICE                EnvType        = "service"
	ROOT                   EnvType        = "root"
)

var ErrInvalidEnvData = errors.New("invalid environment data string")
var ErrInvalidEnvironmentType = errors.New("invalid environment type")
var ErrEnvironmentNotFound = errors.New("no such environment exists")
var ErrManifestPathExistsAlready = errors.New("manifest path exists already")
var ErrManifestNotFound = errors.New("manifest not found")
var ErrWritingManifest = errors.New("error in writing manifest")
var ErrRootExistsAlready = errors.New("root environment exists already")
var ErrInvalidEnvAccessBinding = errors.New("invalid access binding")
var ErrInvalidRSAPublicKey = errors.New("invalid RSA public key")

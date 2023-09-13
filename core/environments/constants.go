package environments

import (
	"errors"

	"github.com/shibme/slv/core/crypto"
)

const (
	userNonceFile                 = ".user_nonce"
	envDefAbbrev                  = "EDS" // Environment Definition Serialized
	EnvironmentKey crypto.KeyType = 'E'
	RootKey        crypto.KeyType = 'R'
	USER           EnvType        = "user"
	SERVICE        EnvType        = "service"
	ROOT           EnvType        = "root"
)

var ErrInvalidEnvDef = errors.New("invalid environment definition")
var ErrInvalidEnvironmentType = errors.New("invalid environment type")
var ErrEnvironmentNotFound = errors.New("no such environment exists")
var ErrManifestPathExistsAlready = errors.New("manifest path exists already")
var ErrManifestNotFound = errors.New("manifest not found")
var ErrWritingManifest = errors.New("error in writing manifest")
var ErrManifestRootExistsAlready = errors.New("error root exists already in manifest")

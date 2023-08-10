package environment

import (
	"errors"

	"github.com/shibme/slv/crypto"
)

const (
	envDefPrefix                  = "SLV_ED_" // Environment Definition
	EnvironmentKey crypto.KeyType = 'E'
	RootKey        crypto.KeyType = 'R'
	USER           EnvType        = "user"
	SERVICE        EnvType        = "service"
	ROOT           EnvType        = "root"
)

var ErrInvalidEnvironmentType = errors.New("invalid environment type")
var ErrEnvironmentNotFound = errors.New("no such environment exists")

var ErrManifestExistsAlready = errors.New("manifest exists already")
var ErrManifestNotFound = errors.New("manifest not found")
var ErrWritingManifest = errors.New("error in writing manifest")
var ErrManifestRootExistsAlready = errors.New("error root exists already in manifest")

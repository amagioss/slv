package groups

import (
	"errors"

	"github.com/shibme/slv/core/crypto"
)

const (
	GroupKey crypto.KeyType = 'G'
)

var ErrGroupNotFound = errors.New("no such environment exists")
var ErrManifestExistsAlready = errors.New("manifest exists already")
var ErrManifestNotFound = errors.New("manifest not found")
var ErrWritingManifest = errors.New("error in writing manifest")

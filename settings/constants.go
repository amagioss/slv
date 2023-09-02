package settings

import (
	"errors"
)

var ErrManifestExistsAlready = errors.New("manifest exists already")
var ErrManifestNotFound = errors.New("manifest not found")
var ErrWritingManifest = errors.New("error in writing manifest")

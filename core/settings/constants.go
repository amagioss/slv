package settings

import (
	"errors"
)

const (
	defaultSyncInterval = 86400
)

var ErrManifestPathExistsAlready = errors.New("manifest path exists already")
var ErrManifestNotFound = errors.New("manifest not found")
var ErrWritingManifest = errors.New("error in writing manifest")

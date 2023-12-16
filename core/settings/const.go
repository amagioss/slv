package settings

import (
	"errors"
)

const (
	defaultSyncInterval = 86400
)

var (
	ErrManifestPathExistsAlready = errors.New("manifest path exists already")
	ErrManifestNotFound          = errors.New("manifest not found")
	ErrWritingManifest           = errors.New("error in writing manifest")
)

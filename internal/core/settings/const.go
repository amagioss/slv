package settings

import (
	"errors"
)

var (
	errManifestPathExistsAlready = errors.New("manifest path exists already")
	errManifestNotFound          = errors.New("manifest not found")
	// errWritingManifest           = errors.New("error in writing manifest")
)

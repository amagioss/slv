package settings

import (
	"errors"
)

var (
	errManifestPathExistsAlready = errors.New("manifest path exists already")
	errManifestNotFound          = errors.New("manifest not found")
)

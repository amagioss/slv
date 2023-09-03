package config

import (
	"errors"
)

// Errors and constants used by config & manifests functions
const (
	// Config Manager constants
	configManagerDirName             = "configs"
	configManagerPreferencesFileName = ".preferences.yml"

	// Config constants
	configFileName            = ".config.yml"
	configDataDirName         = "data"
	settingsFileName          = "settings.yml"
	environmentConfigFileName = "environments.yml"
	groupConfigFileName       = "groups.yml"

	// Settings constants
	defaultSyncInterval = 86400
)

var ErrConfigPathExistsAlready = errors.New("config path exists already")
var ErrCreatingConfigDir = errors.New("error creating config dir")
var ErrWritingManifest = errors.New("error in writing manifest")
var ErrConfigPathDoesNotExist = errors.New("config path does not exist")

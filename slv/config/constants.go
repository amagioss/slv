package config

import (
	"errors"
)

// Errors and constants used by config & manifests functions
const (
	// Config Manager constants
	configManagerDirName             = "configs"
	configManagerPreferencesFileName = ".preferences.slv"

	// Config constants
	configFileName            = ".config.slv"
	configDataDirName         = "data"
	settingsFileName          = "settings.slv"
	environmentConfigFileName = "environments.slv"
	groupConfigFileName       = "groups.slv"

	// Settings constants
	defaultSyncInterval = 86400
)

var ErrConfigManagerInitialization = errors.New("error in initializing config manager")
var ErrSavingConfigManagerPreferences = errors.New("error in saving config manager preferences")
var ErrOpeningConfigManagerDir = errors.New("error in opening manifest config directory")
var ErrNoCurrentConfigFound = errors.New("current manifest not set")
var ErrManifestExistsAlready = errors.New("manifest exists already")
var ErrConfigNotFound = errors.New("manifest not found")
var ErrConfigInitialization = errors.New("unable to initialize config")
var ErrProcessingSettings = errors.New("error in processing settings")
var ErrProcessingEnvironmentsManifest = errors.New("error in processing environments manifest")
var ErrProcessingGroupsManifest = errors.New("error in processing groups manifest")
var ErrEnvironmentNotFound = errors.New("no such environment exists")

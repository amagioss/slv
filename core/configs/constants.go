package configs

import (
	"errors"
)

// Errors and constants used by config & manifests functions
const (
	// Config Manager constants
	configsDirName        = "configs"
	defaultConfigFileName = ".default"

	// Config constants
	configFileName            = ".config.yml"
	configDataDirName         = "data"
	settingsFileName          = "settings.yml"
	environmentConfigFileName = "environments.yml"
	groupConfigFileName       = "groups.yml"
)

var ErrConfigPathExistsAlready = errors.New("config path exists already")
var ErrCreatingConfigDir = errors.New("error creating config dir")
var ErrWritingManifest = errors.New("error in writing manifest")
var ErrConfigPathDoesNotExist = errors.New("config path does not exist")
var ErrCreatingConfigCollectionDir = errors.New("error creating config collection dir in app data dir")
var ErrInitializingConfigManagementDir = errors.New("error initializing config management dir")
var ErrOpeningConfigManagementDir = errors.New("error opening config management dir")
var ErrConfigNotFound = errors.New("config not found")
var ErrConfigExistsAlready = errors.New("config exists already")
var ErrInvalidConfigName = errors.New("invalid config name")
var ErrNoDefaultConfigFound = errors.New("no default config found")
var ErrSettingDefaultConfig = errors.New("error setting default config")

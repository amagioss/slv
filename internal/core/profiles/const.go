package profiles

import (
	"errors"
	"time"
)

// Errors and constants used by profiles
const (
	// Profile Manager constants
	profilesDirName        = "profiles"
	currentProfileFileName = ".current"

	// Profile constants
	profileConfigFileName = ".profile.yaml"
	profileDataDirName    = "data"
	defaultSyncInterval   = time.Hour

	defaultEnvManifestFileName = "environments.yaml"
	defaultSettingsFileName    = "settings.yaml"
)

var (
	defaultRemoteRegistered = false
	envManifestFileNames    = []string{defaultEnvManifestFileName, "environments.yml"}
	settingsFileNames       = []string{defaultSettingsFileName, "settings.yml"}

	errProfilePathExistsAlready         = errors.New("profile path exists already")
	errProfilePathDoesNotExist          = errors.New("profile path does not exist")
	errCreatingProfilesHomeDir          = errors.New("error creating profiles home dir inside app data dir")
	errInitializingProfileManagementDir = errors.New("error initializing profile management dir")
	errOpeningProfileManagementDir      = errors.New("error opening profile management dir")
	errProfileNotFound                  = errors.New("profile not found")
	errProfileExistsAlready             = errors.New("profile exists already")
	errInvalidProfileName               = errors.New("invalid profile name")
	errNoCurrentProfileSet              = errors.New("current profile not set")
	errSettingCurrentProfile            = errors.New("error setting current profile")
	errDeletingCurrentProfile           = errors.New("error deleting current profile")
	errRemoteTypeDoesNotExist           = errors.New("remote type does not exist")
	errRemoteSetupNotImplemented        = errors.New("remote setup not implemented")
	errRemotePullNotImplemented         = errors.New("remote pull not implemented")
	errRemotePushNotSupported           = errors.New("remote push not supported")
)

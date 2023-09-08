package profiles

import (
	"errors"
)

// Errors and constants used by profiles
const (
	// Profile Manager constants
	profilesDirName        = "profiles"
	defaultProfileFileName = ".default"

	// Profile constants
	profileFileName             = ".profile.yml"
	profileDataDirName          = "data"
	profileSettingsFileName     = "settings.yml"
	profileEnvironmentsFileName = "environments.yml"
	profileGroupsFileName       = "groups.yml"
)

var ErrProfilePathExistsAlready = errors.New("profile path exists already")
var ErrCreatingProfileDir = errors.New("error creating profile dir")
var ErrWritingManifest = errors.New("error in writing manifest")
var ErrProfilePathDoesNotExist = errors.New("profile path does not exist")
var ErrCreatingProfileCollectionDir = errors.New("error creating profile collection dir in app data dir")
var ErrInitializingProfileManagementDir = errors.New("error initializing profile management dir")
var ErrOpeningProfileManagementDir = errors.New("error opening profile management dir")
var ErrProfileNotFound = errors.New("profile not found")
var ErrProfileExistsAlready = errors.New("profile exists already")
var ErrInvalidProfileName = errors.New("invalid profile name")
var ErrNoDefaultProfileFound = errors.New("no default profile found")
var ErrSettingDefaultProfile = errors.New("error setting default profile")

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

var (
	ErrProfilePathExistsAlready         = errors.New("profile path exists already")
	ErrCreatingProfileDir               = errors.New("error creating profile dir")
	ErrWritingManifest                  = errors.New("error in writing manifest")
	ErrProfilePathDoesNotExist          = errors.New("profile path does not exist")
	ErrCreatingProfileCollectionDir     = errors.New("error creating profile collection dir in app data dir")
	ErrInitializingProfileManagementDir = errors.New("error initializing profile management dir")
	ErrOpeningProfileManagementDir      = errors.New("error opening profile management dir")
	ErrProfileNotFound                  = errors.New("profile not found")
	ErrProfileExistsAlready             = errors.New("profile exists already")
	ErrInvalidProfileName               = errors.New("invalid profile name")
	ErrNoDefaultProfileFound            = errors.New("no default profile found")
	ErrSettingDefaultProfile            = errors.New("error setting default profile")
)

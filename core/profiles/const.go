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
	errProfilePathExistsAlready         = errors.New("profile path exists already")
	errCreatingProfileDir               = errors.New("error creating profile dir")
	errWritingManifest                  = errors.New("error in writing manifest")
	errProfilePathDoesNotExist          = errors.New("profile path does not exist")
	errCreatingProfileCollectionDir     = errors.New("error creating profile collection dir in app data dir")
	errInitializingProfileManagementDir = errors.New("error initializing profile management dir")
	errOpeningProfileManagementDir      = errors.New("error opening profile management dir")
	errProfileNotFound                  = errors.New("profile not found")
	errProfileExistsAlready             = errors.New("profile exists already")
	errInvalidProfileName               = errors.New("invalid profile name")
	errNoDefaultProfileFound            = errors.New("no default profile found")
	errSettingDefaultProfile            = errors.New("error setting default profile")
)

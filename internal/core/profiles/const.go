package profiles

import (
	"errors"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Errors and constants used by profiles
const (
	// Profile Manager constants
	profilesDirName        = "profiles"
	currentProfileFileName = ".current"

	// Profile constants
	profileDataDirName          = "data"
	profileSettingsFileName     = "settings.yml"
	profileEnvironmentsFileName = "environments.yml"
	profileGroupsFileName       = "groups.yml"
	profileGitSyncInterval      = time.Hour
)

var (
	gitHttpAuth          *http.BasicAuth
	gitHttpAuthProcessed = false

	errProfilePathExistsAlready         = errors.New("profile path exists already")
	errCreatingProfileDir               = errors.New("error creating profile dir")
	errProfilePathDoesNotExist          = errors.New("profile path does not exist")
	errCreatingProfileCollectionDir     = errors.New("error creating profile collection dir in app data dir")
	errInitializingProfileManagementDir = errors.New("error initializing profile management dir")
	errOpeningProfileManagementDir      = errors.New("error opening profile management dir")
	errProfileNotFound                  = errors.New("profile not found")
	errProfileExistsAlready             = errors.New("profile exists already")
	errInvalidProfileName               = errors.New("invalid profile name")
	errNoCurrentProfileSet              = errors.New("current profile not set")
	errSettingCurrentProfile            = errors.New("error setting current profile")
	errDeletingCurrentProfile           = errors.New("error deleting current profile")
	errProfileNotGitRepository          = errors.New("profile is not a git repository")
	errProfileGitPullMarking            = errors.New("error marking profile as pulled")
	errChangesNotAllowedInGitProfile    = errors.New("changes not allowed since the current profile is git based")
)

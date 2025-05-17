package profiles

import (
	"os"
	"path/filepath"
	"time"

	"slv.sh/slv/internal/core/commons"
	"slv.sh/slv/internal/core/config"
)

type profileManager struct {
	dir                string
	profileList        map[string]struct{}
	currentProfileName *string
	currentProfile     *Profile
}

var profileMgr *profileManager = nil
var profileMap map[string]*Profile = make(map[string]*Profile)

func initProfileManager() error {
	if profileMgr == nil {
		RegisterDefaultRemotes()
		var manager profileManager
		manager.dir = filepath.Join(config.GetAppDataDir(), profilesDirName)
		profileManagerDirInfo, err := os.Stat(manager.dir)
		if err != nil {
			err = os.MkdirAll(manager.dir, 0755)
			if err != nil {
				return errCreatingProfilesHomeDir
			}
		} else if !profileManagerDirInfo.IsDir() {
			return errInitializingProfileManagementDir
		}
		profileManagerDir, err := os.Open(manager.dir)
		if err != nil {
			return errOpeningProfileManagementDir
		}
		defer profileManagerDir.Close()
		fileInfoList, err := profileManagerDir.Readdir(-1)
		if err != nil {
			return errOpeningProfileManagementDir
		}
		manager.profileList = make(map[string]struct{})
		for _, fileInfo := range fileInfoList {
			if fileInfo.IsDir() {
				manager.profileList[fileInfo.Name()] = struct{}{}
			}
		}
		profileMgr = &manager
	}
	return nil
}

func List() (profileNames []string, err error) {
	if err = initProfileManager(); err != nil {
		return nil, err
	}
	profileNames = make([]string, 0, len(profileMgr.profileList))
	for profileName := range profileMgr.profileList {
		profileNames = append(profileNames, profileName)
	}
	return
}

func Get(profileName string) (profile *Profile, err error) {
	if err = initProfileManager(); err != nil {
		return nil, err
	}
	if profile = profileMap[profileName]; profile != nil {
		return
	}
	if _, exists := profileMgr.profileList[profileName]; !exists {
		return nil, errProfileNotFound
	}
	if profile, err = getProfile(profileName, filepath.Join(profileMgr.dir, profileName)); err != nil {
		return nil, err
	}
	profile.name = profileName
	profileMap[profileName] = profile
	return
}

func New(profileName, remoteType string, updateInterval time.Duration, remoteConfig map[string]string) error {
	if profileName == "" {
		return errInvalidProfileName
	}
	if err := initProfileManager(); err != nil {
		return err
	}
	if _, exists := profileMgr.profileList[profileName]; exists {
		return errProfileExistsAlready
	}
	if _, err := createProfile(profileName, filepath.Join(profileMgr.dir, profileName), remoteType, updateInterval, remoteConfig); err != nil {
		return err
	}
	profileMgr.profileList[profileName] = struct{}{}
	if profileMgr.currentProfileName == nil {
		return SetCurrentProfile(profileName)
	}
	return nil
}

func SetCurrentProfile(profileName string) error {
	if profileName == "" {
		return errInvalidProfileName
	}
	if err := initProfileManager(); err != nil {
		return err
	}
	if _, exists := profileMgr.profileList[profileName]; !exists {
		return errProfileNotFound
	}
	if commons.WriteToFile(filepath.Join(profileMgr.dir, currentProfileFileName), []byte(profileName)) != nil {
		return errSettingCurrentProfile
	}
	profileMgr.currentProfileName = &profileName
	profileMgr.currentProfile = nil
	return nil
}

func GetCurrentProfileName() (string, error) {
	if err := initProfileManager(); err != nil {
		return "", err
	}
	if profileMgr.currentProfileName != nil {
		return *profileMgr.currentProfileName, nil
	}
	fileContent, err := os.ReadFile(filepath.Join(profileMgr.dir, currentProfileFileName))
	currentProfileName := string(fileContent)
	if err != nil {
		return "", errNoCurrentProfileSet
	}
	profileMgr.currentProfileName = &currentProfileName
	return *profileMgr.currentProfileName, nil
}

func GetCurrentProfile() (profile *Profile, err error) {
	if err = initProfileManager(); err != nil {
		return nil, err
	}
	if profileMgr.currentProfile == nil {
		currentProfileName, err := GetCurrentProfileName()
		if err != nil {
			return nil, err
		}
		if profileMgr.currentProfile, err = Get(currentProfileName); err != nil {
			return nil, err
		}
	}
	return profileMgr.currentProfile, nil
}

func Delete(profileName string) error {
	if profileName == "" {
		return errInvalidProfileName
	}
	if err := initProfileManager(); err != nil {
		return err
	}
	if _, exists := profileMgr.profileList[profileName]; !exists {
		return errProfileNotFound
	}
	if profileMgr.currentProfileName != nil && *profileMgr.currentProfileName == profileName {
		return errDeletingCurrentProfile
	}
	delete(profileMgr.profileList, profileName)
	delete(profileMap, profileName)
	return os.RemoveAll(filepath.Join(profileMgr.dir, profileName))
}

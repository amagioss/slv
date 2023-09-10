package profiles

import (
	"os"
	"path/filepath"

	"github.com/shibme/slv/core/commons"
)

type profileManager struct {
	dir                string
	profileList        map[string]struct{}
	defaultFile        *string
	defaultProfileName *string
	defaultProfile     *Profile
}

var profileMgr *profileManager = nil
var profileMap map[string]*Profile = make(map[string]*Profile)

func initProfileManager() error {
	if profileMgr != nil {
		return nil
	}
	var manager profileManager
	manager.dir = filepath.Join(commons.AppDataDir(), profilesDirName)

	// Creates the profiles directory if it doesn't exist
	profileManagerDirInfo, err := os.Stat(manager.dir)
	if err != nil {
		err = os.MkdirAll(manager.dir, 0755)
		if err != nil {
			return ErrCreatingProfileCollectionDir
		}
	} else if !profileManagerDirInfo.IsDir() {
		return ErrInitializingProfileManagementDir
	}

	// Populating list of all profiles
	profileManagerDir, err := os.Open(manager.dir)
	if err != nil {
		return ErrOpeningProfileManagementDir
	}
	defer profileManagerDir.Close()
	fileInfoList, err := profileManagerDir.Readdir(-1)
	if err != nil {
		return ErrOpeningProfileManagementDir
	}
	manager.profileList = make(map[string]struct{})
	for _, fileInfo := range fileInfoList {
		if fileInfo.IsDir() {
			if f, err := os.Stat(filepath.Join(manager.dir, fileInfo.Name(), profileFileName)); err == nil && !f.IsDir() {
				manager.profileList[fileInfo.Name()] = struct{}{}
			}
		}
	}
	profileMgr = &manager
	GetDefaultProfileName()
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

func GetProfile(profileName string) (profile *Profile, err error) {
	if err = initProfileManager(); err != nil {
		return nil, err
	}
	if profile = profileMap[profileName]; profile != nil {
		return
	}
	if _, exists := profileMgr.profileList[profileName]; !exists {
		return nil, ErrProfileNotFound
	}
	if profile, err = getProfileForPath(filepath.Join(profileMgr.dir, profileName)); err != nil {
		return nil, err
	}
	profileMap[profileName] = profile
	return
}

func New(profileName string) error {
	if profileName == "" {
		return ErrInvalidProfileName
	}
	if err := initProfileManager(); err != nil {
		return err
	}
	if _, exists := profileMgr.profileList[profileName]; exists {
		return ErrProfileExistsAlready
	}
	if _, err := newProfileForPath(filepath.Join(profileMgr.dir, profileName)); err != nil {
		return err
	}
	profileMgr.profileList[profileName] = struct{}{}
	if profileMgr.defaultProfileName == nil {
		return SetDefault(profileName)
	}
	return nil
}

func SetDefault(profileName string) error {
	if profileName == "" {
		return ErrInvalidProfileName
	}
	if err := initProfileManager(); err != nil {
		return err
	}
	if _, exists := profileMgr.profileList[profileName]; !exists {
		return ErrProfileNotFound
	}
	if profileMgr.defaultFile == nil {
		defaultInfoFile := filepath.Join(profileMgr.dir, defaultProfileFileName)
		profileMgr.defaultFile = &defaultInfoFile
	}
	err := os.WriteFile(*profileMgr.defaultFile, []byte(profileName), 0644)
	if err != nil {
		return ErrSettingDefaultProfile
	}
	profileMgr.defaultProfile = nil
	return nil
}

func GetDefaultProfileName() (string, error) {
	if err := initProfileManager(); err != nil {
		return "", err
	}
	if profileMgr.defaultProfileName != nil {
		return *profileMgr.defaultProfileName, nil
	}
	if profileMgr.defaultFile == nil {
		defaultInfoFile := filepath.Join(profileMgr.dir, defaultProfileFileName)
		profileMgr.defaultFile = &defaultInfoFile
	}
	bytes, err := os.ReadFile(*profileMgr.defaultFile)
	if err != nil {
		return "", ErrNoDefaultProfileFound
	}
	defaultProfileName := string(bytes)
	profileMgr.defaultProfileName = &defaultProfileName
	return *profileMgr.defaultProfileName, nil
}

func GetDefaultProfile() (profile *Profile, err error) {
	if err = initProfileManager(); err != nil {
		return nil, err
	}
	if profileMgr.defaultProfile == nil {
		defaultProfileName, err := GetDefaultProfileName()
		if err != nil {
			return nil, err
		}
		if profileMgr.defaultProfile, err = GetProfile(defaultProfileName); err != nil {
			return nil, err
		}
	}
	return profileMgr.defaultProfile, nil
}
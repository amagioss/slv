package profiles

import (
	"os"
	"path/filepath"

	"oss.amagi.com/slv/core/commons"
	"oss.amagi.com/slv/core/config"
)

type profileManager struct {
	dir                string
	profileList        map[string]struct{}
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
	manager.dir = filepath.Join(config.GetAppDataDir(), profilesDirName)
	profileManagerDirInfo, err := os.Stat(manager.dir)
	if err != nil {
		err = os.MkdirAll(manager.dir, 0755)
		if err != nil {
			return errCreatingProfileCollectionDir
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
	if profile, err = getProfileForPath(filepath.Join(profileMgr.dir, profileName)); err != nil {
		return nil, err
	}
	profile.name = commons.StringPtr(profileName)
	profileMap[profileName] = profile
	return
}

func New(profileName, gitURI, gitBranch string) error {
	if profileName == "" {
		return errInvalidProfileName
	}
	if err := initProfileManager(); err != nil {
		return err
	}
	if _, exists := profileMgr.profileList[profileName]; exists {
		return errProfileExistsAlready
	}
	if _, err := newProfile(filepath.Join(profileMgr.dir, profileName), gitURI, gitBranch); err != nil {
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
		return errInvalidProfileName
	}
	if err := initProfileManager(); err != nil {
		return err
	}
	if _, exists := profileMgr.profileList[profileName]; !exists {
		return errProfileNotFound
	}
	if commons.WriteToFile(filepath.Join(profileMgr.dir, defaultProfileFileName), []byte(profileName)) != nil {
		return errSettingDefaultProfile
	}
	profileMgr.defaultProfileName = &profileName
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
	fileContent, err := os.ReadFile(filepath.Join(profileMgr.dir, defaultProfileFileName))
	defaultProfileName := string(fileContent)
	if err != nil {
		return "", errNoDefaultProfileFound
	}
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
		if profileMgr.defaultProfile, err = Get(defaultProfileName); err != nil {
			return nil, err
		}
	}
	return profileMgr.defaultProfile, nil
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
	if profileMgr.defaultProfileName != nil && *profileMgr.defaultProfileName == profileName {
		return errDeletingDefaultProfile
	}
	delete(profileMgr.profileList, profileName)
	delete(profileMap, profileName)
	return os.RemoveAll(filepath.Join(profileMgr.dir, profileName))
}

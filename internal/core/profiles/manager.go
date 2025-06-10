package profiles

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"slv.sh/slv/internal/core/commons"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/keystore"
	"xipher.org/xipher"
)

type profileManagerConfig struct {
	file              string
	activeProfile     *Profile
	ActiveProfileName string `json:"active" yaml:"active"`
}

func (pmc *profileManagerConfig) write() error {
	return commons.WriteToYAML(pmc.file, pmc)
}

func (pmc *profileManagerConfig) getActiveProfile() (*Profile, error) {
	if pmc.activeProfile == nil {
		if pmc.ActiveProfileName == "" {
			return nil, errNoActiveProfileSet
		}
		profile, err := Get(pmc.ActiveProfileName)
		if err != nil {
			return nil, err
		}
		pmc.activeProfile = profile
	}
	return pmc.activeProfile, nil
}

type profileManager struct {
	dir         string
	profileList map[string]struct{}
	config      *profileManagerConfig
}

func (pm *profileManager) getConfig() (*profileManagerConfig, error) {
	if pm.config == nil {
		pmc := &profileManagerConfig{}
		pmcFile := filepath.Join(pm.dir, profileMgrConfigFileName)
		if commons.FileExists(pmcFile) {
			if err := commons.ReadFromYAML(pmcFile, pmc); err != nil {
				return nil, fmt.Errorf("error reading profile manager config file: %w", err)
			}
		} else {
			if err := commons.WriteToYAML(pmcFile, pmc); err != nil {
				return nil, fmt.Errorf("error creating profile manager config file: %w", err)
			}
		}
		pmc.file = pmcFile
		pm.config = pmc
	}
	return pm.config, nil
}

func initProfileManager() error {
	if profileMgr == nil {
		profileMgrMutex.Lock()
		defer profileMgrMutex.Unlock()
		if profileMgr == nil {
			registerDefaultRemotes()
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
					if isValidProfile(filepath.Join(manager.dir, fileInfo.Name())) {
						manager.profileList[fileInfo.Name()] = struct{}{}
					} else {
						if err := os.RemoveAll(filepath.Join(manager.dir, fileInfo.Name())); err != nil {
							return fmt.Errorf("error removing invalid profile directory: %w", err)
						}
					}
				}
			}
			if _, err = manager.getConfig(); err != nil {
				return err
			}
			profileMgr = &manager
		}
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

func New(profileName, remoteType string, readOnly bool, updateInterval time.Duration, remoteConfig map[string]string) error {
	if profileName == "" {
		return errInvalidProfileName
	}
	if err := initProfileManager(); err != nil {
		return err
	}
	if _, exists := profileMgr.profileList[profileName]; exists {
		return errProfileExistsAlready
	}
	if _, err := createProfile(profileName, filepath.Join(profileMgr.dir, profileName), remoteType, readOnly, updateInterval, remoteConfig); err != nil {
		return err
	}
	profileMgr.profileList[profileName] = struct{}{}
	pmc, err := profileMgr.getConfig()
	if err != nil {
		return err
	}
	if pmc.ActiveProfileName == "" {
		return SetActiveProfile(profileName)
	}
	return nil
}

func SetActiveProfile(profileName string) error {
	if profileName == "" {
		return errInvalidProfileName
	}
	if err := initProfileManager(); err != nil {
		return err
	}
	if _, exists := profileMgr.profileList[profileName]; !exists {
		return errProfileNotFound
	}
	pmc, err := profileMgr.getConfig()
	if err != nil {
		return err
	}
	pmc.ActiveProfileName = profileName
	pmc.activeProfile = nil
	if pmc.write() != nil {
		return errSettingActiveProfile
	}
	return nil
}

func GetActiveProfileName() (string, error) {
	if err := initProfileManager(); err != nil {
		return "", err
	}
	pmc, err := profileMgr.getConfig()
	if err != nil {
		return "", err
	}
	if pmc.ActiveProfileName != "" {
		return pmc.ActiveProfileName, nil
	}
	return "", errNoActiveProfileSet
}

func GetActiveProfile() (profile *Profile, err error) {
	if err = initProfileManager(); err != nil {
		return nil, err
	}
	pmc, err := profileMgr.getConfig()
	if err != nil {
		return nil, err
	}
	return pmc.getActiveProfile()
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
	pmc, err := profileMgr.getConfig()
	if err != nil {
		return err
	}
	if pmc.ActiveProfileName == profileName {
		return errDeletingActiveProfile
	}
	delete(profileMgr.profileList, profileName)
	delete(profileMap, profileName)
	return os.RemoveAll(filepath.Join(profileMgr.dir, profileName))
}

func getCryptoKey() (*xipher.SecretKey, error) {
	if profileSK == nil {
		if skBytes, err := keystore.Get([]byte(profileCryptoKeyName), true); err == keystore.ErrNotFound {
			if profileSK, err = xipher.NewSecretKey(); err != nil {
				return nil, fmt.Errorf("error creating new secret key: %w", err)
			} else if skBytes, err = profileSK.Bytes(); err != nil {
				return nil, fmt.Errorf("error getting secret key bytes: %w", err)
			} else if err = keystore.Put([]byte(profileCryptoKeyName), skBytes, true); err != nil {
				return nil, fmt.Errorf("error saving secret key to keystore: %w", err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("error getting secret key from keystore: %w", err)
		} else if profileSK, err = xipher.ParseSecretKey(skBytes); err != nil {
			return nil, fmt.Errorf("error creating secret key from bytes: %w", err)
		}
	}
	return profileSK, nil
}

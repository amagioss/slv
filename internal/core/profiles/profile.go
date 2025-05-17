package profiles

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"slv.sh/slv/internal/core/commons"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/settings"
)

type profileConfig struct {
	RemoteType     string            `json:"remoteType" yaml:"remoteType"`
	UpdatedAt      time.Time         `json:"updatedAt" yaml:"updatedAt"`
	UpdateInterval time.Duration     `json:"updateInterval" yaml:"updateInterval"`
	Config         map[string]string `json:"config" yaml:"config"`
	file           string
}

func (pc *profileConfig) write() error {
	pc.UpdatedAt = time.Now()
	return commons.WriteToYAML(pc.file, "", pc)
}

type Profile struct {
	name          string
	dir           string
	dataDir       string
	profileConfig *profileConfig
	settings      *settings.Settings
	envManifest   *environments.EnvManifest
	remote        *remote
}

func (profile *Profile) Name() string {
	return profile.name
}

func (profile *Profile) getConfig() (*profileConfig, error) {
	if profile.profileConfig == nil {
		profileConfig := &profileConfig{}
		profileConfigFile := filepath.Join(profile.dir, profileConfigFileName)
		if err := commons.ReadFromYAML(profileConfigFile, profileConfig); err != nil {
			return nil, err
		}
		profileConfig.file = profileConfigFile
		profile.profileConfig = profileConfig
		profile.profileConfig.file = filepath.Join(profile.dir, profileConfigFileName)
	}
	return profile.profileConfig, nil
}

func (profile *Profile) getRemote() (*remote, error) {
	if profile.remote == nil {
		profileConfig, err := profile.getConfig()
		if err != nil {
			return nil, err
		}
		profile.remote = remotes[profileConfig.RemoteType]
		if profile.remote == nil {
			return nil, errRemoteTypeDoesNotExist
		}
	}
	return profile.remote, nil
}

func (profile *Profile) pull(setup bool) error {
	remote, err := profile.getRemote()
	if err != nil {
		return err
	}
	if setup {
		err = (remote.setup)(profile.dataDir, profile.profileConfig.Config)
	} else {
		err = (remote.pull)(profile.dataDir, profile.profileConfig.Config)
	}
	if err == nil {
		profile.envManifest = nil
		profile.settings = nil
		profile.profileConfig.UpdatedAt = time.Now()
		err = profile.profileConfig.write()
	}
	return err
}

func (profile *Profile) Pull() error {
	return profile.pull(false)
}

func (profile *Profile) pullOnDue() error {
	if time.Since(profile.profileConfig.UpdatedAt) < profile.profileConfig.UpdateInterval {
		return nil
	}
	return profile.Pull()
}

func (profile *Profile) IsPushSupported() bool {
	remote, err := profile.getRemote()
	if err != nil {
		return false
	}
	return remote.push != nil
}

func (profile *Profile) Push(note string) (err error) {
	if !profile.IsPushSupported() {
		return errRemotePushNotSupported
	}
	if err = (profile.remote.push)(profile.dataDir, profile.profileConfig.Config, note); err == nil {
		profile.profileConfig.UpdatedAt = time.Now()
		err = profile.profileConfig.write()
	}
	return
}

func createProfile(name, dir, remoteType string, updateInterval time.Duration, remoteConfig map[string]string) (profile *Profile, err error) {
	if commons.DirExists(dir) {
		return nil, errProfilePathExistsAlready
	}
	remote := remotes[remoteType]
	if remote == nil {
		return nil, errRemoteTypeDoesNotExist
	}
	if remote.setup == nil {
		return nil, errRemoteSetupNotImplemented
	}
	if remote.pull == nil {
		return nil, errRemotePullNotImplemented
	}
	if err = os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("error creating profile dir: %w", err)
	}
	profileDataDir := filepath.Join(dir, profileDataDirName)
	if err = (remote.setup)(profileDataDir, remoteConfig); err == nil {
		if updateInterval < 0 {
			updateInterval = defaultSyncInterval
		}
		profile = &Profile{
			name:    name,
			dir:     dir,
			dataDir: profileDataDir,
			profileConfig: &profileConfig{
				RemoteType:     remoteType,
				UpdateInterval: updateInterval,
				Config:         remoteConfig,
				file:           filepath.Join(dir, profileConfigFileName),
			},
		}
		err = profile.profileConfig.write()
	}
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	return profile, nil
}

func getProfile(name, dir string) (profile *Profile, err error) {
	if !commons.DirExists(dir) {
		return nil, errProfilePathDoesNotExist
	}
	profile = &Profile{
		name:    name,
		dir:     dir,
		dataDir: filepath.Join(dir, profileDataDirName),
	}
	if _, err = profile.getConfig(); err != nil {
		return nil, err
	}
	if err = profile.pullOnDue(); err != nil {
		return nil, err
	}
	return
}

func getAvailableFilePath(basePath string, fileNames []string) string {
	for _, fileName := range fileNames {
		if commons.FileExists(filepath.Join(basePath, fileName)) {
			return filepath.Join(basePath, fileName)
		}
	}
	return ""
}

func (profile *Profile) GetSettings() (*settings.Settings, error) {
	var err error
	if profile.settings == nil {
		if settingsFile := getAvailableFilePath(profile.dataDir, settingsFileNames); settingsFile != "" {
			profile.settings, err = settings.GetManifest(settingsFile)
		} else {
			profile.settings, err = settings.NewManifest(filepath.Join(profile.dataDir, defaultSettingsFileName))
		}
		if err != nil {
			return nil, err
		}
	}
	return profile.settings, nil
}

func (profile *Profile) getEnvManifest() (*environments.EnvManifest, error) {
	var err error
	if profile.envManifest == nil {
		if envManifestFile := getAvailableFilePath(profile.dataDir, envManifestFileNames); envManifestFile != "" {
			profile.envManifest, err = environments.GetManifest(envManifestFile)
		} else {
			profile.envManifest, err = environments.NewManifest(filepath.Join(profile.dataDir, defaultEnvManifestFileName))
		}
		if err != nil {
			return nil, err
		}
	}
	return profile.envManifest, nil
}

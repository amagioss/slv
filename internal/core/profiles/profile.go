package profiles

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"slv.sh/slv/internal/core/commons"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/settings"
)

type profileConfig struct {
	RemoteType   string            `json:"type" yaml:"type"`
	ReadOnly     bool              `json:"readOnly" yaml:"readOnly"`
	SyncedAt     time.Time         `json:"syncedAt" yaml:"syncedAt"`
	SyncInterval time.Duration     `json:"syncInterval" yaml:"syncInterval"`
	Config       map[string]string `json:"config" yaml:"config"`
	file         string
}

func (pc *profileConfig) decrypt() error {
	sk, err := getCryptoKey()
	if err != nil {
		return err
	}
	for _, arg := range GetRemoteTypeArgs(pc.RemoteType) {
		if value, ok := pc.Config[arg.Name()]; ok && arg.sensitive && value != "" {
			var pt []byte
			if ctBytes, err := base64.StdEncoding.DecodeString(value); err != nil {
				return err
			} else if pt, err = sk.Decrypt(ctBytes); err != nil {
				return err
			}
			pc.Config[arg.Name()] = string(pt)
		}
	}
	return nil
}

func (pc *profileConfig) write() error {
	pc.SyncedAt = time.Now()
	sk, err := getCryptoKey()
	if err != nil {
		return err
	}
	ptMap := make(map[string]string)
	for _, arg := range GetRemoteTypeArgs(pc.RemoteType) {
		if value, ok := pc.Config[arg.Name()]; ok && arg.sensitive && value != "" {
			ptMap[arg.Name()] = value
			var ct []byte
			if ct, err = sk.Encrypt([]byte(value), true, false); err != nil {
				return err
			}
			pc.Config[arg.Name()] = base64.StdEncoding.EncodeToString(ct)
		}
	}
	err = commons.WriteToYAML(pc.file, pc)
	for k, v := range ptMap {
		pc.Config[k] = v
	}
	return err
}

type Profile struct {
	name        string
	dir         string
	dataDir     string
	profConfig  *profileConfig
	settings    *settings.Settings
	envManifest *environments.EnvManifest
	remote      *remote
}

func (profile *Profile) Name() string {
	return profile.name
}

func (profile *Profile) getConfig() (*profileConfig, error) {
	if profile.profConfig == nil {
		profileConfig := &profileConfig{
			file: filepath.Join(profile.dir, profileConfigFileName),
		}
		if err := commons.ReadFromYAML(profileConfig.file, profileConfig); err != nil {
			return nil, err
		}
		if err := profileConfig.decrypt(); err != nil {
			return nil, err
		}
		profile.profConfig = profileConfig
	}
	return profile.profConfig, nil
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

func (profile *Profile) Pull() error {
	remote, err := profile.getRemote()
	if err != nil {
		return err
	}
	profConfig, err := profile.getConfig()
	if err != nil {
		return err
	}
	if commons.DirExists(profile.dataDir) {
		err = (remote.pull)(profile.dataDir, profConfig.Config)
	} else {
		err = (remote.setup)(profile.dataDir, profConfig.Config)
	}
	if err == nil {
		profile.envManifest = nil
		profile.settings = nil
		profConfig.SyncedAt = time.Now()
		err = profConfig.write()
	}
	return err
}

func (profile *Profile) pullOnDue() error {
	profConfig, err := profile.getConfig()
	if err != nil {
		return err
	}
	if commons.DirExists(profile.dataDir) && time.Since(profConfig.SyncedAt) < profConfig.SyncInterval {
		return nil
	}
	return profile.Pull()
}

func (profile *Profile) IsPushSupported() bool {
	if profConfig, err := profile.getConfig(); err == nil && profConfig.ReadOnly {
		return false
	}
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
	profConfig, err := profile.getConfig()
	if err != nil {
		return err
	}
	if err = (profile.remote.push)(profile.dataDir, profConfig.Config, note); err == nil {
		profConfig.SyncedAt = time.Now()
		err = profConfig.write()
	}
	return
}

func createProfile(name, dir, remoteType string, readOnly bool, updateInterval time.Duration, remoteConfig map[string]string) (profile *Profile, err error) {
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
			profConfig: &profileConfig{
				RemoteType:   remoteType,
				ReadOnly:     readOnly,
				SyncInterval: updateInterval,
				Config:       remoteConfig,
				file:         filepath.Join(dir, profileConfigFileName),
			},
		}
		err = profile.profConfig.write()
	}
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	return profile, nil
}

func isValidProfile(dir string) bool {
	if !commons.DirExists(dir) {
		return false
	}
	return commons.FileExists(filepath.Join(dir, profileConfigFileName))
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

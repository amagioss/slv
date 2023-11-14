package profiles

import (
	"os"
	"path/filepath"
	"time"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/environments"
	"github.com/shibme/slv/core/settings"
	"gopkg.in/yaml.v3"
)

type Profile struct {
	dir         *string
	path        *string
	settings    *settings.Settings
	envManifest *environments.EnvManifest
	*profile
}

type profile struct {
	Version uint8 `yaml:"version,omitempty"`
	Repo    struct {
		URI    string `yaml:"uri"`
		Branch string `yaml:"branch,omitempty"`
		Path   string `yaml:"path,omitempty"`
	} `yaml:"repo,omitempty"`
	Diff     bool      `yaml:"diff"`
	LastSync time.Time `yaml:"lastSync,omitempty"`
}

func (profile Profile) MarshalYAML() (interface{}, error) {
	return profile.profile, nil
}

func (profile *Profile) UnmarshalYAML(value *yaml.Node) (err error) {
	return value.Decode(&profile.profile)
}

func (profile *Profile) commit() error {
	if commons.WriteToYAML(*profile.path, "", profile) != nil {
		return ErrWritingManifest
	}
	return nil
}

func newProfileForPath(dir string) (*Profile, error) {
	if commons.DirExists(dir) {
		return nil, ErrProfilePathExistsAlready
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, ErrCreatingProfileDir
	}
	path := filepath.Join(dir, profileFileName)
	profile := &Profile{
		dir:  &dir,
		path: &path,
		profile: &profile{
			Version: commons.Version,
		},
	}
	err = profile.commit()
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func getProfileForPath(dir string) (*Profile, error) {
	if !commons.DirExists(dir) {
		return nil, ErrProfilePathDoesNotExist
	}
	path := filepath.Join(dir, profileFileName)
	profile := &Profile{}
	if err := commons.ReadFromYAML(path, profile); err != nil {
		return nil, err
	}
	profile.dir = &dir
	profile.path = &path
	return profile, nil
}

func (profile *Profile) GetSettings() (*settings.Settings, error) {
	if profile.settings == nil {
		settingsManifest, err := settings.GetManifest(filepath.Join(*profile.dir, profileSettingsFileName))
		if err != nil {
			settingsManifest, err = settings.NewManifest(filepath.Join(*profile.dir, profileSettingsFileName))
			if err != nil {
				return nil, err
			}
		}
		profile.settings = settingsManifest
	}
	return profile.settings, nil
}

func (profile *Profile) GetEnvManifest() (*environments.EnvManifest, error) {
	if profile.envManifest == nil {
		envManifest, err := environments.GetManifest(filepath.Join(*profile.dir, profileEnvironmentsFileName))
		if err != nil {
			envManifest, err = environments.NewManifest(filepath.Join(*profile.dir, profileEnvironmentsFileName))
			if err != nil {
				return nil, err
			}
		}
		profile.envManifest = envManifest
	}
	return profile.envManifest, nil
}

func (profile *Profile) AddEnvDef(envDef string) error {
	envManifest, err := profile.GetEnvManifest()
	if err != nil {
		return err
	}
	envManifest.AddEnvDef(envDef)
	return nil
}

func (profile *Profile) AddEnv(env *environments.Environment) error {
	envManifest, err := profile.GetEnvManifest()
	if err != nil {
		return err
	}
	return envManifest.AddEnv(env)
}

func (profile *Profile) InitRoot() (*crypto.SecretKey, error) {
	envManifest, err := profile.GetEnvManifest()
	if err != nil {
		return nil, err
	}
	return envManifest.InitRoot()
}

func (profile *Profile) Sync() {
	// TODO git operations
}

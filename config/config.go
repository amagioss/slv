package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/shibme/slv/commons"
	"github.com/shibme/slv/environment"
	"github.com/shibme/slv/settings"
	"gopkg.in/yaml.v3"
)

type Config struct {
	dir         *string
	path        *string
	settings    *settings.Settings
	envManifest *environment.EnvManifest
	*manifest
}

type manifest struct {
	Version string `yaml:"version,omitempty"`
	Repo    struct {
		URI    string `yaml:"uri"`
		Branch string `yaml:"branch,omitempty"`
		Path   string `yaml:"path,omitempty"`
	} `yaml:"repo,omitempty"`
	Diff     bool      `yaml:"diff"`
	LastSync time.Time `yaml:"last_sync,omitempty"`
}

func (config Config) MarshalYAML() (interface{}, error) {
	return config.manifest, nil
}

func (config *Config) UnmarshalYAML(value *yaml.Node) (err error) {
	return value.Decode(&config.manifest)
}

func (config *Config) commit() error {
	if commons.WriteToYAML(*config.path, config) != nil {
		return ErrWritingManifest
	}
	return nil
}

func NewConfig(dir string) (*Config, error) {
	if commons.DirExists(dir) {
		return nil, ErrConfigPathExistsAlready
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, ErrCreatingConfigDir
	}
	path := filepath.Join(dir, configFileName)
	config := &Config{
		dir:  &dir,
		path: &path,
		manifest: &manifest{
			Version: commons.Version,
		},
	}
	err = config.commit()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func GetConfig(dir string) (*Config, error) {
	if !commons.DirExists(dir) {
		return nil, ErrConfigPathDoesNotExist
	}
	path := filepath.Join(dir, configFileName)
	config := &Config{}
	if err := commons.ReadFromYAML(path, config); err != nil {
		return nil, err
	}
	config.dir = &dir
	config.path = &path
	return config, nil
}

func (config *Config) GetSettings() (*settings.Settings, error) {
	if config.settings == nil {
		settingsManifest, err := settings.GetManifest(filepath.Join(*config.dir, settingsFileName))
		if err != nil {
			settingsManifest, err = settings.NewManifest(filepath.Join(*config.dir, settingsFileName))
			if err != nil {
				return nil, err
			}
		}
		config.settings = settingsManifest
	}
	return config.settings, nil
}

func (config *Config) GetEnvManifest() (*environment.EnvManifest, error) {
	if config.envManifest == nil {
		envManifest, err := environment.GetManifest(filepath.Join(*config.dir, environmentConfigFileName))
		if err != nil {
			envManifest, err = environment.NewManifest(filepath.Join(*config.dir, environmentConfigFileName))
			if err != nil {
				return nil, err
			}
		}
		config.envManifest = envManifest
	}
	return config.envManifest, nil
}

func (config *Config) Sync() {
	// TODO git operations
}

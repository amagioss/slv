package settings

import (
	"github.com/shibme/slv/core/commons"
	"gopkg.in/yaml.v3"
)

// import (
// 	"os"

// 	"github.com/shibme/slv/core/commons"
// )

type Settings struct {
	path *string
	*manifest
}

type manifest struct {
	Version           string `yaml:"version,omitempty"`
	AllowChanges      bool   `yaml:"allow_changes"`
	AllowCreateEnv    bool   `yaml:"allow_create_env"`
	AllowCreateGroup  bool   `yaml:"allow_create_group"`
	SyncInterval      int    `yaml:"sync_interval"`
	AllowGroups       bool   `yaml:"allow_groups"`
	AllowVaultSharing bool   `yaml:"allow_vault_sharing"`
}

func (settings Settings) MarshalYAML() (interface{}, error) {
	return settings.manifest, nil
}

func (settings *Settings) UnmarshalYAML(value *yaml.Node) (err error) {
	return value.Decode(&settings.manifest)
}

func NewManifest(path string) (settings *Settings, err error) {
	if commons.FileExists(path) {
		return nil, ErrManifestPathExistsAlready
	}
	settings = &Settings{
		path: &path,
		manifest: &manifest{
			Version: commons.Version,
		},
	}
	err = settings.commit()
	if err != nil {
		settings = nil
	}
	return
}

func GetManifest(path string) (settings *Settings, err error) {
	if !commons.FileExists(path) {
		return nil, ErrManifestNotFound
	}
	settings = &Settings{}
	if err = commons.ReadFromYAML(path, settings); err != nil {
		return nil, err
	}
	settings.path = &path
	return
}

func (settings *Settings) commit() error {
	if commons.WriteToYAML(*settings.path, settings) != nil {
		return ErrWritingManifest
	}
	return nil
}

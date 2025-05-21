package settings

import (
	"slv.sh/slv/internal/core/commons"
)

type Settings struct {
	path              *string
	AllowChanges      bool `yaml:"allow_changes"`
	AllowCreateEnv    bool `yaml:"allow_create_env"`
	AllowCreateGroup  bool `yaml:"allow_create_group"`
	SyncInterval      int  `yaml:"sync_interval"`
	AllowGroups       bool `yaml:"allow_groups"`
	AllowVaultSharing bool `yaml:"allow_vault_sharing"`
}

func NewManifest(path string) (settings *Settings, err error) {
	if commons.FileExists(path) {
		return nil, errManifestPathExistsAlready
	}
	settings = &Settings{
		path: &path,
	}
	return
}

func GetManifest(path string) (settings *Settings, err error) {
	if !commons.FileExists(path) {
		return nil, errManifestNotFound
	}
	settings = &Settings{}
	if err = commons.ReadFromYAML(path, settings); err != nil {
		return nil, err
	}
	settings.path = &path
	return
}

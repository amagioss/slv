package settings

import (
	"slv.sh/slv/internal/core/commons"
)

type Settings struct {
	path              *string
	AllowChanges      bool `json:"allowChanges" yaml:"allowChanges"`
	AllowCreateEnv    bool `json:"allowCreateEnv" yaml:"allowCreateEnv"`
	AllowCreateGroup  bool `json:"allowCreateGroup" yaml:"allowCreateGroup"`
	SyncInterval      int  `json:"syncInterval" yaml:"syncInterval"`
	AllowGroups       bool `json:"allowGroups" yaml:"allowGroups"`
	AllowVaultSharing bool `json:"allowVaultSharing" yaml:"allowVaultSharing"`
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

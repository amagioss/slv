package config

import (
	"os"

	"github.com/shibme/slv/slv/access"
	"github.com/shibme/slv/slv/commons"
)

type GroupConfig struct {
	path           string
	pendingChanges bool
	data           *GroupConfigData
}

type GroupConfigData struct {
	Version string                  `yaml:"slv_version"`
	Groups  map[string]access.Group `yaml:"groups"`
}

func initGroupConfig(path string) (groupConfig *GroupConfig, err error) {
	groupConfig = &GroupConfig{}
	groupConfig.path = path
	var groupConfigFileInfo os.FileInfo
	if groupConfigFileInfo, err = os.Stat(groupConfig.path); err == nil {
		if groupConfigFileInfo.IsDir() {
			return nil, ErrProcessingGroupsManifest
		}
		err = commons.ReadFromYAML(groupConfig.path, &groupConfig)
		if err != nil {
			return nil, ErrProcessingGroupsManifest
		}
	} else {
		groupConfig.newGroupConfigData()
	}
	return groupConfig, nil
}

func (groupConfig *GroupConfig) newGroupConfigData() error {
	groupConfig.data = &GroupConfigData{}
	groupConfig.data.Version = commons.Version
	groupConfig.data.Groups = make(map[string]access.Group)
	return groupConfig.Commit()
}

func (groupsManifest *GroupConfig) Commit() error {
	if groupsManifest.pendingChanges {
		if commons.WriteToYAML(groupsManifest.path, &groupsManifest) != nil {
			return ErrProcessingGroupsManifest
		}
		groupsManifest.pendingChanges = false
	}
	return nil
}

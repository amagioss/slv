package config

import (
	"github.com/shibme/slv/slv/commons"
)

type Settings struct {
	path           string
	pendingChanges bool
	data           *SettingsData
}

type SettingsData struct {
	Version           string `yaml:"slv_version"`
	AllowChanges      bool   `yaml:"allow_changes"`
	SyncInterval      int    `yaml:"sync_interval"`
	AllowGroups       bool   `yaml:"allow_groups"`
	AllowVaultSharing bool   `yaml:"allow_vault_sharing"`
}

func initSettings(path string) (settings *Settings, err error) {
	settings = &Settings{}
	settings.path = path
	if err = commons.ReadFromYAML(settings.path, settings.data); err != nil {
		return nil, ErrProcessingSettings
	} else if err = settings.newSettingsData(); err != nil {
		return nil, err
	}
	return settings, nil
}

func (settings *Settings) newSettingsData() error {
	settings.data = &SettingsData{}
	settings.data.Version = commons.Version
	settings.data.AllowChanges = true
	settings.data.SyncInterval = defaultSyncInterval
	settings.data.AllowGroups = true
	settings.data.AllowVaultSharing = true
	settings.pendingChanges = true
	return settings.Commit()
}

func (settings *Settings) Commit() error {
	if settings.pendingChanges {
		if commons.WriteToYAML(settings.path, &settings.data) != nil {
			return ErrProcessingSettings
		}
		settings.pendingChanges = false
	}
	return nil
}

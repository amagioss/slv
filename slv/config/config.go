package config

import (
	"time"
)

type Config struct {
	configManager     *configManager
	dir               string
	configFile        string
	dataDir           string
	configInfo        ConfigInfo
	settings          *Settings
	environmentConfig *EnvironmentConfig
	groupConfig       *GroupConfig
}

type ConfigInfo struct {
	RemoteRepo   string    `yaml:"remoteRepo"`
	RemoteBranch string    `yaml:"remoteBranch"`
	RemoteDir    string    `yaml:"remoteDir"`
	Diff         bool      `yaml:"diff"`
	SyncedAt     time.Time `yaml:"syncedAt"`
}

func (config *Config) Sync() {
	// TODO
}

func (config *Config) GetSettings() (settings Settings) {
	return *config.settings
}

func (config *Config) GetEnvironmentsManifest() (environmentManifest EnvironmentConfig) {
	return *config.environmentConfig
}

func (config *Config) GetGroupsManifest() (groupManager GroupConfig) {
	return *config.groupConfig
}

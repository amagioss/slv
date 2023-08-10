package config

import (
	"time"

	"github.com/shibme/slv/environment"
)

type Config struct {
	configManager *configManager
	dir           string
	configFile    string
	dataDir       string
	data          *ConfigData
	settings      *Settings
	envManifest   *environment.EnvManifest
	groupConfig   *GroupConfig
}

type ConfigData struct {
	RemoteRepo   string    `yaml:"remoteRepo"`
	RemoteBranch string    `yaml:"remoteBranch"`
	RemoteDir    string    `yaml:"remoteDir"`
	Diff         bool      `yaml:"diff"`
	SyncedAt     time.Time `yaml:"syncedAt"`
}

func (config *Config) Sync() {
	// TODO
}

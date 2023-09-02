package config

import (
	"time"

	"github.com/shibme/slv/environment"
	"github.com/shibme/slv/settings"
)

type Config struct {
	configManager *configManager
	dir           string
	configFile    string
	dataDir       string
	data          *ConfigData
	settings      *settings.Settings
	envManifest   *environment.EnvManifest
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

package config

import (
	"os"
	"path/filepath"

	"github.com/shibme/slv/slv/commons"
)

type configManager struct {
	dir           string
	prefsFile     string
	preferences   ConfigPreferences
	configList    map[string]struct{}
	currentConfig *Config
}

type ConfigPreferences struct {
	CurrentConfigName string `yaml:"current"`
}

var confManager *configManager = nil
var comfigMap map[string]*Config = make(map[string]*Config)

func initConfigManager() (err error) {
	if confManager != nil {
		return nil
	}
	var confMgr configManager
	confMgr.dir = filepath.Join(commons.AppDataDir(), configManagerDirName)
	confMgr.prefsFile = filepath.Join(confMgr.dir, configManagerPreferencesFileName)

	// Creates the configs directory if it doesn't exist
	configManagerDirInfo, err := os.Stat(confMgr.dir)
	if err != nil {
		err = os.MkdirAll(confMgr.dir, 0755)
		if err != nil {
			return ErrConfigManagerInitialization
		}
	} else if !configManagerDirInfo.IsDir() {
		return ErrConfigManagerInitialization
	}

	// Create the config preferences file if it doesn't exist
	preferencesFileInfo, err := os.Stat(confMgr.prefsFile)
	if err != nil {
		confMgr.preferences = ConfigPreferences{}
		err = confMgr.savePreferences()
		if err != nil {
			return ErrConfigManagerInitialization
		}
	} else if preferencesFileInfo.IsDir() {
		return ErrConfigManagerInitialization
	} else {
		err = commons.ReadFromYAML(confMgr.prefsFile, &confMgr.preferences)
		if err != nil {
			return ErrConfigManagerInitialization
		}
	}

	// Read the manifest config directory and populate the config list
	confManagerDir, err := os.Open(confMgr.dir)
	if err != nil {
		return ErrOpeningConfigManagerDir
	}
	defer confManagerDir.Close()
	fileInfoList, err := confManagerDir.Readdir(-1)
	if err != nil {
		return ErrOpeningConfigManagerDir
	}
	confMgr.configList = make(map[string]struct{})
	for _, fileInfo := range fileInfoList {
		if fileInfo.IsDir() {
			if f, err := os.Stat(filepath.Join(confMgr.dir, fileInfo.Name(), configFileName)); err == nil && !f.IsDir() {
				confMgr.configList[fileInfo.Name()] = struct{}{}
			}
		}
	}
	confManager = &confMgr
	return nil
}

func (confMgr *configManager) savePreferences() error {
	if commons.WriteToYAML(confMgr.prefsFile, confMgr.preferences) != nil {
		return ErrSavingConfigManagerPreferences
	}
	return nil
}

func GetAllConfigNames() (configNames []string, err error) {
	if err = initConfigManager(); err != nil {
		return nil, err
	}
	configNames = make([]string, 0, len(confManager.configList))
	for configName := range confManager.configList {
		configNames = append(configNames, configName)
	}
	return
}

func GetConfig(configName string) (config *Config, err error) {
	if err = initConfigManager(); err != nil {
		return nil, err
	}
	if config = comfigMap[configName]; config != nil {
		return
	}
	if _, exists := confManager.configList[configName]; !exists {
		return nil, ErrConfigNotFound
	}
	if config, err = confManager.initConfig(configName); err != nil {
		return nil, ErrConfigInitialization
	}
	comfigMap[configName] = config
	return
}

func SetCurrentConfig(configName string) (config *Config, err error) {
	config, err = GetConfig(configName)
	if err == nil {
		confManager.preferences.CurrentConfigName = configName
		confManager.currentConfig = config
		err = confManager.savePreferences()
		if err != nil {
			return nil, err
		}
	}
	return
}

func GetCurrentConfig() (config *Config, err error) {
	if err = initConfigManager(); err != nil {
		return nil, err
	}
	if confManager.currentConfig != nil {
		return confManager.currentConfig, nil
	}
	if confManager.preferences.CurrentConfigName == "" {
		return nil, ErrNoCurrentConfigFound
	}
	return GetConfig(confManager.preferences.CurrentConfigName)
}

func NewConfig(configName string) (config *Config, err error) {
	if err = initConfigManager(); err != nil {
		return nil, err
	}
	if _, exists := confManager.configList[configName]; exists {
		return nil, ErrManifestExistsAlready
	}
	if config, err = confManager.initConfig(configName); err != nil {
		return nil, ErrConfigInitialization
	}
	confManager.configList[configName] = struct{}{}
	return
}

func (configManager *configManager) initConfig(configName string) (config *Config, err error) {
	config = &Config{}
	config.configManager = configManager
	config.dir = filepath.Join(configManager.dir, configName)
	config.dataDir = filepath.Join(config.dir, configDataDirName)
	config.configFile = filepath.Join(config.dir, configFileName)

	// Creating config data dir if it doesn't exists
	if configDataDirInfo, err := os.Stat(config.dataDir); err == nil {
		if !configDataDirInfo.IsDir() {
			return nil, ErrConfigInitialization
		}
	} else {
		err = os.MkdirAll(config.dataDir, 0755)
		if err != nil {
			return nil, ErrConfigInitialization
		}
	}

	// Attempting to read the config file, if not creating one
	if configFileInfo, err := os.Stat(config.configFile); err == nil {
		if configFileInfo.IsDir() {
			return nil, ErrConfigInitialization
		} else {
			if commons.ReadFromYAML(config.configFile, &config.configInfo) != nil {
				return nil, ErrConfigInitialization
			}
		}
	} else {
		config.configInfo = ConfigInfo{}
		if commons.WriteToYAML(config.configFile, &config.configInfo) != nil {
			return nil, ErrConfigInitialization
		}
	}

	// Attempting to initialize settings
	settingsFile := filepath.Join(config.dataDir, settingsFileName)
	if config.settings, err = initSettings(settingsFile); err != nil {
		return nil, err
	}

	// Attempting to initialize environments manifest
	environmentConfigFile := filepath.Join(config.dataDir, environmentConfigFileName)
	if config.environmentConfig, err = initEnvironmentConfig(environmentConfigFile); err != nil {
		return nil, err
	}

	// Attempting to initialize groups manifest
	groupConfigFile := filepath.Join(config.dataDir, groupConfigFileName)
	if config.groupConfig, err = initGroupConfig(groupConfigFile); err != nil {
		return nil, err
	}
	return
}

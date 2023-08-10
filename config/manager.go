package config

import (
	"os"
	"path/filepath"

	"github.com/shibme/slv/commons"
)

type configManager struct {
	dir       string
	prefsFile string
	pref      Pref
	list      map[string]struct{}
	current   *Config
}

type Pref struct {
	Current string `yaml:"current"`
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
		confMgr.pref = Pref{}
		err = confMgr.savePreferences()
		if err != nil {
			return ErrConfigManagerInitialization
		}
	} else if preferencesFileInfo.IsDir() {
		return ErrConfigManagerInitialization
	} else {
		err = commons.ReadFromYAML(confMgr.prefsFile, &confMgr.pref)
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
	confMgr.list = make(map[string]struct{})
	for _, fileInfo := range fileInfoList {
		if fileInfo.IsDir() {
			if f, err := os.Stat(filepath.Join(confMgr.dir, fileInfo.Name(), configFileName)); err == nil && !f.IsDir() {
				confMgr.list[fileInfo.Name()] = struct{}{}
			}
		}
	}
	confManager = &confMgr
	return nil
}

func (confMgr *configManager) savePreferences() error {
	if commons.WriteToYAML(confMgr.prefsFile, confMgr.pref) != nil {
		return ErrSavingConfigManagerPreferences
	}
	return nil
}

func List() (configNames []string, err error) {
	if err = initConfigManager(); err != nil {
		return nil, err
	}
	configNames = make([]string, 0, len(confManager.list))
	for configName := range confManager.list {
		configNames = append(configNames, configName)
	}
	return
}

func getConfig(configName string) (config *Config, err error) {
	if err = initConfigManager(); err != nil {
		return nil, err
	}
	if config = comfigMap[configName]; config != nil {
		return
	}
	if _, exists := confManager.list[configName]; !exists {
		return nil, ErrConfigNotFound
	}
	if config, err = confManager.initConfig(configName); err != nil {
		return nil, ErrConfigInitialization
	}
	comfigMap[configName] = config
	return
}

func Set(configName string) error {
	if err := initConfigManager(); err != nil {
		return err
	}
	if _, exists := confManager.list[configName]; !exists {
		return ErrConfigNotFound
	}
	confManager.pref.Current = configName
	return confManager.savePreferences()
}

func GetCurrentConfigName() (currentConfigName string, err error) {
	if err = initConfigManager(); err != nil {
		return "", err
	}
	if confManager.pref.Current != "" {
		return confManager.pref.Current, nil
	} else {
		return "", ErrNoCurrentConfigFound
	}
}

func GetCurrentConfig() (config *Config, err error) {
	if err = initConfigManager(); err != nil {
		return nil, err
	}
	if confManager.current != nil {
		return confManager.current, nil
	}
	if confManager.pref.Current == "" {
		return nil, ErrNoCurrentConfigFound
	}
	return getConfig(confManager.pref.Current)
}

func New(configName string) error {
	if err := initConfigManager(); err != nil {
		return err
	}
	if _, exists := confManager.list[configName]; exists {
		return ErrConfigExistsAlready
	}
	if _, err := confManager.initConfig(configName); err != nil {
		return err
	}
	confManager.list[configName] = struct{}{}
	if confManager.pref.Current == "" {
		return Set(configName)
	}
	return nil
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
			if commons.ReadFromYAML(config.configFile, &config.data) != nil {
				return nil, ErrConfigInitialization
			}
		}
	} else {
		config.data = &ConfigData{}
		if commons.WriteToYAML(config.configFile, &config.data) != nil {
			return nil, ErrConfigInitialization
		}
	}

	// Attempting to initialize settings
	settingsFile := filepath.Join(config.dataDir, settingsFileName)
	if config.settings, err = initSettings(settingsFile); err != nil {
		return nil, err
	}

	// Attempting to initialize environments manifest
	// environmentConfigFile := filepath.Join(config.dataDir, environmentConfigFileName)
	// if config.envConfig, err = initEnvironmentConfig(environmentConfigFile); err != nil {
	// 	return nil, err
	// }

	// Attempting to initialize groups manifest
	groupConfigFile := filepath.Join(config.dataDir, groupConfigFileName)
	if config.groupConfig, err = initGroupConfig(groupConfigFile); err != nil {
		return nil, err
	}
	return
}

// func AddEnvToConfig(configName, envDef string) error {
// 	cfg, err := getConfig(configName)
// 	if err != nil {
// 		return err
// 	}
// 	envConf := cfg.GetEnvConfig()
// 	envConf.AddEnvironment(envDef)
// }

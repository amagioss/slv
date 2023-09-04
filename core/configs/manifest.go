package configs

import (
	"os"
	"path/filepath"

	"github.com/shibme/slv/core/commons"
)

type configManager struct {
	dir               string
	confList          map[string]struct{}
	defaultFile       *string
	defaultConfigName *string
	defaultConfig     *Config
}

var configMgr *configManager = nil
var comfigMap map[string]*Config = make(map[string]*Config)

func initConfigManager() error {
	if configMgr != nil {
		return nil
	}
	var manager configManager
	manager.dir = filepath.Join(commons.AppDataDir(), configsDirName)

	// Creates the configs directory if it doesn't exist
	configManagerDirInfo, err := os.Stat(manager.dir)
	if err != nil {
		err = os.MkdirAll(manager.dir, 0755)
		if err != nil {
			return ErrCreatingConfigCollectionDir
		}
	} else if !configManagerDirInfo.IsDir() {
		return ErrInitializingConfigManagementDir
	}

	// Populating list of all configs
	confManagerDir, err := os.Open(manager.dir)
	if err != nil {
		return ErrOpeningConfigManagementDir
	}
	defer confManagerDir.Close()
	fileInfoList, err := confManagerDir.Readdir(-1)
	if err != nil {
		return ErrOpeningConfigManagementDir
	}
	manager.confList = make(map[string]struct{})
	for _, fileInfo := range fileInfoList {
		if fileInfo.IsDir() {
			if f, err := os.Stat(filepath.Join(manager.dir, fileInfo.Name(), configFileName)); err == nil && !f.IsDir() {
				manager.confList[fileInfo.Name()] = struct{}{}
			}
		}
	}
	configMgr = &manager
	GetDefaultConfigName()
	return nil
}

func List() (configNames []string, err error) {
	if err = initConfigManager(); err != nil {
		return nil, err
	}
	configNames = make([]string, 0, len(configMgr.confList))
	for configName := range configMgr.confList {
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
	if _, exists := configMgr.confList[configName]; !exists {
		return nil, ErrConfigNotFound
	}
	if config, err = getConfigForPath(filepath.Join(configMgr.dir, configName)); err != nil {
		return nil, err
	}
	comfigMap[configName] = config
	return
}

func New(configName string) error {
	if configName == "" {
		return ErrInvalidConfigName
	}
	if err := initConfigManager(); err != nil {
		return err
	}
	if _, exists := configMgr.confList[configName]; exists {
		return ErrConfigExistsAlready
	}
	if _, err := newConfigForPath(filepath.Join(configMgr.dir, configName)); err != nil {
		return err
	}
	configMgr.confList[configName] = struct{}{}
	if configMgr.defaultConfigName == nil {
		return SetDefault(configName)
	}
	return nil
}

func SetDefault(configName string) error {
	if configName == "" {
		return ErrInvalidConfigName
	}
	if err := initConfigManager(); err != nil {
		return err
	}
	if _, exists := configMgr.confList[configName]; !exists {
		return ErrConfigNotFound
	}
	if configMgr.defaultFile == nil {
		defaultInfoFile := filepath.Join(configMgr.dir, defaultConfigFileName)
		configMgr.defaultFile = &defaultInfoFile
	}
	err := os.WriteFile(*configMgr.defaultFile, []byte(configName), 0644)
	if err != nil {
		return ErrSettingDefaultConfig
	}
	configMgr.defaultConfig = nil
	return nil
}

func GetDefaultConfigName() (string, error) {
	if err := initConfigManager(); err != nil {
		return "", err
	}
	if configMgr.defaultConfigName != nil {
		return *configMgr.defaultConfigName, nil
	}
	if configMgr.defaultFile == nil {
		defaultInfoFile := filepath.Join(configMgr.dir, defaultConfigFileName)
		configMgr.defaultFile = &defaultInfoFile
	}
	bytes, err := os.ReadFile(*configMgr.defaultFile)
	if err != nil {
		return "", ErrNoDefaultConfigFound
	}
	defaultConfigName := string(bytes)
	configMgr.defaultConfigName = &defaultConfigName
	return *configMgr.defaultConfigName, nil
}

func GetDefaultConfig() (config *Config, err error) {
	if err = initConfigManager(); err != nil {
		return nil, err
	}
	if configMgr.defaultConfig == nil {
		defaultConfigName, err := GetDefaultConfigName()
		if err != nil {
			return nil, err
		}
		if configMgr.defaultConfig, err = GetConfig(defaultConfigName); err != nil {
			return nil, err
		}
	}
	return configMgr.defaultConfig, nil
}

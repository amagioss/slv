package config

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

package config

import (
	"strings"

	"github.com/shibme/slv/slv/commons"
	"github.com/shibme/slv/slv/crypto"
	"github.com/shibme/slv/slv/environment"
)

type EnvironmentConfig struct {
	path           string
	pendingChanges bool
	data           *EnvironmentConfigData
}

type EnvironmentConfigData struct {
	Version         string                              `yaml:"slv_version"`
	RootEnvironment *environment.Environment            `yaml:"root"`
	Environments    map[string]*environment.Environment `yaml:"environments"`
}

func initEnvironmentConfig(path string) (envConfig *EnvironmentConfig, err error) {
	envConfig = &EnvironmentConfig{}
	envConfig.path = path
	if err = commons.ReadFromYAML(envConfig.path, envConfig.data); err != nil {
		return nil, ErrProcessingSettings
	} else if err = envConfig.newEnvConfigData(); err != nil {
		return nil, err
	}
	return envConfig, nil
}

func (envConfig *EnvironmentConfig) newEnvConfigData() error {
	envConfig.data = &EnvironmentConfigData{}
	envConfig.data.Version = commons.Version
	envConfig.data.Environments = make(map[string]*environment.Environment)
	return envConfig.Commit()
}

func (envConfig *EnvironmentConfig) Commit() error {
	if envConfig.pendingChanges {
		if commons.WriteToYAML(envConfig.path, &envConfig) != nil {
			return ErrProcessingEnvironmentsManifest
		}
		envConfig.pendingChanges = false
	}
	return nil
}

func (envConfig *EnvironmentConfig) GetRootEnv() (environment environment.Environment) {
	return *envConfig.data.RootEnvironment
}

func (envConfig *EnvironmentConfig) GetEnvironments() (environments []*environment.Environment) {
	for _, environment := range envConfig.data.Environments {
		environments = append(environments, environment)
	}
	return
}

func (envConfig *EnvironmentConfig) GetEnvironment(id string) (environment *environment.Environment, err error) {
	if environment, ok := envConfig.data.Environments[id]; ok {
		return environment, nil
	}
	return nil, ErrEnvironmentNotFound
}

// Lists environments that match a given query by searching for parts of name, email and tags
func (envConfig *EnvironmentConfig) SearchEnvironments(query string) (environments []*environment.Environment) {
	query = strings.ToLower(query)
	for _, environment := range envConfig.data.Environments {
		if environment.Search(query) {
			environments = append(environments, environment)
		}
	}
	return
}

func (envConfig *EnvironmentConfig) NewEnvironment(name, email string, envType environment.EnvType) (env *environment.Environment, privateKey *crypto.PrivateKey, err error) {
	env, privateKey, err = environment.New(name, email, envType)
	if err != nil {
		return
	}
	envConfig.updateEnvironment(env)
	return
}

func (envConfig *EnvironmentConfig) updateEnvironment(env *environment.Environment) {
	envConfig.data.Environments[env.Id()] = env
	envConfig.pendingChanges = true
}

func (envConfig *EnvironmentConfig) AddEnvironment(envString string) (err error) {
	environment, err := environment.EnvFromSLVFormat(envString)
	if err != nil {
		return
	}
	envConfig.updateEnvironment(&environment)
	return
}

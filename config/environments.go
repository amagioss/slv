package config

import (
	"os"
	"strings"

	"github.com/shibme/slv/commons"
	"github.com/shibme/slv/crypto"
	"github.com/shibme/slv/environment"
)

type EnvConfig struct {
	path           string
	pendingChanges bool
	data           *EnvironmentConfigData
}

type EnvironmentConfigData struct {
	Version         string                              `yaml:"version"`
	RootEnvironment *environment.Environment            `yaml:"root"`
	Environments    map[string]*environment.Environment `yaml:"environments"`
}

func initEnvironmentConfig(path string) (envConfig *EnvConfig, err error) {
	envConfig = &EnvConfig{}
	envConfig.path = path
	if _, err := os.Stat(envConfig.path); err == nil {
		if err = commons.ReadFromYAML(envConfig.path, envConfig.data); err != nil {
			return nil, ErrProcessingEnvironmentsManifest
		}
	} else if err = envConfig.newEnvConfigData(); err != nil {
		return nil, err
	}
	return envConfig, nil
}

func (envConfig *EnvConfig) newEnvConfigData() error {
	envConfig.data = &EnvironmentConfigData{}
	envConfig.data.Version = commons.Version
	envConfig.data.Environments = make(map[string]*environment.Environment)
	return envConfig.Commit()
}

func (envConfig *EnvConfig) Commit() error {
	if envConfig.pendingChanges {
		if commons.WriteToYAML(envConfig.path, &envConfig) != nil {
			return ErrProcessingEnvironmentsManifest
		}
		envConfig.pendingChanges = false
	}
	return nil
}

func (envConfig *EnvConfig) GetRootEnv() (environment environment.Environment) {
	return *envConfig.data.RootEnvironment
}

func (envConfig *EnvConfig) GetEnvironments() (environments []*environment.Environment) {
	for _, environment := range envConfig.data.Environments {
		environments = append(environments, environment)
	}
	return
}

func (envConfig *EnvConfig) GetEnvironment(id string) (environment *environment.Environment, err error) {
	if environment, ok := envConfig.data.Environments[id]; ok {
		return environment, nil
	}
	return nil, ErrEnvironmentNotFound
}

// Lists environments that match a given query by searching for parts of name, email and tags
func (envConfig *EnvConfig) SearchEnvironments(query string) (environments []*environment.Environment) {
	query = strings.ToLower(query)
	for _, environment := range envConfig.data.Environments {
		if environment.Search(query) {
			environments = append(environments, environment)
		}
	}
	return
}

func (envConfig *EnvConfig) NewEnvironment(name, email string, envType environment.EnvType) (env *environment.Environment, privateKey *crypto.PrivateKey, err error) {
	env, privateKey, err = environment.New(name, email, envType)
	if err != nil {
		return
	}
	envConfig.updateEnvironment(env)
	return
}

func (envConfig *EnvConfig) updateEnvironment(env *environment.Environment) {
	envConfig.data.Environments[env.Id()] = env
	envConfig.pendingChanges = true
}

func (envConfig *EnvConfig) AddEnvironment(enfDef string) (err error) {
	environment, err := environment.FromEnvDef(enfDef)
	if err != nil {
		return
	}
	envConfig.updateEnvironment(&environment)
	return
}

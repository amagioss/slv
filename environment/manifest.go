package environment

import (
	"strings"

	"github.com/shibme/slv/commons"
	"github.com/shibme/slv/crypto"
)

type EnvManifest struct {
	path string
	data *EnvManifestData
}

type EnvManifestData struct {
	Version      string                  `yaml:"version"`
	Root         *Root                   `yaml:"root"`
	Environments map[string]*Environment `yaml:"environments"`
}

func NewManifest(path string) (envManifest *EnvManifest, err error) {
	if commons.FileExists(path) {
		return nil, ErrManifestExistsAlready
	}
	envManifest = &EnvManifest{}
	envManifest.path = path
	envManifest.data = &EnvManifestData{}
	envManifest.data.Version = commons.Version
	envManifest.data.Environments = make(map[string]*Environment)
	err = envManifest.commit()
	if err != nil {
		envManifest = nil
	}
	return
}

func GetManifest(path string) (envManifest *EnvManifest, err error) {
	if !commons.FileExists(path) {
		return nil, ErrManifestNotFound
	}
	envManifestData := &EnvManifestData{}
	if err = commons.ReadFromYAML(path, envManifestData); err != nil {
		return nil, err
	}
	envManifest = &EnvManifest{}
	envManifest.path = path
	envManifest.data = envManifestData
	return
}

func (envManifest *EnvManifest) commit() error {
	if commons.WriteToYAML(envManifest.path, &envManifest.data) != nil {
		return ErrWritingManifest
	}
	return nil
}

func (envManifest *EnvManifest) HasRoot() bool {
	return envManifest.data.Root != nil
}

func (envManifest *EnvManifest) InitRoot() (*crypto.PrivateKey, error) {
	if envManifest.HasRoot() {
		return nil, ErrManifestRootExistsAlready
	}
	root, rootPrivateKey, err := newRoot()
	if err != nil {
		return nil, err
	}
	envManifest.data.Root = root
	err = envManifest.commit()
	if err != nil {
		return nil, err
	}
	return rootPrivateKey, nil
}

func (envConfig *EnvManifest) ListEnv() (environments []*Environment) {
	for _, environment := range envConfig.data.Environments {
		environments = append(environments, environment)
	}
	return
}

func (envConfig *EnvManifest) GetEnv(id string) (environment *Environment, err error) {
	if environment, ok := envConfig.data.Environments[id]; ok {
		return environment, nil
	}
	return nil, ErrEnvironmentNotFound
}

// Lists environments that match a given query by searching for parts of name, email and tags
func (envConfig *EnvManifest) SearchEnv(query string) (environments []*Environment) {
	query = strings.ToLower(query)
	for _, env := range envConfig.data.Environments {
		if env.Search(query) {
			environments = append(environments, env)
		}
	}
	return
}

func (envConfig *EnvManifest) updateEnvironment(env *Environment) error {
	envConfig.data.Environments[env.Id()] = env
	return envConfig.commit()
}

func (envConfig *EnvManifest) AddEnv(envString string) (err error) {
	environment, err := FromEnvDef(envString)
	if err != nil {
		return
	}
	envConfig.updateEnvironment(&environment)
	return
}

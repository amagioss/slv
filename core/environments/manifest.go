package environments

import (
	"strings"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"gopkg.in/yaml.v3"
)

type EnvManifest struct {
	path *string
	*manifest
}

type manifest struct {
	Version      string                  `yaml:"version"`
	Root         *Root                   `yaml:"root,omitempty"`
	Environments map[string]*Environment `yaml:"environments,omitempty"`
}

func (envManifest EnvManifest) MarshalYAML() (interface{}, error) {
	return envManifest.manifest, nil
}

func (envManifest *EnvManifest) UnmarshalYAML(value *yaml.Node) (err error) {
	return value.Decode(&envManifest.manifest)
}

func NewManifest(path string) (envManifest *EnvManifest, err error) {
	if commons.FileExists(path) {
		return nil, ErrManifestPathExistsAlready
	}
	envManifest = &EnvManifest{
		path: &path,
		manifest: &manifest{
			Version: commons.Version,
		},
	}
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
	envManifest = &EnvManifest{}
	if err = commons.ReadFromYAML(path, envManifest); err != nil {
		return nil, err
	}
	envManifest.path = &path
	return
}

func (envManifest *EnvManifest) commit() error {
	if commons.WriteToYAML(*envManifest.path, envManifest) != nil {
		return ErrWritingManifest
	}
	return nil
}

func (envManifest *EnvManifest) RootPublicKey() *crypto.PublicKey {
	if envManifest.Root != nil {
		return &envManifest.Root.PublicKey
	}
	return nil
}

func (envManifest *EnvManifest) InitRoot() (*crypto.SecretKey, error) {
	if envManifest.Root != nil {
		return nil, ErrManifestRootExistsAlready
	}
	root, rootSecretKey, err := newRoot()
	if err != nil {
		return nil, err
	}
	envManifest.Root = root
	err = envManifest.commit()
	if err != nil {
		return nil, err
	}
	return rootSecretKey, nil
}

func (envManifest *EnvManifest) ListEnv() (environments []*Environment) {
	if envManifest.Environments != nil {
		for _, environment := range envManifest.Environments {
			environments = append(environments, environment)
		}
	}
	return
}

func (envManifest *EnvManifest) GetEnv(id string) (environment *Environment, err error) {
	if environment, ok := envManifest.Environments[id]; ok {
		return environment, nil
	}
	return nil, ErrEnvironmentNotFound
}

// Lists environments that match a given query by searching for parts of name, email and tags
func (envManifest *EnvManifest) SearchEnv(query string) (environments []*Environment) {
	query = strings.ToLower(query)
	for _, env := range envManifest.Environments {
		if env.Search(query) {
			environments = append(environments, env)
		}
	}
	return
}

func (envManifest *EnvManifest) updateEnvironment(env *Environment) error {
	if envManifest.Environments == nil {
		envManifest.Environments = make(map[string]*Environment)
	}
	envManifest.Environments[env.Id()] = env
	return envManifest.commit()
}

func (envManifest *EnvManifest) AddEnv(envString string) (err error) {
	environment, err := FromEnvDef(envString)
	if err != nil {
		return
	}
	envManifest.updateEnvironment(environment)
	return
}

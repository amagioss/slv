package environments

import (
	"strings"

	"slv.sh/slv/internal/core/commons"
)

type EnvManifest struct {
	path         *string
	Root         *Environment            `json:"root,omitempty" yaml:"root,omitempty"`
	Environments map[string]*Environment `json:"environments,omitempty" yaml:"environments,omitempty"`
}

func NewManifest(path string) (envManifest *EnvManifest, err error) {
	if commons.FileExists(path) {
		return nil, errManifestPathExistsAlready
	}
	envManifest = &EnvManifest{
		path: &path,
	}
	return
}

func GetManifest(path string) (envManifest *EnvManifest, err error) {
	if !commons.FileExists(path) {
		return nil, errManifestNotFound
	}
	envManifest = &EnvManifest{}
	if err = commons.ReadFromYAML(path, envManifest); err != nil {
		return nil, err
	}
	envManifest.path = &path
	return
}

func (envManifest *EnvManifest) write() error {
	if commons.WriteToYAML(*envManifest.path, envManifest) != nil {
		return errWritingManifest
	}
	return nil
}

func (envManifest *EnvManifest) GetRoot() (*Environment, error) {
	if envManifest.Root != nil {
		return envManifest.Root, nil
	}
	return nil, nil
}

func (envManifest *EnvManifest) SetRoot(env *Environment) error {
	if envManifest.Root != nil {
		return errRootExistsAlready
	}
	envManifest.Root = env
	return envManifest.write()
}

func (envManifest *EnvManifest) ListEnvs() (environments []*Environment) {
	if envManifest.Environments != nil {
		for _, environment := range envManifest.Environments {
			environments = append(environments, environment)
		}
	}
	return
}

func (envManifest *EnvManifest) GetEnv(id string) (env *Environment) {
	if envManifest.Environments != nil {
		env = envManifest.Environments[id]
	}
	return
}

func (envManifest *EnvManifest) searchEnvs(query string) (environments []*Environment) {
	query = strings.ToLower(query)
	for _, env := range envManifest.Environments {
		if env.Search(query) {
			environments = append(environments, env)
		}
	}
	return
}

func (envManifest *EnvManifest) SearchEnvs(queries []string) (environments []*Environment) {
	for _, query := range queries {
		if query != "" {
			environments = append(environments, envManifest.searchEnvs(query)...)
		}
	}
	return
}

func (envManifest *EnvManifest) DeleteEnv(id string) (env *Environment, err error) {
	if envManifest.Environments == nil {
		return nil, errEnvNotFound
	}
	if envManifest.Environments[id] == nil {
		return nil, errEnvNotFound
	}
	env = envManifest.Environments[id]
	delete(envManifest.Environments, id)
	return env, envManifest.write()
}

func (envManifest *EnvManifest) PutEnv(env *Environment) (err error) {
	if envManifest.Environments == nil {
		envManifest.Environments = make(map[string]*Environment)
	}
	envManifest.Environments[env.PublicKey] = env
	return envManifest.write()
}

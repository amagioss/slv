package environments

import (
	"strings"

	"gopkg.in/yaml.v3"
	"savesecrets.org/slv/core/commons"
	"savesecrets.org/slv/core/crypto"
)

type EnvManifest struct {
	path *string
	*manifest
}

type manifest struct {
	Root         *Environment            `yaml:"root,omitempty"`
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
		return nil, errManifestPathExistsAlready
	}
	envManifest = &EnvManifest{
		path:     &path,
		manifest: new(manifest),
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

func (envManifest *EnvManifest) commit() error {
	if commons.WriteToYAML(*envManifest.path, "", envManifest) != nil {
		return errWritingManifest
	}
	return nil
}

func (envManifest *EnvManifest) RootPublicKey() (*crypto.PublicKey, error) {
	if envManifest.Root != nil {
		return envManifest.Root.getPublicKey()
	}
	return nil, nil
}

func (envManifest *EnvManifest) SetRoot(env *Environment) error {
	if envManifest.Root != nil {
		return errRootExistsAlready
	}
	envManifest.Root = env
	return envManifest.commit()
}

func (envManifest *EnvManifest) ListEnvs() (environments []*Environment) {
	if envManifest.Environments != nil {
		for _, environment := range envManifest.Environments {
			environments = append(environments, environment)
		}
	}
	return
}

func (envManifest *EnvManifest) SearchEnvs(query string) (environments []*Environment) {
	query = strings.ToLower(query)
	for _, env := range envManifest.Environments {
		if env.Search(query) {
			environments = append(environments, env)
		}
	}
	return
}

func (envManifest *EnvManifest) PutEnv(env *Environment) (err error) {
	if envManifest.Environments == nil {
		envManifest.Environments = make(map[string]*Environment)
	}
	envManifest.Environments[env.Id()] = env
	return envManifest.commit()
}

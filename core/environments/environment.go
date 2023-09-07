package environments

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"gopkg.in/yaml.v3"
)

type EnvType string

type Environment struct {
	*environment
}

type environment struct {
	PublicKey crypto.PublicKey `yaml:"publicKey"`
	Name      string           `yaml:"name"`
	Email     string           `yaml:"email"`
	EnvType   EnvType          `yaml:"type"`
	Tags      []string         `yaml:"tags"`
}

func (eType *EnvType) isValid() bool {
	return *eType == SERVICE || *eType == USER || *eType == ROOT
}

func New(name, email string, envType EnvType) (*Environment, *crypto.SecretKey, error) {
	if !envType.isValid() {
		return nil, nil, ErrInvalidEnvironmentType
	}
	pKey, sKey, err := crypto.NewKeyPair(EnvironmentKey)
	if err != nil {
		return nil, nil, err
	}
	return &Environment{
		environment: &environment{
			PublicKey: *pKey,
			Name:      name,
			Email:     email,
			EnvType:   envType,
		},
	}, sKey, nil
}

func (env *Environment) Id() string {
	return env.PublicKey.String()
}

func (env *Environment) AddTags(tags ...string) {
	env.Tags = append(env.Tags, tags...)
}

func FromEnvDef(envDef string) (env Environment, err error) {
	if !strings.HasPrefix(envDef, envDefPrefix) {
		return
	}
	serializedEnvString := strings.TrimPrefix(envDef, envDefPrefix)
	err = commons.Deserialize(serializedEnvString, &env)
	return
}

func (env *Environment) ToEnvDef() (string, error) {
	data, err := commons.Serialize(env)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", envDefPrefix, data), nil
}

func (env Environment) MarshalYAML() (interface{}, error) {
	return env.environment, nil
}

func (env *Environment) UnmarshalYAML(value *yaml.Node) (err error) {
	return value.Decode(&env.environment)
}

func (env *Environment) UnmarshalJSON(data []byte) (err error) {
	var environment *environment = new(environment)
	err = json.Unmarshal(data, environment)
	if err == nil {
		env.environment = environment
	}
	return
}

func (env *Environment) Search(query string) bool {
	return commons.SearchStruct(env, query)
}

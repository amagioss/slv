package environment

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shibme/slv/commons"
	"github.com/shibme/slv/crypto"
	"gopkg.in/yaml.v3"
)

type EnvType string

type Environment struct {
	*environment
}

type environment struct {
	PublicKey crypto.PublicKey `yaml:"public_key"`
	Name      string           `yaml:"name"`
	Email     string           `yaml:"email"`
	EnvType   EnvType          `yaml:"type"`
	Tags      []string         `yaml:"tags"`
}

func (eType *EnvType) isValid() bool {
	return *eType == SERVICE || *eType == USER || *eType == ROOT
}

func New(name, email string, envType EnvType) (env *Environment, privateKey *crypto.PrivateKey, err error) {
	if !envType.isValid() {
		return nil, nil, ErrInvalidEnvironmentType
	}
	kp, err := crypto.NewKeyPair(EnvironmentKey)
	if err != nil {
		return nil, nil, err
	}
	privateKey = new(crypto.PrivateKey)
	env = &Environment{
		environment: &environment{
			PublicKey: kp.PublicKey(),
			Name:      name,
			Email:     email,
			EnvType:   envType,
		},
	}
	*privateKey = kp.PrivateKey()
	return
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

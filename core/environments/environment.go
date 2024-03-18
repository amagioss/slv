package environments

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"savesecrets.org/slv/core/commons"
	"savesecrets.org/slv/core/config"
	"savesecrets.org/slv/core/crypto"
)

type EnvType string

type Environment struct {
	*environment
}

type environment struct {
	PublicKey     string   `yaml:"publicKey"`
	Name          string   `yaml:"name"`
	Email         string   `yaml:"email"`
	EnvType       EnvType  `yaml:"type"`
	Tags          []string `yaml:"tags"`
	SecretBinding string   `yaml:"binding,omitempty"`
	publicKey     *crypto.PublicKey
}

func (eType *EnvType) isValid() bool {
	return *eType == SERVICE || *eType == USER || *eType == ROOT
}

func NewEnvironmentForPublicKey(name string, envType EnvType, publicKey *crypto.PublicKey) (*Environment, error) {
	if !envType.isValid() {
		return nil, errInvalidEnvironmentType
	}
	publicKeyStr, err := publicKey.String()
	if err != nil {
		return nil, err
	}
	return &Environment{
		environment: &environment{
			PublicKey: publicKeyStr,
			Name:      name,
			EnvType:   envType,
		},
	}, nil
}

func NewEnvironment(name string, envType EnvType, pq bool) (*Environment, *crypto.SecretKey, error) {
	secretKey, err := crypto.NewSecretKey(EnvironmentKey)
	if err == nil {
		publicKey, err := secretKey.PublicKey(pq)
		if err != nil {
			return nil, nil, err
		}
		env, err := NewEnvironmentForPublicKey(name, envType, publicKey)
		return env, secretKey, err
	}
	return nil, nil, err
}

func (env *Environment) getPublicKey() (publicKey *crypto.PublicKey, err error) {
	if env.publicKey == nil {
		if env.PublicKey == "" {
			return nil, errEnvironmentPublicKeyNotFound
		}
		publicKey, err = crypto.PublicKeyFromString(env.PublicKey)
		if err == nil {
			env.publicKey = publicKey
		}
	}
	return env.publicKey, nil
}

func (env *Environment) Id() string {
	return env.PublicKey
}

func (env *Environment) SetEmail(email string) {
	env.Email = email
}

func (env *Environment) AddTags(tags ...string) {
	env.Tags = append(env.Tags, tags...)
}

func FromEnvDef(envDef string) (env *Environment, err error) {
	sliced := strings.Split(envDef, "_")
	if len(sliced) != 3 || sliced[0] != slvPrefix || sliced[1] != envDefStringAbbrev {
		return nil, errInvalidEnvDef
	}
	err = commons.Deserialize(sliced[2], &env)
	return
}

func (env *Environment) ToEnvDef() (string, error) {
	data, err := commons.Serialize(env)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s_%s", slvPrefix, envDefStringAbbrev, data), nil
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
	return strings.Contains(strings.ToLower(fmt.Sprintf("%s\n%s\n%s\n%s", env.Name, env.Email,
		env.EnvType, strings.Join(env.Tags, "\n"))),
		strings.ToLower(query))
}

func GetSelf() *Environment {
	selfEnvFilePath := filepath.Join(config.GetAppDataDir(), selfEnvFileName)
	env := &Environment{
		environment: new(environment),
	}
	if err := commons.ReadFromYAML(selfEnvFilePath, env.environment); err != nil {
		return nil
	}
	return env
}

func (env *Environment) MarkAsSelf() error {
	if env.SecretBinding == "" {
		return errMarkingSelfEnvBindingNotFound
	}
	if env.EnvType != USER {
		return errMarkingSelfNonUserEnv
	}
	selfEnvFilePath := filepath.Join(config.GetAppDataDir(), selfEnvFileName)
	return commons.WriteToYAML(selfEnvFilePath, "", env)
}

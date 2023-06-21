package environment

import (
	"fmt"
	"strings"

	"github.com/shibme/slv/slv/commons"
	"github.com/shibme/slv/slv/crypto"
)

type EnvType string

type Environment struct {
	data       EnvData
	searchData map[string]struct{}
}

type EnvData struct {
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
	env = &Environment{}
	kp, err := crypto.NewKeyPair(EnvironmentKey)
	if err != nil {
		return nil, nil, err
	}
	privateKey = new(crypto.PrivateKey)
	*privateKey = kp.PrivateKey()
	env.data.PublicKey = kp.PublicKey()
	env.data.Name = name
	env.data.Email = email
	env.data.EnvType = envType
	return env, privateKey, nil
}

func (env *Environment) Id() string {
	return env.PublicKey().String()
}

func (env *Environment) PublicKey() crypto.PublicKey {
	return env.data.PublicKey
}

func (env *Environment) Type() EnvType {
	return env.data.EnvType
}

func (env *Environment) Name() string {
	return env.data.Name
}

func (env *Environment) Email() string {
	return env.data.Email
}

func (env *Environment) AddTags(tags ...string) {
	env.data.Tags = append(env.data.Tags, tags...)
}

func (env *Environment) Tags() (tags []string) {
	return env.data.Tags
}

func EnvFromSLVFormat(slvFormatEnv string) (env Environment, err error) {
	if !strings.HasPrefix(slvFormatEnv, slvFormatEnvironmentPrefix) {
		return
	}
	serializedEnvString := strings.TrimPrefix(slvFormatEnv, slvFormatEnvironmentPrefix)
	data := commons.Decode(serializedEnvString)
	err = commons.Deserialize(data, &env)
	return
}

func (env *Environment) ToSLVFormat() (string, error) {
	data, err := commons.Serialize(env.data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", slvFormatEnvironmentPrefix, commons.Encode(data)), nil
}

func (env *Environment) Search(query string) bool {
	if env.searchData == nil || len(env.searchData) == 0 {
		env.searchData = make(map[string]struct{})
		env.searchData[strings.ToLower(env.data.Name)] = struct{}{}
		env.searchData[strings.ToLower(env.data.Email)] = struct{}{}
		for _, tag := range env.data.Tags {
			env.searchData[strings.ToLower(tag)] = struct{}{}
		}
	}
	for searchStr := range env.searchData {
		if strings.Contains(searchStr, query) {
			return true
		}
	}
	return false
}

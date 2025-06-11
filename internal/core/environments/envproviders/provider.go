package envproviders

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"slv.sh/slv/internal/core/commons"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments"
)

type bind func(skBytes []byte, inputs map[string]string) (ref map[string][]byte, err error)
type unbind func(ref map[string][]byte) (secretKeyBytes []byte, err error)

const (
	envSecretBindingAbbrev = "ESB"
)

var (
	providerMap         = make(map[string]*provider)
	providerInitializer sync.Once

	errInvalidEnvSecretBindingFormat = errors.New("invalid environment secret binding format")
	errEnvSecretBindingUnspecified   = errors.New("environment secret binding unspecified")
)

type provider struct {
	id          string
	name        string
	desc        string
	bind        bind
	unbind      unbind
	refRequired bool
	args        []arg
}

type envSecretBinding struct {
	Provider string            `json:"p"`
	Ref      map[string][]byte `json:"r"`
}

func (esb *envSecretBinding) string() (string, error) {
	data, err := commons.Serialize(*esb)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s_%s", config.AppNameUpperCase, envSecretBindingAbbrev, data), nil
}

func envSecretBindingFromString(envSecretBindingStr string) (*envSecretBinding, error) {
	sliced := strings.Split(envSecretBindingStr, "_")
	if len(sliced) != 3 || sliced[0] != config.AppNameUpperCase || sliced[1] != envSecretBindingAbbrev {
		return nil, errInvalidEnvSecretBindingFormat
	}
	binding := new(envSecretBinding)
	if err := commons.Deserialize(sliced[2], &binding); err != nil {
		return nil, err
	}
	return binding, nil
}

func Register(id, name, desc string, bind bind, unbind unbind, refRequired bool, args []arg) {
	providerMap[id] = &provider{
		id:          id,
		name:        name,
		desc:        desc,
		bind:        bind,
		unbind:      unbind,
		refRequired: refRequired,
		args:        args,
	}
}

func registerDefaultProviders() {
	providerInitializer.Do(func() {
		Register(PasswordProviderId, passwordProviderName, passwordProviderDesc, bindWithPassword, unBindWithPassword, true, pwdArgs)
		Register(awsProviderId, awsProviderName, awsProviderDesc, bindWithAWSKMS, unBindFromAWSKMS, true, awsArgs)
		Register(gcpProviderId, gcpProviderName, gcpProviderDesc, bindWithGCP, unBindWithGCP, true, gcpArgs)
		Register(azureProviderId, azureProviderName, azureProviderDesc, bindWithAzure, unBindFromAzure, true, azureArgs)
	})
}

func ListIds() []string {
	registerDefaultProviders()
	providerIds := make([]string, 0, len(providerMap))
	for providerId := range providerMap {
		providerIds = append(providerIds, providerId)
	}
	return providerIds
}

func GetName(providerId string) string {
	registerDefaultProviders()
	if provider, ok := providerMap[providerId]; ok {
		return provider.name
	}
	return ""
}

func GetDesc(providerId string) string {
	registerDefaultProviders()
	if provider, ok := providerMap[providerId]; ok {
		return provider.desc
	}
	return ""
}

func GetArgs(providerId string) []arg {
	registerDefaultProviders()
	if provider, ok := providerMap[providerId]; ok {
		return provider.args
	}
	return nil
}

func NewEnv(providerId, envName string, envType environments.EnvType,
	inputs map[string]string, quantumSafe bool) (*environments.Environment, error) {
	registerDefaultProviders()
	provider, ok := providerMap[providerId]
	if !ok {
		return nil, fmt.Errorf("unknown environment provider: %s", providerId)
	}
	for _, arg := range provider.args {
		if arg.required && inputs[arg.id] == "" {
			return nil, fmt.Errorf("missing required input: %s", arg.id)
		}
	}
	env, sk, err := environments.New(envName, envType, quantumSafe)
	if err != nil {
		return nil, err
	}
	skBytes, err := sk.Bytes()
	if err != nil {
		return nil, err
	}
	ref, err := (provider.bind)(skBytes, inputs)
	if err != nil {
		return nil, err
	}
	if provider.refRequired {
		esb := &envSecretBinding{
			Provider: providerId,
			Ref:      ref,
		}
		if env.SecretBinding, err = esb.string(); err != nil {
			return nil, err
		}
	}
	return env, nil
}

func getSecretKeyFromBytesForBinding(skBytes []byte) (secretKey *crypto.SecretKey, err error) {
	if secretKey, err = crypto.SecretKeyFromBytes(skBytes); err == nil {
		secretKey.RestrictSerialization()
	}
	return
}

func GetSecretKeyFromSecretBinding(envSecretBindingStr string) (secretKey *crypto.SecretKey, err error) {
	registerDefaultProviders()
	var secretKeyBytes []byte
	var esb *envSecretBinding
	if envSecretBindingStr == "" {
		var providersWithoutRef []provider
		for _, provider := range providerMap {
			if !provider.refRequired {
				providersWithoutRef = append(providersWithoutRef, *provider)
			}
		}
		for _, provider := range providersWithoutRef {
			if secretKeyBytes, err = (provider.unbind)(nil); err == nil {
				return getSecretKeyFromBytesForBinding(secretKeyBytes)
			}
		}
		return nil, errEnvSecretBindingUnspecified
	} else if esb, err = envSecretBindingFromString(envSecretBindingStr); err == nil {
		if provider, ok := providerMap[esb.Provider]; !ok {
			return nil, fmt.Errorf("unknown environment provider: %s", esb.Provider)
		} else if secretKeyBytes, err = (provider.unbind)(esb.Ref); err == nil {
			return getSecretKeyFromBytesForBinding(secretKeyBytes)
		}
	}
	return nil, err
}

type arg struct {
	id          string
	name        string
	required    bool
	description string
}

func (a *arg) Id() string {
	return a.id
}

func (a *arg) Name() string {
	return a.name
}

func (a *arg) Required() bool {
	return a.required
}

func (a *arg) Description() string {
	return a.description
}

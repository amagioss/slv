package envproviders

import (
	"fmt"
	"strings"
	"sync"

	"slv.sh/slv/internal/core/commons"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments"
)

type bind func(skBytes []byte, inputs map[string][]byte) (ref map[string][]byte, err error)
type unbind func(ref map[string][]byte) (secretKeyBytes []byte, err error)

var (
	providerMap         = make(map[string]*provider)
	providerInitializer sync.Once
)

type provider struct {
	name        string
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
	return fmt.Sprintf("%s_%s_%s", slvPrefix, envSecretBindingAbbrev, data), nil
}

func envSecretBindingFromString(envSecretBindingStr string) (*envSecretBinding, error) {
	sliced := strings.Split(envSecretBindingStr, "_")
	if len(sliced) != 3 || sliced[0] != slvPrefix || sliced[1] != envSecretBindingAbbrev {
		return nil, errInvalidEnvSecretBindingFormat
	}
	binding := new(envSecretBinding)
	if err := commons.Deserialize(sliced[2], &binding); err != nil {
		return nil, err
	}
	return binding, nil
}

func Register(name string, bind bind, unbind unbind, refRequired bool) {
	providerMap[name] = &provider{
		name:        name,
		bind:        bind,
		unbind:      unbind,
		refRequired: refRequired,
	}
}

func registerDefaultProviders() {
	providerInitializer.Do(func() {
		Register(passwordProviderName, bindWithPassword, unBindWithPassword, true)
		Register(awsProviderName, bindWithAWSKMS, unBindFromAWSKMS, true)
		Register(gcpProviderName, bindWithGCP, unBindWithGCP, true)
	})
}

func ListNames() []string {
	registerDefaultProviders()
	providerNames := make([]string, 0, len(providerMap))
	for name := range providerMap {
		providerNames = append(providerNames, name)
	}
	return providerNames
}

func GetArgs(providerName string) []arg {
	if provider, ok := providerMap[providerName]; ok {
		return provider.args
	}
	return nil
}

func NewEnv(providerName, envName string, envType environments.EnvType,
	inputs map[string][]byte, quantumSafe bool) (*environments.Environment, error) {
	registerDefaultProviders()
	provider, ok := providerMap[providerName]
	if !ok {
		return nil, errProviderUnknown
	}
	env, sk, err := environments.NewEnvironment(envName, envType, quantumSafe)
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
			Provider: providerName,
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
			return nil, errProviderUnknown
		} else if secretKeyBytes, err = (provider.unbind)(esb.Ref); err == nil {
			return getSecretKeyFromBytesForBinding(secretKeyBytes)
		}
	}
	return nil, err
}

type arg struct {
	name        string
	required    bool
	sensitive   bool
	description string
}

func (a *arg) Name() string {
	return a.name
}

func (a *arg) Required() bool {
	return a.required
}

func (a *arg) Sensitive() bool {
	return a.sensitive
}

func (a *arg) Description() string {
	return a.description
}

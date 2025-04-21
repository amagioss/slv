package providers

import (
	"fmt"
	"strings"

	"slv.sh/slv/internal/core/commons"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments"
)

type Bind func(skBytes []byte, inputs map[string][]byte) (ref map[string][]byte, err error)
type UnBind func(ref map[string][]byte) (secretKeyBytes []byte, err error)

var providerMap = make(map[string]*provider)

type provider struct {
	Name        string
	bind        *Bind
	unbind      *UnBind
	refRequired bool
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

func registerProvider(name string, bind Bind, unbind UnBind, refRequired bool) {
	providerMap[name] = &provider{
		Name:        name,
		bind:        &bind,
		unbind:      &unbind,
		refRequired: refRequired,
	}
}

func RegisterEnvSecretProvider(name string, bind Bind, unbind UnBind, refRequired bool) error {
	loadDefaultProviders()
	if _, ok := providerMap[name]; ok {
		return errProviderRegisteredAlready
	}
	registerProvider(name, bind, unbind, refRequired)
	return nil
}

func NewEnvForProvider(providerName, envName string, envType environments.EnvType,
	inputs map[string][]byte, quantumSafe bool) (*environments.Environment, error) {
	loadDefaultProviders()
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
	ref, err := (*provider.bind)(skBytes, inputs)
	if err != nil {
		return nil, err
	}
	if provider.refRequired {
		esb := &envSecretBinding{
			Provider: providerName,
			Ref:      ref,
		}
		env.SecretBinding, err = esb.string()
		if err != nil {
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
	loadDefaultProviders()
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
			if secretKeyBytes, err = (*provider.unbind)(nil); err == nil {
				return getSecretKeyFromBytesForBinding(secretKeyBytes)
			}
		}
		return nil, errEnvSecretBindingUnspecified
	} else if esb, err = envSecretBindingFromString(envSecretBindingStr); err == nil {
		if provider, ok := providerMap[esb.Provider]; !ok {
			return nil, errProviderUnknown
		} else if secretKeyBytes, err = (*provider.unbind)(esb.Ref); err == nil {
			return getSecretKeyFromBytesForBinding(secretKeyBytes)
		}
	}
	return nil, err
}

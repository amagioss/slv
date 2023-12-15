package environments

import (
	"fmt"
	"strings"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
)

type Bind func(inputs map[string][]byte) (publicKey *crypto.PublicKey, ref map[string][]byte, err error)
type UnBind func(ref map[string][]byte) (secretKeyBytes []byte, err error)

var providerMap = make(map[string]*provider)

type provider struct {
	Name        string
	bind        *Bind
	unbind      *UnBind
	refRequired bool
}

type providerAccessBinding struct {
	Provider string            `json:"p"`
	Ref      map[string][]byte `json:"r"`
}

func (pab *providerAccessBinding) string() (string, error) {
	data, err := commons.Serialize(*pab)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s_%s", commons.SLV, providerAccessBindingAbbrev, data), nil
}

func providerAccessBindingFromString(providerAccessBindingStr string) (*providerAccessBinding, error) {
	sliced := strings.Split(providerAccessBindingStr, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV || sliced[1] != providerAccessBindingAbbrev {
		return nil, ErrInvalidProviderAccessBindingFormat
	}
	binding := new(providerAccessBinding)
	if err := commons.Deserialize(sliced[2], &binding); err != nil {
		return nil, err
	}
	return binding, nil
}

func RegisterAccessProvider(name string, bind Bind, unbind UnBind, refRequired bool) {
	providerMap[name] = &provider{
		Name:        name,
		bind:        &bind,
		unbind:      &unbind,
		refRequired: refRequired,
	}
}

func NewEnvForProvider(providerName, envName string, envType EnvType,
	inputs map[string][]byte) (*Environment, error) {
	provider, ok := providerMap[providerName]
	if !ok {
		return nil, ErrProviderUnknown
	}
	publicKey, ref, err := (*provider.bind)(inputs)
	if err != nil {
		return nil, err
	}
	env, err := NewEnvironmentForPublicKey(envName, envType, publicKey)
	if err != nil {
		return nil, err
	}
	if provider.refRequired {
		pab := &providerAccessBinding{
			Provider: providerName,
			Ref:      ref,
		}
		env.ProviderBinding, err = pab.string()
		if err != nil {
			return nil, err
		}
	}
	return env, nil
}

func GetSecretKeyFromAccessBinding(providerAccessBindingStr string) (secretKey *crypto.SecretKey, err error) {
	if providerAccessBindingStr == "" {
		var providersWithoutRef []provider
		for _, provider := range providerMap {
			if !provider.refRequired {
				providersWithoutRef = append(providersWithoutRef, *provider)
			}
		}
		for _, provider := range providersWithoutRef {
			secretKeyBytes, err := (*provider.unbind)(nil)
			if err == nil {
				return crypto.SecretKeyFromBytes(secretKeyBytes)
			}
		}
		return nil, ErrProviderAccessBindingUnspecified
	}
	pab, err := providerAccessBindingFromString(providerAccessBindingStr)
	if err != nil {
		return nil, err
	}
	provider, ok := providerMap[pab.Provider]
	if !ok {
		return nil, ErrProviderUnknown
	}
	secretKeyBytes, err := (*provider.unbind)(pab.Ref)
	if err == nil {
		return crypto.SecretKeyFromBytes(secretKeyBytes)
	}
	return nil, err
}

package providers

import (
	"fmt"
	"strings"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/environments"
)

type Bind func(inputs map[string]any) (publicKey *crypto.PublicKey, ref map[string]any, err error)
type UnBind func(ref map[string]any) (secretKeyBytes []byte, err error)

func NewEnvForProvider(name, email string, envType environments.EnvType, providerName string,
	inputs map[string]any) (*environments.Environment, error) {
	provider, ok := providerMap[providerName]
	if !ok {
		return nil, ErrEnvProviderUnknown
	}
	publicKey, ref, err := (*provider.bind)(inputs)
	if err != nil {
		return nil, err
	}
	env, err := environments.NewEnvironmentForPublicKey(name, email, envType, publicKey)
	if err != nil {
		return nil, err
	}
	if provider.refRequired {
		eab := &envProviderBinding{
			Provider: providerName,
			Ref:      ref,
		}
		env.ProviderBinding, err = eab.string()
		if err != nil {
			return nil, err
		}
	}
	return env, nil
}

func GetSecretKeyFromProvider(envProviderBindingStr string) (secretKey *crypto.SecretKey, err error) {
	if envProviderBindingStr == "" {
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
		return nil, ErrEnvProviderBindingUnspecified
	}
	envAccessBinding, err := providerBindingFromString(envProviderBindingStr)
	if err != nil {
		return nil, err
	}
	provider, ok := providerMap[envAccessBinding.Provider]
	if !ok {
		return nil, ErrEnvProviderUnknown
	}
	secretKeyBytes, err := (*provider.unbind)(envAccessBinding.Ref)
	if err == nil {
		return crypto.SecretKeyFromBytes(secretKeyBytes)
	}
	return nil, err
}

type envProviderBinding struct {
	Provider string         `json:"p"`
	Ref      map[string]any `json:"r"`
}

func (epb *envProviderBinding) string() (string, error) {
	data, err := commons.Serialize(*epb)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s_%s", commons.SLV, envAccessBindingAbbrev, data), nil
}

func providerBindingFromString(envProviderBindingStr string) (eab *envProviderBinding, err error) {
	sliced := strings.Split(envProviderBindingStr, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV || sliced[1] != envAccessBindingAbbrev {
		return nil, ErrEnvProviderBindingInvalid
	}
	if err = commons.Deserialize(sliced[2], eab); err != nil {
		return nil, err
	}
	return
}

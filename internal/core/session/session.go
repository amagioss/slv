package session

import (
	"fmt"
	"os"
	"sync"

	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/environments/envproviders"
	"slv.sh/slv/internal/core/profiles"
)

const (
	envar_SLV_ENV_SECRET_KEY     = "SLV_ENV_SECRET_KEY"
	envar_SLV_ENV_SECRET_BINDING = "SLV_ENV_SECRET_BINDING"
)

var (
	session          *Session
	sessionInitMutex sync.Mutex
	slvK8sSecret     = func() string {
		if val := os.Getenv("SLV_K8S_ENV_SECRET"); val != "" {
			return val
		}
		return config.AppNameLowerCase
	}()
)

type Session struct {
	secretKey   *crypto.SecretKey
	env         *environments.Environment
	pubKeyEC    string
	pubKeyPQ    string
	initialized bool
	initMutex   sync.Mutex
}

func (s *Session) SecretKey() *crypto.SecretKey {
	return s.secretKey
}

func (s *Session) Env() (*environments.Environment, error) {
	if !s.initialized && s.secretKey != nil {
		s.initMutex.Lock()
		defer s.initMutex.Unlock()
		if !s.initialized {
			profile, err := profiles.GetActiveProfile()
			if err != nil {
				return nil, err
			}
			if s.env == nil {
				if env, err := profile.GetEnv(s.pubKeyEC); err == nil {
					s.env = env
				}
				if s.env == nil {
					if env, err := profile.GetEnv(s.pubKeyPQ); err == nil {
						s.env = env
					}
				}
			}
			s.initialized = true
		}
	}
	return s.env, nil
}

func (s *Session) PublicKeyEC() string {
	if s.pubKeyEC == "" {
		if s.secretKey != nil {
			if pkEC, _ := s.secretKey.PublicKey(false); pkEC != nil {
				s.pubKeyEC, _ = pkEC.String()
			}
		}
	}
	return s.pubKeyEC
}

func (s *Session) PublicKeyPQ() string {
	if s.pubKeyPQ == "" {
		if s.secretKey != nil {
			if pkPQ, _ := s.secretKey.PublicKey(true); pkPQ != nil {
				s.pubKeyPQ, _ = pkPQ.String()
			}
		}
	}
	return s.pubKeyPQ
}

func GetSession() (*Session, error) {
	if session == nil {
		sessionInitMutex.Lock()
		defer sessionInitMutex.Unlock()
		if session == nil {
			session = &Session{}
			var err error
			if secretKeyStr := os.Getenv(envar_SLV_ENV_SECRET_KEY); secretKeyStr != "" {
				if session.secretKey, err = crypto.SecretKeyFromString(secretKeyStr); err != nil {
					session = nil
					return nil, err
				}
			}
			if session.secretKey == nil {
				envSecretBindingStr := os.Getenv(envar_SLV_ENV_SECRET_BINDING)
				if envSecretBindingStr == "" {
					if selfEnv := environments.GetSelf(); selfEnv != nil {
						session.env = selfEnv
						session.initialized = true
						envSecretBindingStr = selfEnv.SecretBinding
					}
				}
				if envSecretBindingStr != "" {
					if session.secretKey, err = envproviders.GetSecretKeyFromSecretBinding(envSecretBindingStr); err != nil {
						session = nil
						return nil, err
					}
				}
			}
			if session.secretKey == nil && isInKubernetesCluster() {
				if kubeClientSet, _ := getKubeClientSet(); kubeClientSet != nil {
					session.secretKey, _ = getSecretKeyFor(kubeClientSet, GetK8sNamespace())
				}
			}
			if isInKubernetesCluster() {
				if kubeClientSet, _ := getKubeClientSet(); kubeClientSet != nil && (session.PublicKeyEC() != "" || session.PublicKeyPQ() != "") {
					putPublicKeyToConfigMap(kubeClientSet, session.PublicKeyEC(), session.PublicKeyPQ())
				}
			}
		}
	}
	return session, nil
}

func GetSecretKey() (*crypto.SecretKey, error) {
	if session, err := GetSession(); err != nil {
		return nil, err
	} else if session.SecretKey() == nil {
		return nil, fmt.Errorf("environment access not configured")
	} else {
		return session.SecretKey(), nil
	}
}

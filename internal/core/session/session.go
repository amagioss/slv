package session

import (
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
	return s.pubKeyEC
}

func (s *Session) PublicKeyPQ() string {
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
			if session.secretKey == nil {
				if kubeClientSet, _ := getKubeClientSet(); kubeClientSet != nil {
					session.secretKey, _ = getSecretKeyFor(kubeClientSet, GetK8sNamespace())
				}
			}
			if session.secretKey != nil {
				var pkEC, pkPQ *crypto.PublicKey
				if pkEC, _ = session.secretKey.PublicKey(false); pkEC != nil {
					session.pubKeyEC, _ = pkEC.String()
				}
				if pkPQ, _ = session.secretKey.PublicKey(true); pkPQ != nil {
					session.pubKeyPQ, _ = pkPQ.String()
				}
			}
			if kubeClientSet, _ := getKubeClientSet(); kubeClientSet != nil && (session.pubKeyEC != "" || session.pubKeyPQ != "") {
				putPublicKeyToConfigMap(kubeClientSet, session.pubKeyEC, session.pubKeyPQ)
			}
		}
	}
	return session, nil
}

func GetSecretKey() (*crypto.SecretKey, error) {
	if session, err := GetSession(); err == nil {
		return session.SecretKey(), nil
	} else {
		return nil, err
	}
}

package session

import (
	"sync"

	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/environments/envproviders"
	"slv.sh/slv/internal/core/profiles"
)

var (
	session          *Session
	sessionInitMutex sync.Mutex
)

type Session struct {
	secretKey      *crypto.SecretKey
	env            *environments.Environment
	envInitialized bool
}

func (s *Session) SecretKey() *crypto.SecretKey {
	return s.secretKey
}

func (s *Session) Env() *environments.Environment {
	if !s.envInitialized && s.secretKey != nil {
		profile, err := profiles.GetActiveProfile()
		if err != nil {
			return nil
		}
		if pubKeyEC, err := s.secretKey.PublicKey(false); err == nil && pubKeyEC != nil {
			if pubKeyStr, _ := pubKeyEC.String(); pubKeyStr != "" {
				if env, err := profile.GetEnv(pubKeyStr); err == nil {
					s.env = env
				}
			}
		}
		if s.env == nil {
			if pubKeyPQ, err := s.secretKey.PublicKey(true); err == nil && pubKeyPQ != nil {
				if pubKeyStr, _ := pubKeyPQ.String(); pubKeyStr != "" {
					if env, err := profile.GetEnv(pubKeyStr); err == nil {
						s.env = env
					}
				}
			}
		}
		s.envInitialized = true
	}
	return s.env
}

func GetSession() (*Session, error) {
	if session == nil {
		sessionInitMutex.Lock()
		defer sessionInitMutex.Unlock()
		if session == nil {
			session = &Session{}
			var err error
			if secretKeyStr := config.GetEnvSecretKey(); secretKeyStr != "" {
				if session.secretKey, err = crypto.SecretKeyFromString(secretKeyStr); err == nil {
					return session, nil
				} else {
					session = nil
					return nil, err
				}
			}
			envSecretBindingStr := config.GetEnvSecretBinding()
			if envSecretBindingStr == "" {
				if selfEnv := environments.GetSelf(); selfEnv != nil {
					session.env = selfEnv
					session.envInitialized = true
					envSecretBindingStr = selfEnv.SecretBinding
				}
			}
			if envSecretBindingStr != "" {
				if session.secretKey, err = envproviders.GetSecretKeyFromSecretBinding(envSecretBindingStr); err == nil {
					return session, nil
				} else {
					session = nil
					return nil, err
				}
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

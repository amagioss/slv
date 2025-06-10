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
				var pubKeyEC, pubKeyPQ *crypto.PublicKey
				if pubKeyEC, err = s.secretKey.PublicKey(false); err == nil && pubKeyEC != nil {
					if s.pubKeyEC, _ = pubKeyEC.String(); s.pubKeyEC != "" {
						if env, err := profile.GetEnv(s.pubKeyEC); err == nil {
							s.env = env
						}
					}
				}
				if s.env == nil {
					if pubKeyPQ, err = s.secretKey.PublicKey(true); err == nil && pubKeyPQ != nil {
						if s.pubKeyPQ, _ = pubKeyPQ.String(); s.pubKeyPQ != "" {
							if env, err := profile.GetEnv(s.pubKeyPQ); err == nil {
								s.env = env
							}
						}
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
		s.Env()
	}
	return s.pubKeyEC
}

func (s *Session) PublicKeyPQ() string {
	if s.pubKeyPQ == "" {
		s.Env()
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
					session.initialized = true
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

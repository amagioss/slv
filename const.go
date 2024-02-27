package slv

import (
	"errors"

	"savesecrets.org/slv/core/crypto"
)

var (
	Version                      = "dev"
	Commit                       = "none"
	BuildDate                    = ""
	secretKey                    *crypto.SecretKey
	errEnvironmentAccessNotFound = errors.New("environment doesn't have access. please set the required environment variables")
)

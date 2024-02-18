package slv

import (
	"errors"

	"savesecrets.org/slv/core/crypto"
)

const (
	// AppName = config.AppName
	Prefix = "SLV"
)

var (
	Version                      = "dev"
	secretKey                    *crypto.SecretKey
	errEnvironmentAccessNotFound = errors.New("environment doesn't have access. please set the required environment variables")
)
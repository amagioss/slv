package secretkeystore

import "errors"

const (
	slvSecreKeyEnvarName      = "SLV_SECRET_KEY"
	slvAccessBindingEnvarName = "SLV_ACCESS_BINDING"
)

var (
	ErrInvalidEnvAccessProvider  = errors.New("invalid environment access provider")
	ErrEnvironmentAccessNotFound = errors.New("environment doesn't have access. please set one of the environment variables: " + slvSecreKeyEnvarName + " or " + slvAccessBindingEnvarName)
)

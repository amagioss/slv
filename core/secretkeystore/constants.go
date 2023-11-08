package secretkeystore

import "errors"

const (
	slvSecreKeyEnvarName        = "SLV_SECRET_KEY"
	slvProviderContextEnvarName = "SLV_PROVIDER_CONTEXT"
)

var ErrInvalidEnvProviderType = errors.New("invalid environment provider type")
var ErrEnvironmentAccessNotFound = errors.New("environment doesn't have access. please set one of the environment variables: " + slvSecreKeyEnvarName + " or " + slvProviderContextEnvarName)

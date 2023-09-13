package keystore

import "errors"

const (
	slvSecreKeyEnvarName   = "SLV_SECRET_KEY"
	slvPassphraseEnvarName = "SLV_USER_PASSWORD"
	slvKeyringServiceName  = "slv"
	slvKeyringItemKey      = "slv_password"
)

var ErrEnvSecretNotSet = errors.New(slvSecreKeyEnvarName + " environment variable is not set")

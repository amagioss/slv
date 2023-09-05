package keyreader

import "errors"

var ErrEnvSecretNotSet = errors.New("SLV_SECRET_KEY environment variable is not set")

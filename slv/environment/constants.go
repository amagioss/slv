package environment

import (
	"errors"

	"github.com/shibme/slv/slv/crypto"
)

const (
	slvFormatEnvironmentPrefix                = "SLV_ED_" // Environment Definition
	EnvironmentKey             crypto.KeyType = 'E'
	USER                       EnvType        = "user"
	SERVICE                    EnvType        = "service"
	ROOT                       EnvType        = "root"
)

var ErrInvalidEnvironmentType = errors.New("invalid environment type")

var ErrProcessingEnvironmentsManifest = errors.New("error in processing environments manifest")
var ErrEnvironmentNotFound = errors.New("no such environment exists")

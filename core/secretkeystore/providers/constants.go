package providers

import "errors"

const (
	envAccessBindingAbbrev = "EAB" // Environment Access Binding
)

var ErrEnvProviderUnknown = errors.New("unknown environment provider")
var ErrEnvProviderBindingInvalid = errors.New("invalid environment provider binding")
var ErrEnvProviderBindingUnspecified = errors.New("unspecified environment provider binding")

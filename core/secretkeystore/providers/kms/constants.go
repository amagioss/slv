package kms

import "errors"

var ErrInvalidRSAPublicKey = errors.New("invalid RSA public key")
var ErrMissingProviderBindingRef = errors.New("missing provider binding ref")

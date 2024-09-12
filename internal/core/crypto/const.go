package crypto

import (
	"errors"

	"oss.amagi.com/slv/internal/core/config"
)

const (
	publicKeyAbbrev          = "PK" // PK = Public Key
	secretKeyAbbrev          = "SK" // SK = Secret Key
	wrappedKeyAbbrev         = "WK" // WK = Wrapped Key
	sealedSecretAbbrev       = "SS" // SS = Sealed Secret
	slvPrefix                = config.AppNameUpperCase
	cryptoVersion      uint8 = 1

	hashMaxLength = 4
)

var (
	errUnsupportedCryptoVersion = errors.New("unsupported cryptography version")
	errGeneratingSecretKey      = errors.New("error generating a new secret key")
	errDerivingPublicKey        = errors.New("error deriving public key from the secret key")
	errInvalidPublicKeyFormat   = errors.New("invalid public key format")
	errInvalidSecretKeyFormat   = errors.New("invalid secret key format")
	errEncryptionFailed         = errors.New("encryption failed")
	errDecryptionFailed         = errors.New("decryption failed")
	errSecretKeyMismatch        = errors.New("given secret key cannot decrypt the data")
	errInvalidCiphertextFormat  = errors.New("invalid ciphertext format")
)

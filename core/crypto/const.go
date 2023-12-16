package crypto

import (
	"errors"

	"gopkg.shib.me/gociphers/ecc"
)

const (
	publicKeyAbbrev    = "PK" // PK = Public Key
	secretKeyAbbrev    = "SK" // SK = Secret Key
	wrappedKeyAbbrev   = "WK" // WK = Wrapped Key
	sealedSecretAbbrev = "SS" // SS = Sealed Secret

	keyLength           = ecc.KeyLength + 3
	cipherTextMinLength = ecc.CipherTextMinLength + keyLength + 2

	argon2HashMaxLength = 4
)

var (
	ErrUnsupportedCryptoVersion = errors.New("unsupported cryptography version")
	ErrGeneratingSecretKey      = errors.New("error generating a new secret key")
	ErrDerivingPublicKey        = errors.New("error deriving public key from the secret key")
	ErrInvalidPublicKeyFormat   = errors.New("invalid public key format")
	ErrInvalidSecretKeyFormat   = errors.New("invalid secret key format")
	ErrEncryptionFailed         = errors.New("encryption failed")
	ErrDecryptionFailed         = errors.New("decryption failed")
	ErrSecretKeyMismatch        = errors.New("given secret key cannot decrypt the data")
	ErrInvalidCiphertextFormat  = errors.New("invalid ciphertext format")
)

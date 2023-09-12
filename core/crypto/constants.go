package crypto

import (
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

const (
	publicKeyAbbrev             = "PK"  // PK = Public Key
	secretKeyAbbrev             = "SK"  // SK = Secret Key
	wrappedKeyAbbrev            = "WK"  // WK = Wrapped Key
	sealedSecretAbbrev          = "SS"  // SS = Sealed Secret
	passwordProtectedDataAbbrev = "PPD" // PP = Password Protected Data
	passwordProtectedKeyAbbrev  = "PPK" // PP = Password Protected Key

	minimumPasswordLength = 8
	shortKeyIdLength      = 4
	keyBaseLength         = curve25519.ScalarSize + 3
	cipherBytesMinLength  = chacha20poly1305.NonceSize + curve25519.ScalarSize + shortKeyIdLength + 3

	argon2SaltLength    uint8 = 16
	argon2Iterations          = 12
	argon2Memory              = 16 * 1024
	argon2Threads             = 1
	argon2HashMaxLength       = 4
)

var ErrUnsupportedCryptoVersion = errors.New("unsupported cryptographic version")
var ErrGeneratingKey = errors.New("error generating a new key")
var ErrInvalidKeyFormat = errors.New("invalid key format")
var ErrDecryptionFailed = errors.New("decryption failed")
var ErrSecretKeyMismatch = errors.New("given secret key cannot decrypt the data")
var ErrKeyTypeMismatch = errors.New("given key type cannot decrypt the data")
var ErrInvalidCiphertextFormat = errors.New("invalid ciphertext format")

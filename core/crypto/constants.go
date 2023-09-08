package crypto

import (
	"errors"
)

const (
	publicKeyAbbrev    = "PK" // PK = Public Key
	secretKeyAbbrev    = "SK" // SK = Secret Key
	wrappedKeyAbbrev   = "WK" // WK = Wrapped Key
	sealedSecretAbbrev = "SS" // SS = Sealed Secret

	keyIdLength     = 8
	keyLength       = 32
	publicKeyLength = keyLength + 1
	secretKeyLength = publicKeyLength + keyLength
	nonceLength     = 24

	secretHashTime      = 3
	secretHashMemory    = 16 * 1024
	secretHashThreads   = 1
	secretHashMaxLength = 4
)

var ErrGeneratingKeyPair = errors.New("error generating a new key pair")
var ErrInvalidKeyFormat = errors.New("invalid key format")
var ErrDecryptionFailed = errors.New("decryption failed")
var ErrSealedKeyFormat = errors.New("invalid wrapped key format")
var ErrSealedDataFormat = errors.New("invalid sealed secret format")
var ErrSecretKeyMismatch = errors.New("given secret key cannot decrypt the data")
var ErrSecretKeyTypeMismatch = errors.New("given key type cannot decrypt the data")
var ErrInvalidPassphraseEncryptedData = errors.New("invalid passphrase encrypted data")
var ErrCiphertextFormat = errors.New("invalid ciphertext format")
var ErrSymmetricCiphertextShort = errors.New("symmetric ciphertext too short")
var ErrInvalidSymmetricCipherdataFormat = errors.New("invalid symmetric cipherdata format")

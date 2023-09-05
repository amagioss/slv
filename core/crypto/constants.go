package crypto

import "errors"

const (
	publicKeyPrefix           = "SLV_PK_" // PK = Public Key
	privateKeyPrefix          = "SLV_SK_" // SK = Secret Key
	sealedKeyPrefix           = "SLV_WK_" // WK = Wrapped Key
	sealedDataPrefix          = "SLV_SS_" // SS = Sealed Secret
	passphraseProtectedPrefix = "SLV_PP_" // PP = Password Protected
)

var ErrGeneratingKeyPair = errors.New("error generating a new key pair")
var ErrInvalidKeyFormat = errors.New("invalid key format")
var ErrDecryptionFailed = errors.New("decryption failed")
var ErrSealedKeyFormat = errors.New("invalid wrapped key format")
var ErrSealedDataFormat = errors.New("invalid sealed secret format")
var ErrAccessKeyMismatch = errors.New("given access key cannot decrypt the payload")
var ErrInvalidPassphraseEncryptedData = errors.New("invalid passphrase encrypted data")

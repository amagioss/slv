package crypto

import (
	"strings"

	"dev.shib.me/xipher"
	"oss.amagi.com/slv/internal/core/commons"
)

type KeyType byte

type PublicKey struct {
	version *uint8
	keyType *KeyType
	pubKey  *xipher.PublicKey
}

func (publicKey *PublicKey) toBytes() ([]byte, error) {
	pubKeyBytes, err := publicKey.pubKey.Bytes()
	if err != nil {
		return nil, err
	}
	return append([]byte{*publicKey.version, 1, byte(*publicKey.keyType)}, pubKeyBytes...), nil
}

func (publicKey *PublicKey) Type() KeyType {
	return *publicKey.keyType
}

func publicKeyFromBytes(bytes []byte) (*PublicKey, error) {
	if bytes[1] != 1 {
		return nil, errInvalidPublicKeyFormat
	}
	if bytes[0] > cryptoVersion {
		return nil, errUnsupportedCryptoVersion
	}
	var version uint8 = bytes[0]
	var keyType KeyType = KeyType(bytes[2])
	pubKey, err := xipher.ParsePublicKey(bytes[3:])
	if err != nil {
		return nil, errInvalidPublicKeyFormat
	}
	return &PublicKey{
		version: &version,
		keyType: &keyType,
		pubKey:  pubKey,
	}, nil
}

func (publicKey PublicKey) String() (string, error) {
	publicKeyBytes, err := publicKey.toBytes()
	if err != nil {
		return "", err
	}
	return slvPrefix + "_" + string(*publicKey.keyType) + publicKeyAbbrev + "_" + commons.Encode(publicKeyBytes), err
}

func PublicKeyFromString(publicKeyStr string) (*PublicKey, error) {
	sliced := strings.Split(publicKeyStr, "_")
	if len(sliced) != 3 || sliced[0] != slvPrefix {
		return nil, errInvalidPublicKeyFormat
	}
	decoded, err := commons.Decode(sliced[2])
	if err != nil {
		return nil, errInvalidPublicKeyFormat
	}
	publicKey, err := publicKeyFromBytes(decoded)
	if err != nil {
		return nil, err
	}
	if len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*publicKey.keyType)) ||
		!strings.HasSuffix(sliced[1], publicKeyAbbrev) {
		return nil, errInvalidPublicKeyFormat
	}
	return publicKey, nil
}

type SecretKey struct {
	version      *uint8
	keyType      *KeyType
	privKey      *xipher.SecretKey
	pqPublicKey  *PublicKey
	eccPublicKey *PublicKey
	restricted   bool
}

func NewSecretKey(keyType KeyType) (secretKey *SecretKey, err error) {
	privKey, err := xipher.NewSecretKey()
	if err != nil {
		return nil, errGeneratingSecretKey
	}
	return newSecretKey(privKey, keyType, false), nil
}

func NewSecretKeyForPassword(password []byte, keyType KeyType) (secretKey *SecretKey, err error) {
	privKey, err := xipher.NewSecretKeyForPassword(password)
	if err != nil {
		return nil, err
	}
	return newSecretKey(privKey, keyType, true), nil
}

func newSecretKey(privKey *xipher.SecretKey, keyType KeyType, restricted bool) *SecretKey {
	version := cryptoVersion
	return &SecretKey{
		version:    &version,
		keyType:    &keyType,
		privKey:    privKey,
		restricted: restricted,
	}
}

func (secretKey *SecretKey) getPublicKey(postQuantum bool) (*PublicKey, error) {
	pubKey, err := secretKey.privKey.PublicKey(postQuantum)
	if err != nil {
		return nil, errDerivingPublicKey
	}
	return &PublicKey{
		version: secretKey.version,
		keyType: secretKey.keyType,
		pubKey:  pubKey,
	}, nil
}

func (secretKey *SecretKey) PublicKey(postQuantum bool) (publicKey *PublicKey, err error) {
	if postQuantum {
		if secretKey.pqPublicKey == nil {
			if secretKey.pqPublicKey, err = secretKey.getPublicKey(postQuantum); err != nil {
				return nil, err
			}
		}
		return secretKey.pqPublicKey, nil
	} else {
		if secretKey.eccPublicKey == nil {
			if secretKey.eccPublicKey, err = secretKey.getPublicKey(postQuantum); err != nil {
				return nil, err
			}
		}
		return secretKey.eccPublicKey, nil

	}
}

func SecretKeyFromBytes(bytes []byte) (*SecretKey, error) {
	if bytes[1] != 0 {
		return nil, errInvalidSecretKeyFormat
	}
	if bytes[0] > cryptoVersion {
		return nil, errUnsupportedCryptoVersion
	}
	var version uint8 = bytes[0]
	var keyType KeyType = KeyType(bytes[2])
	privKey, err := xipher.ParseSecretKey(bytes[3:])
	if err != nil {
		return nil, errInvalidSecretKeyFormat
	}
	return &SecretKey{
		version: &version,
		keyType: &keyType,
		privKey: privKey,
	}, nil
}

func SecretKeyFromString(secretKeyStr string) (*SecretKey, error) {
	sliced := strings.Split(secretKeyStr, "_")
	if len(sliced) != 3 || sliced[0] != slvPrefix {
		return nil, errInvalidSecretKeyFormat
	}
	decoded, err := commons.Decode(sliced[2])
	if err != nil {
		return nil, err
	}
	secretKey, err := SecretKeyFromBytes(decoded)
	if err != nil {
		return nil, err
	}
	if len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*secretKey.keyType)) ||
		!strings.HasSuffix(sliced[1], secretKeyAbbrev) {
		return nil, errInvalidSecretKeyFormat
	}
	return secretKey, nil
}

func (secretKey *SecretKey) Bytes() ([]byte, error) {
	if privKeyBytes, err := secretKey.privKey.Bytes(); err != nil {
		return nil, err
	} else {
		return append([]byte{*secretKey.version, 0, byte(*secretKey.keyType)}, privKeyBytes...), nil
	}
}

func (secretKey *SecretKey) RestrictSerialization() {
	secretKey.restricted = true
}

func (secretKey *SecretKey) IsSerializationRestricted() bool {
	return secretKey.restricted
}

func (secretKey SecretKey) String() string {
	if secretKey.restricted {
		return ""
	}
	if secretKeyBytes, err := secretKey.Bytes(); err != nil {
		return ""
	} else {
		return slvPrefix + "_" + string(*secretKey.keyType) + secretKeyAbbrev + "_" + commons.Encode(secretKeyBytes)
	}
}

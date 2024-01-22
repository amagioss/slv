package crypto

import (
	"strings"

	"dev.shib.me/xipher"
	"github.com/amagimedia/slv/core/commons"
)

type KeyType byte

type PublicKey struct {
	version *uint8
	keyType *KeyType
	pubKey  *xipher.PublicKey
}

func (publicKey *PublicKey) toBytes() []byte {
	return append([]byte{*publicKey.version, 1, byte(*publicKey.keyType)}, publicKey.pubKey.Bytes()...)
}

func (publicKey *PublicKey) Type() KeyType {
	return *publicKey.keyType
}

func publicKeyFromBytes(bytes []byte) (*PublicKey, error) {
	if bytes[1] != 1 {
		return nil, errInvalidPublicKeyFormat
	}
	if bytes[0] > commons.Version {
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

func (publicKey PublicKey) String() string {
	return commons.SLV + "_" + string(*publicKey.keyType) + publicKeyAbbrev + "_" + commons.Encode(publicKey.toBytes())
}

func PublicKeyFromString(publicKeyStr string) (*PublicKey, error) {
	sliced := strings.Split(publicKeyStr, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV {
		return nil, errInvalidPublicKeyFormat
	}
	decoded, err := commons.Decode(sliced[2])
	if err != nil {
		return nil, err
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
	version   *uint8
	keyType   *KeyType
	privKey   *xipher.PrivateKey
	publicKey *PublicKey
}

func NewSecretKey(keyType KeyType) (secretKey *SecretKey, err error) {
	privKey, err := xipher.NewPrivateKey()
	if err != nil {
		return nil, errGeneratingSecretKey
	}
	return newSecretKey(privKey, keyType), nil
}

func NewSecretKeyForPassword(password []byte, keyType KeyType) (secretKey *SecretKey, err error) {
	privKey, err := xipher.NewPrivateKeyForPassword(password)
	if err != nil {
		return nil, err
	}
	return newSecretKey(privKey, keyType), nil
}

func newSecretKey(privKey *xipher.PrivateKey, keyType KeyType) *SecretKey {
	version := commons.Version
	return &SecretKey{
		version: &version,
		keyType: &keyType,
		privKey: privKey,
	}
}

func (secretKey *SecretKey) PublicKey() (*PublicKey, error) {
	if secretKey.publicKey == nil {
		pubKey, err := secretKey.privKey.PublicKey()
		if err != nil {
			return nil, errDerivingPublicKey
		}
		secretKey.publicKey = &PublicKey{
			version: secretKey.version,
			keyType: secretKey.keyType,
			pubKey:  pubKey,
		}
	}
	return secretKey.publicKey, nil
}

func SecretKeyFromBytes(bytes []byte) (*SecretKey, error) {
	if bytes[1] != 0 {
		return nil, errInvalidSecretKeyFormat
	}
	if bytes[0] > commons.Version {
		return nil, errUnsupportedCryptoVersion
	}
	var version uint8 = bytes[0]
	var keyType KeyType = KeyType(bytes[2])
	privKey, err := xipher.ParsePrivateKey(bytes[3:])
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
	if len(sliced) != 3 || sliced[0] != commons.SLV {
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

func (secretKey *SecretKey) Bytes() []byte {
	if privKeyBytes, err := secretKey.privKey.Bytes(); err != nil {
		return nil
	} else {
		return append([]byte{*secretKey.version, 0, byte(*secretKey.keyType)}, privKeyBytes...)
	}
}

func (secretKey SecretKey) String() string {
	if secretKeyBytes := secretKey.Bytes(); secretKeyBytes == nil {
		return ""
	} else {
		return commons.SLV + "_" + string(*secretKey.keyType) + secretKeyAbbrev + "_" + commons.Encode(secretKeyBytes)
	}
}

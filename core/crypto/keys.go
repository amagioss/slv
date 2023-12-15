package crypto

import (
	"strings"

	"github.com/shibme/slv/core/commons"
	"gopkg.shib.me/gociphers/argon2"
	"gopkg.shib.me/gociphers/ecc"
)

type KeyType byte

type PublicKey struct {
	version *uint8
	keyType *KeyType
	pubKey  *ecc.PublicKey
}

func (publicKey *PublicKey) toBytes() []byte {
	return append([]byte{*publicKey.version, 1, byte(*publicKey.keyType)}, publicKey.pubKey.Bytes()...)
}

func (publicKey *PublicKey) Type() KeyType {
	return *publicKey.keyType
}

func publicKeyFromBytes(bytes []byte) (*PublicKey, error) {
	if len(bytes) != keyLength || bytes[1] != 1 {
		return nil, ErrInvalidPublicKeyFormat
	}
	if bytes[0] > commons.Version {
		return nil, ErrUnsupportedCryptoVersion
	}
	var version uint8 = bytes[0]
	var keyType KeyType = KeyType(bytes[2])
	pubKey, err := ecc.GetPublicKeyForBytes(bytes[3:])
	if err != nil {
		return nil, ErrInvalidPublicKeyFormat
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
		return nil, ErrInvalidPublicKeyFormat
	}
	publicKey, err := publicKeyFromBytes(commons.Decode(sliced[2]))
	if err != nil {
		return nil, err
	}
	if len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*publicKey.keyType)) ||
		!strings.HasSuffix(sliced[1], publicKeyAbbrev) {
		return nil, ErrInvalidPublicKeyFormat
	}
	return publicKey, nil
}

type SecretKey struct {
	version   *uint8
	keyType   *KeyType
	privKey   *ecc.PrivateKey
	publicKey *PublicKey
}

func NewSecretKey(keyType KeyType) (secretKey *SecretKey, err error) {
	privKey, err := ecc.NewPrivateKey()
	if err != nil {
		return nil, ErrGeneratingSecretKey
	}
	return newSecretKey(privKey, keyType), nil
}

func NewSecretKeyForPassword(password []byte, keyType KeyType) (secretKey *SecretKey, salt []byte, err error) {
	key, salt, err := argon2.GenerateKeyForPassword(password)
	if err != nil {
		return nil, nil, err
	}
	privKey, err := ecc.GetPrivateKeyForBytes(key)
	if err != nil {
		return nil, nil, err
	}
	return newSecretKey(privKey, keyType), salt, nil
}

func NewSecretKeyForPasswordAndSalt(password, salt []byte, keyType KeyType) (secretKey *SecretKey, err error) {
	key, err := argon2.GenerateKeyForPasswordAndSalt(password, salt)
	if err != nil {
		return nil, err
	}
	privKey, err := ecc.GetPrivateKeyForBytes(key)
	if err != nil {
		return nil, err
	}
	return newSecretKey(privKey, keyType), nil
}

func newSecretKey(privKey *ecc.PrivateKey, keyType KeyType) *SecretKey {
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
			return nil, ErrDerivingPublicKey
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
	if len(bytes) != keyLength || bytes[1] != 0 {
		return nil, ErrInvalidSecretKeyFormat
	}
	if bytes[0] > commons.Version {
		return nil, ErrUnsupportedCryptoVersion
	}
	var version uint8 = bytes[0]
	var keyType KeyType = KeyType(bytes[2])
	privKey, err := ecc.GetPrivateKeyForBytes(bytes[3:])
	if err != nil {
		return nil, ErrInvalidSecretKeyFormat
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
		return nil, ErrInvalidSecretKeyFormat
	}
	secretKey, err := SecretKeyFromBytes(commons.Decode(sliced[2]))
	if err != nil {
		return nil, err
	}
	if len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*secretKey.keyType)) ||
		!strings.HasSuffix(sliced[1], secretKeyAbbrev) {
		return nil, ErrInvalidSecretKeyFormat
	}
	return secretKey, nil
}

func (secretKey *SecretKey) Bytes() []byte {
	return append([]byte{*secretKey.version, 0, byte(*secretKey.keyType)}, secretKey.privKey.Bytes()...)
}

func (secretKey SecretKey) String() string {
	return commons.SLV + "_" + string(*secretKey.keyType) + secretKeyAbbrev + "_" + commons.Encode(secretKey.Bytes())
}

package crypto

import (
	"fmt"
	"strings"

	"github.com/shibme/slv/core/commons"
)

type SecretKey struct {
	key *[keyLength]byte
	*PublicKey
}

func (secretKey *SecretKey) toBytes() []byte {
	if secretKey.PublicKey == nil || secretKey.key == nil {
		return nil
	}
	return append(secretKey.key[:], secretKey.PublicKey.toBytes()...)
}

func secretKeyFromBytes(secretKeyBytes []byte) (*SecretKey, error) {
	if len(secretKeyBytes) != secretKeyLength {
		return nil, ErrInvalidKeyFormat
	}
	publicKeyBytes := secretKeyBytes[:publicKeyLength]
	key := new([keyLength]byte)
	copy(key[:], secretKeyBytes[len(secretKeyBytes)-keyLength:])
	var publicKey *PublicKey
	var err error
	if publicKey, err = publicKeyFromBytes(publicKeyBytes); err != nil {
		return nil, err
	}
	return &SecretKey{
		key:       key,
		PublicKey: publicKey,
	}, nil
}

func (secretKey SecretKey) String() string {
	return fmt.Sprintf("%s_%s%s_%s", commons.SLV, secretKey.keyType,
		secretKeyAbbrev, commons.Encode(secretKey.toBytes()))
}

func SecretKeyFromString(secretKeyStr string) (secretKey *SecretKey, err error) {
	sliced := strings.Split(secretKeyStr, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV {
		return nil, ErrInvalidKeyFormat
	}
	if secretKey, err = secretKeyFromBytes(commons.Decode(sliced[2])); err != nil {
		return nil, err
	}
	if len(sliced[1]) != 3 || strings.HasPrefix(sliced[1], string(*secretKey.keyType)) ||
		strings.HasSuffix(sliced[1], secretKeyAbbrev) {
		return nil, ErrInvalidKeyFormat
	}
	return
}

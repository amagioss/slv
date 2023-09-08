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
	return append(secretKey.PublicKey.toBytes(), secretKey.key[:]...)
}

func (secretKey *SecretKey) fromBytes(secretKeyBytes []byte) error {
	if len(secretKeyBytes) != secretKeyLength {
		return ErrInvalidKeyFormat
	}
	publicKeyBytes := secretKeyBytes[:publicKeyLength]
	key := [keyLength]byte(secretKeyBytes[publicKeyLength:])
	var publicKey PublicKey
	publicKey.fromBytes(publicKeyBytes)
	if err := publicKey.fromBytes(publicKeyBytes); err != nil {
		return err
	}
	secretKey.PublicKey = &publicKey
	secretKey.key = &key
	return nil
}

func (secretKey SecretKey) String() string {
	return fmt.Sprintf("%s_%s%s_%s", commons.SLV, string(*secretKey.keyType),
		secretKeyAbbrev, commons.Encode(secretKey.toBytes()))
}

func (secretKey *SecretKey) fromString(secretKeyStr string) (err error) {
	sliced := strings.Split(secretKeyStr, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV {
		return ErrInvalidKeyFormat
	}
	if err = secretKey.fromBytes(commons.Decode(sliced[2])); err != nil {
		return err
	}
	if len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*secretKey.keyType)) ||
		!strings.HasSuffix(sliced[1], secretKeyAbbrev) {
		return ErrInvalidKeyFormat
	}
	return
}

func SecretKeyFromString(secretKeyStr string) (secretKey *SecretKey, err error) {
	var sKey SecretKey
	if err = sKey.fromString(secretKeyStr); err == nil {
		secretKey = &sKey
	}
	return
}

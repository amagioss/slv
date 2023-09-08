package crypto

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shibme/slv/core/commons"
	"gopkg.in/yaml.v3"
)

type KeyType byte

type PublicKey struct {
	key       *[keyLength]byte
	keyType   *KeyType
	encrypter *encrypter
}

func (publicKey *PublicKey) Id() string {
	return commons.Encode(publicKey.toBytes())
}

func (publicKey *PublicKey) ShortId() *[keyIdLength]byte {
	if publicKey.key == nil {
		return nil
	}
	sum := sha1.Sum(publicKey.key[:])
	keyId := [keyIdLength]byte(sum[len(sum)-keyIdLength:])
	return &keyId
}

func (publicKey *PublicKey) toBytes() []byte {
	return append([]byte{byte(*publicKey.keyType)}, publicKey.key[:]...)
}

func (publicKey *PublicKey) fromBytes(publicKeyBytes []byte) error {
	if len(publicKeyBytes) != publicKeyLength {
		return ErrInvalidKeyFormat
	}
	keyType := KeyType(publicKeyBytes[0])
	key := [keyLength]byte(publicKeyBytes[1:])
	publicKey.keyType = &keyType
	publicKey.key = &key
	return nil
}

func (publicKey PublicKey) Type() KeyType {
	return *publicKey.keyType
}

func (publicKey PublicKey) String() string {
	return fmt.Sprintf("%s_%s%s_%s", commons.SLV, string(*publicKey.keyType),
		publicKeyAbbrev, publicKey.Id())
}

func (publicKey *PublicKey) fromString(publicKeyStr string) (err error) {
	sliced := strings.Split(publicKeyStr, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV {
		return ErrInvalidKeyFormat
	}
	if err = publicKey.fromBytes(commons.Decode(sliced[2])); err != nil {
		return err
	}
	if len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*publicKey.keyType)) ||
		!strings.HasSuffix(sliced[1], publicKeyAbbrev) {
		return ErrInvalidKeyFormat
	}
	return
}

func PublicKeyFromString(publicKeyStr string) (publicKey *PublicKey, err error) {
	var pKey PublicKey
	if err = pKey.fromString(publicKeyStr); err == nil {
		publicKey = &pKey
	}
	return
}

func (publicKey PublicKey) MarshalYAML() (interface{}, error) {
	return publicKey.String(), nil
}

func (publicKey *PublicKey) UnmarshalYAML(value *yaml.Node) (err error) {
	var pubKeyStr string
	err = value.Decode(&pubKeyStr)
	if err == nil {
		publicKey.fromString(pubKeyStr)
	}
	return
}

func (publicKey PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(publicKey.String())
}

func (publicKey *PublicKey) UnmarshalJSON(data []byte) (err error) {
	var pubKeyStr string
	err = json.Unmarshal(data, &pubKeyStr)
	if err == nil {
		publicKey.fromString(pubKeyStr)
	}
	return
}

type encrypter struct {
	encryptionKeyId    *[keyIdLength]byte
	ephemeralPublicKey *[keyLength]byte
	sharedKey          *[keyLength]byte
}

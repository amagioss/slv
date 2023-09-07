package crypto

import (
	"bytes"
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

func (publicKey *PublicKey) Id() *[keyIdLength]byte {
	if publicKey.keyType == nil || publicKey.key == nil {
		return nil
	}
	sum := sha1.Sum(publicKey.key[:])
	keyId := [keyIdLength]byte(sum[len(sum)-keyIdLength:])
	return &keyId
}

func (publicKey *PublicKey) toBytes() []byte {
	bytes := append(publicKey.Id()[:], publicKey.key[:]...)
	return append([]byte{byte(*publicKey.keyType)}, bytes...)
}

func publicKeyFromBytes(publicKeyBytes []byte) (*PublicKey, error) {
	if len(publicKeyBytes) != publicKeyLength {
		return nil, ErrInvalidKeyFormat
	}
	keyType := KeyType(publicKeyBytes[0])
	sumFromBytes := publicKeyBytes[1 : keyIdLength+1]
	key := new([keyLength]byte)
	copy(key[:], publicKeyBytes[len(publicKeyBytes)-keyLength:])
	publicKey := &PublicKey{
		keyType: &keyType,
		key:     key,
	}
	if !bytes.Equal(sumFromBytes, publicKey.Id()[:]) {
		return nil, ErrInvalidKeyFormat
	}
	return &PublicKey{
		keyType: &keyType,
		key:     key,
	}, nil
}

func (publicKey PublicKey) Type() KeyType {
	return *publicKey.keyType
}

func (publicKey PublicKey) String() string {
	return fmt.Sprintf("%s_%s%s_%s", commons.SLV, publicKey.keyType,
		publicKeyAbbrev, commons.Encode(publicKey.toBytes()))
}

func PublicKeyFromString(publicKeyStr string) (publicKey *PublicKey, err error) {
	sliced := strings.Split(publicKeyStr, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV {
		return nil, ErrInvalidKeyFormat
	}
	if publicKey, err = publicKeyFromBytes(commons.Decode(sliced[2])); err != nil {
		return nil, err
	}
	if len(sliced[1]) != 3 || strings.HasPrefix(sliced[1], string(*publicKey.keyType)) ||
		strings.HasSuffix(sliced[1], publicKeyAbbrev) {
		return nil, ErrInvalidKeyFormat
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
		publicKey, err = PublicKeyFromString(pubKeyStr)
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
		publicKey, err = PublicKeyFromString(pubKeyStr)
	}
	return
}

type encrypter struct {
	encryptionKeyId    *[keyIdLength]byte
	ephemeralPublicKey *[keyLength]byte
	sharedKey          *[keyLength]byte
}

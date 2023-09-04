package crypto

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/shibme/slv/commons"
)

type KeyType byte

type key struct {
	id      [8]byte
	keyData [32]byte
	public  bool
	keyType KeyType
	idStr   *string
	str     *string
}

func (k *key) Id() string {
	if k.idStr == nil {
		k.idStr = new(string)
		*k.idStr = commons.Encode(k.id[:])
	}
	return *k.idStr
}

func (k *key) Public() bool {
	return k.public
}

func (k *key) Type() KeyType {
	return k.keyType
}

func (k key) String() string {
	if k.str == nil {
		k.str = new(string)
		keyPart := []byte{byte(k.keyType), 0}
		if k.public {
			keyPart = []byte{byte(k.keyType), 1}
		}
		keyPart = append(keyPart, k.keyData[:]...)
		var prefix string
		if k.public {
			prefix = publicKeyPrefix
		} else {
			prefix = privateKeyPrefix
		}
		*k.str = fmt.Sprintf("%s%s_%s", prefix, commons.Encode(k.id[:]), commons.Encode(keyPart))
	}
	return *k.str
}

func isValidKey(k key) bool {
	var checksumFromId [4]byte
	if k.public {
		checksumFromId = [4]byte(k.id[0:4])
	} else {
		checksumFromId = [4]byte(k.id[4:])
	}
	calculatedChecksum := sha1.Sum(k.keyData[:])
	checksumFromKey := [4]byte(calculatedChecksum[len(calculatedChecksum)-4:])
	return checksumFromId == checksumFromKey
}

func keyFromString(keyString string) (*key, error) {
	var public bool
	var trimmedKeyString string
	if strings.HasPrefix(keyString, privateKeyPrefix) {
		public = false
		trimmedKeyString = strings.TrimPrefix(keyString, privateKeyPrefix)
	} else if strings.HasPrefix(keyString, publicKeyPrefix) {
		public = true
		trimmedKeyString = strings.TrimPrefix(keyString, publicKeyPrefix)
	} else {
		return nil, ErrInvalidKeyFormat
	}
	sliced := strings.Split(trimmedKeyString, "_")
	if len(sliced) != 2 {
		return nil, ErrInvalidKeyFormat
	}
	idPart := commons.Decode(sliced[0])
	keyPart := commons.Decode(sliced[1])
	if len(idPart) != 8 || len(keyPart) != 34 {
		return nil, ErrInvalidKeyFormat
	}
	k := new(key)
	if keyPart[1] == 0 && !public {
		k.public = false
	} else if keyPart[1] == 1 && public {
		k.public = true
	} else {
		return nil, ErrInvalidKeyFormat
	}
	k.keyType = KeyType(keyPart[0])
	k.id = [8]byte(idPart)
	k.keyData = [32]byte(keyPart[2:])
	if !isValidKey(*k) {
		return nil, ErrInvalidKeyFormat
	}
	return k, nil
}

func (k *key) toBytes() []byte {
	keyPart := []byte{byte(k.keyType), 0}
	if k.public {
		keyPart = []byte{byte(k.keyType), 1}
	}
	keyPart = append(keyPart, k.keyData[:]...)
	return append(k.id[:], keyPart...)
}

func keyFromBytes(keyBytes []byte) (k *key, err error) {
	k = new(key)
	if len(keyBytes) != 42 {
		return nil, ErrInvalidKeyFormat
	}
	k.id = [8]byte(keyBytes[0:8])
	if keyBytes[9] == 0 {
		k.public = false
	} else if keyBytes[9] == 1 {
		k.public = true
	} else {
		return nil, ErrInvalidKeyFormat
	}
	k.keyType = KeyType(keyBytes[8])
	k.keyData = [32]byte(keyBytes[10:])
	if !isValidKey(*k) {
		return nil, ErrInvalidKeyFormat
	}
	return k, nil
}

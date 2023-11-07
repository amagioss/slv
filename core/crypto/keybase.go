package crypto

import (
	"strings"

	"github.com/shibme/slv/core/commons"
)

type KeyType byte

type keyBase struct {
	version *uint8
	public  bool
	keyType *KeyType
	key     *[]byte
}

func (kb *keyBase) Type() KeyType {
	return *kb.keyType
}

func (keyBase *keyBase) toBytes() []byte {
	if keyBase.public {
		return append([]byte{*keyBase.version, 1, byte(*keyBase.keyType)}, *keyBase.key...)
	} else {
		return append([]byte{*keyBase.version, 0, byte(*keyBase.keyType)}, *keyBase.key...)
	}
}

func keyBaseFromBytes(keyBaseBytes []byte) (*keyBase, error) {
	if len(keyBaseBytes) != keyBaseLength || (keyBaseBytes[1] != 0 && keyBaseBytes[1] != 1) {
		return nil, ErrInvalidKeyFormat
	}
	if keyBaseBytes[0] > commons.Version {
		return nil, ErrUnsupportedCryptoVersion
	}
	var version uint8 = keyBaseBytes[0]
	public := false
	if keyBaseBytes[1] == 1 {
		public = true
	}
	var keyType KeyType = KeyType(keyBaseBytes[2])
	key := keyBaseBytes[3:]
	return &keyBase{
		version: &version,
		public:  public,
		keyType: &keyType,
		key:     &key,
	}, nil
}

func (keyBase *keyBase) String() string {
	if keyBase.public {
		return commons.SLV + "_" + string(*keyBase.keyType) + publicKeyAbbrev + "_" + commons.Encode(keyBase.toBytes())
	} else {
		return commons.SLV + "_" + string(*keyBase.keyType) + secretKeyAbbrev + "_" + commons.Encode(keyBase.toBytes())
	}
}

func keyBaseFromString(keyBaseStr string) (*keyBase, error) {
	sliced := strings.Split(keyBaseStr, "_")
	if len(sliced) != 3 || sliced[0] != commons.SLV {
		return nil, ErrInvalidKeyFormat
	}
	keyBase, err := keyBaseFromBytes(commons.Decode(sliced[2]))
	if err != nil {
		return nil, err
	}
	var keyAbbrev string
	if keyBase.public {
		keyAbbrev = publicKeyAbbrev
	} else {
		keyAbbrev = secretKeyAbbrev
	}
	if len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*keyBase.keyType)) ||
		!strings.HasSuffix(sliced[1], keyAbbrev) {
		return nil, ErrInvalidKeyFormat
	}
	return keyBase, nil
}

package crypto

import (
	"strings"

	"github.com/shibme/slv/core/commons"
	"gopkg.in/yaml.v3"
)

type ciphered struct {
	ciphertext *[]byte
	keyId      *[keyIdLength]byte
	keyType    *KeyType
	hash       *[]byte
}

func (ciph ciphered) toString(ciphertextType string) string {
	str := commons.SLV + "_" + string(*ciph.keyType) + ciphertextType + "_" + commons.Encode(ciph.keyId[:]) + "_"
	encrypted := commons.Encode(append([]byte{byte(*ciph.keyType)}, *ciph.ciphertext...))
	if ciph.hash == nil {
		return str + encrypted
	} else {
		return str + commons.Encode(*ciph.hash) + "_" + encrypted
	}
}

func cipheredFromString(cipheredStr, ciphertextType string) (*ciphered, error) {
	sliced := strings.Split(cipheredStr, "_")
	if len(sliced) != 4 && len(sliced) != 5 {
		return nil, ErrCiphertextFormat
	}
	encryptionKeyId := [keyIdLength]byte(commons.Decode(sliced[2]))
	encrypted := commons.Decode(sliced[len(sliced)-1])
	var keyType KeyType = KeyType(encrypted[0])
	encrypted = encrypted[1:]
	if sliced[0] != commons.SLV || len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(keyType)) ||
		!strings.HasSuffix(sliced[1], ciphertextType) || len(encryptionKeyId) != keyIdLength {
		return nil, ErrCiphertextFormat
	}
	ciph := &ciphered{
		keyType:    &keyType,
		keyId:      &encryptionKeyId,
		ciphertext: &encrypted,
	}
	if len(sliced) == 5 {
		hash := commons.Decode(sliced[3])
		if len(hash) > secretHashMaxLength {
			return nil, ErrCiphertextFormat
		}
		ciph.hash = &hash
	}
	return ciph, nil
}

func (ciph *ciphered) GetKeyId() *[keyIdLength]byte {
	return ciph.keyId
}

type SealedSecret struct {
	*ciphered
}

func (sealedSecret SealedSecret) String() string {
	return sealedSecret.toString(sealedSecretAbbrev)
}

func (sealedSecret *SealedSecret) FromString(sealedSecretStr string) (err error) {
	var ciphered *ciphered
	if ciphered, err = cipheredFromString(sealedSecretStr, sealedSecretAbbrev); err == nil {
		sealedSecret.ciphered = ciphered
	}
	return
}

func (sealedSecret *ciphered) GetHash() string {
	return commons.Encode(*sealedSecret.hash)
}

func (sealedSecret SealedSecret) MarshalYAML() (interface{}, error) {
	return sealedSecret.String(), nil
}

func (sealedSecret *SealedSecret) UnmarshalYAML(value *yaml.Node) (err error) {
	var sealedSecretStr string
	if value.Decode(&sealedSecretStr) == nil {
		var ciphered *ciphered
		if ciphered, err = cipheredFromString(sealedSecretStr, sealedSecretAbbrev); err == nil {
			sealedSecret.ciphered = ciphered
		}
	}
	return
}

type WrappedKey struct {
	*ciphered
}

func (wrappedKey WrappedKey) String() string {
	return wrappedKey.toString(wrappedKeyAbbrev)
}

func (wrappedKey WrappedKey) MarshalYAML() (interface{}, error) {
	return wrappedKey.String(), nil
}

func (wrappedKey *WrappedKey) UnmarshalYAML(value *yaml.Node) (err error) {
	var wrappedKeyStr string
	if value.Decode(&wrappedKeyStr) == nil {
		var ciphered *ciphered
		if ciphered, err = cipheredFromString(wrappedKeyStr, wrappedKeyAbbrev); err == nil {
			wrappedKey.ciphered = ciphered
		}
	}
	return
}

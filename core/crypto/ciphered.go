package crypto

import (
	"strings"

	"github.com/shibme/slv/core/commons"
	"gopkg.in/yaml.v3"
)

type ciphered struct {
	version    *uint8
	keyType    *KeyType
	ciphertext *[]byte
	shortKeyId *[]byte
}

func (ciph ciphered) toBytes() []byte {
	cipheredBytes := append([]byte{*ciph.version, byte(*ciph.keyType)}, *ciph.shortKeyId...)
	cipheredBytes = append(cipheredBytes, *ciph.ciphertext...)
	return cipheredBytes
}

func cipheredFromBytes(cipheredBytes []byte) (*ciphered, error) {
	if len(cipheredBytes) < cipherBytesMinLength {
		return nil, ErrInvalidCiphertextFormat
	}
	var version byte = cipheredBytes[0]
	var keyType KeyType = KeyType(cipheredBytes[1])
	cipheredBytes = cipheredBytes[2:]
	shortKeyId := cipheredBytes[:shortKeyIdLength]
	ciphertext := cipheredBytes[shortKeyIdLength:]
	return &ciphered{
		version:    &version,
		keyType:    &keyType,
		shortKeyId: &shortKeyId,
		ciphertext: &ciphertext,
	}, nil
}

func (ciph *ciphered) GetKeyId() *[]byte {
	return ciph.shortKeyId
}

type SealedSecret struct {
	*ciphered
	hash *[]byte
}

func (sealedSecret SealedSecret) String() string {
	if sealedSecret.hash == nil {
		return commons.SLV + "_" + string(*sealedSecret.keyType) + sealedSecretAbbrev + "_" +
			commons.Encode(sealedSecret.toBytes())
	} else {
		return commons.SLV + "_" + string(*sealedSecret.keyType) + sealedSecretAbbrev + "_" +
			commons.Encode(*sealedSecret.hash) + "_" + commons.Encode(sealedSecret.toBytes())
	}
}

func (sealedSecret *SealedSecret) fromString(sealedSecretStr string) (err error) {
	sliced := strings.Split(sealedSecretStr, "_")
	if len(sliced) != 3 && len(sliced) != 4 {
		return ErrInvalidCiphertextFormat
	}
	ciphered, err := cipheredFromBytes(commons.Decode(sliced[len(sliced)-1]))
	if err != nil {
		return err
	}
	if sliced[0] != commons.SLV || len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*ciphered.keyType)) ||
		!strings.HasSuffix(sliced[1], sealedSecretAbbrev) || len(*ciphered.shortKeyId) != shortKeyIdLength {
		return ErrInvalidCiphertextFormat
	}
	if len(sliced) == 4 {
		hash := commons.Decode(sliced[2])
		if len(hash) > argon2HashMaxLength {
			return ErrInvalidCiphertextFormat
		}
		sealedSecret.hash = &hash
	}
	sealedSecret.ciphered = ciphered
	return
}

func (sealedSecret *SealedSecret) GetHash() string {
	return commons.Encode(*sealedSecret.hash)
}

func (sealedSecret SealedSecret) MarshalYAML() (interface{}, error) {
	return sealedSecret.String(), nil
}

func (sealedSecret *SealedSecret) UnmarshalYAML(value *yaml.Node) (err error) {
	var sealedSecretStr string
	if value.Decode(&sealedSecretStr) == nil {
		return sealedSecret.fromString(sealedSecretStr)
	}
	return
}

type WrappedKey struct {
	*ciphered
}

func (wrappedKey WrappedKey) String() string {
	return commons.SLV + "_" + string(*wrappedKey.keyType) + wrappedKeyAbbrev + "_" +
		commons.Encode(wrappedKey.toBytes())
}

func (wrappedKey *WrappedKey) fromString(wrappedKeyStr string) (err error) {
	sliced := strings.Split(wrappedKeyStr, "_")
	if len(sliced) != 3 {
		return ErrInvalidCiphertextFormat
	}
	ciphered, err := cipheredFromBytes(commons.Decode(sliced[len(sliced)-1]))
	if err != nil {
		return err
	}
	if sliced[0] != commons.SLV || len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*ciphered.keyType)) ||
		!strings.HasSuffix(sliced[1], wrappedKeyAbbrev) || len(*ciphered.shortKeyId) != shortKeyIdLength {
		return ErrInvalidCiphertextFormat
	}
	wrappedKey.ciphered = ciphered
	return
}

func (wrappedKey WrappedKey) MarshalYAML() (interface{}, error) {
	return wrappedKey.String(), nil
}

func (wrappedKey *WrappedKey) UnmarshalYAML(value *yaml.Node) (err error) {
	var wrappedKeyStr string
	if value.Decode(&wrappedKeyStr) == nil {
		return wrappedKey.fromString(wrappedKeyStr)
	}
	return
}

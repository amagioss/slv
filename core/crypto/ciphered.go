package crypto

import (
	"bytes"
	"encoding/binary"
	"strings"
	"time"

	"oss.amagi.com/slv/core/commons"
)

type ciphered struct {
	version     *uint8
	keyType     *KeyType
	encryptedAt *time.Time
	ciphertext  []byte
	encryptedBy []byte
}

func timeToBytes(t time.Time) []byte {
	tsBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tsBytes, uint32(t.Unix()))
	return tsBytes
}

func bytesToTime(timeBytes []byte) time.Time {
	return time.Unix(int64(binary.BigEndian.Uint32(timeBytes)), 0)
}

func (ciph ciphered) toBytes() []byte {
	keyLenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(keyLenBytes, uint16(len(ciph.encryptedBy)))
	cipheredBytes := append([]byte{*ciph.version, byte(*ciph.keyType)}, timeToBytes(*ciph.encryptedAt)...)
	cipheredBytes = append(cipheredBytes, keyLenBytes...)
	cipheredBytes = append(cipheredBytes, ciph.encryptedBy...)
	cipheredBytes = append(cipheredBytes, ciph.ciphertext...)
	return cipheredBytes
}

func cipheredFromBytes(cipheredBytes []byte) (*ciphered, error) {
	var version byte = cipheredBytes[0]
	var keyType KeyType = KeyType(cipheredBytes[1])
	cipheredBytes = cipheredBytes[2:]
	encryptedAt := bytesToTime(cipheredBytes[:4])
	cipheredBytes = cipheredBytes[4:]
	keyLen := binary.BigEndian.Uint16(cipheredBytes[:2])
	cipheredBytes = cipheredBytes[2:]
	if len(cipheredBytes) < int(keyLen) {
		return nil, errInvalidCiphertextFormat
	}
	encryptedBy := cipheredBytes[:keyLen]
	ciphertext := cipheredBytes[keyLen:]
	return &ciphered{
		version:     &version,
		keyType:     &keyType,
		encryptedAt: &encryptedAt,
		encryptedBy: encryptedBy,
		ciphertext:  ciphertext,
	}, nil
}

func (ciph *ciphered) IsEncryptedBy(publicKey *PublicKey) bool {
	if publicKeyBytes, err := publicKey.toBytes(); err == nil {
		return bytes.Equal(ciph.encryptedBy, publicKeyBytes)
	}
	return false
}

func (ciph *ciphered) EncryptedBy() []byte {
	return ciph.encryptedBy
}

func (ciph *ciphered) EncryptedByPublicKey() (*PublicKey, error) {
	return publicKeyFromBytes(ciph.encryptedBy)
}

func (ciph *ciphered) EncryptedAt() time.Time {
	return *ciph.encryptedAt
}

type SealedSecret struct {
	*ciphered
	hash *[]byte
}

func (sealedSecret SealedSecret) String() string {
	if sealedSecret.hash == nil {
		return slvPrefix + "_" + string(*sealedSecret.keyType) + sealedSecretAbbrev + "_" +
			commons.Encode(sealedSecret.toBytes())
	} else {
		return slvPrefix + "_" + string(*sealedSecret.keyType) + sealedSecretAbbrev + "_" +
			commons.Encode(*sealedSecret.hash) + "_" + commons.Encode(sealedSecret.toBytes())
	}
}

func (sealedSecret *SealedSecret) FromString(sealedSecretStr string) (err error) {
	sliced := strings.Split(sealedSecretStr, "_")
	if len(sliced) != 3 && len(sliced) != 4 {
		return errInvalidCiphertextFormat
	}
	decoded, err := commons.Decode(sliced[len(sliced)-1])
	if err != nil {
		return err
	}
	ciphered, err := cipheredFromBytes(decoded)
	if err != nil {
		return err
	}
	if sliced[0] != slvPrefix || len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*ciphered.keyType)) ||
		!strings.HasSuffix(sliced[1], sealedSecretAbbrev) {
		return errInvalidCiphertextFormat
	}
	if len(sliced) == 4 {
		hash, err := commons.Decode(sliced[2])
		if err != nil {
			return err
		}
		if len(hash) > hashMaxLength {
			return errInvalidCiphertextFormat
		}
		sealedSecret.hash = &hash
	}
	sealedSecret.ciphered = ciphered
	return
}

func (sealedSecret *SealedSecret) Hash() string {
	if sealedSecret.hash == nil {
		return ""
	}
	return commons.Encode(*sealedSecret.hash)
}

type WrappedKey struct {
	*ciphered
}

func (wrappedKey WrappedKey) String() string {
	return slvPrefix + "_" + string(*wrappedKey.keyType) + wrappedKeyAbbrev + "_" +
		commons.Encode(wrappedKey.toBytes())
}

func (wrappedKey *WrappedKey) FromString(wrappedKeyStr string) (err error) {
	sliced := strings.Split(wrappedKeyStr, "_")
	if len(sliced) != 3 {
		return errInvalidCiphertextFormat
	}
	decoded, err := commons.Decode(sliced[len(sliced)-1])
	if err != nil {
		return err
	}
	ciphered, err := cipheredFromBytes(decoded)
	if err != nil {
		return err
	}
	if sliced[0] != slvPrefix || len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*ciphered.keyType)) ||
		!strings.HasSuffix(sliced[1], wrappedKeyAbbrev) {
		return errInvalidCiphertextFormat
	}
	wrappedKey.ciphered = ciphered
	return
}

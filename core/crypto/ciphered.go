package crypto

import (
	"bytes"
	"strings"

	"github.com/amagimedia/slv/core/commons"
)

type ciphered struct {
	version     *uint8
	keyType     *KeyType
	ciphertext  *[]byte
	pubKeyBytes *[]byte
}

func (ciph ciphered) toBytes() []byte {
	cipheredBytes := append([]byte{*ciph.version, byte(*ciph.keyType)}, *ciph.pubKeyBytes...)
	cipheredBytes = append(cipheredBytes, *ciph.ciphertext...)
	return cipheredBytes
}

func cipheredFromBytes(cipheredBytes []byte) (*ciphered, error) {
	if len(cipheredBytes) < cipherTextMinLength {
		return nil, errInvalidCiphertextFormat
	}
	var version byte = cipheredBytes[0]
	var keyType KeyType = KeyType(cipheredBytes[1])
	cipheredBytes = cipheredBytes[2:]
	pubKeyBytes := cipheredBytes[:publicKeyLength]
	ciphertext := cipheredBytes[publicKeyLength:]
	return &ciphered{
		version:     &version,
		keyType:     &keyType,
		pubKeyBytes: &pubKeyBytes,
		ciphertext:  &ciphertext,
	}, nil
}

func (ciph *ciphered) IsEncryptedBy(publicKey *PublicKey) bool {
	return bytes.Equal(*ciph.pubKeyBytes, publicKey.toBytes())
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
	if sliced[0] != commons.SLV || len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*ciphered.keyType)) ||
		!strings.HasSuffix(sliced[1], sealedSecretAbbrev) || len(*ciphered.pubKeyBytes) != publicKeyLength {
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

func (sealedSecret *SealedSecret) GetHash() string {
	return commons.Encode(*sealedSecret.hash)
}

type WrappedKey struct {
	*ciphered
}

func (wrappedKey WrappedKey) String() string {
	return commons.SLV + "_" + string(*wrappedKey.keyType) + wrappedKeyAbbrev + "_" +
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
	if sliced[0] != commons.SLV || len(sliced[1]) != 3 || !strings.HasPrefix(sliced[1], string(*ciphered.keyType)) ||
		!strings.HasSuffix(sliced[1], wrappedKeyAbbrev) || len(*ciphered.pubKeyBytes) != publicKeyLength {
		return errInvalidCiphertextFormat
	}
	wrappedKey.ciphered = ciphered
	return
}

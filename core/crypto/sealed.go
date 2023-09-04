package crypto

import (
	"fmt"
	"strings"

	"github.com/shibme/slv/core/commons"
	"gopkg.in/yaml.v3"
)

type SealedData struct {
	checksum        [8]byte
	encryptionKeyId [8]byte
	data            []byte
}

func (sealedData *SealedData) FromString(sealedDataString string) (err error) {
	err = ErrSealedDataFormat
	if !strings.HasPrefix(sealedDataString, sealedDataPrefix) {
		return
	}
	trimmedCipheredSecretString := strings.TrimPrefix(sealedDataString, sealedDataPrefix)
	sliced := strings.Split(trimmedCipheredSecretString, "_")
	if len(sliced) != 3 {
		return
	}
	encryptionKeyId := commons.Decode(sliced[0])
	checksum := commons.Decode(sliced[1])
	data := commons.Decode(sliced[2])
	if len(checksum) != 8 || len(encryptionKeyId) != 8 || len(data) == 0 {
		return
	}
	sealedData.checksum = [8]byte(checksum)
	sealedData.encryptionKeyId = [8]byte(encryptionKeyId)
	sealedData.data = data
	return nil
}

func (sealedData *SealedData) GetChecksum() string {
	return commons.Encode(sealedData.checksum[:])
}

func (sealedData *SealedData) UnmarshalYAML(value *yaml.Node) (err error) {
	var sealedDataStr string
	err = value.Decode(&sealedDataStr)
	if err == nil {
		err = sealedData.FromString(sealedDataStr)
	}
	return
}

func (sealedData SealedData) MarshalYAML() (interface{}, error) {
	return sealedData.String(), nil
}

func (sealedData SealedData) String() string {
	return fmt.Sprintf("%s%s_%s_%s", sealedDataPrefix, commons.Encode(sealedData.encryptionKeyId[:]),
		commons.Encode(sealedData.checksum[:]), commons.Encode(sealedData.data))
}

type SealedKey struct {
	encryptionKeyId [8]byte
	data            []byte
}

func (sealedKey *SealedKey) fromString(sealedKeyString string) (err error) {
	err = ErrSealedKeyFormat
	if !strings.HasPrefix(sealedKeyString, sealedKeyPrefix) {
		return
	}
	trimmedSealedKeyString := strings.TrimPrefix(sealedKeyString, sealedKeyPrefix)
	sliced := strings.Split(trimmedSealedKeyString, "_")
	if len(sliced) != 2 {
		return
	}
	encryptionKeyId := commons.Decode(sliced[0])
	data := commons.Decode(sliced[1])
	if len(encryptionKeyId) != 8 || len(data) == 0 {
		return
	}
	sealedKey.encryptionKeyId = [8]byte(encryptionKeyId)
	sealedKey.data = data
	return nil
}

func (sealedKey *SealedKey) GetAccessKeyId() string {
	return commons.Encode(sealedKey.encryptionKeyId[:])
}

func (sealedKey *SealedKey) UnmarshalYAML(value *yaml.Node) (err error) {
	var sealedKeyStr string
	err = value.Decode(&sealedKeyStr)
	if err == nil {
		err = sealedKey.fromString(sealedKeyStr)
	}
	return
}

func (sealedKey SealedKey) MarshalYAML() (interface{}, error) {
	return sealedKey.String(), nil
}

func (sealedKey SealedKey) String() (sealedKeyStr string) {
	return fmt.Sprintf("%s%s_%s", sealedKeyPrefix, commons.Encode(sealedKey.encryptionKeyId[:]), commons.Encode(sealedKey.data))
}

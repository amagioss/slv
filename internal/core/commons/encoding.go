package commons

import (
	"encoding/base32"
	"encoding/json"
)

var (
	base32Encoding = base32.StdEncoding.WithPadding(base32.NoPadding)
)

func jsonSerialize(data interface{}) (dataBytes []byte, err error) {
	dataBytes, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return dataBytes, nil
}

func jsonDeserialize(dataBytes []byte, data interface{}) (err error) {
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return err
	}
	return nil
}

func Serialize(data interface{}) (string, error) {
	serialized, err := jsonSerialize(data)
	if err != nil {
		return "", err
	}
	serialized, err = Compress(serialized)
	if err != nil {
		return "", err
	}
	return Encode(serialized), nil
}

func Deserialize(serialized string, data interface{}) (err error) {
	serializedBytes, err := Decode(serialized)
	if err != nil {
		return err
	}
	serializedBytes, err = Decompress(serializedBytes)
	if err != nil {
		return err
	}
	return jsonDeserialize(serializedBytes, &data)
}

func Encode(data []byte) (encoded string) {
	return base32Encoding.EncodeToString(data)
}

func Decode(encoded string) (data []byte, err error) {
	return base32Encoding.DecodeString(encoded)
}

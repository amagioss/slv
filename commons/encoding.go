package commons

import (
	"encoding/json"

	"github.com/btcsuite/btcutil/base58"
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
	serializedBytes := Decode(serialized)
	serializedBytes, err = Decompress(serializedBytes)
	if err != nil {
		return err
	}
	return jsonDeserialize(serializedBytes, &data)
}

func Encode(data []byte) (encoded string) {
	return base58.Encode(data)
}

func Decode(encoded string) (data []byte) {
	return base58.Decode(encoded)
}

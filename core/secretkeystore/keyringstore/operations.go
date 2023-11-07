package keyringstore

import "github.com/99designs/keyring"

var ring keyring.Keyring

func initKeyring() (err error) {
	if ring == nil {
		ring, err = keyring.Open(keyring.Config{
			ServiceName: slvKeyringServiceName,
		})
	}
	return
}

func putUserSecretToKeyring(secretBytes []byte) error {
	if err := initKeyring(); err != nil {
		return err
	}
	_ = ring.Set(keyring.Item{
		Key:  slvKeyringItemKey,
		Data: (secretBytes),
	})
	return nil
}

func getUserSecretFromKeyring() ([]byte, error) {
	if err := initKeyring(); err != nil {
		return nil, err
	}
	item, err := ring.Get(slvKeyringItemKey)
	if err != nil {
		return nil, err
	}
	return item.Data, nil
}

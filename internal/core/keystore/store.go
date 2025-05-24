package keystore

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
	"slv.sh/slv/internal/core/commons"
	"slv.sh/slv/internal/core/config"
	"xipher.org/xipher"
)

const (
	keyringServiceName = config.AppNameLowerCase
	keystoreDirName    = "keystore"
)

var (
	keyStoreDir = filepath.Join(config.GetAppDataDir(), keystoreDirName)

	ErrNotFound = fmt.Errorf("data not found in keystore")
)

func putToKeyStore(storeId string, payloadId, payload []byte) error {
	seed := sha512.Sum512(payloadId)
	if xsk, err := xipher.SecretKeyFromSeed(seed); err != nil {
		return fmt.Errorf("error writing data to keystore: %w", err)
	} else {
		if !commons.DirExists(keyStoreDir) {
			if err := os.MkdirAll(keyStoreDir, 0755); err != nil {
				return fmt.Errorf("error creating keystore directory: %w", err)
			}
		}
		if ct, err := xsk.Encrypt(payload, true, false); err != nil {
			return fmt.Errorf("error encrypting data for keystore: %w", err)
		} else {
			storePath := filepath.Join(keyStoreDir, storeId)
			if err := os.WriteFile(storePath, ct, 0644); err != nil {
				return fmt.Errorf("error writing data to keystore: %w", err)
			}
			return nil
		}
	}
}

func getFromKeyStore(storeId string, payloadId []byte) ([]byte, error) {
	seed := sha512.Sum512(payloadId)
	if xsk, err := xipher.SecretKeyFromSeed(seed); err != nil {
		return nil, fmt.Errorf("error reading data from keystore: %w", err)
	} else {
		storePath := filepath.Join(keyStoreDir, storeId)
		if commons.FileExists(storePath) {
			if payload, err := os.ReadFile(storePath); err == nil && len(payload) > 0 {
				if decryptedPayload, err := xsk.Decrypt(payload); err != nil {
					return nil, fmt.Errorf("error reading data from keystore: %w", err)
				} else {
					return decryptedPayload, nil
				}
			} else if err != nil {
				return nil, fmt.Errorf("error reading data from keystore: %w", err)
			}
		}
		return nil, ErrNotFound
	}
}

func Put(id, payload []byte, useLocalStore bool) error {
	sha256sum := sha256.Sum256(id)
	storeId := commons.Encode(sha256sum[:])
	if err := keyring.Set(keyringServiceName, storeId, string(payload)); err == nil {
		return nil
	} else if useLocalStore {
		return putToKeyStore(storeId, id, payload)
	} else {
		return fmt.Errorf("error saving data to keystore: %w", err)
	}
}

func Get(id []byte, useLocalStore bool) ([]byte, error) {
	sha256sum := sha256.Sum256(id)
	storeId := commons.Encode(sha256sum[:])
	if payloadStr, err := keyring.Get(keyringServiceName, storeId); err == nil {
		return []byte(payloadStr), nil
	} else if err == keyring.ErrNotFound {
		return nil, ErrNotFound
	} else if useLocalStore {
		return getFromKeyStore(storeId, id)
	} else {
		return nil, fmt.Errorf("error reading data from keystore: %w", err)
	}
}

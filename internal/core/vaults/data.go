package vaults

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"slv.sh/slv/internal/core/crypto"
)

func (vlt *Vault) putWithoutCommit(name string, value []byte, encrypt bool) (err error) {
	if !secretNameRegex.MatchString(name) {
		return errInvalidVaultItemName
	}
	var finalValue string
	if encrypt {
		var vaultPublicKey *crypto.PublicKey
		if vaultPublicKey, err = vlt.getPublicKey(); err == nil {
			var sealedSecret *crypto.SealedSecret
			if sealedSecret, err = vaultPublicKey.EncryptSecret(value, vlt.Spec.Config.Hash); err == nil {
				finalValue = sealedSecret.String()
			}
		}
	} else {
		finalValue = string(value)
	}
	if err == nil {
		if vlt.Spec.Data == nil {
			vlt.Spec.Data = make(map[string]string)
		}
		vlt.Spec.Data[name] = finalValue
	}
	return
}

func (vlt *Vault) Put(name string, value []byte, encrypt bool) (err error) {
	if !vlt.Spec.writable {
		return errVaultNotWritable
	}
	if err = vlt.putWithoutCommit(name, value, encrypt); err == nil {
		err = vlt.commit()
	}
	return
}

func (vlt *Vault) Import(importData []byte, force, encrypt bool) (err error) {
	if !vlt.Spec.writable {
		return errVaultNotWritable
	}
	dataMap := make(map[string]string)
	if err = yaml.Unmarshal(importData, &dataMap); err != nil {
		if envMap, err := godotenv.Unmarshal(string(importData)); err != nil {
			return errInvalidImportDataFormat
		} else {
			dataMap = envMap
		}
	}
	if !force {
		for name := range dataMap {
			if vlt.ItemExists(name) {
				return fmt.Errorf("the name %s already exists", name)
			}
		}
	}
	for name, value := range dataMap {
		if err = vlt.putWithoutCommit(name, []byte(value), encrypt); err != nil {
			return err
		}
	}
	return vlt.commit()
}

func (vlt *Vault) ItemExists(name string) (exists bool) {
	if vlt.Spec.Data != nil {
		_, exists = vlt.Spec.Data[name]
	}
	return
}

func (vlt *Vault) GetItemNames() (itemNames []string) {
	if vlt.Spec.Data != nil {
		for name := range vlt.Spec.Data {
			itemNames = append(itemNames, name)
		}
	}
	return
}

func (vlt *Vault) GetAllItems() (map[string]*VaultItem, error) {
	itemMap := make(map[string]*VaultItem)
	for name := range vlt.Spec.Data {
		if data, err := vlt.Get(name); err == nil {
			itemMap[name] = data
		} else {
			return nil, fmt.Errorf("error retrieving '%s': %w", name, err)
		}
	}
	return itemMap, nil
}

func (vlt *Vault) GetAllValues() (map[string][]byte, error) {
	itemValueMap := make(map[string][]byte)
	for name := range vlt.Spec.Data {
		if item, err := vlt.Get(name); err == nil {
			if itemValue, err := item.Value(); err == nil {
				itemValueMap[name] = itemValue
			} else {
				return nil, fmt.Errorf("error retrieving '%s': %w", name, err)
			}
		} else {
			return nil, fmt.Errorf("error retrieving '%s': %w", name, err)
		}
	}
	return itemValueMap, nil
}

func (vlt *Vault) Get(name string) (*VaultItem, error) {
	rawValue, ok := vlt.Spec.Data[name]
	if !ok {
		return nil, errVaultItemNotFound
	}
	item := vlt.getFromCache(name)
	if item == nil {
		item = &VaultItem{
			vlt:      vlt,
			rawValue: rawValue,
		}
		sealedSecret := &crypto.SealedSecret{}
		if err := sealedSecret.FromString(rawValue); err == nil {
			item.encryptedAt = new(time.Time)
			*item.encryptedAt = sealedSecret.EncryptedAt()
			if sealedSecret.Hash() != "" {
				item.hash = sealedSecret.Hash()
			}
		} else {
			item.plaintext = true
		}
		vlt.putToCache(name, item)
	}
	return item, nil
}

func (vlt *Vault) GetValue(name string) ([]byte, error) {
	item, err := vlt.Get(name)
	if err != nil {
		return nil, err
	}
	return item.Value()
}

func (vlt *Vault) DeleteItem(name string) error {
	return vlt.DeleteItems([]string{name})
}

func (vlt *Vault) DeleteItems(names []string) error {
	if !vlt.Spec.writable {
		return errVaultNotWritable
	}
	for _, name := range names {
		delete(vlt.Spec.Data, name)
		vlt.deleteFromCache(name)
	}
	return vlt.commit()
}

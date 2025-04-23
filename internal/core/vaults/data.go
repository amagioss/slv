package vaults

import (
	"fmt"
	"time"

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
	if err = vlt.putWithoutCommit(name, value, encrypt); err == nil {
		err = vlt.commit()
	}
	return
}

func (vlt *Vault) Import(importData []byte, force, encrypt bool) (err error) {
	dataMap := make(map[string]string)
	if err = yaml.Unmarshal(importData, &dataMap); err != nil {
		return errInvalidImportDataFormat
	}
	if !force {
		for name := range dataMap {
			if vlt.Exists(name) {
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

func (vlt *Vault) Exists(name string) (exists bool) {
	if vlt.Spec.Data != nil {
		_, exists = vlt.Spec.Data[name]
	}
	return exists
}

func (vlt *Vault) List(decrypt bool) (map[string]*VaultItem, error) {
	itemMap := make(map[string]*VaultItem)
	for name := range vlt.Spec.Data {
		if data, err := vlt.get(name, decrypt); err == nil {
			itemMap[name] = data
		} else {
			return nil, fmt.Errorf("error retrieving '%s': %w", name, err)
		}
	}
	return itemMap, nil
}

func (vlt *Vault) GetAllValues() (map[string][]byte, error) {
	if vaultItemMap, err := vlt.List(true); err == nil {
		valuesMap := make(map[string][]byte)
		for name, data := range vaultItemMap {
			valuesMap[name] = data.value
		}
		return valuesMap, nil
	} else {
		return nil, err
	}
}

func (vlt *Vault) get(name string, decrypt bool) (*VaultItem, error) {
	rawValue := vlt.Spec.Data[name]
	if rawValue == "" {
		return nil, errVaultItemNotFound
	}
	item := vlt.getFromCache(name)
	if item == nil {
		item = &VaultItem{}
		sealedSecret := &crypto.SealedSecret{}
		if err := sealedSecret.FromString(rawValue); err == nil {
			item.isSecret = true
			if decrypt {
				if vlt.IsLocked() {
					return nil, errVaultLocked
				}
				if item.value, err = vlt.Spec.secretKey.DecryptSecret(*sealedSecret); err != nil {
					return nil, err
				}
			}
			item.updatedAt = new(time.Time)
			*item.updatedAt = sealedSecret.EncryptedAt()
			if sealedSecret.Hash() != "" {
				item.hash = sealedSecret.Hash()
			}
		}
		if !item.isSecret {
			item.value = []byte(rawValue)
		}
		if decrypt || !item.isSecret {
			vlt.putToCache(name, item)
		}
	}
	return item, nil
}

func (vlt *Vault) Get(name string) (item *VaultItem, err error) {
	return vlt.get(name, true)
}

func (vlt *Vault) IsSecret(name string) (isSecret bool, err error) {
	if data, err := vlt.get(name, false); err == nil {
		isSecret = data.isSecret
	}
	return
}

func (vlt *Vault) DeleteItem(name string) error {
	return vlt.DeleteItems([]string{name})
}

func (vlt *Vault) DeleteItems(names []string) error {
	for _, name := range names {
		delete(vlt.Spec.Data, name)
		vlt.deleteFromCache(name)
	}
	return vlt.commit()
}

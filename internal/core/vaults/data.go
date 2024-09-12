package vaults

import (
	"fmt"

	"gopkg.in/yaml.v3"
	"oss.amagi.com/slv/internal/core/crypto"
)

func (vlt *Vault) putWithoutCommit(name string, value []byte, encrypt bool) (err error) {
	if !secretNameRegex.MatchString(name) {
		return errInvalidVaultDataName
	}
	var finalValue string
	if encrypt {
		var vaultPublicKey *crypto.PublicKey
		if vaultPublicKey, err = vlt.getPublicKey(); err == nil {
			var sealedSecret *crypto.SealedSecret
			if sealedSecret, err = vaultPublicKey.EncryptSecret(value, vlt.Config.Hash); err == nil {
				finalValue = sealedSecret.String()
			}
		}
	} else {
		finalValue = string(value)
	}
	if err == nil {
		if vlt.Data == nil {
			vlt.Data = make(map[string]string)
		}
		vlt.Data[name] = finalValue
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
	if vlt.Data != nil {
		_, exists = vlt.Data[name]
	}
	return exists
}

func (vlt *Vault) List() (secretsMap map[string]*crypto.SealedSecret, plaintextMap map[string]string) {
	secretsMap = make(map[string]*crypto.SealedSecret)
	plaintextMap = make(map[string]string)
	for name, value := range vlt.Data {
		sealedSecret := crypto.SealedSecret{}
		if err := sealedSecret.FromString(value); err == nil {
			secretsMap[name] = &sealedSecret
		} else {
			plaintextMap[name] = value
		}
	}
	return secretsMap, plaintextMap
}

func (vlt *Vault) GetAll() (map[string]*VaultData, error) {
	dataMap := make(map[string]*VaultData)
	for name := range vlt.Data {
		if data, err := vlt.Get(name); err == nil {
			dataMap[name] = data
		} else {
			return nil, fmt.Errorf("error retrieving secret %s: %w", name, err)
		}
	}
	return dataMap, nil
}

func (vlt *Vault) GetAllValues() (map[string][]byte, error) {
	if vaultDataMap, err := vlt.GetAll(); err == nil {
		valuesMap := make(map[string][]byte)
		for name, data := range vaultDataMap {
			valuesMap[name] = data.value
		}
		return valuesMap, nil
	} else {
		return nil, err
	}
}

func (vlt *Vault) Get(name string) (data *VaultData, err error) {
	rawValue := vlt.Data[name]
	if rawValue == "" {
		return nil, errVaultDataNotFound
	}
	data = vlt.getFromCache(name)
	if data == nil {
		data = &VaultData{}
		sealedSecret := &crypto.SealedSecret{}
		if err = sealedSecret.FromString(rawValue); err == nil {
			data.isSecret = true
			if vlt.IsLocked() {
				return nil, errVaultLocked
			}
			if data.value, err = vlt.secretKey.DecryptSecret(*sealedSecret); err != nil {
				return nil, err
			}
		}
		if data.value == nil || err != nil {
			data.value = []byte(rawValue)
			err = nil
		}
		vlt.putToCache(name, data)
	}
	return
}

func (vlt *Vault) DeleteItem(name string) error {
	return vlt.DeleteItems([]string{name})
}

func (vlt *Vault) DeleteItems(names []string) error {
	for _, name := range names {
		delete(vlt.Data, name)
		vlt.deleteFromCache(name)
	}
	return vlt.commit()
}

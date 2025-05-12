package vaults

import (
	"time"

	"slv.sh/slv/internal/core/crypto"
)

type VaultItem struct {
	value       []byte     `json:"-"`
	rawValue    string     `json:"-"`
	isSecret    bool       `json:"-"`
	encryptedAt *time.Time `json:"-"`
	hash        string     `json:"-"`
	vlt         *Vault     `json:"-"`
}

func (vi *VaultItem) Vault() *Vault {
	return vi.vlt
}

func (vi *VaultItem) Value() (value []byte, err error) {
	if vi.value == nil {
		if vi.isSecret {
			if vi.vlt.IsLocked() {
				return nil, errVaultLocked
			}
			sealedSecret := &crypto.SealedSecret{}
			if err = sealedSecret.FromString(vi.rawValue); err == nil {
				vi.value, err = vi.vlt.Spec.secretKey.DecryptSecret(*sealedSecret)
			}
		} else {
			vi.value = []byte(vi.rawValue)
		}
	}
	return vi.value, err
}

func (vi *VaultItem) ValueString() (value string, err error) {
	if valueBytes, err := vi.Value(); err == nil {
		value = string(valueBytes)
	}
	return
}

func (vi *VaultItem) IsSecret() bool {
	return vi.isSecret
}

func (vi *VaultItem) EncryptedAt() *time.Time {
	return vi.encryptedAt
}

func (vi *VaultItem) Hash() string {
	return vi.hash
}

func (vi *VaultItem) String() string {
	return vi.rawValue
}

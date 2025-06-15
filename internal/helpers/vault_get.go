package helpers

import (
	"encoding/base64"
	"fmt"
	"time"

	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
	"slv.sh/slv/internal/core/vaults"
)

type VaultInfo struct {
	Name      string                        `json:"name,omitempty"`
	PublicKey string                        `json:"publicKey,omitempty"`
	Locked    bool                          `json:"locked,omitempty"`
	Accessors map[string]*VaultAccessorInfo `json:"accessors,omitempty"`
	Items     map[string]VaultItemInfo      `json:"items,omitempty"`
}

type VaultAccessorInfo struct {
	Name  string   `json:"name,omitempty"`
	Email string   `json:"email,omitempty"`
	Tags  []string `json:"tags,omitempty"`
	Ref   string   `json:"ref,omitempty"`
}

type VaultItemInfo struct {
	Value       string `json:"value,omitempty"`
	Plaintext   bool   `json:"plaintext,omitempty"`
	EncryptedAt string `json:"encryptedAt,omitempty"`
	Hash        string `json:"hash,omitempty"`
}

func GetVaultInfo(vault *vaults.Vault, includeAccessors, encodeValuesToBase64 bool) *VaultInfo {
	info := &VaultInfo{
		Name:      vault.Name,
		PublicKey: vault.Spec.Config.PublicKey,
		Locked:    vault.IsLocked(),
	}
	if includeAccessors {
		self := environments.GetSelf()
		profile, _ := profiles.GetActiveProfile()
		var root *environments.Environment
		if profile != nil {
			root, _ = profile.GetRoot()
		}
		accessors, err := vault.ListAccessors()
		if err != nil {
			return nil
		}
		if len(accessors) > 0 {
			info.Accessors = make(map[string]*VaultAccessorInfo)
		}
		for _, accessor := range accessors {
			accessorPubKey, err := accessor.String()
			if err != nil {
				return nil
			}
			var accessorInfo *VaultAccessorInfo
			if root != nil && root.PublicKey == accessorPubKey {
				accessorInfo = &VaultAccessorInfo{
					Name:  root.Name,
					Ref:   fmt.Sprintf("Root (from %s)", profile.Name()),
					Email: root.Email,
					Tags:  root.Tags,
				}
			} else if self != nil && self.PublicKey == accessorPubKey {
				accessorInfo = &VaultAccessorInfo{
					Name:  self.Name,
					Ref:   "Self (Current User)",
					Email: self.Email,
					Tags:  self.Tags,
				}
			} else if profile != nil {
				env, _ := profile.GetEnv(accessorPubKey)
				if env != nil {
					accessorInfo = &VaultAccessorInfo{
						Name:  env.Name,
						Email: env.Email,
						Tags:  env.Tags,
					}
					if env.EnvType == environments.USER {
						accessorInfo.Ref = fmt.Sprintf("User (from %s)", profile.Name())
					} else {
						accessorInfo.Ref = fmt.Sprintf("Service (from %s)", profile.Name())
					}
				}
			}
			if accessorInfo == nil {
				accessorInfo = &VaultAccessorInfo{
					Name: accessorPubKey,
					Ref:  "Unknown Accessor",
				}
			}
			info.Accessors[accessorPubKey] = accessorInfo
		}
	}
	items, err := vault.GetAllItems()
	if err != nil {
		return nil
	}
	if len(items) > 0 {
		info.Items = make(map[string]VaultItemInfo)
	}
	for itemName, item := range items {
		itemInfo := VaultItemInfo{
			Plaintext: item.IsPlaintext(),
			Hash:      item.Hash(),
		}
		if !item.IsPlaintext() {
			itemInfo.EncryptedAt = item.EncryptedAt().Format(time.RFC3339)
		}
		if item.IsPlaintext() || !vault.IsLocked() {
			itemValue, err := item.Value()
			if err != nil {
				return nil
			}
			if encodeValuesToBase64 {
				itemInfo.Value = base64.StdEncoding.EncodeToString(itemValue)
			} else {
				itemInfo.Value = string(itemValue)
			}
		} else {
			itemInfo.Value = "(Locked)"
		}
		info.Items[itemName] = itemInfo
	}
	return info
}

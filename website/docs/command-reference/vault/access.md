---
sidebar_position: 6
---

# Manage Vault Access
Add or Remove access to the vault.

> **Before you begin:** You need to have a vault created and an environment with access to that vault. The environment managing access to a vault must be able to access the vault in the first place.

#### General Usage:
```bash
slv vault --vault <PATH_TO_VAULT> access [flags]
slv vault --vault <PATH_TO_VAULT> access [flags] [command]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --env-self | None | NA | NA | Modify vault access for the environment set to `self` |
| --env-k8s | None | NA | NA | Modify vault access for the environment in current kubernetes context |
| --env-pubkey | String(s) | False | None | Modify vault access for the environment with given Public Keys |
| --env-search | String(s) | False | None | Share vault with environment based on search string |
| --quantum-safe | None | NA | NA | Use Quantum Resistant Cryptography (Kyber1024) |
| --vault | String | True | NA | Path to the SLV Vault file |
| --help | None | NA | NA | Help text for `slv vault access` |

---
## Grant Access to a Vault
The canonical command is `grant` (aliases: `add`, `allow`, `share`).
#### Usage:
```bash
slv vault --vault <PATH_TO_VAULT> access --env-search <SEARCH_STRING> grant
```
#### Example:
```bash
$ slv vault --vault test.slv.yaml access --env-search alice grant
Added vault access: test.slv.yaml
```
---
## Revoke Access to a Vault
The canonical command is `rm` (aliases: `remove`, `revoke`, `deny`, `del`, `delete`).
#### Usage:
```bash
slv vault --vault <PATH_TO_VAULT> access --env-search <SEARCH_STRING> rm
```
#### Example:
```bash
$ slv vault --vault test.slv.yaml access --env-search bob@example.com rm
Revoked vault access: test.slv.yaml
```

---

## See Also

- [Create a New Vault](/docs/command-reference/vault/new) - Create a new vault
- [Get a Secret](/docs/command-reference/vault/get) - Retrieve secrets from your vault
- [Environment Component](/docs/components/environment) - Learn about environments
- [Vault Component](/docs/components/vault) - Learn more about vaults

---
sidebar_position: 6
---

# Manage Vault Access
Add or Remove access to the vault.\
**Important Condition:** The environment managing access to a vault must be able to access the vault in the first place.

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
## Add Access to a Vault
#### Usage:
```bash
slv vault --vault <PATH_TO_VAULT> access --env-search <SEARCH_STRING> add
```
#### Example:
```bash
$ slv vault --vault test.slv.yaml access --env-search alice add
Shared vault: test.slv.yaml
```
---
## Remove Access to a Vault
#### Usage:
```bash
slv vault --vault <PATH_TO_VAULT> access --env-search <SEARCH_STRING> remove
```
#### Example:
```bash
$ slv vault --vault test.slv.yaml access --env-search bob@example.com remove
Shared vault: test.slv.yaml
```

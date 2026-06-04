---
sidebar_position: 5
---
# Remove Secrets 
Delete one or more items from a vault.
#### General Usage:
```bash
slv vault --vault <PATH_TO_VAULT> rm [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --name | String (comma-separated or repeated) | True | NA | Name(s) of the item(s) (keys) to delete |
| --vault | String | True | NA | Path to the SLV Vault file |
| --help | None | NA | NA | Help text for `slv vault rm` |

#### Usage:
```bash
slv vault --vault <PATH_TO_VAULT> rm --name <ITEM_KEY>[,<ITEM_KEY>...]
```
#### Example:
```bash
$ slv vault --vault test.slv.yaml rm --name my_secret
Successfully deleted the secrets: [my_secret] from the vault: test.slv.yaml

$ slv vault --vault test.slv.yaml rm --name secret1,secret2
Successfully deleted the secrets: [secret1 secret2] from the vault: test.slv.yaml
```

---

## See Also

- [Put a Secret](/docs/command-reference/vault/put) - Add secrets to your vault
- [Get a Secret](/docs/command-reference/vault/get) - Retrieve secrets from your vault
- [Update Vault Attributes](/docs/command-reference/vault/update) - Update vault metadata
- [Vault Component](/docs/components/vault) - Learn more about vaults

---
sidebar_position: 5
---
# Remove a Secret 
Delete an item from a vault.
#### General Usage:
```bash
slv vault --vault <PATH_TO_VAULT> rm [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --name | String | False | None | Name of the item (key) to delete |
| --vault | String | True | NA | Path to the SLV Vault file |
| --help | None | NA | NA | Help text for `slv vault rm` |

#### Usage:
```bash
slv vault --vault <PATH_TO_VAULT> rm --name <ITEM_KEY>
```
#### Example:
```bash
$ slv vault --vault test.slv.yaml rm --name my_secret
Successfully deleted the secrets: [my_secret] from the vault: test.slv.yaml
```

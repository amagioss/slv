---
sidebar_position: 7
---

# Referencing
This is particularly useful when you have files with secrets. The command would replace the actual secrets with references to SLV secrets and add the SLV secrets to an SLV vault.
#### General Usage:
```bash
slv vault --vault <PATH_TO_VAULT> ref [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --file | String | True | NA | The reference YAML/JSON/BLOB file (Needs to be flat) |
| --name | String | False | None | Name of the item (key) to reference. References all keys if not provided |
| --force | None | NA | NA | Overwrite the item if it already exists |
| --preview | None | NA | NA | Dry Run - Show what the referenced file would look like after referencing |
| --vault | String | True | NA | Path to the SLV Vault file |
| --help | None | NA | NA | Help text for `slv vault ref` |

#### Usage:
```bash
slv vault --vault <PATH_TO_VAULT> ref --file <PATH_TO_FILE_TO_BE_REFERENCED>
```
#### Example:
```bash
$ cat secrets.yaml 
username: johndoe
password: super_secret_password

$ slv vault --vault test.slv.yaml ref --file secrets.yaml
Auto referenced secrets.yaml (YAML) with vault test.slv.yaml

$ cat secrets.yaml 
password: '{{SLV.test.password}}'
username: '{{SLV.test.username}}'
```

---

## See Also

- [Dereferencing](/docs/command-reference/vault/deref) - Replace references with actual secrets
- [Put a Secret](/docs/command-reference/vault/put) - Add secrets to your vault
- [Get a Secret](/docs/command-reference/vault/get) - Retrieve secrets from your vault
- [Vault Component](/docs/components/vault) - Learn more about vaults

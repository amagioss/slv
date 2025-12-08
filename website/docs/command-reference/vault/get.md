---
sidebar_position: 4
---

# Get a Secret 
Retrieve one or all items in a vault.

> **Before you begin:** You need to have a vault with secrets and an environment with access to that vault. If you haven't set this up yet, see [Create a New Vault](/docs/command-reference/vault/new) and [Create a New Environment](/docs/command-reference/environment/new), or follow the [Quick Start Guide](/docs/quick-start).
#### General Usage:
```bash
slv vault --vault <PATH_TO_VAULT> get [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --name | String | False | None | Name of the item (key) |
| --format | String | False | None | List secrets as one of [`json`, `yaml`, `envar`] |
| --with-metadata | None | NA | NA | Print metadata of items when using `--format` |
| --base64 | None | NA | NA | Encode the item values as base64 |
| --vault | String | True | NA | Path to the SLV Vault file or Vault URL|
| --help | None | NA | NA | Help text for `slv vault get` |
---
## Get everything from the Vault
#### Usage:
```bash
slv vault --vault <PATH_TO_VAULT> get
```
#### Example:
```bash
slv vault --vault test.slv.yaml get
Enter Password: 
Do you want to save the password in keyring? (y/n): y
Vault ID:  SLV_VPK_AEAVMAAAACYH33FBSJWDB7R4QUGQMPBX2F4DZLWC5LLZIAWSA7EQPDEYEP7A6
Vault Data:
+-----------+-------------------------+--------+----------------------+
|   NAME    |          VALUE          |  TYPE  |      UPDATED AT      |
+-----------+-------------------------+--------+----------------------+
| my_secret | this_is_super_sensitive | Secret | 25-Apr-2025 14:52:42 |
| password  | super_secret_password   | Secret | 25-Apr-2025 14:55:45 |
| username  | johndoe                 | Secret | 25-Apr-2025 14:55:45 |
+-----------+-------------------------+--------+----------------------+
Accessible by:
+-----------------------------------------------------------------------+------+----------------------+
|                              PUBLIC KEY                               | TYPE |         NAME         |
+-----------------------------------------------------------------------+------+----------------------+
| SLV_EPK_AEAUKAAAACRIHZIK3U46HKV7PML7VIY4JXO2FYNTNCVKNN23U2LNTZCYTJQGY | Self | John Doe |
| SLV_EPK_AEAUKAAAABQHMUEM6YBE6D63FYAPYNZXJD3LSUJRVPMG7SOGAYTUQ4XAFJ6EQ | Root | root           |
+-----------------------------------------------------------------------+------+----------------------+
```
---
## Get a specific secret
#### Usage:
```bash
slv vault --vault <PATH_TO_VAULT> get --with-metadata --name <ITEM_KEY>
```
#### Example:
```bash
$ slv vault --vault test.slv.yaml get --name password
Enter Password: 
Do you want to save the password in keyring? (y/n): y
super_secret_password
```
---
## Get items in a specific format with metadata
#### Usage:
```bash
slv vault --vault <PATH_TO_VAULT> get --format [yaml/json] --with-metadata
```
#### Example:
```bash
$ slv vault --vault test.slv.yaml get --format yaml --with-metadata
my_secret:
    value: this_is_super_sensitive
    secret: true
    updatedAt: "2025-04-25T14:52:42+05:30"
password:
    value: super_secret_password
    secret: true
    updatedAt: "2025-04-25T14:55:45+05:30"
username:
    value: johndoe
    secret: true
    updatedAt: "2025-04-25T14:55:45+05:30"
```

---

## See Also

- [Put a Secret](/docs/command-reference/vault/put) - Add secrets to your vault
- [Create a New Vault](/docs/command-reference/vault/new) - Create a new vault
- [Vault Component](/docs/components/vault) - Learn more about vaults
- [Load Vault as Environment Variables](/docs/command-reference/vault/run) - Use the vault secrets as environment variables
- [Dereference Secrets](/docs/command-reference/vault/deref) - Replace references with actual values

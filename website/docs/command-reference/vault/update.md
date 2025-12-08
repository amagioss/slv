---
sidebar_position: 10
---

# Update Vault Attributes
Update the metadata of a vault
#### General Usage:
```bash
slv vault --vault <PATH_TO_VAULT> update [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --name | String | False | None | The name to update the vault name with |
| --k8s-namespace | String | False | None | Namespace for the K8S Custom Resource |
| --vault | String | True | NA | Path to the SLV Vault file |
| --help | None | NA | NA | Help text for `slv vault update` |

When none of the flags above are given, the command would update the vault structure into a K8S compatible resource YAML.

#### Usage
```bash
slv vault --vault test.slv.yaml update --name <NEW_VAULT_NAME> --k8s-namespace <NAMESPACE>
```

#### Example:
```bash
$ slv vault --vault test.slv.yaml update --name new_vault --k8s-namespace slv
Vault test.slv.yaml transformed to K8s resource new_vault
```

---

## See Also

- [Create a New Vault](/docs/command-reference/vault/new) - Create a new vault
- [Put a Secret](/docs/command-reference/vault/put) - Add secrets to your vault
- [Kubernetes Operator](/docs/extensions/slv-in-kubernetes/operator) - Use SLV with Kubernetes
- [Vault Component](/docs/components/vault) - Learn more about vaults

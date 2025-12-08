---
sidebar_position: 9
---

# Load Vault to Environment Variables
Launch a shell or run a command with the secrets loaded as Environment Variables in it.
#### General Usage:
```bash
slv vault --vault <PATH_TO_VAULT> run [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --command | String | False | $SHELL | The command to be run |
| --prefix | String | False | None | Prefix to add to the ENVAR |
| --vault | String | True | NA | Path to the SLV Vault file or Vault URL |
| --help | None | NA | NA | Help text for `slv vault run` |

#### Usage:
```bash
slv vault --vault <PATH_TO_VAULT> run -c <COMMAND_TO_RUN>
```

#### Example:
```bash
$ slv vault --vault test.slv.yaml run -c bash --prefix SLV_ENV_VAR_
Running command [bash] with secrets loaded into environment variables from the vault test.slv.yaml...
Please note that the secret names are prefixed with SLV_ENV_VAR_

$ env | grep SLV_ENV_VAR
SLV_ENV_VAR_username=johndoe
SLV_ENV_VAR_password=super_secret_password

$ exit
exit
Command [bash] ended successfully
```

---

## See Also

- [Get a Secret](/docs/command-reference/vault/get) - Retrieve secrets from your vault
- [Vault Component](/docs/components/vault) - Learn more about vaults
- [GitHub Actions Integration](/docs/extensions/slv-in-github-actions) - Use SLV in CI/CD pipelines

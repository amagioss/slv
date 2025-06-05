---
sidebar_position: 7
---

# Github Actions

## Inputs
- `version`: The version of slv to install, defaulting to `latest`
- `vault`: Path to the vault file
- `env-secret-key`: The SLV environment secret (key/binding) to use for the action
- `prefix`: Prefix to use for the environment variable names along with the SLV secret name

## Use Cases
### Set Up SLV CLI
You can use the action to set up SLV CLI on the runner.
```yaml
steps:
- name: Setup SLV
  uses: amagioss/slv@v<MAJOR_VERSION>
```

#### Install a Specific Version
```yaml
steps:
- name: Setup SLV
  uses: amagioss/slv@v<MAJOR_VERSION>
  with:
    version: 0.16.3
```

### Load SLV Secrets Into Environment Variables
You can use the action to load secrets from a vault into environment variables that can further be consumed by other actions or programs.
```yaml
steps:
- name: Load SLV Secrets
  uses: amagioss/slv@v<MAJOR_VERSION>
  with:
    vault: creds.slv.yaml
    env-secret-key: ${{ secrets.SLV_ENV_SECRET_KEY }}
```

#### Set a Prefix for Variables
If you'd like to set a prefix across all the environment variables created by SLV, you can do so by specifying the `prefix` parameter.
```yaml
steps:
- name: Load SLV Secrets - PROD
  uses: amagioss/slv@v<MAJOR_VERSION>
  with:
    vault: creds.slv.yaml
    env-secret-key: ${{ secrets.SLV_ENV_SECRET_KEY }}
    prefix: "PROD_"
```

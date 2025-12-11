---
sidebar_position: 1
---

# List Vaults

List all vault files in a directory.

## Overview

The `list` command scans a directory for SLV vault files (`.slv.yaml` or `.slv.yml`) and displays them as a simple list or detailed table. Useful for:

- Finding vaults in a new project
- Getting an overview of all vaults
- Searching nested directories
- Checking vault metadata without unlocking

## General Usage

```bash
slv vault list [flags]
```

### Aliases

The following aliases are available:
- `slv vault list`
- `slv vault ls`
- `slv vault find`
- `slv vault search`

## Flags

| Flag | Shorthand | Arguments | Default | Description |
| -- | -- | -- | -- | -- |
| --dir | -d | String | Current directory | Directory to search for vaults |
| --recursive | -r | None | false | Search recursively in subdirectories |
| --details | -l | None | false | Show detailed information (name, namespace, secret count, accessors) |
| --help | -h | None | NA | Help text for `slv vault list` |

---

## How It Works

### Basic List

By default, the command lists vault files in the current directory:

```bash
slv vault list
```

**Output:**
```
Found 3 vault(s):

  • dev.slv.yaml
  • staging.slv.yaml
  • prod.slv.yaml

Use 'slv vault list --details' to see more information.
```

### Recursive Search

Search for vaults in all subdirectories:

```bash
slv vault list --recursive
```

**Output:**
```
Found 5 vault(s):

  • dev.slv.yaml
  • prod.slv.yaml
  • team/backend.slv.yaml
  • team/frontend.slv.yaml
  • secrets/api-keys.slv.yaml

Use 'slv vault list --details' to see more information.
```

### Detailed View

Show metadata for each vault:

```bash
slv vault list --recursive --details
```

**Short form:**
```bash
slv vault ls -r -l
```

**Output:**
```
Found 5 vault(s):

┌──────────────────────────┬───────────┬───────────┬─────────┬───────────┐
│ VAULT FILE               │ NAME      │ NAMESPACE │ SECRETS │ ACCESSORS │
├──────────────────────────┼───────────┼───────────┼─────────┼───────────┤
│ dev.slv.yaml             │ dev-env   │ default   │      12 │         3 │
│ prod.slv.yaml            │ prod-env  │ prod      │       8 │         2 │
│ team/backend.slv.yaml    │ backend   │ team      │      15 │         5 │
│ team/frontend.slv.yaml   │ frontend  │ team      │      10 │         4 │
│ secrets/api-keys.slv.yaml│ api-keys  │ -         │       6 │         2 │
└──────────────────────────┴───────────┴───────────┴─────────┴───────────┘
```

:::tip
To view the list of accessors for a specific vault, use:
```bash
slv vault -v /path/to/vault.slv.yaml
```
:::

### Specify Directory

List vaults in a specific directory:

```bash
slv vault list --dir /path/to/project
```

**Or with a relative path:**
```bash
slv vault ls -d ~/projects/myapp -r
```

---

## Examples

### Example 1: Quick Vault Discovery

Find all vaults in your current project:

```bash
cd /path/to/project
slv vault ls -r
```

```
Found 3 vault(s):

  • config/dev.slv.yaml
  • config/prod.slv.yaml
  • infrastructure/secrets.slv.yaml

Use 'slv vault list --details' to see more information.
```

### Example 2: Audit Vault Statistics

Get detailed statistics about all vaults:

```bash
slv vault list -r -l
```

```
Found 3 vault(s):

┌─────────────────────────┬──────────┬───────────┬─────────┬───────────┐
│ VAULT FILE              │ NAME     │ NAMESPACE │ SECRETS │ ACCESSORS │
├─────────────────────────┼──────────┼───────────┼─────────┼───────────┤
│ config/dev.slv.yaml     │ dev      │ default   │      25 │         8 │
│ config/prod.slv.yaml    │ prod     │ prod      │      12 │         3 │
│ infrastructure/secrets..│ infra    │ -         │       5 │         2 │
└─────────────────────────┴──────────┴───────────┴─────────┴───────────┘
```

The output shows:
- 25 secrets in dev vault, accessible by 8 environments
- 12 secrets in prod vault, accessible by 3 environments
- 5 secrets in infrastructure vault, accessible by 2 environments

:::tip
The **ACCESSORS** column shows the count of environments/profiles that can decrypt this vault.
To see the actual accessor public keys, use `slv vault -v /path/to/vault.slv.yaml`
:::

### Example 3: Search Specific Directory

List vaults in a specific project without changing directories:

```bash
slv vault find --dir ~/projects/my-app --recursive
```

### Example 4: CI/CD Validation

Check if vaults exist in your repository:

```bash
if slv vault list -r | grep -q "Found 0 vault"; then
  echo "No vaults found - please create vaults"
  exit 1
fi
```

---

## Use Cases

### 1. New Project Onboarding

List all vaults when joining a new project:

```bash
slv vault list -r -l
```

Shows all vaults with secret counts and environment access.

### 2. Vault Organization

Find all vaults in the project:

```bash
slv vault find -r
```

Useful for identifying duplicate vaults and understanding vault structure.

### 3. Security Audit

Audit vault distribution and access:

```bash
slv vault ls -r -l
```

Review secret counts, accessor counts, and vault organization.

### 4. Documentation

Generate vault inventory for documentation:

```bash
slv vault list -r -l > vault-inventory.txt
```

---

## Features

- Automatically finds `.slv.yaml` and `.slv.yml` files
- Shows vault statistics without unlocking
- Recursive directory scanning
- Fast operation (no vault unlocking required)
- Works with absolute and relative paths
- Only displays metadata, never secret values

---

## Notes

- No unlocking required - the command shows basic information without decrypting vaults
- Locked vaults can still display name, namespace, secret count, and accessor count
- Corrupted vault files will be shown as "(Error loading)"
- Fast performance since it only reads metadata

---

## Related Commands

- [Create a New Vault](/docs/command-reference/vault/new) - Create vaults to list
- [Get a Secret](/docs/command-reference/vault/get) - View secrets in a specific vault
- [Vault Overview](/docs/components/vault) - Learn more about vaults

---

## Tips

- Use `ls -r -l` for a complete vault infrastructure overview
- Combine with grep for filtering: `slv vault ls -r | grep prod`
- Validate vault presence in scripts: `slv vault list -d ./config -r`

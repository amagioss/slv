---
sidebar_position: 4
---
# Sync Profile with Remote
Sync the local cache of your profile with the remote profile (git/http). SLV automatically syncs periodically as specified by the `--sync-interval` flag. You can use this command to manually sync profiles.
#### General Usage:
```bash
slv profile sync [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --help | None | NA | NA | Help text for `slv profile sync` |
#### Usage:
```bash
slv profile sync
```
#### Example:
```bash
$ slv profile sync
Profile test is updated from remote successfully
```

---

## See Also

- [New Profile](/docs/command-reference/profile/new) - Create a new profile
- [List Profiles](/docs/command-reference/profile/list) - View all available profiles
- [Profile Component](/docs/components/profile) - Learn more about profiles

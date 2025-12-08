---
sidebar_position: 3
---

# Set Active Profile
Set an already added SLV profile as the active profile.
#### General usage:
```bash
slv profile activate [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- | 
| --name | String | True | NA | Name of the profile to set |
| --help | None | NA | NA|Help text for `slv profile activate` |

#### Usage:
```bash
slv profile activate --name <SLV_PROFILE_NAME>
```
#### Example:
```bash
$ slv profile activate --name my_other_profile
Successfully set my_other_profile as active profile
```

---

## See Also

- [List Profiles](/docs/command-reference/profile/list) - View all available profiles
- [New Profile](/docs/command-reference/profile/new) - Create a new profile
- [Profile Component](/docs/components/profile) - Learn more about profiles

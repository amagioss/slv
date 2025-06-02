---
sidebar_position: 5
---
# Delete a profile
Deletes an existing profile. Note that you **cannot delete the active profile**.
#### General Usage:
```bash
slv profile rm [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --name | String | True | NA | Name of the profile to delete |
| --help | None | NA | NA|Help text for `slv profile delete` |

#### Usage:
```bash
slv profile rm --name <PROFILE_NAME>
```
#### Example:
```bash
$ slv profile rm --name my_org
Deleted profile:  my_org
```

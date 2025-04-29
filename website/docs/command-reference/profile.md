---
sidebar_position: 1
---

# Profile (Beta)
Documentation for the `slv profile` command.

Aliases: `profile`, `profiles`

#### Commands Available: 
- [`new`](#new-profile)
- [`list`](#list-profiles)
- [`set`](#set-default-profile)
- [`pull`](#pull-remote-changes-into-local-profile)
- [`push`](#push-local-changes-to-remote-profile)
- [`delete`](#delete-a-profile)

---
## New Profile
Used to create a new SLV profile or add an existing one from a git repository.

#### General Usage:
```bash
slv profile new [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --git | String | False | None |Git URI to clone the profile from |
| --git-branch | String | False | main | Git branch corresponding to the git URI |
| --name | String | True | NA | Name of the profile (Scoped Locally) |
| --help | None | NA | NA|Help text for `slv profile new` |

### Create a new local profile
- Typically used when you want to manage the profile locally
- Ideal for individuals (One Man Army)
#### Usage:
```bash
slv profile new --name <SLV_PROFILE_NAME>
```
#### Example:
```bash
$ slv profile new --name my_slv_profile
Created profile:  my_slv_profile
```

### Add a profile from a git repository
- Syncs the profile info with a git repository and environments are managed from Git
- Ideal for groups or organisations
#### Usage:
```bash
slv profile new --name <SLV_PROFILE_NAME> --git <GIT_URI_TO_CLONE_PROFILE_FROM> --git-branch <GIT_BRANCH_CORRESPONDING_TO_GIT_URI>
```
#### Example:
```bash
slv profile new --name my_org --git https://github.com/my_org/slvprofile.git --git-branch main
Enter the git username       : johndoe
Enter the git token/password : 
Created profile:  my_org
```

---

## List Profiles
Shows all the SLV profiles available on local

#### General Usage:
```bash
slv profile list [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --help | None | NA | NA|Help text for `slv profile list` |

#### Usage:
```bash
slv profile list
```
#### Example:
```bash
$ slv profile list
  my_org_slv_profile
* my_local_slv_profile
```
The `*` in the output represents the current profile.

---

## Set Default Profile
Set an already added SLV profile as the current profile.
#### General usage:
```bash
slv profile set [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- | 
| --name | String | True | NA | Name of the profile to set |
| --help | None | NA | NA|Help text for `slv profile set` |

#### Usage:
```bash
slv profile set --name <SLV_PROFILE_NAME>
```
#### Example:
```bash
$ slv profile set --name my_other_profile
Successfully set my_other_profile as current profile
```

---

## Pull remote changes into local profile
Pulls the latest changes for the current profile from remote repository. (Only works for profiles with remote repository set using `--git` flag.)
#### General Usage:
```bash
slv profile pull [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --help | None | NA | NA | Help text for `slv profile pull` |

#### Usage:
```bash
slv profile pull
```
#### Example:
```bash
$ slv profile pull
Enter the git username       : johndoe
Enter the git token/password : 
Successfully pulled changes into profile: my_org
```

---

## Push local changes to remote profile
Pushes the changes in the current profile to the pre-configured remote repository. (Only works for profiles with remote repository set using `--git` flag.)
#### General Usage:
```bash
slv profile push [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --help | None | NA | NA | Help text for `slv profile push` |
#### Usage:
```bash
slv profile push
```
#### Example:
```bash
$ slv profile push
Enter the git username       : johndoe
Enter the git token/password : 
Successfully pushed changes from profile: my_org
```

---

## Delete a profile
Deletes an existing profile. Note that you **cannot delete the current profile**.
#### General Usage:
```bash
slv profile delete [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --name | String | True | NA | Name of the profile to delete |
| --help | None | NA | NA|Help text for `slv profile delete` |

#### Usage:
```bash
slv profile delete --name <PROFILE_NAME>
```
#### Example:
```bash
$ slv profile delete --name my_org
Deleted profile:  my_org
```


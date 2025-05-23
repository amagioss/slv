---
sidebar_position: 1
---

# Profile
Documentation for the `slv profile` command.

Aliases: `profile`, `profiles`

#### Commands Available: 
- [`new`](#new-profile)
- [`list`](#list-profiles)
- [`activate`](#set-active-profile)
- [`sync`](#sync-profile-with-remote)
- [`delete`](#delete-a-profile)

---
## New Profile
Used to create a new SLV profile or add an existing one from a git repository.

#### General Usage:
```bash
slv profile new [command]
```
#### Commands Available:
- [`git`](#creating-a-git-based-profile)
- [`http`](#creating-a-http-url-based-profile)

### Creating a Git Based Profile
Use a remote git reposioty to maintain the profile. 
#### Usage:
```bash
slv profile new git [flags] 
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --repo | String | False | None | Git URI to clone the profile from |
| --branch | String | False | main | Git branch corresponding to the git URI |
| --token | String | False | None | The token to authenticate with the git repository over HTTP |
| --username | String | False | None | The username to authenticate with the git repository over HTTP |
| --name | String | True | NA | Name of the profile (Scoped Locally) |
| --ssh-key | String | False | None | The Path to private key to use for SSH |
| --read-only | None | NA | NA | Set profile as read-only |
| --sync-interval | Duration | False | `1h0m0s` | Profile sync interval |
| --help | None | NA | NA|Help text for `slv profile new git` |

#### Example:
```bash
$ slv profile new git --name test --repo git@github.com:username/reponame.git
Created profile test from remote (git)
```

### Creating a HTTP URL based Profile 
- Syncs the profile info with a remote HTTP URL
- Much faster than Git based profiles
#### Usage:
```bash
slv profile new http [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --auth-header | String | False | None | The header to be used for HTTP URLs protected by authentication |
| --name | String | True | NA | Name of the profile (Scoped Locally) |
| --url | String | True | NA | The HTTP base URL of the remote profile |
| --sync-interval | Duration | False | `1h0m0s` | Profile sync interval |
| --help | None | NA | NA|Help text for `slv profile new http` |
#### Example:
```bash
$ slv profile new http -n test --url https://example.com/slvprofile/
Created profile test from remote (http)
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
my_local_slv_profile
```
The active profile is shown in a different color.

---

## Set Active Profile
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

## Sync Profile with Remote
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

## Delete a profile
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


---
sidebar_position: 1
---

# New Profile
Used to create a new SLV profile or add an existing one from a git repository.

#### General Usage:
```bash
slv profile new [command]
```
#### Commands Available:
- [`git`](#creating-a-git-based-profile)
- [`http`](#creating-a-http-url-based-profile)

## Creating a Git Based Profile
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
---
## Creating a HTTP URL based Profile 
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

## See Also

- [List Profiles](/docs/command-reference/profile/list) - View all available profiles
- [Set Active Profile](/docs/command-reference/profile/set-active) - Switch between profiles
- [Sync Profile](/docs/command-reference/profile/sync) - Sync profile with remote
- [Profile Component](/docs/components/profile) - Learn more about profiles
- [Git Profile Integration](/docs/integration-guide/git-profile) - Use profiles with Git repositories

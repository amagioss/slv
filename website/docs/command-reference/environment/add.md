---
sidebar_position: 2
---

# Add an existing Environment

Used to add an environment that is created elsewhere to the existing machine. The Environment Definition String (`EDS`) can be used to do the same. 
> **Note:** This is not available when using **read-only** profiles.

#### General usage:
```bash
slv env add [flags]
```

#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --env-def | String(s) | True | NA | EDS for the environment to be added |
| --root | None | NA | NA| Set the environment as root environment for the active profile |
| --help | None | NA | NA| Help text for `slv env add` |

#### Usage:
```bash
slv env add --env-def <ENVIRONMENT_DEFINITION_STRING>
```

#### Example:
```bash
$ slv env add --env-def SLV_EDS_AF4JYRGM35FMGMA4YXYXOOOXYGFS2U6ISUEWMZUJWMS52H4HJCE7DRYIGS23IWRM4K5YWUGZ5X4RZPW7FA7V7GYUDVGRBKA6B22S4XJNWXODWKNFXTN24NXFU3YEPO7AYXETY4K33ENX7LP4WZMKZNBR67NLXNK5F22XE7D5HTWLSUJHGA6IMTAQUCXZBO4G5KA7UMKFAKB44IJVCCMJPW7ZOEK57447W3RW52XI4JQNRBPTADY7ZX6CBNBULMNHB6K5VN6UTYQYBH67AAAAB777TR3T5CA
Successfully added 1 environments to profile my_org
```

---

## See Also

- [Create a New Environment](/docs/command-reference/environment/new) - Create a new environment
- [List Environments](/docs/command-reference/environment/list) - View all available environments
- [Show Environment](/docs/command-reference/environment/show) - View environment details
- [Profile Component](/docs/components/profile) - Learn about profiles
- [Environment Component](/docs/components/environment) - Learn more about environments

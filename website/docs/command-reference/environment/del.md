---
sidebar_position: 6
---
# Delete an Environment

Delete one or more environments based on search parameters such as `name`, `email`, `tags`.

#### General Usage:
```bash
slv env rm [flags]
```

#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --env-search | String(s) | False | None | Search for environments based on `tag`/`email`/`name` to delete |
| --help | None | NA | NA| Help text for `slv env del` |

#### Usage:
```bash
slv env rm --env-search <SEARCH_STRING>
```
#### Example:
```bash
$ slv env rm --env-search alice
Public Key:      SLV_EPK_AEAUKAAAABUEMSPQ4BJIIWMSAKFUUXUV4THOP3ERH25CY4HR54W25HUJQR6XK
Name:            alice
Email:           alice@example.com
Tags:            [example_env]
Secret Binding:  SLV_ESB_AF4JYBGA2FXKUMAYADQHP6LPPMJM6JQDILREKS3KIHUQJWRZZKO5FQKCM4FHIRGR7DXPWHRQIACMHSNZVOOTJ7EDBGRAOODHEABHYIYY5O4Q23WD7TOZWSKIF66XVPZPLNXNTHJVNLVG3DH4N3237LZ7QMOJLGKSHHS6F7JQKOWCW3QLCYXLDYBBZFJ5ZRGVCEXXZVZYLR2ER33X3JLNJNHYICODVMPQ5VREN5GDSLSDLENJ6PUMFXKHZ5EHGOIGIT4TEW6LOYW6XMYR452BRPZSKKXLM5ZT7KHAVL64LPKNK45R3HAPH6IXAAAP772O3VBWC

Are you sure you wish to delete the above environment(s) [yes/no]: yes
Environment alice deleted successfully
```

---

## See Also

- [List Environments](/docs/command-reference/environment/list) - View all available environments
- [Create a New Environment](/docs/command-reference/environment/new) - Create a new environment
- [Environment Component](/docs/components/environment) - Learn more about environments

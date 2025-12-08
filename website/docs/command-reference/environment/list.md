---
sidebar_position: 3
---

# List Environments

Used to list all the profiles that are present in the active profile. The command can also be used to filter down results or search based on name, tag and email.

#### General usage:
```bash
slv env list [flags]
```

#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --env-search | String(s) | False | None | Search for environments based on `tag`/`email`/`name` |
| --help | None | NA | NA| Help text for `slv env get` |

#### Usage:
```bash
slv env list --env-search <SEARCH_STRING>
```

#### Example:
##### Without search
```bash
$ slv env list
Public Key:      SLV_EPK_AEAUKAAAABUEMSPQ4BJIIWMSAKFUUXUV4THOP3ERH25CY4HR54W25HUJQR6XK
Name:            alice
Email:           alice@example.com
Tags:            [example_env]
Secret Binding:  SLV_ESB_AF4JYBGA2FXKUMAYADQHP6LPPMJM6JQDILREKS3KIHUQJWRZZKO5FQKCM4FHIRGR7DXPWHRQIACMHSNZVOOTJ7EDBGRAOODHEABHYIYY5O4Q23WD7TOZWSKIF66XVPZPLNXNTHJVNLVG3DH4N3237LZ7QMOJLGKSHHS6F7JQKOWCW3QLCYXLDYBBZFJ5ZRGVCEXXZVZYLR2ER33X3JLNJNHYICODVMPQ5VREN5GDSLSDLENJ6PUMFXKHZ5EHGOIGIT4TEW6LOYW6XMYR452BRPZSKKXLM5ZT7KHAVL64LPKNK45R3HAPH6IXAAAP772O3VBWC

Public Key:  SLV_EPK_AEAUKAAAADM5IPIORWJ24OYHX4JSJ7R6BRMO25EHHGERKFJ33EBK4FWVU4HBY
Name:        example_service
Email:       service@example.com
Tags:        [example]
```
##### With search
```bash
$ slv env list --env-search alice
Public Key:      SLV_EPK_AEAUKAAAABUEMSPQ4BJIIWMSAKFUUXUV4THOP3ERH25CY4HR54W25HUJQR6XK
Name:            alice
Email:           alice@example.com
Tags:            [example_env]
Secret Binding:  SLV_ESB_AF4JYBGA2FXKUMAYADQHP6LPPMJM6JQDILREKS3KIHUQJWRZZKO5FQKCM4FHIRGR7DXPWHRQIACMHSNZVOOTJ7EDBGRAOODHEABHYIYY5O4Q23WD7TOZWSKIF66XVPZPLNXNTHJVNLVG3DH4N3237LZ7QMOJLGKSHHS6F7JQKOWCW3QLCYXLDYBBZFJ5ZRGVCEXXZVZYLR2ER33X3JLNJNHYICODVMPQ5VREN5GDSLSDLENJ6PUMFXKHZ5EHGOIGIT4TEW6LOYW6XMYR452BRPZSKKXLM5ZT7KHAVL64LPKNK45R3HAPH6IXAAAP772O3VBWC
```

---

## See Also

- [Create a New Environment](/docs/command-reference/environment/new) - Create a new environment
- [Show Environment](/docs/command-reference/environment/show) - View environment details
- [Add Environment](/docs/command-reference/environment/add) - Add an existing environment
- [Environment Component](/docs/components/environment) - Learn more about environments

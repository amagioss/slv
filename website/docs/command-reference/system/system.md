---
sidebar_position: 1
---

# Reset SLV 
Wipe clean of all profile configurations and delete environments in your local SLV installation. This does not touch remote profile configurations.

#### General Usage:
```bash
slv system reset [flags]
```

#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --yes | None | NA | NA | Confirm action without prompt |
| --help | None | NA | NA | Help text for `slv system reset` |

#### Usage:
```bash
slv system reset
```

#### Example:
```bash
$ slv system reset
You have a configured environment which you might have to consider backing up:
Public Key:      SLV_EPK_AEAUKAAAABUEMSPQ4BJIIWMSAKFUUXUV4THOP3ERH25CY4HR54W25HUJQR6XK
Name:            alice
Email:           alice@example.com
Tags:            [example_env]
Secret Binding:  SLV_ESB_AF4JYBGA2FXKUMAYADQHP6LPPMJM6JQDILREKS3KIHUQJWRZZKO5FQKCM4FHIRGR7DXPWHRQIACMHSNZVOOTJ7EDBGRAOODHEABHYIYY5O4Q23WD7TOZWSKIF66XVPZPLNXNTHJVNLVG3DH4N3237LZ7QMOJLGKSHHS6F7JQKOWCW3QLCYXLDYBBZFJ5ZRGVCEXXZVZYLR2ER33X3JLNJNHYICODVMPQ5VREN5GDSLSDLENJ6PUMFXKHZ5EHGOIGIT4TEW6LOYW6XMYR452BRPZSKKXLM5ZT7KHAVL64LPKNK45R3HAPH6IXAAAP772O3VBWC
------------------------------------------------------------
Env Definition:  SLV_EDS_AF4JYNGKIFF4GMAYQ7Y674R7A4HTH5MQSOUUJMWFNHLDJW2EUTCPQMQCJEOJWKYW6G5UWQMPB7H66G7W6KLFGNBUIHAD23AHMXGSBFPUIZFCSW5P23HG46F3LIXHO2ZHZW67O6574W6X5MWXVLJYXZXOTHV25YN3IWR722WT3XWA2GA6IMQQRBE4EKAUDFMQ6J757USXFDDUYV7RUPGK7DX5OSOPLZKME4YPJYLQQZ4MCHY3VCHHQZLQCRH7JWNG7KPOUAIC7D4Q2AAA77745LZ2FQ
Are you sure you wish to proceed? (yes/no): yes
System reset successful
```

---

## See Also

- [Interactive Mode (TUI)](/docs/interactive-mode) - Use the interactive TUI to manage SLV
- [Create a New Environment](/docs/command-reference/environment/new) - Create a new environment after reset
- [Quick Start Guide](/docs/quick-start) - Get started with SLV
- [Overview](/docs/overview) - Learn about SLV

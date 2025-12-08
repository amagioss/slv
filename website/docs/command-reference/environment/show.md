---
sidebar_position: 5
---

# Show Requested Environment
Used to print details about environments from the current context

#### General usage:
```bash
slv env show [flags]
slv env show [command]
```

#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --help | None | NA | NA| Help text for `slv env show` |

#### Commands:
- [`self`](#show-self-environment) - Shows the self environment.
- [`root`](#show-root-environment) - Shows the root environment in the active profile.
- [`k8s`](#show-k8s-environment) - Takes the current kubernetes context and shows the environment present in the context.
---
## Show Self Environment
Print details about the self environment. This is the environment that will be used by default while opening vaults and sharing vaults with others.
#### General Usage:
```bash
slv env show self [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --help | None | NA | NA| Help text for `slv env show self` |
#### Usage:
```bash
slv env show self
```
#### Example:
```bash
$ slv env show self
Public Key:      SLV_EPK_AEAUKAAAABUEMSPQ4BJIIWMSAKFUUXUV4THOP3ERH25CY4HR54W25HUJQR6XK
Name:            alice
Email:           alice@example.com
Tags:            [example_env]
Secret Binding:  SLV_ESB_AF4JYBGA2FXKUMAYADQHP6LPPMJM6JQDILREKS3KIHUQJWRZZKO5FQKCM4FHIRGR7DXPWHRQIACMHSNZVOOTJ7EDBGRAOODHEABHYIYY5O4Q23WD7TOZWSKIF66XVPZPLNXNTHJVNLVG3DH4N3237LZ7QMOJLGKSHHS6F7JQKOWCW3QLCYXLDYBBZFJ5ZRGVCEXXZVZYLR2ER33X3JLNJNHYICODVMPQ5VREN5GDSLSDLENJ6PUMFXKHZ5EHGOIGIT4TEW6LOYW6XMYR452BRPZSKKXLM5ZT7KHAVL64LPKNK45R3HAPH6IXAAAP772O3VBWC
------------------------------------------------------------
Env Definition:  SLV_EDS_AF4JYNGKIFF4GMAYQ7Y674R7A4HTH5MQSOUUJMWFNHLDJW2EUTCPQMQCJEOJWKYW6G5UWQMPB7H66G7W6KLFGNBUIHAD23AHMXGSBFPUIZFCSW5P23HG46F3LIXHO2ZHZW67O6574W6X5MWXVLJYXZXOTHV25YN3IWR722WT3XWA2GA6IMQQRBE4EKAUDFMQ6J757USXFDDUYV7RUPGK7DX5OSOPLZKME4YPJYLQQZ4MCHY3VCHHQZLQCRH7JWNG7KPOUAIC7D4Q2AAA77745LZ2FQ
```
---
## Show Root Environment
Print details about the root environment. This is applicable only when there is an active profile.
#### General Usage:
```bash
slv env show root [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --help | None | NA | NA| Help text for `slv env show root` |
#### Usage:
```bash
slv env show root
```
#### Example:
```bash
$ slv env show root
Public Key:  SLV_EPK_AEAUKAAAADM5IPIORWJ24OYHX4JSJ7R6BRMO25EHHGERKFJ33EBK4FWVU4HBY
Name:        example_service
Email:       service@example.com
Tags:        [example]
------------------------------------------------------------
Env Definition:  SLV_EDS_AF4JYRGM35FMGMA4YXYXOOOXYGFS2U6ISUEWMZUJWMS52H4HJCE7DRYIGS23IWRM4K5YWUGZ5X4RZPW7FA7V7GYUDVGRBKA6B22S4XJNWXODWKNFXTN24NXFU3YEPO7AYXETY4K33ENX7LP4WZMKZNBR67NLXNK5F22XE7D5HTWLSUJHGA6IMTAQUCXZBO4G5KA7UMKFAKB44IJVCCMJPW7ZOEK57447W3RW52XI4JQNRBPTADY7ZX6CBNBULMNHB6K5VN6UTYQYBH67AAAAB777TR3T5CA
```
---
## Show K8S Environment
Print details about the environment present in the cluster corresponding to the current kubernete context. The current context can be found by running the command `kubectl config get-contexts`.
#### General Usage:
```bash
slv env show k8s [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --help | None | NA | NA| Help text for `slv env show k8s` |
#### Usage:
```bash
slv env show k8s
```
#### Example:
```bash
$ slv env show k8s
Public Key: SLV_EPK_AEAUKAAAAD6XTJCYBCIHYKDPPHQN3YNDEVBOFCOIVDMGESLJFH65KG3VULVBK

K8s Cluster Info:
Name   : slv-test-cluster
Address: https://127.0.0.1:52358
User   : slv-test-cluster
```

---

## See Also

- [List Environments](/docs/command-reference/environment/list) - View all available environments
- [Create a New Environment](/docs/command-reference/environment/new) - Create a new environment
- [Environment Component](/docs/components/environment) - Learn more about environments
- [Kubernetes Integration](/docs/extensions/slv-in-kubernetes/operator) - Use SLV in Kubernetes

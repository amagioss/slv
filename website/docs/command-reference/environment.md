---
sidebar_position: 2
---

# Environment
Documentation for the `slv env` command.

Aliases: `env`, `envs`, `environment`, `environments`

Commands:
- [`new`](#create-a-new-environment)
- [`add`](#add-an-existing-environment)
- [`get`](#get-environments-in-active-profile)
- [`set-self`](#set-self-environment) 
- [`show`](#show-requested-environment)
- [`del`](#delete-an-environment)

---

## Create a new Environment
Used to create a new environment from scratch.
#### General Usage:
```bash
slv env new [command] [flags]
slv env new [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --quantum-safe | None | NA | NA| Use Quantum Resistant Cryptography (Kyber1024) |
| --help | None | NA | NA| Help text for `slv env new` |

#### Commands Available:
- [`self`](#create-a-new-self-environment)
- [`service`](#create-a-new-service-environment)

### Create a new self environment
Used to create a new self environment for end devices such as users' laptops.
#### General Usage:
```bash
slv env new self [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --email | String | True | None | Email Address for the environment being created |
| --name | String | True | None | Name of the environment to be created |
| --tags | String(s) | False | None | Tags to be set for the environment |
| --add | None | False | NA | Add the created environment to the active profile |
| --help | None | NA | NA| Help text for `slv env new self` |

#### Usage:
```bash
# For quantum safe environment
slv env new -q self --email <EMAIL_ADDRESS> --name <ENVIRONMENT_NAME> --add --tags <TAGS>

# For a regular environment
slv env new self --email <EMAIL_ADDRESS> --name <ENVIRONMENT_NAME> --add --tags <TAGS>
```

#### Example:
```bash
$ slv env new self --email alice@example.com --name alice --add --tags example_env
Enter a Password: 
Confirm Password: 
Public Key:      SLV_EPK_AEAUKAAAABUEMSPQ4BJIIWMSAKFUUXUV4THOP3ERH25CY4HR54W25HUJQR6XK
Name:            alice
Email:           alice@example.com
Tags:            [example_env]
Secret Binding:  SLV_ESB_AF4JYBGA2FXKUMAYADQHP6LPPMJM6JQDILREKS3KIHUQJWRZZKO5FQKCM4FHIRGR7DXPWHRQIACMHSNZVOOTJ7EDBGRAOODHEABHYIYY5O4Q23WD7TOZWSKIF66XVPZPLNXNTHJVNLVG3DH4N3237LZ7QMOJLGKSHHS6F7JQKOWCW3QLCYXLDYBBZFJ5ZRGVCEXXZVZYLR2ER33X3JLNJNHYICODVMPQ5VREN5GDSLSDLENJ6PUMFXKHZ5EHGOIGIT4TEW6LOYW6XMYR452BRPZSKKXLM5ZT7KHAVL64LPKNK45R3HAPH6IXAAAP772O3VBWC
------------------------------------------------------------
Env Definition:  SLV_EDS_AF4JYNGKIFF4GMAYQ7Y674R7A4HTH5MQSOUUJMWFNHLDJW2EUTCPQMQCJEOJWKYW6G5UWQMPB7H66G7W6KLFGNBUIHAD23AHMXGSBFPUIZFCSW5P23HG46F3LIXHO2ZHZW67O6574W6X5MWXVLJYXZXOTHV25YN3IWR722WT3XWA2GA6IMQQRBE4EKAUDFMQ6J757USXFDDUYV7RUPGK7DX5OSOPLZKME4YPJYLQQZ4MCHY3VCHHQZLQCRH7JWNG7KPOUAIC7D4Q2AAA77745LZ2FQ
Successfully registered as self environment
Please note down the "Secret Binding" somewhere safe so that you don't lose it.
It is required to access your registered environment.
```

### Create a new service environment
Used to create a new service environment for production kubernetes clusters, Github Actions, CI Pipelines.
#### General Usage:
```bash
slv env new service [command] 
slv env new service [flags]
```
#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --email | String | True | None | Email Address for the environment being created |
| --name | String | True | None | Name of the environment to be created |
| --tags | String(s) | False | None | Tags to be set for the environment |
| --add | None | False | NA | Add the created environment to the active profile |
| --help | None | NA | NA| Help text for `slv env new service` |

#### Commands available:
- [`aws`](#creating-aws-kms-based-service-environments)
- [`gcp`](#creating-gcp-kms-based-service-environments)

3 types Environemts can be created
- [Regular Service](#creating-regular-service-environments) - Uses a conventional secret key (not recommended)
- [AWS KMS](#creating-aws-kms-based-service-environments) - Uses AWS KMS for secret key
- [GCP KMS](#creating-gcp-kms-based-service-environments) - Uses GCP KMS for secret key

#### Creating regular service environments
##### Usage: 
```bash
slv env new service --email <EMAIL_ADDRESS> --name <ENVIRONMENT_NAME> --tags <TAGS>
```
##### Example:
```bash
$ slv env new service --email service@example.com --name example_service --tags example --add
Public Key:  SLV_EPK_AEAUKAAAADM5IPIORWJ24OYHX4JSJ7R6BRMO25EHHGERKFJ33EBK4FWVU4HBY
Name:        example_service
Email:       service@example.com
Tags:        [example]
------------------------------------------------------------
Env Definition:  SLV_EDS_AF4JYRGM35FMGMA4YXYXOOOXYGFS2U6ISUEWMZUJWMS52H4HJCE7DRYIGS23IWRM4K5YWUGZ5X4RZPW7FA7V7GYUDVGRBKA6B22S4XJNWXODWKNFXTN24NXFU3YEPO7AYXETY4K33ENX7LP4WZMKZNBR67NLXNK5F22XE7D5HTWLSUJHGA6IMTAQUCXZBO4G5KA7UMKFAKB44IJVCCMJPW7ZOEK57447W3RW52XI4JQNRBPTADY7ZX6CBNBULMNHB6K5VN6UTYQYBH67AAAAB777TR3T5CA

Secret Key:      SLV_ESK_AEAEKAAA44QDOTDOAV5B5TZCVWYYVPMM6RGG65RQHEZMJHL4ZMVQJ3TZTV7772ZHKXLXOV5L4G2Y6NKRY4AYMXUBBXQRFGP5BSGGJVAJFIHRROI
```
The secret key is to be stored safely. SLV does not store the secret key anywhere. Since this method deals with the secret key directly, it is **not recommended to create regular service environments**.

#### Creating GCP KMS based service environments
##### General Usage:
```bash
slv env new service gcp [flags]
```

##### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --resource-name | String | True | None | GCP KMS resource name |
| --rsa-pubkey | String | True | None | KMS public key [RSA 4096] as pem file (Recommended to perform offline access binding) |
| --email | String | True | None | Email Address for the environment being created |
| --name | String | True | None | Name of the environment to be created |
| --tags | String(s) | False | None | Tags to be set for the environment |
| --add | None | False | NA | Add the created environment to the active profile |
| --help | None | NA | NA| Help text for `slv env new service gcp` |

##### Usage:
```bash
slv env new service gcp --email <EMAIL_ADDRESS> --name <ENVIRONMENT_NAME> --tags <TAGS> --add --resource-name <GCP_RESOURCE_NAME> --rsa-pubkey <PATH_TO_PUBKEY_PEM>
```

#### Example:
```bash
TBA
```

#### Creating AWS KMS based service environments
##### General Usage:
```bash
slv env new service aws [flags]
```

##### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --arn | String | True | None | AWS KMS arn |
| --rsa-pubkey | String | True | None | KMS public key [RSA 4096] as pem file (Recommended to perform offline access binding) |
| --email | String | True | None | Email Address for the environment being created |
| --name | String | True | None | Name of the environment to be created |
| --tags | String(s) | False | None | Tags to be set for the environment |
| --add | None | False | NA | Add the created environment to the active profile |
| --help | None | NA | NA| Help text for `slv env new service aws` |

##### Usage:
```bash
slv env new service aws --email <EMAIL_ADDRESS> --name <ENVIRONMENT_NAME> --tags <TAGS> --add --arn <AWS_RESOURCE_ARN> --rsa-pubkey <PATH_TO_PUBKEY_PEM>
```

#### Example:
```bash
TBA
```
---

## Add an existing Environment

Used to add an environment that is created elsewhere to the existing machine. The Environment Definition String (`EDS`) can be used to do the same.

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

## Get Environments in Active Profile
Used to list all the profiles that are present in the active profile. The command can also be used to filter down results or search based on name, tag and email.

#### General usage:
```bash
slv env get [flags]
```

#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --env-search | String(s) | False | None | Search for environments based on `tag`/`email`/`name` |
| --help | None | NA | NA| Help text for `slv env get` |

#### Usage:
```bash
slv env get --env-search <SEARCH_STRING>
```

#### Example:
##### Without search
```bash
$ slv env get
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
$ slv env get --env-search alice
Public Key:      SLV_EPK_AEAUKAAAABUEMSPQ4BJIIWMSAKFUUXUV4THOP3ERH25CY4HR54W25HUJQR6XK
Name:            alice
Email:           alice@example.com
Tags:            [example_env]
Secret Binding:  SLV_ESB_AF4JYBGA2FXKUMAYADQHP6LPPMJM6JQDILREKS3KIHUQJWRZZKO5FQKCM4FHIRGR7DXPWHRQIACMHSNZVOOTJ7EDBGRAOODHEABHYIYY5O4Q23WD7TOZWSKIF66XVPZPLNXNTHJVNLVG3DH4N3237LZ7QMOJLGKSHHS6F7JQKOWCW3QLCYXLDYBBZFJ5ZRGVCEXXZVZYLR2ER33X3JLNJNHYICODVMPQ5VREN5GDSLSDLENJ6PUMFXKHZ5EHGOIGIT4TEW6LOYW6XMYR452BRPZSKKXLM5ZT7KHAVL64LPKNK45R3HAPH6IXAAAP772O3VBWC
```
---
## Set Self Environment
Self Environment is the environment that will be used while unlocking vaults or sharing them. It can be seen as the environment that is currently being used.

#### General usage:
```bash
slv env set-self [flags]
```

#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --env-def | String | True | None | The `EDS` for the environment to be set as self (Need not be a part of the profile)  |
| --help | None | NA | NA| Help text for `slv env set-self` |

#### Usage:
```bash
slv env set-self --env-def <ENVIRONMENT_DEFINITION_STRING>
```

#### Example:
```
$ slv env set-self --env-def SLV_EDS_AF4JYNGKIFF4GMAYQ7Y674R7A4HTH5MQSOUUJMWFNHLDJW2EUTCPQMQCJEOJWKYW6G5UWQMPB7H66G7W6KLFGNBUIHAD23AHMXGSBFPUIZFCSW5P23HG46F3LIXHO2ZHZW67O6574W6X5MWXVLJYXZXOTHV25YN3IWR722WT3XWA2GA6IMQQRBE4EKAUDFMQ6J757USXFDDUYV7RUPGK7DX5OSOPLZKME4YPJYLQQZ4MCHY3VCHHQZLQCRH7JWNG7KPOUAIC7D4Q2AAA77745LZ2FQ
Enter the secret binding: SLV_ESB_AF4JYBGA2FXKUMAYADQHP6LPPMJM6JQDILREKS3KIHUQJWRZZKO5FQKCM4FHIRGR7DXPWHRQIACMHSNZVOOTJ7EDBGRAOODHEABHYIYY5O4Q23WD7TOZWSKIF66XVPZPLNXNTHJVNLVG3DH4N3237LZ7QMOJLGKSHHS6F7JQKOWCW3QLCYXLDYBBZFJ5ZRGVCEXXZVZYLR2ER33X3JLNJNHYICODVMPQ5VREN5GDSLSDLENJ6PUMFXKHZ5EHGOIGIT4TEW6LOYW6XMYR452BRPZSKKXLM5ZT7KHAVL64LPKNK45R3HAPH6IXAAAP772O3VBWC
Public Key:      SLV_EPK_AEAUKAAAABUEMSPQ4BJIIWMSAKFUUXUV4THOP3ERH25CY4HR54W25HUJQR6XK
Name:            alice
Email:           alice@example.com
Tags:            [example_env]
Secret Binding:  SLV_ESB_AF4JYBGA2FXKUMAYADQHP6LPPMJM6JQDILREKS3KIHUQJWRZZKO5FQKCM4FHIRGR7DXPWHRQIACMHSNZVOOTJ7EDBGRAOODHEABHYIYY5O4Q23WD7TOZWSKIF66XVPZPLNXNTHJVNLVG3DH4N3237LZ7QMOJLGKSHHS6F7JQKOWCW3QLCYXLDYBBZFJ5ZRGVCEXXZVZYLR2ER33X3JLNJNHYICODVMPQ5VREN5GDSLSDLENJ6PUMFXKHZ5EHGOIGIT4TEW6LOYW6XMYR452BRPZSKKXLM5ZT7KHAVL64LPKNK45R3HAPH6IXAAAP772O3VBWC
------------------------------------------------------------
Env Definition:  SLV_EDS_AF4JYNGKIFF4GMAYQ7Y674R7A4HTH5MQSOUUJMWFNHLDJW2EUTCPQMQCJEOJWKYW6G5UWQMPB7H66G7W6KLFGNBUIHAD23AHMXGSBFPUIZFCSW5P23HG46F3LIXHO2ZHZW67O6574W6X5MWXVLJYXZXOTHV25YN3IWR722WT3XWA2GA6IMQQRBE4EKAUDFMQ6J757USXFDDUYV7RUPGK7DX5OSOPLZKME4YPJYLQQZ4MCHY3VCHHQZLQCRH7JWNG7KPOUAIC7D4Q2AAA77745LZ2FQ
Successfully registered self environment
```
---

## Show Requested Environment
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

### Show Self Environment
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

### Show Root Environment
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

### Show K8S Environment
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
## Delete an Environment

Delete one or more environments based on search parameters such as `name`, `email`, `tags`.

#### General Usage:
```bash
slv env del [flags]
```

#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --env-search | String(s) | False | None | Search for environments based on `tag`/`email`/`name` to delete |
| --help | None | NA | NA| Help text for `slv env del` |

#### Usage:
```bash
slv env del --env-search <SEARCH_STRING>
```
#### Example:
```bash
$ slv env del --env-search alice
Public Key:      SLV_EPK_AEAUKAAAABUEMSPQ4BJIIWMSAKFUUXUV4THOP3ERH25CY4HR54W25HUJQR6XK
Name:            alice
Email:           alice@example.com
Tags:            [example_env]
Secret Binding:  SLV_ESB_AF4JYBGA2FXKUMAYADQHP6LPPMJM6JQDILREKS3KIHUQJWRZZKO5FQKCM4FHIRGR7DXPWHRQIACMHSNZVOOTJ7EDBGRAOODHEABHYIYY5O4Q23WD7TOZWSKIF66XVPZPLNXNTHJVNLVG3DH4N3237LZ7QMOJLGKSHHS6F7JQKOWCW3QLCYXLDYBBZFJ5ZRGVCEXXZVZYLR2ER33X3JLNJNHYICODVMPQ5VREN5GDSLSDLENJ6PUMFXKHZ5EHGOIGIT4TEW6LOYW6XMYR452BRPZSKKXLM5ZT7KHAVL64LPKNK45R3HAPH6IXAAAP772O3VBWC

Are you sure you wish to delete the above environment(s) [yes/no]: yes
Environment alice deleted successfully
```






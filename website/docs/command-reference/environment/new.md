---
sidebar_position: 1
---
# Create a New Environment
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

## Create a new self environment
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
| --help | None | NA | NA| Help text for `slv env new self` |

#### Usage:
```bash
# For quantum safe environment
slv env new -q self --email <EMAIL_ADDRESS> --name <ENVIRONMENT_NAME> --tags <TAGS>

# For a regular environment
slv env new self --email <EMAIL_ADDRESS> --name <ENVIRONMENT_NAME> --tags <TAGS>
```

#### Example:
```bash
$ slv env new self --email alice@example.com --name alice --tags example_env
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
---
## Create a new service environment
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
| --help | None | NA | NA| Help text for `slv env new service` |

#### Commands available:
- [`direct`](#creating-regular-service-environments)
- [`aws`](#creating-aws-kms-based-service-environments)
- [`gcp`](#creating-gcp-kms-based-service-environments)
- [`azure`](#creating-azure-kms-based-service-environments)

4 types Environemts can be created
- [Regular Service](#creating-regular-service-environments) - Uses a conventional secret key (not recommended)
- [AWS KMS](#creating-aws-kms-based-service-environments) - Uses AWS KMS for secret key
- [GCP KMS](#creating-gcp-kms-based-service-environments) - Uses GCP KMS for secret key
- [Azure KMS](#creating-azure-kms-based-service-environments) - Uses GCP KMS for secret key


### Creating regular service environments
#### Usage: 
```bash
slv env new service direct --email <EMAIL_ADDRESS> --name <ENVIRONMENT_NAME> --tags <TAGS>
```
#### Example:
```bash
$ slv env new service direct --email service@example.com --name example_service --tags example 
Public Key:  SLV_EPK_AEAUKAAAADM5IPIORWJ24OYHX4JSJ7R6BRMO25EHHGERKFJ33EBK4FWVU4HBY
Name:        example_service
Email:       service@example.com
Tags:        [example]
------------------------------------------------------------
Env Definition:  SLV_EDS_AF4JYRGM35FMGMA4YXYXOOOXYGFS2U6ISUEWMZUJWMS52H4HJCE7DRYIGS23IWRM4K5YWUGZ5X4RZPW7FA7V7GYUDVGRBKA6B22S4XJNWXODWKNFXTN24NXFU3YEPO7AYXETY4K33ENX7LP4WZMKZNBR67NLXNK5F22XE7D5HTWLSUJHGA6IMTAQUCXZBO4G5KA7UMKFAKB44IJVCCMJPW7ZOEK57447W3RW52XI4JQNRBPTADY7ZX6CBNBULMNHB6K5VN6UTYQYBH67AAAAB777TR3T5CA

Secret Key:      SLV_ESK_AEAEKAAA44QDOTDOAV5B5TZCVWYYVPMM6RGG65RQHEZMJHL4ZMVQJ3TZTV7772ZHKXLXOV5L4G2Y6NKRY4AYMXUBBXQRFGP5BSGGJVAJFIHRROI
```
The secret key is to be stored safely. SLV does not store the secret key anywhere. Since this method deals with the secret key directly, it is **not recommended to create regular service environments**.

### Creating GCP KMS based service environments
#### General Usage:
```bash
slv env new service gcp [flags]
```

#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --resource-name | String | True | None | GCP KMS resource name |
| --rsa-pubkey | String | True | None | KMS public key [RSA 4096] as pem file (Recommended to perform offline access binding) |
| --email | String | True | None | Email Address for the environment being created |
| --name | String | True | None | Name of the environment to be created |
| --tags | String(s) | False | None | Tags to be set for the environment |
| --help | None | NA | NA| Help text for `slv env new service gcp` |

##### Usage:
```bash
slv env new service gcp --email <EMAIL_ADDRESS> --name <ENVIRONMENT_NAME> --tags <TAGS>  --resource-name <GCP_RESOURCE_NAME> --rsa-pubkey <PATH_TO_PUBKEY_PEM>
```

#### Example:
```bash
TBA
```

---

### Creating AWS KMS based service environments
#### General Usage:
```bash
slv env new service aws [flags]
```

#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --arn | String | True | None | AWS KMS arn |
| --rsa-pubkey | String | True | None | KMS public key [RSA 4096] as pem file (Recommended to perform offline access binding) |
| --email | String | True | None | Email Address for the environment being created |
| --name | String | True | None | Name of the environment to be created |
| --tags | String(s) | False | None | Tags to be set for the environment |
| --help | None | NA | NA| Help text for `slv env new service aws` |

#### Usage:
```bash
slv env new service aws --email <EMAIL_ADDRESS> --name <ENVIRONMENT_NAME> --tags <TAGS>  --arn <AWS_RESOURCE_ARN> --rsa-pubkey <PATH_TO_PUBKEY_PEM>
```

#### Example:
```bash
TBA
```

---

### Creating Azure KMS based service environments
#### General Usage:
```bash
slv env new service azure [flags]
```

#### Flags:
| Flag | Arguments | Required | Default | Description |
| -- | -- | -- | -- | -- |
| --vault-url | String | True | None | Azure key vault URL |
| --key-name | String | True | None | Name of the key in Azure Key Vault to use |
| --key-version | String | False | None |  Version of the key in Azure Key Vault to use (optional, latest version will be used if not specified) |
| --rsa-pubkey | String | True | None | KMS public key [RSA 4096] as pem file (Recommended to perform offline access binding) |
| --email | String | True | None | Email Address for the environment being created |
| --name | String | True | None | Name of the environment to be created |
| --tags | String(s) | False | None | Tags to be set for the environment |
| --help | None | NA | NA| Help text for `slv env new service azure` |

#### Usage:
```bash
slv env new service azure --email <EMAIL_ADDRESS> --name <ENVIRONMENT_NAME> --tags <TAGS>  --vault-url <AZURE_VAULT_URL> --rsa-pubkey <PATH_TO_PUBKEY_PEM> --key-name <KEY_NAME>
```

#### Example:
```bash
TBA
```

---

## See Also

- [List Environments](/docs/command-reference/environment/list) - View all available environments
- [Show Environment](/docs/command-reference/environment/show) - View environment details
- [Add Environment](/docs/command-reference/environment/add) - Add an existing environment
- [Environment Component](/docs/components/environment) - Learn more about environments
- [Quick Start Guide](/docs/quick-start) - Get started with SLV

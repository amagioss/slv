# SLV - Secure Local Vault
Securely store, share, and access secrets alongside the codebase.

SLV is designed based on the following **key principles**
 - Anyone can add or update secrets, however will not be able to read them unless they have access to the vault
 - An environment should have a single identity that will give access to all necessary secrets from any vault shared with it

## Installation
Download the latest SLV binary from the [releases](https://github.com/amagioss/slv/releases/latest) page and add it to your path.

### Homebrew
SLV can be installed with brew using the following command on macOS
```zsh
brew install amagioss/slv/slv
```

### Install Script

#### Install Latest Version
**With Shell (MacOs/Linux):**
```sh
curl -fsSL https://slv.sh/scripts/install.sh | sh
```
**With PowerShell (Windows):**
```powershell
irm https://slv.sh/scripts/install.ps1 | iex
```

#### Install Specific Version
**With Shell (MacOs/Linux):**
```sh
curl -fsSL https://slv.sh/scripts/install.sh | sh -s v0.1.7
```
**With PowerShell (Windows):**
```powershell
$v="0.1.7"; irm https://slv.sh/scripts/install.ps1 | iex
```

### Docker
You can also run SLV without installing using Docker:
```zsh
docker run -it --rm -v $PWD:/workspace ghcr.io/amagioss/slv:latest version
```

## Usage

### Basic CLI Commands

#### Create a new profile
```sh
$ slv profile new -n my_org

Created profile: my_org
```

#### Create a new environment
```sh
$ slv env new service -n alice -e alice@example.com --add

Public Key:       SLV_EPK_AEAUKAAAAD6XTJCYBCIHYKDPPHQN3YNDEVBOFCOIVDMGESLJFH65KG3VULVBK
Name:             alice
Email:            alice@example.com
Tags:             []
------------------------------------------------------------
Env Definition:  SLV_EDS_AF4JYNGKJ5FYMMA4YDY7P4R3JOLYPHWDJZWW57U35FBB26MSWV7MQYC3UIUUT5G6IOROHF7P44N5J7XGTWKXQAUBV3LJGUDSUKBA5ESSJL473NNP2KI2KZJRJKXFJ4OS3TDIMC6N3IWG2S6NT5Z5DVKVK3OB6ZL62NB23GMEAQNBGEAIDDXSYQQCEIMOP773BG7UYWB4H3MI64F5PD2OO4XJBXL6HT7XM3PIBRG57MCDVNBLPYZBPX25TSAQB7H4AYAAB777D2YDPOA

Secret Key:	 SLV_ESK_AEAEKAAATI5CXB7QMFSUGY4RUT6UTUSK7SGMIECTJKRTQBFY6BN5ZV5M5XGF6DWLV2RVCJJSMXH43DJ6A5TK7Y6L6PYEMCDGQRBX46GUQPUIYUQ
```

#### Create a vault
- To create a vault and share it with the environment `alice`, use the following command
```sh
$ slv vault new -v test.slv.yaml -s alice

Created vault: test.slv.yaml
```
- To create a K8s compatible vault, use the following command
```sh
$ slv vault new -v test.slv.yaml -s alice --k8s production

Created vault: test.slv.yaml
```

#### Add secrets to the vault
```sh
$ slv vault put -v test.slv.yaml -n db_password -s "super_secret_pwd"

Added secret: db_password to vault: test.slv.yaml
```

#### Get secrets from the vault
Set the environment variable `SLV_ENV_SECRET_KEY` to the secret key generated in the previous step
```sh
$ export SLV_ENV_SECRET_KEY=SLV_ESK_AEAEKAAATI5CXB7QMFSUGY4RUT6UTUSK7SGMIECTJKRTQBFY6BN5ZV5M5XGF6DWLV2RVCJJSMXH43DJ6A5TK7Y6L6PYEMCDGQRBX46GUQPUIYUQ
$ slv vault get -v test.slv.yaml -n db_password

super_secret_pwd
```

#### Share the vault with other environments
Ensure that the current environment has access to the vault in order to share it with other environments
```sh
$ slv vault share -v test.slv.yaml -s bob

Shared vault: test.slv.yaml
```
Once shared, the other environments can access the vault using their respective secret keys

### Using SLV as a library
SLV can also be used as a library in your Go projects. The following is an example of how to use SLV to read secrets from a vault. Ensure that the vault is shared with the environment and the environment is configured by setting the environment variable `SLV_ENV_SECRET_KEY` or `SLV_ENV_SECRET_BINDING` before executing it.
```go
package main

import "slv.sh/slv"

func main() {
	viMap, err := slv.GetAllVaultItems("demo.slv.yaml")
	if err != nil {
		panic(err)
	}
	for k, v := range viMap {
		print(k, " : ", string(v.Value()), "\n")
	}
}
```

## Integrations
Some of the integrations that SLV currently supports are:
- [Kubernetes](/docs/KUBERNETES.md)
- [GitHub Actions](https://github.com/amagioss/slv-action)
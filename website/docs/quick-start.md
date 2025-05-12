---
sidebar_position: 3
---
# Quick Start

## Try SLV yourself

### Create a new profile
A profile is a collection of environments or identities. For example, you can create a profile for your organisation.
```bash
slv profile new -n my_org
```
#### Output
```
Created profile:  my_org
```

### Create a new environment
An environment here refers to an identity. The identity can be a person or a service. For example, let's create an identity for the service alice within the organisation.
```bash
slv env new service -n alice -e alice@example.com --add
```
#### Output
```
Public Key:       SLV_EPK_AEAUKAAAAD6XTJCYBCIHYKDPPHQN3YNDEVBOFCOIVDMGESLJFH65KG3VULVBK
Name:             alice
Email:            alice@example.com
Tags:             []
------------------------------------------------------------
Env Definition:  SLV_EDS_AF4JYNGKJ5FYMMA4YDY7P4R3JOLYPHWDJZWW57U35FBB26MSWV7MQYC3UIUUT5G6IOROHF7P44N5J7XGTWKXQAUBV3LJGUDSUKBA5ESSJL473NNP2KI2KZJRJKXFJ4OS3TDIMC6N3IWG2S6NT5Z5DVKVK3OB6ZL62NB23GMEAQNBGEAIDDXSYQQCEIMOP773BG7UYWB4H3MI64F5PD2OO4XJBXL6HT7XM3PIBRG57MCDVNBLPYZBPX25TSAQB7H4AYAAB777D2YDPOA

Secret Key:	 SLV_ESK_AEAEKAAATI5CXB7QMFSUGY4RUT6UTUSK7SGMIECTJKRTQBFY6BN5ZV5M5XGF6DWLV2RVCJJSMXH43DJ6A5TK7Y6L6PYEMCDGQRBX46GUQPUIYUQ
```
Have the secret key stored somewhere safe, as you may have to use it later.

### Create a new vault
A vault is a collection of secrets (as the name suggests). You can share vaults with other environments within the profile. Let's create a secret and share it with alice.
```bash
slv vault new -v test.slv.yaml -s alice
```
#### Output
```
Created vault: test.slv.yaml
```
The vault has been created and also shared with alice.

### Add a secret to the vault
Now that the vault has been created and shared with alice, let's try and add a secret to the vault.
```bash
slv vault put -v test.slv.yaml -n db_password -s "super_secret_pwd"
```
#### Output
```
Added secret: db_password to vault: test.slv.yaml
```

### Retrieve a secret from the vault
Now that we have a vault with a secret, let's try and extract the secret.
Before doing this, you will have to set the Environment variable `SLV_ENV_SECRET_KEY` to the Secret key generated while creating the environment.
```bash
export SLV_ENV_SECRET_KEY=SLV_ESK_AEAEKAAATI5CXB7QMFSUGY4RUT6UTUSK7SGMIECTJKRTQBFY6BN5ZV5M5XGF6DWLV2RVCJJSMXH43DJ6A5TK7Y6L6PYEMCDGQRBX46GUQPUIYUQ
slv vault get -v test.slv.yaml -n db_password
```
#### Output
```
super_secret_pwd
```

### Share a vault with another environment
If you already have a vault with you and would like to share it with someone else, you can do so by running
```bash
slv vault share -v test.slv.yaml -s bob
```
#### Output
```
Shared vault: test.slv.yaml
```
You can share a vault with another environment **only if the current environment has access to it**.



---
sidebar_position: 3
---
# Quick Start

### Create a new environment
An [environment](/docs/components/environment) represents an accessing entity, which can be either a `user` or a `service`. For now, we'll create an environment associated with the current user and register it to the development machine that is in use.

> You’ll be prompted to create and confirm a password. This password will be required to access the environment later, so be sure to remember it.
```bash
slv env new self --name kuwan --email kuwan@example.com --tags cats --tags pets
```
#### Output
```
Public Key:      SLV_EPK_AEAUKAAAACBGNPQ5ORZWPBM6AG3C7IYC57NMQIXIXCTN2GWE422PBAQMWQ5AE
Name:            kuwan
Email:           kuwan@example.com
Tags:            [cats pets]
Secret Binding:  SLV_ESB_AF4JYBGA3VXIEMAUADQHOONXLOROHL2K4LCVCHE3QGHATLV2HOQTRHLFVUW24DOGO733XAMCDDKN4GBLGWDWPUAQ37AJQM6EQBXI5WH6IS6N44BLP3U4TNBJSY2D2XHUSPL3C3USE7RYZETLLN4HVXJFUSHZEIHPUXE5KDX5YDQH4JMSM3MKUZTRTCDNKFLN2253MHCRXGVN4SYSJ23H6R735JILAD2WRQWETZZ6C5V7TQXFVG7JZM7FKKJSQYWT43Z36CPMILK2JWEE2W4PXTSYTL4F4MNY3VNXKFC7NYUDS5IGM4Z3RXY7AEAAB777TCWUJCY
------------------------------------------------------------
Env Definition:  SLV_EDS_AF4JYNGKIFF4GMAYQ7Y674R7A4H5KOXIZG3SLFCSDMJVP3KUMTCPQMUCJUWEXKYO6G5UWDZ3HY6L6X7I4VW7JLXFCMFGY3Y765JLO64S6TIBEEKVMWW3JSPP52PQOXLW25KF6VU342U4UN5KGPG25WKVXXFOUQK6MWMS5SLUQPEUSQSA3HACR4FRPTNQQAIZVQP467ODH43EYI27XDH3BPXY2WP2MVJPRGHRB2HNEGQXRANTOOBMBRDTYKV4BFW5SHT5FR3XD4HSRAF774AAAAH776UJAOLP
Successfully registered as self environment
Please note down the "Secret Binding" somewhere safe so that you don't lose it.
It is required to access your registered environment.
```

### Create a new vault
A [vault](/docs/components/vault) stores secrets in encrypted format in SLV. Let's quickly create one that's accessible to the environment we just set up.
```bash
slv vault new --vault test.slv.yaml --env-self
```
#### Output
```
Created vault: test.slv.yaml
```
The vault has been created and can be accessed by the environment we just created.

### Add a secret to the vault
Now that the vault has been created, let's try to add a secret to the vault.
```bash
slv vault put --vault test.slv.yaml --name my_cat --value "is_pawsome"
```
#### Output
```
Successfully added/updated secret my_cat into the vault test.slv.yaml
```

### Retrieve the secret from the vault
Now that we have a secret in the vault, let's try to read it.
> You will be prompted to enter the password you set when creating the environment, and asked if you’d like to save it to your keychain for future use.
```bash
slv vault get --vault test.slv.yaml --name my_cat
```
#### Output
```
is_pawsome
```

> Alternatively, you can inject the secret into your environment variables. This approach is particularly useful when you want to access secrets in your code locally during testing.
```bash
slv vault shell --vault test.slv.yaml
```
#### Output
```
Initialized SLV session with secrets loaded as environment variables from the vault test.slv.yaml.
```
> Once set, the secrets will be available as environment variables and accessible to your application at runtime.
```bash
echo $my_cat
```
#### Output
```
is_pawsome
```

---

## Next Steps

Now that you've completed the quick start, explore these resources:

- **[Interactive Mode (TUI)](/docs/interactive-mode)** - Try the interactive TUI for managing vaults
- **[Vault Component](/docs/components/vault)** - Learn more about how vaults work
- **[Profile Component](/docs/components/profile)** - Understand profiles for team collaboration
- **[Environment Component](/docs/components/environment)** - Learn about environments
- **[Command Reference](/docs/command-reference/vault/get)** - Explore all available commands
- **[Kubernetes Integration](/docs/extensions/slv-in-kubernetes/operator)** - Use SLV in Kubernetes
# SLV - Secure Local Vault
Secure Local Vault - SLV (a.k.a Secrets Launch Vehicle 🔐🚀) is a tool to manage secrets locally in a secure manner. It is designed to be used by developers to manage secrets along with their code so that they can be shared with other developers and services in a secure manner.

SLV is designed based on the following **key principles**
 - Anyone can add or update secrets, however will not be allowed to read them unless they have access to the vault
 - An environment should have a single key that will give access to all necessary secrets irrespective of the number of vaults shared with it

 ## How to install
 SLV can be installed using brew using the following command
```zsh
brew install shibme/tap/slv
```
Alternatively, you can download the SLV binary from the [releases](https://github.com/shibme/slv-beta/releases/latest) page and add it to your path.
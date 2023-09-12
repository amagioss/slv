# SLV - Secure Local Vault
Secure Local Vault - SLV (a.k.a Secrets Launch Vehicle) is a tool to manage secrets locally in a secure manner. It is designed to be used by developers to manage secrets along with their code so that they can be shared with other developers and services in a secure manner.

SLV is designed based on the following **key principles**
 - Anyone can add or update secrets, however may not be allowed to view them
 - The receiving environment of the secret should have a single key that gives access to all necessary secrets irrespective of the number of vaults shared with it
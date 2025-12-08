---
sidebar_position: 3
---

# Vault

## What is a vault?

A **vault** in SLV is a container for storing a secret. To create a vault, see [Create a New Vault](/docs/command-reference/vault/new).

As the name suggests, a vault represents a secure space where sensitive data — such as tokens, credentials, or API keys — can be stored and shared with specific environments. It holds the secret value, as well as a list of environments it is shared with.

Each vault can be shared with multiple environments. These environments are identified by their public keys, which are used to determine who can access the secret.

From a user’s perspective, a vault is:
- A **named container** for storing multiple secrets.
- A **sharing mechanism** for distributing secrets to specific environments.

---

## Write Without Access

One of the defining features of SLV is that **you do not need access to a vault to put a secret into it**.
This means:
- Anyone can write to a vault, even if they cannot read from it.
- Access to read secrets is strictly limited to environments the vault is shared with.

This makes SLV particularly powerful in collaborative workflows.

---

## Related Topics

- [List Vaults](/docs/command-reference/vault/list) - Discover and list all vaults in your project
- [Create a New Vault](/docs/command-reference/vault/new) - Learn how to create vaults
- [Put a Secret](/docs/command-reference/vault/put) - Add secrets to your vault
- [Get a Secret](/docs/command-reference/vault/get) - Retrieve secrets from your vault
- [Quick Start Guide](/docs/quick-start) - Get started with SLV

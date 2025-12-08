---
sidebar_position: 1
---

# Overview

# What is SLV?

**SLV (Secure Local Vault)** is a free and open-source software designed for securely storing, sharing, and accessing secrets directly within your codebase. It ensures secrets remain encrypted at rest and in transit, allowing access only within explicitly authorized environments, thus simplifying secure collaboration and enhancing overall security.

---

## What problem does SLV solve?

Secrets such as API keys and tokens are often scattered across Git repositories, CI/CD systems, and cloud configurations, increasing the risk of leaks. SLV addresses this challenge by providing a simple, secure, and decentralized solution for managing secrets without relying on a centralized vault. This approach enables seamless cross-collaboration across hybrid cloud environments and even supports air-gapped systems.

---

## Core Principles of SLV

- **Secrets can be added or updated by anyone**, but only authorized environments are permitted to decrypt them.
- **Each environment is treated as a distinct entity**, with the ability to access secrets distributed across multiple vaults.

This approach ensures that SLV remains secure, auditable, and easy to use.

---

## Next Steps

- **[Quick Start Guide](/docs/quick-start)** - Get up and running with SLV in minutes
- **[Installation](/docs/installation)** - Install SLV on your system
- **[Components](/docs/components/vault)** - Learn about vaults, profiles, and environments
- **[Command Reference](/docs/command-reference/vault/get)** - Explore all available commands
- **[Contributing](/docs/contributing)** - Contribute to SLV development

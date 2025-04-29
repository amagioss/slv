---
sidebar_position: 1
---

# Overview

# What is SLV?

**SLV (Secure Local Vault)** is an open-source tool by Amagi OSS that enables secure storage, sharing, and access to secrets directly alongside your codebase. It keeps secrets encrypted and only accessible to authorized environments.

---

## What Problem Does SLV Solve?

Secrets like API keys, credentials, and tokens are often spread across Git repos, CI/CD systems, and cloud configs — increasing the risk of leaks and complexity. SLV solves this by offering a simple, secure, and decentralized way to manage secrets without relying on a central vault.

---

## Core Idea

- **Anyone can add or update secrets**, but only authorized environments can decrypt them.
- **Each environment has a unique identity**, which controls access to secrets stored within a shared encrypted vault.

This makes SLV secure, auditable, and friendly to GitOps workflows.

---

## Get Started

Jump right into using SLV with setup instructions and usage examples - [Quick Start](/docs/quick-start)

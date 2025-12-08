---
sidebar_position: 2
---
# Profile

## What is an SLV Profile?

An **SLV profile** is a foundational concept that brings together a group of [environments](/docs/components/environment) under a single logical unit. To create a profile, see [New Profile](/docs/command-reference/profile/new). It defines **who** can access secrets by associating multiple environments that are intended to work together — such as teams, deployment stages, or organizational boundaries.

Profiles are designed to support scalable and secure secret sharing across individual developers, automation workflows, and infrastructure environments — while keeping access isolated and auditable.

---

## Purpose

An SLV profile primarily serves to **group environments** that belong to the same logical context. These may represent team members, CI pipelines, production systems, or any combination of entities that should have coordinated access to a set of secrets.

By grouping related environments under a single profile, SLV simplifies secret sharing, enforces access boundaries, and enables recovery through shared access mechanisms like the root environment.

---

## Environments

### Environment Grouping

A profile can include any number of environments. Each environment is uniquely identified by its own public-private key pair and can access secrets based on how those secrets are shared within the profile.

This design enables multiple environments to operate independently while securely accessing the same vaults, as long as they are authorized participants in the profile.

### Root Environment

A profile can optionally define a **root environment**. This environment is automatically granted access to all vaults created under the profile by default.

The root environment serves as a secure fallback, ensuring that secrets are not permanently lost if an individual environment's secret key or binding is lost. It is especially useful for recovery scenarios and administrative oversight.

---

## Use Cases

- A **developer team** working on a project can be grouped into a single profile, with each developer having their own environment key.
- A **CI/CD pipeline** and associated preview environments can be added to a profile for shared access to build-time secrets.
- A **production infrastructure** profile may include both automation environments and secure operator access, with a root environment acting as a recovery path.

---

## Related Topics

- [Create a New Profile](/docs/command-reference/profile/new) - Learn how to create profiles
- [List Profiles](/docs/command-reference/profile/list) - View all available profiles
- [Set Active Profile](/docs/command-reference/profile/set-active) - Switch between profiles
- [Environment Component](/docs/components/environment) - Learn about environments
- [Git Profile Integration](/docs/integration-guide/git-profile) - Use profiles with Git repositories

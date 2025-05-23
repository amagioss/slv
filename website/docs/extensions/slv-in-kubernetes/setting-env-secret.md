---
sidebar_position: 2
---

# Setup K8s Environment

## Overview

Accessing a vault requires the secret key of the environment with which the vault is shared. When using SLV locally, you can set the secret key via environment variables, store it in the system keyring, or manage it through secret bindings integrated with passwords or a KMS.

However, when using SLV inside Kubernetes, the approach to setting the environment secret key is different. This document explains how to securely configure it.

---

## Setting the Secret Key 

You can configure the secret key directly in Kubernetes by creating a **Kubernetes Secret** in the release namespace. Then, specify the name of that secret using the Helm chart value `k8sSecret`.

- **Default Behavior:** If `k8sSecret` is not provided, SLV will look for a secret named `slv`.
- **Key Requirement:** The Kubernetes Secret must contain a key named `SecretKey`. SLV will **only** look for this specific key.

> ⚠️ **Warning:**
> Directly handling secret keys is **not recommended** because it involves managing raw cryptographic material, increasing the risk of exposure. When using cloud based deployments, using **secret bindings** is strongly recommended for enhanced security.

---

## Setting the Secret Binding

A safer approach involves using a **Secret Binding**, which encapsulates key material securely.

### 1. Through Helm

You can provide the secret binding string directly via the Helm chart value `secretBinding`:

```bash
helm upgrade --install slv slv/slv-operator --set secretBinding=<your-secret-binding>
```

### 2. Through a Kubernetes Secret

Alternatively, you can store the secret binding inside a Kubernetes Secret, similar to the secret key method:

- Create a Kubernetes Secret containing a key named `SecretBinding`.
- Set the secret's name using the Helm chart value `k8sSecret`.

Example:

```bash
kubectl create secret generic slv -n slv --from-literal=SecretBinding=<your-secret-binding>
```

```bash
helm upgrade --install slv slv/slv-operator --set k8sSecret=slv
```

---

## Recommended Practices

- **Always prefer using Secret Bindings** instead of directly handling Secret Keys.
- **Integrate Secret Bindings with a KMS (Key Management System)** to secure the binding and unbinding process.

> You can learn more about integrating KMS into service environments [here](#).



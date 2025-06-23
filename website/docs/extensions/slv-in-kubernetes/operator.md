---
sidebar_position: 3
---
# Operator

## Overview

The **SLV Operator** manages `SLV` custom resources within a Kubernetes cluster, ensuring that secure vaults are automatically reconciled into Kubernetes Secrets for consumption by applications.

---

## Behavior

The operator continuously watches the `SLV` custom resources and performs automatic reconciliation:

- **Creation:** When a new `SLV` object is created, the operator generates a corresponding Kubernetes Secret in the same namespace.
- **Update:** When an `SLV` object is updated with new entries, the corresponding Secret is updated accordingly.
- **Deletion:** When an `SLV` object is deleted, the associated Kubernetes Secret is also removed.

This ensures a seamless, real-time synchronization between vault data and Kubernetes Secrets.

---

## Quick Install (For Testing)

We recommend using the official [Helm chart method](#installation) for deploying the SLV Operator in Production. The following method is for quick testing purposes only.

Create a namespace for the SLV Operator
```
kubectl create namespace slv
```

Preload the environment secret with the `SecretKey` or `SecretBinding` as shown below:
```bash
kubectl create secret generic slv -n slv --from-literal=SecretKey=<your_slv_env_secret_key>
```

Quickly install the SLV Operator using the provided YAML:
```bash
kubectl apply -f https://slv.sh/k8s/samples/deploy/operator.yaml
```

## Installation 

The SLV Operator can be easily deployed using the official Helm chart:

```bash
helm repo add slv https://slv.sh/charts
helm repo update
helm upgrade --install slv slv/slv-operator --namespace slv --create-namespace
```

> **Note:** By default, the operator expects a secret named `slv` in the release namespace containing either the `SecretKey` or `SecretBinding`.

---

## Helm Chart Values

### SLV Helm Chart Configuration

| Parameter | Description | Default |
|---|---|---|
| `secretBinding` | Secret binding string used for the environment. Either this or `k8sSecret` must be specified. | `""` |
| `k8sSecret` | Name of the Kubernetes Secret that contains either the `SecretKey` or `SecretBinding` under keys `SecretKey` or `SecretBinding` respectively. Must be in the release namespace. | `""` |
| `image` | Full image URL including tag. Must match the Helm chart version. | `"ghcr.io/amagioss/slv:<CHART_VERSION>"` |
| `resource` | CPU and memory resource limits/requests for the operator pods. Use standard `limits` and `requests` structure. | Refer Helm |
| `labels` | Additional labels to add to the Deployment | `{}` |
| `podLabels` | Additional labels to add to individual SLV pods. | `{}` |
| `podAnnotations` | Additional annotations to add to individual SLV pods. | `{}` |
| `nodeSelector` | Node selector for scheduling SLV pods on specific nodes. | `{}` |
| `affinity` | Pod affinity/anti-affinity rules for scheduling SLV pods. | `{}` |
| `tolerations` | Tolerations for scheduling SLV pods on tainted nodes. | `[]` |
| `env` | Environment variables to be set for the SLV operator container. | `{}` |
| `serviceAccount.labels` | Labels to be added to the ServiceAccount. | `{}` |
| `serviceAccount.annotations` | Annotations to be added to the ServiceAccount. | `{}` |
| `volumes` | Additional volumes to mount in the SLV pods. | `[]` |
| `volumeMounts` | Volume mounts for the volumes specified above. | `[]` |
| `replicas` | Number of SLV operator replicas. | `1` |
| `webhook.disableAutomaticCertManagement` | If `false`, SLV manages TLS certs for the webhook using the built-in mechanism. If `true`, you must manually manage TLS and caBundle injection. | `false` |
| `webhook.serviceName` | Name of the Kubernetes service pointing to the webhook server. | `"slv-webhook-service"` |
| `webhook.validatingWebhookConfigName` | Name of the `ValidatingWebhookConfiguration` resource. | `"slv-operator-validating-webhook"` |
| `webhook.certSecretName` | Name of the Kubernetes Secret where SLV stores TLS certs for the webhook. | `"slv-webhook-server-cert"` |
| `webhook.vwhAnnotations` | Additional annotations to add to the `ValidatingWebhookConfiguration` (e.g., for CA injection). | `{}` |

---

## Minimum Permissions for the Operator

If you prefer to create your own roles instead of using the default ClusterRole provided by the Helm chart, here are the minimum required permissions:

```yaml
rules:
  - apiGroups: ["slv.sh"]
    resources: ["slvs"]
    verbs: ["get", "list", "watch", "update"]

  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["create", "get", "list", "update", "delete", "watch"]

  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get", "create", "update"]
```

> **Tip:** By default, permissions are granted cluster-wide. You can scope them down to a namespace if needed.
---

## Example

### Preload the Environment Secret

```bash
kubectl create secret generic slv -n slv --from-literal=SecretKey=SLV_ESK_AEAEKAAATI5CXB7QMFSUGY4RUT6UTUSK7SGMIECTJKRTQBFY6BN5ZV5M5XGF6DWLV2RVCJJSMXH43DJ6A5TK7Y6L6PYEMCDGQRBX46GUQPUIYUQ
```
#### Output:
```bash
secret/slv created
```

### Install the SLV Operator
```bash
helm upgrade --install slv slv/slv-operator --namespace slv --create-namespace
```
#### Output:
```bash
Release "slv" has been upgraded. Happy Helming!
NAME: slv
LAST DEPLOYED: Mon Apr 28 14:45:45 2025
NAMESPACE: slv
STATUS: deployed
REVISION: 2
TEST SUITE: None
NOTES:
SLV Install Successful.
WARNING: You have not set the value for ".Values.secretBinding" or "Values.slvEnvironment.k8sSecret".
SLV will now look for a secret named "slv" in the "slv" namespace.
If a secret is not found, SLV will not run as expected and return an error.

Ensure that you have set atleast one of the following
- secret key for the environment (under key "SecretKey") 
- secret binding for the environment (under key "SecretBinding"),
under the secret name "slv" 
in namespace "slv"
```

### Apply an SLV object
```bash
kubectl apply -f https://slv.sh/k8s/samples/pets.slv.yaml
```
#### Output:
```bash
slv.slv.sh/pets created
```

### Retrieve the Corresponding Secret
```bash
kubectl get secret pets -o jsonpath='{.data.mycat}' | base64 --decode
```
#### Output:
```bash
Kuwan
```









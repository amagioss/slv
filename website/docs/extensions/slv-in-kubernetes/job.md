---
sidebar_position: 4
---

# Job / CronJob

## Overview

The **SLV Job** is a one-time runner that reconciles all existing SLV secrets. It is particularly useful when you prefer more control over the reconciliation process or when running a continuous operator is not feasible.

When a **schedule** is provided, the same reconciliation process is executed periodically as a **CronJob**.

---

## Behavior

- **Creation:** If an `SLV` object exists but its corresponding secret does not, a new secret is created.
- **Update:** If both the `SLV` object and the secret exist but differ, the secret is updated.
- **Deletion:** If a secret exists but its corresponding `SLV` object does not, the secret is deleted.

When run as a CronJob, the same behavior is applied during each scheduled execution.

---

## Installation

Deploy the SLV Job easily using the official Helm chart:

```bash
helm repo add slv https://slv.sh/charts
helm repo update
helm upgrade --install slv slv/slv-job --set jobName=my-job-$(date +%s)
```

> **Note:** The `jobName` is overridden to ensure uniqueness for each Helm upgrade.
> 
> By default, SLV expects a secret named `slv` in the release namespace, containing either a `SecretKey` or a `SecretBinding`.

---

## Helm Chart Values

| Parameter | Description | Default |
| --- | --- | --- |
| `secretBinding` | Secret binding string for the environment. | None |
| `k8sSecret` | Name of the Kubernetes Secret containing the `SecretKey` or `SecretBinding`. | `slv` |
| `image` | Full image URL including tag. Tag must match the chart version. | `ghcr.io/amagioss/slv:<CHART_VERSION>` |
| `resource` | CPU and memory resource limits. | `250m` CPU, `250Mi` Memory |
| `labels` | Additional labels for the deployment. | None |
| `podLabels` | Additional labels for the pods. | None |
| `podAnnotations` | Additional annotations for the pods. | `{}` |
| `nodeSelector` | Node selector for scheduling SLV pods on specific nodes. | `{}` |
| `affinity` | Pod affinity/anti-affinity rules for scheduling SLV pods. | `{}` |
| `tolerations` | Tolerations for scheduling SLV pods on tainted nodes. | `[]` |
| `env` | Environment variables to be set for the SLV job container. | `{}` |
| `serviceAccount.labels` | Labels to be added to the ServiceAccount. | `{}` |
| `serviceAccount.annotations` | Annotations to be added to the ServiceAccount. | `{}` |
| `backoffLimit` | Number of retries if the job fails. | `4` |
| `ttlSecondsAfterFinished` | Time to retain job resource after completion (seconds). | `3600` |
| `schedule` | Cron expression to run as a CronJob. | None |

---

## Minimum Permissions for the Job

If you are managing RBAC manually, here are the minimal permissions required:

```yaml
rules:
  - apiGroups: ["slv.sh"]
    resources: ["slvs"]
    verbs: ["get", "list", "update"]
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "create", "list", "update", "delete"]
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get", "create", "update"]
```

> **Tip:** The default permissions are cluster-wide but can be scoped to a specific namespace as needed.

---

## Example Usage

### Run as a Job

#### Preload the Environment Secret

```bash
kubectl create secret generic slv -n slv --from-literal=SecretKey=SLV_ESK_AEAEKAAATI5CXB7QMFSUGY4RUT6UTUSK7SGMIECTJKRTQBFY6BN5ZV5M5XGF6DWLV2RVCJJSMXH43DJ6A5TK7Y6L6PYEMCDGQRBX46GUQPUIYUQ
```

#### Deploy the SLV Job

```bash
helm upgrade --install slv slv/slv-job --namespace slv --create-namespace --set jobName=my-job-$(date +%s)
```

#### Apply an SLV Object

```bash
kubectl apply -f https://slv.sh/k8s/samples/pets.slv.yaml
```

#### Run the Job Again (if needed)

```bash
helm upgrade --install slv slv/slv-job --namespace slv --create-namespace --set jobName=my-job-$(date +%s)
```

#### Retrieve the Corresponding Secret

```bash
kubectl get secret pets -n slv -o jsonpath='{.data.mycat}' | base64 --decode
```
Expected Output:
```
Kuwan
```

---

### Run as a CronJob

#### Preload the Environment Secret

```bash
kubectl create secret generic slv -n slv --from-literal=SecretKey=SLV_ESK_AEAEKAAATI5CXB7QMFSUGY4RUT6UTUSK7SGMIECTJKRTQBFY6BN5ZV5M5XGF6DWLV2RVCJJSMXH43DJ6A5TK7Y6L6PYEMCDGQRBX46GUQPUIYUQ
```

#### Deploy as CronJob

Specify a cron schedule (e.g., every hour):

```bash
helm upgrade --install slv slv/slv-job --set schedule="0 * * * *"
```

#### Apply an SLV Object

```bash
kubectl apply -f https://slv.sh/k8s/samples/pets.slv.yaml
```

#### Retrieve the Corresponding Secret

After waiting for the schedule to trigger reconciliation:

```bash
kubectl get secret pets -n slv -o jsonpath='{.data.mycat}' | base64 --decode
```
Expected Output:
```
Kuwan
```





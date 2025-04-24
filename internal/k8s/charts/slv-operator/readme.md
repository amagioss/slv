# SLV Helm Chart Configuration

This document provides an overview of the configurable parameters available in the `values.yaml` file for the SLV Helm chart. These parameters allow you to customize the behavior and deployment of SLV in a Kubernetes cluster.

---

## Configuration Parameters

### SLV Environment Configuration (`slvEnvironment`)

This section is responsible for configuring the SLV environment, which is required to decrypt SLV secrets.

- **`secretBinding`**:  
  The secret binding string passed in as plain text.


- **`k8sSecret`**:  
  The name of the Kubernetes secret that contains either the `secretKey` or the `secretBinding`.  
  - The `secretKey` should be stored under the key name `SecretKey`.  
  - The `secretBinding` should be stored under the key name `SecretBinding`.  
  **Note**: Ensure this secret exists in the same namespace as the release namespace.
  
**Note**: Ensure that at least one of `secretBinding` or `k8sSecret` is specified. SLV may not work as expected without a secret key or binding.

---

### SLV Operation Configuration (`config`)

This section defines how SLV operates within the Kubernetes cluster.

- **`mode`**:  
  Specifies the mode in which SLV should run.  
  **Possible Values**:  
  - `operator`: Runs SLV as a Kubernetes operator (Deployment). It watches for changes in the SLV CRD and takes action accordingly.  
  - `job`: Runs SLV as a Kubernetes job. It executes once and exits. Assumes that the CRDs are already created.  
  - `cronjob`: Runs SLV as a Kubernetes cronjob. It executes at the specified schedule and exits.  
  **Default**: `operator`

- **`enableWebhook`**:  
  Determines whether to enable the SLV webhook.  
  **Default**: `false`  
  **Note**: This feature is only applicable in `operator` mode and is a work in progress.

- **`replicas`**:  
  The number of replicas to be used by the deployment for SLV pods.  
  **Default**: `1`  
  **Note**: This is only applicable when the mode is set to `operator`.

- **`backoffLimit`**:  
  The number of times Kubernetes must retry the job in case of failure.  
  **Default**: `4`  
  **Note**: This is only applicable when the mode is set to `job` or `cronjob`.

- **`ttlSecondsAfterFinished`**:  
  The number of seconds after which the job must be deleted after it has finished.  
  **Default**: `3600` (1 hour)  
  **Note**: This is only applicable when the mode is set to `job` or `cronjob`.

- **`schedule`**:  
  Specifies how frequently the CronJob should run.  
  **Default**: `"0 * * * *"` (Runs at the top of every hour)  
  **Note**: This is only applicable when the mode is set to `cronjob`.

---

### SLV Pod Configuration (`runnerConfig`)

This section contains configuration parameters that are shared across all modes (`operator`, `job`, and `cronjob`).

- **`image`**:  
  The container image to be used for the SLV pod.  
  **Default**: `"ghcr.io/amagioss/slv:<chart_version>"`
  **Note**: Ensure that the image tag is the same as the chart version if being overridden. Helm may not allow installation if it differs.

- **`imagePullPolicy`**:  
  The image pull policy for the SLV pod.  
  **Default**: `"IfNotPresent"`  
  **Possible Values**:  
  - `Always`: Always pull the image.  
  - `IfNotPresent`: Pull the image only if it is not already present.  
  - `Never`: Never pull the image.

- **`resources`**:  
  Resource requests and limits for the SLV pod.  
  **Default**:  
  ```yaml
  resources:
    requests:
      memory: "256Mi"
      cpu: "200m"
    limits:
      memory: "256Mi"
      cpu: "200m"

- **`labels`**:  
  Custom labels to be applied to Kubernetes resources created by the SLV Helm chart. <br>
  **Default**: `{}` (empty object)  

- **`podLabels`**:
  Custom labels to be applied specifically to the SLV pods. <br>
  **Default**: `{}` (empty object)

- **`serviceAccountName`**:  
  The name of the Kubernetes ServiceAccount to be used by the SLV pods.  
  **Default**: `"slv-sa"`  

---

### Additional Notes

- Ensure that the `values.yaml` file is properly configured before deploying the SLV Helm chart.
- For more details on specific parameters, refer to the inline comments in the `values.yaml` file.

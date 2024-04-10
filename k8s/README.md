# SLV - Kuberenetes Integration
SLV Kubernetes Integration helps in reconciling SLV vaults as kubernetes secrets into namespaces.
SLV can create SLV's kuberenetes compatible vaults with a `--k8s` flag. Doing this will create vaults that are technically custom resources based on SLV's [CRD](https://oss.amagi.com/slv/k8s/crd.yaml).

```sh
slv vault new -v [vault-file.yaml] --k [public_key] --search [env_key_word] --k8s [k8s-secret-file-path | - | k8s-slv-object-name]
```

The `--k8s` flag takes in any of the following arguments and validates them in the following order
- An existing K8s Secret object stored as plaintext K8s Secret yaml file
- An input `-` signifies that you'd like to input the content of the K8s Secret object through stdin
- Name of the SLV's K8s custom resource which directly translates to the name of the K8s Secret Object. This creates an empty K8s compatible SLV vault file.


## Getting Started
To get started apply the SLV [CRD](https://oss.amagi.com/slv/k8s/crd.yaml) using the following command.
```sh
kubectl apply -f https://oss.amagi.com/slv/k8s/crd.yaml
```

SLV support two ways to reconcile SLV vaults as kuberenetes secrets:
1. [Operator](#operator)
2. [Job](#job)

## Operator
SLV operator is a kubenetes controller that runs inside a given cluster to write secrets into given namespaces based on changes in SLV objects.

The following example shows how it is achieved using the operator.

- Create a namespace and add SLV environment secret key as a secret (recommended to use Access Binding using KMS for cloud environments)
```sh
kubectl create ns slv
# Disclaimer: The below secret key is only for demonstration purposes. Please avoid using it in production.
kubectl create secret generic slv -n slv --from-literal=SecretKey=SLV_ESK_AEAEKAAATI5CXB7QMFSUGY4RUT6UTUSK7SGMIECTJKRTQBFY6BN5ZV5M5XGF6DWLV2RVCJJSMXH43DJ6A5TK7Y6L6PYEMCDGQRBX46GUQPUIYUQ
```
- Install the SLV Kubernetes Operator into your cluster (modify the values in the yaml file based on your requirement)
```sh
kubectl apply -f https://oss.amagi.com/slv/k8s/samples/deploy/operator.yaml
```
- Download this vault and keep it locally
```sh
curl -s https://oss.amagi.com/slv/k8s/samples/pets.slv.yaml > pets.slv.yaml
```
- Apply the downloaded vault
```sh
kubectl apply -f pets.slv.yaml
```
- Try reading the processed secret
```sh
kubectl get secret pets -o jsonpath='{.data.mycat}' | base64 --decode
```
- Add any secret value using the following command and apply again
```sh
slv vault secret put -v pets.slv.yaml -n hi --secret "Hello World"
kubectl apply -f pets.slv.yaml
```
- Try again by reading the updated secret
```sh
kubectl get secret pets -o jsonpath='{.data.hi}' | base64 --decode
```

## Job
SLV job is a one-time job that can reconcile any existing SLV objects as kubernetes secrets. This is useful in environments that can't afford to run a persistent operator or there aren't many secrets to deal with.

The following example shows how SLV objects are reconciled to secrets using the job.

- Create a namespace and add SLV environment secret key as a secret (recommended to use Access Binding using KMS for cloud environments)
```sh
kubectl create ns samplespace
# Disclaimer: The below secret key is only for demonstration purposes. Please avoid using it in production.
kubectl create secret generic slv -n samplespace --from-literal=SecretKey=SLV_ESK_AEAEKAAATI5CXB7QMFSUGY4RUT6UTUSK7SGMIECTJKRTQBFY6BN5ZV5M5XGF6DWLV2RVCJJSMXH43DJ6A5TK7Y6L6PYEMCDGQRBX46GUQPUIYUQ
```
- Download this vault and keep it locally
```sh
curl -s https://oss.amagi.com/slv/k8s/samples/pets.slv.yaml > pets.slv.yaml
```
- Apply the downloaded vault to the namespace
```sh
kubectl apply -f pets.slv.yaml -n samplespace
```
- Run the job in your namespace (modify the values in the yaml file based on your requirement)
```sh
kubectl apply -f https://oss.amagi.com/slv/k8s/samples/deploy/job.yaml -n samplespace
```
- Try reading the processed secret
```sh
kubectl get secret pets -o jsonpath='{.data.mycat}' -n samplespace | base64 --decode
```
- Add any secret value using the following command and apply again
```sh
slv vault secret put -v pets.slv.yaml -n hi --secret "Hello World"
kubectl apply -f pets.slv.yaml -n samplespace
```
- Run the job again
```sh
kubectl apply -f https://oss.amagi.com/slv/k8s/samples/deploy/job.yaml -n samplespace
```
- Try again by reading the updated secret
```sh
kubectl get secret pets -o jsonpath='{.data.hi}' -n samplespace | base64 --decode
```

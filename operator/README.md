# SLV Operator
SLV operator is a Kubernetes controller that helps in reconciling SLV vaults as kubernetes secrets into namespaces. This can be achieved by creating a vault with a `--k8s` flag. Doing so will create vaults that are custom resources managed by the SLV operator.

### A working example
- Create a namespace and add SLV_ENV_SECRET_KEY as a secret (recommended to use Access Binding using KMS for cloud environments)
```sh
kubectl create ns slv
# Disclaimer: The below secret key is only for demonstration purposes. Please avoid using it in production.
kubectl create secret generic slv -n slv --from-literal=secretkey=SLV_ESK_AEAEKAAATI5CXB7QMFSUGY4RUT6UTUSK7SGMIECTJKRTQBFY6BN5ZV5M5XGF6DWLV2RVCJJSMXH43DJ6A5TK7Y6L6PYEMCDGQRBX46GUQPUIYUQ
```
- Install the Kubernetes operator into your cluster
```sh
kubectl apply -f https://oss.amagi.com/slv/operator/samples/deploy.yaml
```
- Download this vault and keep it locally
```sh
curl -s https://oss.amagi.com/slv/operator/samples/pets.slv.yaml > pets.slv.yaml
```
- Apply the downloaded vault to the cluster
```sh
kubectl apply -f pets.slv.yaml
```
- Try reading SLV controller reconciled secrets from the cluster
```sh
kubectl get secret pets -o jsonpath='{.data.supercat}' | base64 --decode
```
- Add any secret using the following commands
```sh
slv vault put -v pets.slv.yaml -n hi -s "Hello World"
kubectl apply -f pets.slv.yaml
```
- Try reading newly created secret from the cluster
```sh
kubectl get secret pets -o jsonpath='{.data.hi}' | base64 --decode
```
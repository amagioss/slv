package main

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"oss.amagi.com/slv/internal/core/secretkey"
)

var namespace *string

func getNamespace() string {
	if namespace == nil {
		ns := os.Getenv("NAMESPACE")
		if ns == "" {
			namespaceBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
			if err != nil {
				panic(err)
			}
			ns = string(namespaceBytes)
		}
		namespace = &ns
	}
	return *namespace
}

func main() {
	secretKey, err := secretkey.Get()
	if err != nil {
		panic(err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	slvObjs, err := listSLVs(config)
	if err != nil {
		panic(err)
	}

	if err = slvsToSecrets(clientset, secretKey, slvObjs); err != nil {
		panic(err)
	}
}

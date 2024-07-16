package job

import (
	"k8s.io/client-go/kubernetes"
	"oss.amagi.com/slv/internal/k8s/utils"
)

func Run() {
	secretKey, err := utils.SecretKey()
	if err != nil {
		panic(err)
	}

	config, err := utils.GetKubeClientConfig()
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

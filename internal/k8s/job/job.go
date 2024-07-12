package job

import (
	"k8s.io/client-go/kubernetes"
	"oss.amagi.com/slv/internal/k8s/utils"
)

func Run() {
	if err := utils.InitSecretKey(); err != nil {
		panic(err.Error())
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

	if err = slvsToSecrets(clientset, utils.SecretKey(), slvObjs); err != nil {
		panic(err)
	}
}

package job

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/k8s/utils"
)

var logger = log.Default()

func Run() {
	logger.Println("Starting SLV job...")
	logger.Println(config.VersionInfo())

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

package job

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/session"
)

var logger = log.Default()

func Run() {
	logger.Println("Starting SLV job...")
	logger.Println(config.VersionInfo())

	secretKey, err := session.GetSecretKey()
	if err != nil {
		logger.Fatal(err)
	}

	kubeCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	config, err := kubeCfg.ClientConfig()
	if err != nil {
		logger.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatal(err)
	}

	slvObjs, err := listSLVs(config)
	if err != nil {
		logger.Fatal(err)
	}

	if err = slvsToSecrets(clientset, secretKey, slvObjs); err != nil {
		logger.Fatal(err)
	}
}

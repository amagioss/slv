package utils

import (
	"fmt"
	"os"
	"strings"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"slv.sh/slv/internal/core/config"
)

var (
	currentNamespace *string
	kubeConfig       *clientcmd.ClientConfig
	nameSpacedMode   = strings.ToLower(os.Getenv("SLV_K8S_NAMESPACED_MODE")) == "true"
	slvK8sEnvSecret  = func() string {
		if val := os.Getenv("SLV_K8S_ENV_SECRET"); val != "" {
			return val
		}
		return config.AppNameLowerCase
	}()
)

func IsNamespacedMode() bool {
	return nameSpacedMode
}

func getKubeConfig() clientcmd.ClientConfig {
	if kubeConfig == nil {
		kubeCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
		)
		kubeConfig = &kubeCfg
	}
	return *kubeConfig
}

func GetKubeClientConfig() (*rest.Config, error) {
	return getKubeConfig().ClientConfig()
}

func GetClusterInfo() (name, address, user string, err error) {
	kubeconfig := getKubeConfig()
	rawConfig, err := kubeconfig.RawConfig()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get raw kubeconfig: %w", err)
	}
	currentContext := rawConfig.CurrentContext
	contextDetails, exists := rawConfig.Contexts[currentContext]
	if !exists {
		return "", "", "", fmt.Errorf("context %s not found in kubeconfig", currentContext)
	}
	clusterDetails, exists := rawConfig.Clusters[contextDetails.Cluster]
	if !exists {
		return "", "", "", fmt.Errorf("cluster %s not found in kubeconfig", contextDetails.Cluster)
	}
	return contextDetails.Cluster, clusterDetails.Server, contextDetails.AuthInfo, nil
}

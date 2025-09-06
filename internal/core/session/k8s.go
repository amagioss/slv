package session

import (
	"context"
	"fmt"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments/envproviders"
)

const (
	slvK8sConfigMap         = config.AppNameLowerCase
	publicKeyNameEC         = "PublicKeyEC"
	publicKeyNamePQ         = "PublicKeyPQ"
	envar_NAMESPACE         = "NAMESPACE"
	envar_SLV_K8S_NAMESPACE = "SLV_K8S_NAMESPACE"
	namespaceFile           = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

var (
	currentNamespace *string
	kubeConfig       *clientcmd.ClientConfig
)

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

func getKubeClientConfig() (*rest.Config, error) {
	return getKubeConfig().ClientConfig()
}

func getKubeClientSet() (*kubernetes.Clientset, error) {
	config, err := getKubeClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}
	return kubernetes.NewForConfig(config)
}

func getSecretKeyFor(clientset *kubernetes.Clientset, namespace string) (secretKey *crypto.SecretKey, err error) {
	if clientset == nil {
		return nil, fmt.Errorf("failed to get k8s clientset: %w", err)
	}
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), slvK8sSecret, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if secretKey, err = extractSecretKeyFromSecret(secret); secretKey != nil {
		return secretKey, err
	}
	if configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), slvK8sSecret, metav1.GetOptions{}); err == nil {
		if secretKey, err = extractSecretKeyFromConfigMapBinding(configMap); secretKey != nil {
			return secretKey, err
		}
	}
	return nil, fmt.Errorf("secret key not found")
}

func extractSecretKeyFromSecret(slvSecret *corev1.Secret) (*crypto.SecretKey, error) {
	for k, v := range slvSecret.Data {
		lowerCaseKey := strings.ToLower(k)
		if lowerCaseKey == "secretkey" || lowerCaseKey == "secret_key" {
			return crypto.SecretKeyFromString(string(v))
		}
		if lowerCaseKey == "secretbinding" || lowerCaseKey == "secret_binding" {
			return envproviders.GetSecretKeyFromSecretBinding(string(v))
		}
	}
	return nil, nil
}

func extractSecretKeyFromConfigMapBinding(configMap *corev1.ConfigMap) (*crypto.SecretKey, error) {
	for k, v := range configMap.Data {
		lowerCaseKey := strings.ToLower(k)
		if lowerCaseKey == "secretbinding" || lowerCaseKey == "secret_binding" {
			return envproviders.GetSecretKeyFromSecretBinding(v)
		}
	}
	return nil, nil
}

func putPublicKeyToConfigMap(clientset *kubernetes.Clientset, publicKeyStrEC, publicKeyStrPQ string) error {
	namespace := GetK8sNamespace()
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), slvK8sConfigMap, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			configMap = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      slvK8sConfigMap,
					Namespace: namespace,
				},
				Data: map[string]string{
					publicKeyNameEC: publicKeyStrEC,
					publicKeyNamePQ: publicKeyStrPQ,
				},
			}
			_, err = clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), configMap, metav1.CreateOptions{})
		}
	} else {
		if configMap.Data == nil {
			configMap.Data = make(map[string]string)
		}
		if configMap.Data[publicKeyNameEC] != publicKeyStrEC || configMap.Data[publicKeyNamePQ] != publicKeyStrPQ {
			configMap.Data[publicKeyNameEC] = publicKeyStrEC
			configMap.Data[publicKeyNamePQ] = publicKeyStrPQ
			_, err = clientset.CoreV1().ConfigMaps(namespace).Update(context.Background(), configMap, metav1.UpdateOptions{})
		}
	}
	return err
}

func GetPublicKeyFromK8s(namespace string, pq bool) (string, error) {
	clientset, err := getKubeClientSet()
	if err != nil {
		return "", fmt.Errorf("failed to get k8s clientset: %w", err)
	}
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), slvK8sConfigMap, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	configMapKey := publicKeyNameEC
	if pq {
		configMapKey = publicKeyNamePQ
	}
	publicKeyStr, ok := configMap.Data[configMapKey]
	if !ok {
		return "", fmt.Errorf("public key not found in configmap: %s", slvK8sConfigMap)
	}
	return publicKeyStr, nil
}

func GetK8sNamespace() string {
	if currentNamespace == nil {
		ns := os.Getenv(envar_NAMESPACE)
		if ns == "" {
			ns = os.Getenv(envar_SLV_K8S_NAMESPACE)
		}
		if ns == "" {
			if namespaceBytes, err := os.ReadFile(namespaceFile); err == nil {
				ns = string(namespaceBytes)
			}
		}
		currentNamespace = &ns
	}
	return *currentNamespace
}

func isInKubernetesCluster() bool {
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		return true
	}
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" || os.Getenv("KUBERNETES_SERVICE_PORT") != "" {
		return true
	}
	if _, err := os.Stat("/var/run/secrets/kubernetes.io"); err == nil {
		return true
	}
	return false
}

func GetK8sClusterInfo() (name, address, user string, err error) {
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

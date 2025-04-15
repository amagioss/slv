package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"oss.amagi.com/slv/internal/core/crypto"
	"oss.amagi.com/slv/internal/core/environments/providers"
)

const (
	secretKeyName           = "SecretKey"
	publicKeyNameEC         = "PublicKeyEC"
	publicKeyNamePQ         = "PublicKeyPQ"
	envar_NAMESPACE         = "NAMESPACE"
	envar_SLV_K8S_NAMESPACE = "SLV_K8S_NAMESPACE"
)

func getKubeClientSet() (*kubernetes.Clientset, error) {
	config, err := GetKubeClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}
	return kubernetes.NewForConfig(config)
}

func GetSecretKeyFor(clientset *kubernetes.Clientset, namespace string) (secretKey *crypto.SecretKey, err error) {
	if clientset == nil {
		if clientset, err = getKubeClientSet(); err != nil {
			return nil, fmt.Errorf("failed to get k8s clientset: %w", err)
		}
	}
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), slvK8sEnvSecret, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if secretKey, err = ExtractSecretKeyFromSecret(secret); secretKey != nil {
		return secretKey, err
	}
	if configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), slvK8sEnvSecret, metav1.GetOptions{}); err == nil {
		if secretKey, err = ExtractSecretKeyFromConfigMapBinding(configMap); secretKey != nil {
			return secretKey, err
		}
	}
	return nil, fmt.Errorf("secret key not found")
}

func ExtractSecretKeyFromSecret(slvSecret *corev1.Secret) (*crypto.SecretKey, error) {
	for k, v := range slvSecret.Data {
		lowerCaseKey := strings.ToLower(k)
		if lowerCaseKey == "secretkey" || lowerCaseKey == "secret_key" {
			return crypto.SecretKeyFromString(string(v))
		}
		if lowerCaseKey == "secretbinding" || lowerCaseKey == "secret_binding" {
			return providers.GetSecretKeyFromSecretBinding(string(v))
		}
	}
	return nil, nil
}

func ExtractSecretKeyFromConfigMapBinding(configMap *corev1.ConfigMap) (*crypto.SecretKey, error) {
	for k, v := range configMap.Data {
		lowerCaseKey := strings.ToLower(k)
		if lowerCaseKey == "secretbinding" || lowerCaseKey == "secret_binding" {
			return providers.GetSecretKeyFromSecretBinding(v)
		}
	}
	return nil, nil
}

func putPublicKeyToConfigMap(clientset *kubernetes.Clientset, publicKeyStrEC, publicKeyStrPQ string) error {
	namespace := GetCurrentNamespace()
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), slvK8sEnvSecret, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			configMap = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      slvK8sEnvSecret,
					Namespace: GetCurrentNamespace(),
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

func GetCurrentNamespace() string {
	if currentNamespace == nil {
		ns := os.Getenv(envar_NAMESPACE)
		if ns == "" {
			ns = os.Getenv(envar_SLV_K8S_NAMESPACE)
		}
		if ns == "" {
			namespaceBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
			if err != nil {
				panic(err)
			}
			ns = string(namespaceBytes)
		}
		currentNamespace = &ns
	}
	return *currentNamespace
}

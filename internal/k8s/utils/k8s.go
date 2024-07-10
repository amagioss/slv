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
	"oss.amagi.com/slv/internal/core/config"
	"oss.amagi.com/slv/internal/core/crypto"
	"oss.amagi.com/slv/internal/core/environments/providers"
)

const (
	resourceName    = config.AppNameLowerCase
	secretKeyName   = "SecretKey"
	publicKeyNameEC = "PublicKeyEC"
	publicKeyNamePQ = "PublicKeyPQ"
)

func getKubeClientSet() (*kubernetes.Clientset, error) {
	config, err := GetKubeClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}
	return kubernetes.NewForConfig(config)
}

func getSecretKeyFromCluster(clientset *kubernetes.Clientset) (*crypto.SecretKey, error) {
	namespace := GetCurrentNamespace()
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	for k, v := range secret.Data {
		lowerCaseKey := strings.ToLower(k)
		if lowerCaseKey == "secretkey" || lowerCaseKey == "secret_key" {
			return crypto.SecretKeyFromString(string(v))
		}
		if lowerCaseKey == "secretbinding" || lowerCaseKey == "secret_binding" {
			return providers.GetSecretKeyFromSecretBinding(string(v))
		}
	}
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	for k, v := range configMap.Data {
		lowerCaseKey := strings.ToLower(k)
		if lowerCaseKey == "secretbinding" || lowerCaseKey == "secret_binding" {
			return providers.GetSecretKeyFromSecretBinding(v)
		}
	}
	return nil, fmt.Errorf("secret key not found")
}

func putSecretKeyToSecret(clientset *kubernetes.Clientset, secretKeyStr string) error {
	namespace := GetCurrentNamespace()
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			secret = &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: namespace,
				},
				Data: map[string][]byte{
					secretKeyName: []byte(secretKeyStr),
				},
			}
			_, err = clientset.CoreV1().Secrets(namespace).Create(context.Background(), secret, metav1.CreateOptions{})
		}
	} else {
		updated := false
		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}
		for k, v := range secret.Data {
			lowerCaseKey := strings.ToLower(k)
			if lowerCaseKey == "secretkey" || lowerCaseKey == "secret_key" {
				if string(v) != secretKeyStr {
					secret.Data[k] = []byte(secretKeyStr)
					updated = true
				}
			}
		}
		if string(secret.Data[secretKeyName]) != secretKeyStr {
			secret.Data[secretKeyName] = []byte(secretKeyStr)
			updated = true
		}
		if updated {
			_, err = clientset.CoreV1().Secrets(namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
		}
	}
	return err
}

func putPublicKeyToConfigMap(clientset *kubernetes.Clientset, publicKeyStrEC, publicKeyStrPQ string) error {
	namespace := GetCurrentNamespace()
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			configMap = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
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

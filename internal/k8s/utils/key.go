package utils

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"oss.amagi.com/slv/internal/core/crypto"
	"oss.amagi.com/slv/internal/core/environments"
	"oss.amagi.com/slv/internal/core/secretkey"
)

var sKey *crypto.SecretKey

func InitSecretKey() (err error) {
	if sKey == nil {
		sKey, _ = secretkey.Get()
		if clientset, _ := getKubeClientSet(); clientset != nil {
			if sKey == nil {
				if sKey, err = getSecretKeyFromCluster(clientset); err != nil && isEnvGenEnabled() {
					sKey, err = crypto.NewSecretKey(environments.EnvironmentKey)
				}
			}
			if err == nil && sKey != nil {
				if err = putSecretKeyToSecret(clientset, sKey.String()); err == nil {
					var pKeyEC, pKeyPQ *crypto.PublicKey
					if pKeyEC, err = sKey.PublicKey(false); err == nil {
						if pKeyPQ, err = sKey.PublicKey(true); err == nil {
							var pkStrEC, pkStrPQ string
							if pkStrEC, err = pKeyEC.String(); err == nil {
								if pkStrPQ, err = pKeyPQ.String(); err == nil {
									err = putPublicKeyToConfigMap(clientset, pkStrEC, pkStrPQ)
								}
							}
						}
					}
				}
			}
		}
		if err != nil {
			sKey = nil
		}
	}
	return err
}

func SecretKey() *crypto.SecretKey {
	if err := InitSecretKey(); err != nil {
		panic(err.Error())
	}
	return sKey
}

func GetPublicKeyFromK8s(namespace string, pq bool) (string, error) {
	clientset, err := getKubeClientSet()
	if err != nil {
		return "", fmt.Errorf("failed to get k8s clientset: %w", err)
	}
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	configMapKey := publicKeyNameEC
	if pq {
		configMapKey = publicKeyNamePQ
	}
	publicKeyStr, ok := configMap.Data[configMapKey]
	if !ok {
		return "", fmt.Errorf("public key not found in configmap: %s", resourceName)
	}
	return publicKeyStr, nil
}
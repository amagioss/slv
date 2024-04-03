package main

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"oss.amagi.com/slv/core/crypto"
	slvv1 "oss.amagi.com/slv/k8s/api/v1"
)

func slvsToSecrets(clientset *kubernetes.Clientset, secretKey *crypto.SecretKey, slvObjs []slvv1.SLV) error {
	for _, slvObj := range slvObjs {
		fmt.Println("Attempting to unlock SLV vault:", slvObj.Name)
		if err := slvObj.Spec.Unlock(*secretKey); err != nil {
			return err
		}
		slvSecretMap, err := slvObj.Spec.GetAllSecrets()
		if err != nil {
			return err
		}
		secret, err := clientset.CoreV1().Secrets(getNamespace()).Get(context.Background(), slvObj.Name, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: slvObj.Name,
					},
					Data: slvSecretMap,
				}
				if _, err = clientset.CoreV1().Secrets(getNamespace()).Create(context.Background(), secret, metav1.CreateOptions{}); err != nil {
					return err
				}
				log.Println("Created secret", secret.Name)
			} else {
				return err
			}
		} else {
			secret.Data = slvSecretMap
			_, err = clientset.CoreV1().Secrets(getNamespace()).Update(context.Background(), secret, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
			log.Println("Updated secret", secret.Name)
		}
	}
	return nil
}

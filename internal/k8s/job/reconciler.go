package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"oss.amagi.com/slv/internal/core/config"
	"oss.amagi.com/slv/internal/core/crypto"
	slvv1 "oss.amagi.com/slv/internal/k8s/api/v1"
)

const (
	slvVersionAnnotationKey = slvv1.Group + "/version"
)

func isAnnotationUpdateRequred(slvAnnotations, secretAnnotations map[string]string) bool {
	if len(secretAnnotations) != (len(slvAnnotations) + 1) {
		return true
	}
	for k, v := range slvAnnotations {
		if secretAnnotations[k] != v {
			return true
		}
	}
	return secretAnnotations[slvVersionAnnotationKey] != config.Version
}

func toSecret(clientset *kubernetes.Clientset, secretKey *crypto.SecretKey, slvObj slvv1.SLV) error {
	if err := slvObj.Spec.Unlock(*secretKey); err != nil {
		return err
	}
	slvSecretMap, err := slvObj.Spec.GetAllSecrets()
	if err != nil {
		return err
	}
	secret, err := clientset.CoreV1().Secrets(slvObj.Namespace).Get(context.Background(), slvObj.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:        slvObj.Name,
					Namespace:   slvObj.Namespace,
					Annotations: slvObj.Annotations,
				},
				Type: slvObj.Type,
				Data: slvSecretMap,
			}
			if secret.Annotations == nil {
				secret.Annotations = make(map[string]string)
			}
			secret.Annotations[slvVersionAnnotationKey] = config.Version
			if _, err = clientset.CoreV1().Secrets(slvObj.Namespace).Create(context.Background(), secret, metav1.CreateOptions{}); err != nil {
				return err
			}
			log.Println("Created secret:", secret.Name)
		} else {
			return err
		}
	} else {
		updateRequired := false
		if len(secret.Data) != len(slvSecretMap) {
			updateRequired = true
			secret.Data = slvSecretMap
		} else {
			for k, v := range slvSecretMap {
				if !bytes.Equal(secret.Data[k], v) {
					updateRequired = true
					secret.Data = slvSecretMap
					break
				}
			}
		}
		if isAnnotationUpdateRequred(slvObj.Annotations, secret.Annotations) {
			secret.Annotations = slvObj.Annotations
			if secret.Annotations == nil {
				secret.Annotations = make(map[string]string)
			}
			secret.Annotations[slvVersionAnnotationKey] = config.Version
			updateRequired = true
		}
		if secret.Type != slvObj.Type {
			secret.Type = slvObj.Type
			updateRequired = true
		}
		var msg string
		if updateRequired {
			if _, err = clientset.CoreV1().Secrets(slvObj.Namespace).Update(context.Background(), secret, metav1.UpdateOptions{}); err != nil {
				return err
			}
			msg = "Updated secret: " + secret.Name
		} else {
			msg = "No update required for secret: " + secret.Name
		}
		log.Println(msg)
	}
	return nil
}

func slvsToSecrets(clientset *kubernetes.Clientset, secretKey *crypto.SecretKey, slvObjs []slvv1.SLV) error {
	var errors []error
	for _, slvObj := range slvObjs {
		if err := toSecret(clientset, secretKey, slvObj); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) == 1 {
		return errors[0]
	} else if len(errors) > 1 {
		return fmt.Errorf("multiple errors occurred: %v", errors)
	}
	return nil
}

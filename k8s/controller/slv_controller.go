/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/amagimedia/slv"
	"github.com/amagimedia/slv/core/secretkeystore"
	k8samagicomv1 "github.com/amagimedia/slv/k8s/api/v1"
)

const (
	secretManagedByAnnotationKey   = k8samagicomv1.Group + "/managed-by"
	secretManagedByAnnotationValue = slv.AppName
	secretSLVVersionAnnotationKey  = k8samagicomv1.Group + "/slv-version"
)

var secretSLVVersionAnnotationValue = slv.Version

// SLVReconciler reconciles a SLV object
type SLVReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=k8s.amagi.com,resources=slv,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.amagi.com,resources=slv/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.amagi.com,resources=slv/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SLV object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *SLVReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	logger := log.FromContext(ctx)

	var slvCR k8samagicomv1.SLV
	if err := r.Get(ctx, req.NamespacedName, &slvCR); err != nil {
		if !errors.IsNotFound(err) {
			logger.Error(err, "Failed to get SLV")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	secretKey, err := secretkeystore.GetSecretKey()
	if err != nil {
		logger.Error(err, "SLV has no configured environment")
		return ctrl.Result{}, err
	}
	vault := slvCR.Vault
	if err = vault.Unlock(*secretKey); err != nil {
		logger.Error(err, "Failed to unlock vault", "Vault", vault)
		return ctrl.Result{}, err
	}
	slvSecretMap, err := vault.GetAllSecrets()
	if err != nil {
		logger.Error(err, "Failed to get all secrets from vault", "Vault", vault)
		return ctrl.Result{}, err
	}

	// Check if the secret exists
	secret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: slvCR.Name, Namespace: req.Namespace}, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create secret
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      slvCR.Name,
					Namespace: req.Namespace,
					Annotations: map[string]string{
						secretManagedByAnnotationKey:  secretManagedByAnnotationValue,
						secretSLVVersionAnnotationKey: secretSLVVersionAnnotationValue,
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: slvCR.APIVersion,
							Kind:       slvCR.Kind,
							Name:       slvCR.Name,
							UID:        slvCR.UID,
						},
					},
				},
				Data: slvSecretMap,
			}
			if err := r.Create(ctx, secret); err != nil {
				logger.Error(err, "Failed to create secret", "Secret", secret)
				return ctrl.Result{}, err
			}
			logger.Info("Created secret", "Secret", slvCR.Name)
		} else {
			logger.Error(err, "Failed to get secret", "Secret", secret)
			return ctrl.Result{}, err
		}
	} else {
		// Update secret
		secret.Data = slvSecretMap
		if secret.Annotations == nil {
			secret.Annotations = make(map[string]string)
		}
		secret.Annotations[secretManagedByAnnotationKey] = secretManagedByAnnotationValue
		secret.Annotations[secretSLVVersionAnnotationKey] = secretSLVVersionAnnotationValue
		secret.SetOwnerReferences([]metav1.OwnerReference{
			{
				APIVersion: slvCR.APIVersion,
				Kind:       slvCR.Kind,
				Name:       slvCR.Name,
				UID:        slvCR.UID,
			},
		})
		if err := r.Update(ctx, secret); err != nil {
			logger.Error(err, "Failed to update secret", "Secret", secret)
			return ctrl.Result{}, err
		}
		logger.Info("Updated secret", "Secret", slvCR.Name)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SLVReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8samagicomv1.SLV{}).
		Complete(r)
}

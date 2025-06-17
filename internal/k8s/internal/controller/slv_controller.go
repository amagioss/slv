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
	"bytes"
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/session"
	slvv1 "slv.sh/slv/internal/k8s/api/v1"
)

const (
	slvVersionAnnotationKey = config.K8SLVAnnotationVersionKey
	slvResourceName         = config.AppNameLowerCase
)

// SLVReconciler reconciles a SLV object
type SLVReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *SLVReconciler) returnError(ctx context.Context,
	slvObj *slvv1.SLV, logger *logr.Logger, err error, msg string) (ctrl.Result, error) {
	logger.Error(err, msg, slvv1.Kind, slvObj.Name)
	if slvObj.Status.Error == "" {
		slvObj.Status.Error = err.Error()
		if err := r.Status().Update(ctx, slvObj); err != nil {
			logger.Error(err, "Failed to update status", slvv1.Kind, slvObj.Name)
			return ctrl.Result{}, err
		}
	} else {
		err = nil
	}
	return ctrl.Result{}, err
}

func (r *SLVReconciler) success(ctx context.Context,
	slvObj *slvv1.SLV, logger *logr.Logger, msg string) error {
	logger.Info(msg, "Secret", slvObj.Name)
	if slvObj.Status.Error != "" {
		slvObj.Status = slvv1.SLVStatus{}
		if err := r.Status().Update(ctx, slvObj); err != nil {
			logger.Error(err, "Failed to update status", slvv1.Kind, slvObj.Name)
			return err
		}
	}
	return nil
}

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

//+kubebuilder:rbac:groups=slv.sh,resources=slvs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=slv.sh,resources=slvs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=slv.sh,resources=slvs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SLV object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *SLVReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling SLV")

	var slvObj slvv1.SLV
	if err := r.Get(ctx, req.NamespacedName, &slvObj); err != nil {
		if !errors.IsNotFound(err) {
			logger.Error(err, "Failed to get SLV")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	vault := slvObj.Vault
	if secretKey, err := session.GetSecretKey(); err != nil {
		return r.returnError(ctx, &slvObj, &logger, err, "Failed to get secret key")
	} else {
		if err = vault.Unlock(secretKey); err != nil {
			return r.returnError(ctx, &slvObj, &logger, err, "Failed to unlock vault")
		}
	}
	slvSecretMap, err := vault.GetAllValues()
	if err != nil {
		return r.returnError(ctx, &slvObj, &logger, err, "Failed to get all secrets from vault")
	}

	// Check if the secret exists
	secret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: slvObj.Name, Namespace: req.Namespace}, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create secret
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:        slvObj.Name,
					Namespace:   req.Namespace,
					Annotations: slvObj.Annotations,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: slvObj.APIVersion,
							Kind:       slvObj.Kind,
							Name:       slvObj.Name,
							UID:        slvObj.UID,
							Controller: &[]bool{true}[0],
						},
					},
				},
				Type: corev1.SecretType(slvObj.Type),
				Data: slvSecretMap,
			}
			if secret.Annotations == nil {
				secret.Annotations = make(map[string]string)
			}
			secret.Annotations[slvVersionAnnotationKey] = config.Version
			if err = controllerutil.SetControllerReference(&slvObj, secret, r.Scheme); err != nil {
				return r.returnError(ctx, &slvObj, &logger, err, "Failed to set controller reference for secret")
			}
			if err := r.Create(ctx, secret); err != nil {
				return r.returnError(ctx, &slvObj, &logger, err, "Failed to create secret")
			}
			if err := r.success(ctx, &slvObj, &logger, "Created secret"); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			return r.returnError(ctx, &slvObj, &logger, err, "Error getting secret for SLV Object")
		}
	} else {
		// Update secret
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
		if string(secret.Type) != slvObj.Type {
			secret.Type = corev1.SecretType(slvObj.Type)
			updateRequired = true
		}
		var msg string
		if updateRequired {
			if !controllerutil.HasControllerReference(secret) {
				if err = controllerutil.SetControllerReference(&slvObj, secret, r.Scheme); err != nil {
					return r.returnError(ctx, &slvObj, &logger, err, "Failed to set controller reference for secret")
				}
			}
			if err = r.Update(ctx, secret); err != nil {
				return r.returnError(ctx, &slvObj, &logger, err, "Failed to update secret")
			}
			msg = "Updated secret"
		} else {
			msg = "No update required for secret"
		}
		if err := r.success(ctx, &slvObj, &logger, msg); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SLVReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&slvv1.SLV{}).
		Watches(&corev1.Secret{}, handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &slvv1.SLV{})).
		Complete(r)
}

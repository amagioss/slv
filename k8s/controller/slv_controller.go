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

	"github.com/amagimedia/slv"
	"github.com/amagimedia/slv/core/crypto"
	"github.com/amagimedia/slv/core/secretkeystore"
	k8samagicomv1 "github.com/amagimedia/slv/k8s/api/v1"
	"github.com/go-logr/logr"
)

const (
	secretManagedByAnnotationKey   = k8samagicomv1.Group + "/managed-by"
	secretManagedByAnnotationValue = slv.AppName
	secretSLVVersionAnnotationKey  = k8samagicomv1.Group + "/slv-version"
)

var secretSLVVersionAnnotationValue = slv.Version
var secretKey *crypto.SecretKey

func InitSLVSecretKey() error {
	if secretKey == nil {
		sk, err := secretkeystore.GetSecretKey()
		if err != nil {
			return err
		}
		secretKey = sk
	}
	return nil
}

// SLVReconciler reconciles a SLV object
type SLVReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *SLVReconciler) returnError(ctx context.Context, req ctrl.Request,
	slvObj *k8samagicomv1.SLV, logger *logr.Logger, err error, msg string) (ctrl.Result, error) {
	logger.Error(err, msg, "SLV", slvObj.Name)
	if slvObj.Status.Error == "" {
		slvObj.Status.Error = err.Error()
		if err := r.Status().Update(ctx, slvObj); err != nil {
			logger.Error(err, "Failed to update status", "SLV", slvObj.Name)
			return ctrl.Result{}, err
		}
	} else {
		err = nil
	}
	return ctrl.Result{}, err
}

func (r *SLVReconciler) success(ctx context.Context, req ctrl.Request,
	slvObj *k8samagicomv1.SLV, logger *logr.Logger, msg string) error {
	logger.Info(msg, "Secret", slvObj.Name)
	if slvObj.Status.Error != "" {
		slvObj.Status = k8samagicomv1.SLVStatus{}
		if err := r.Status().Update(ctx, slvObj); err != nil {
			logger.Error(err, "Failed to update status", "SLV", slvObj.Name)
			return err
		}
	}
	return nil
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

	logger.Info("Reconciling SLV")

	var slvObj k8samagicomv1.SLV
	if err := r.Get(ctx, req.NamespacedName, &slvObj); err != nil {
		if !errors.IsNotFound(err) {
			logger.Error(err, "Failed to get SLV")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	vault := slvObj.Vault
	if err := vault.Unlock(*secretKey); err != nil {
		return r.returnError(ctx, req, &slvObj, &logger, err, "Failed to unlock vault")
	}
	slvSecretMap, err := vault.GetAllSecrets()
	if err != nil {
		return r.returnError(ctx, req, &slvObj, &logger, err, "Failed to get all secrets from vault")
	}

	// Check if the secret exists
	secret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: slvObj.Name, Namespace: req.Namespace}, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create secret
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      slvObj.Name,
					Namespace: req.Namespace,
					Annotations: map[string]string{
						secretManagedByAnnotationKey:  secretManagedByAnnotationValue,
						secretSLVVersionAnnotationKey: secretSLVVersionAnnotationValue,
					},
				},
				Data: slvSecretMap,
			}
			if err = controllerutil.SetControllerReference(&slvObj, secret, r.Scheme); err != nil {
				return r.returnError(ctx, req, &slvObj, &logger, err, "Failed to set controller reference for secret")
			}
			if err = controllerutil.SetOwnerReference(&slvObj, secret, r.Scheme); err != nil {
				return r.returnError(ctx, req, &slvObj, &logger, err, "Failed to set owner reference for secret")
			}
			if err := r.Create(ctx, secret); err != nil {
				return r.returnError(ctx, req, &slvObj, &logger, err, "Failed to create secret")
			}
			if err := r.success(ctx, req, &slvObj, &logger, "Created secret"); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			return r.returnError(ctx, req, &slvObj, &logger, err, "Error getting secret for SLV Object")
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
		if secret.Annotations == nil {
			secret.Annotations = make(map[string]string)
		}
		if secret.Annotations[secretManagedByAnnotationKey] != secretManagedByAnnotationValue {
			secret.Annotations[secretManagedByAnnotationKey] = secretManagedByAnnotationValue
			updateRequired = true
		}
		if secret.Annotations[secretSLVVersionAnnotationKey] != secretSLVVersionAnnotationValue {
			secret.Annotations[secretSLVVersionAnnotationKey] = secretSLVVersionAnnotationValue
			updateRequired = true
		}
		var msg string
		if updateRequired {
			if err = controllerutil.SetControllerReference(&slvObj, secret, r.Scheme); err != nil {
				return r.returnError(ctx, req, &slvObj, &logger, err, "Failed to set controller reference for secret")
			}
			if err = controllerutil.SetOwnerReference(&slvObj, secret, r.Scheme); err != nil {
				return r.returnError(ctx, req, &slvObj, &logger, err, "Failed to set owner reference for secret")
			}
			if err = r.Update(ctx, secret); err != nil {
				return r.returnError(ctx, req, &slvObj, &logger, err, "Failed to update secret")
			}
			msg = "Updated secret"
		} else {
			msg = "No update required for secret"
		}
		if err := r.success(ctx, req, &slvObj, &logger, msg); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SLVReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8samagicomv1.SLV{}).
		Watches(&corev1.Secret{}, handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &k8samagicomv1.SLV{})).
		Complete(r)
}

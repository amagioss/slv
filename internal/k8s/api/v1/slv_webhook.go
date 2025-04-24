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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"slv.sh/slv/internal/core/config"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/k8s/utils"
)

// log is for logging in this package.
var slvlog = logf.Log.WithName(config.AppNameLowerCase)

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *SLV) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate,mutating=false,failurePolicy=fail,sideEffects=None,groups=slv.sh,resources=slvs,verbs=create;update,versions=v1,name=validate-slv,admissionReviewVersions=v1

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *SLV) Default() {
	slvlog.Info("default", "name", r.Name)
	// Set SLV default values - not required for now
}

func (r *SLV) validateSLV() (err error) {
	vault := r
	if utils.IsNamespacedMode() && r.Namespace != utils.GetCurrentNamespace() {
		if namespacedSecretKey, _ := utils.GetSecretKeyFor(nil, r.Namespace); namespacedSecretKey != nil {
			vault.Unlock(namespacedSecretKey)
		}
	}
	if vault.IsLocked() {
		var secretKey *crypto.SecretKey
		if secretKey, err = utils.SecretKey(); err != nil {
			slvlog.Error(err, "failed to get secret key", "name", r.Name)
		}
		if err = vault.Unlock(secretKey); err != nil {
			slvlog.Error(err, "failed to unlock vault", "name", r.Name)
			return err
		}
	}
	if _, err = vault.List(true); err != nil {
		slvlog.Error(err, "failed to get all secrets", "name", r.Name)
		return err
	}
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *SLV) ValidateCreate() (admission.Warnings, error) {
	slvlog.Info("Validating create", "name", r.Name)
	return nil, r.validateSLV()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *SLV) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	slvlog.Info("Validating update", "name", r.Name)
	return nil, r.validateSLV()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *SLV) ValidateDelete() (admission.Warnings, error) {
	slvlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

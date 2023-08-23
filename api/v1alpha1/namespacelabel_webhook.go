/*
Copyright 2023.

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

package v1alpha1

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// list of disallowed prefixes
var disallowedPrefixes = []string{
	"kubernetes.io/",
}

// log is for logging in this package.
var namespacelabellog = logf.Log.WithName("namespacelabel-resource")

func (r *NamespaceLabel) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-dana-io-dana-io-v1alpha1-namespacelabel,mutating=false,failurePolicy=fail,sideEffects=None,groups=dana.io.dana.io,resources=namespacelabels,verbs=create;update,versions=v1alpha1,name=vnamespacelabel.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &NamespaceLabel{}

// ValidateCreate implements webhook.Validator to validate the creation of NamespaceLabel objects.
// It checks if any label has a disallowed prefix and returns an error if so.
func (r *NamespaceLabel) ValidateCreate() (admission.Warnings, error) {
	namespacelabellog.Info("validate create", "name", r.Name)

	// Iterating through all the labels in the spec
	for key := range r.Spec.Labels {
		// Check if the label key has any disallowed prefix
		for _, prefix := range disallowedPrefixes {
			if strings.HasPrefix(key, prefix) {
				return nil, fmt.Errorf("label with key %q is not allowed to have the '%s' prefix", key, prefix)
			}
		}
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator to validate the update of NamespaceLabel objects.
// It reuses the logic from ValidateCreate since the validation criteria are the same.
func (r *NamespaceLabel) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	namespacelabellog.Info("validate update", "name", r.Name)

	// Reuse the validation logic from ValidateCreate
	return r.ValidateCreate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *NamespaceLabel) ValidateDelete() (admission.Warnings, error) {
	namespacelabellog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

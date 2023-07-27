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

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"sigs.k8s.io/controller-runtime/pkg/handler"

	danaiodanaiov1alpha1 "dana.io/hello-world/api/v1alpha1"
)

// NamespaceLabelReconciler reconciles a NamespaceLabel object
type NamespaceLabelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dana.io.dana.io,resources=namespacelabels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dana.io.dana.io,resources=namespacelabels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dana.io.dana.io,resources=namespacelabels/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NamespaceLabel object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func namespaceReconcile(ctx context.Context, req ctrl.Request, r *NamespaceLabelReconciler) (bool, error) {
	logger := log.FromContext(ctx)

	var namespace corev1.Namespace

	if err := r.Get(ctx, types.NamespacedName{Name: req.Name}, &namespace); err != nil {
		// requeue the request if we could not get the namespace
		return false, err
	}

	if namespace.ObjectMeta.Labels == nil {
		namespace.ObjectMeta.Labels = make(map[string]string)
	}

	var namespaceLabelList danaiodanaiov1alpha1.NamespaceLabelList

	if err := r.List(ctx, &namespaceLabelList, client.InNamespace(req.Namespace)); err != nil {
		return false, nil
	}

	logger.Info("hey", "hello", namespaceLabelList)

	for _, namespaceLabel := range namespaceLabelList.Items {
		for key, value := range namespaceLabel.Spec.Labels {
			namespace.ObjectMeta.Labels[key] = value
			logger.Info("loop key values", "key", key, "value", value)
		}
		logger.Info("loop namespace label", "namespaceLabel", namespaceLabel)
	}
	// update the namespace with the new labels
	if err := r.Update(ctx, &namespace); err != nil {
		logger.Info("this is it", "error", err)
		return false, err
	} else {
		return true, nil
	}
}

func (r *NamespaceLabelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger := log.FromContext(ctx)

	logger.Info("Reconcile invoked", "req", req)

	// if this is reconcile because change in namespace object, reconcile accordingly and return
	isNamespaceReconcile, _ := namespaceReconcile(ctx, req, r)

	logger.Info("namespace.ObjectMeta.Labels", "req", isNamespaceReconcile)

	if isNamespaceReconcile {
		return ctrl.Result{}, nil
	}

	namespaceLabel := danaiodanaiov1alpha1.NamespaceLabel{}
	if err := r.Get(ctx, req.NamespacedName, &namespaceLabel); err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// the namespace we'll apply labels to will have the same name as the NamespaceLabel object

	var namespace corev1.Namespace
	if err := r.Get(ctx, types.NamespacedName{Name: req.Namespace}, &namespace); err != nil {
		// requeue the request if we could not get the namespace
		return ctrl.Result{}, err
	}

	logger.Info("namespace.ObjectMeta.Labels", "req", namespace.ObjectMeta.Labels)

	// update the labels
	if namespace.ObjectMeta.Labels == nil {
		namespace.ObjectMeta.Labels = make(map[string]string)
	}
	logger.Info("namespacelabel status", "namespacelabel", namespaceLabel.Status)
	for key, value := range namespaceLabel.Spec.Labels {
		namespace.ObjectMeta.Labels[key] = value
	}

	// update the namespace with the new labels
	if err := r.Update(ctx, &namespace); err != nil {
		return ctrl.Result{}, err
	}

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&danaiodanaiov1alpha1.NamespaceLabel{}).
		Watches(&corev1.Namespace{}, &handler.EnqueueRequestForObject{}).
		Complete(r)

}

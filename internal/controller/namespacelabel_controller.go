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
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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

func (r *NamespaceLabelReconciler) deleteExternalResources(ctx context.Context,
	namespaceLabel *danaiodanaiov1alpha1.NamespaceLabel,
	namespace *corev1.Namespace) error {

	logger := log.FromContext(ctx)

	logger.Info("deleteExternalResources", "namespaceLabel", namespaceLabel, "namespace", namespace)
	for key, _ := range namespaceLabel.Spec.Labels {
		delete(namespace.ObjectMeta.Labels, key)
		logger.Info("namespace.ObjectMeta.Labels", "HEY", namespace.ObjectMeta.Labels)
	}

	const ControllerUpdateAnnotation = "namespacelabeler.dana.io/controller-update"
	namespace.ObjectMeta.Annotations = make(map[string]string)
	namespace.ObjectMeta.Annotations[ControllerUpdateAnnotation] = "true"
	// update the namespace with the new labels
	if err := r.Update(ctx, namespace); err != nil {
		return err
	}

	return nil

}

func (r *NamespaceLabelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger := log.FromContext(ctx)

	logger.Info("Reconcile invoked", "req", req)

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

	namespaceLabelFinalizerName := "namespacelabeller.dana.io/finalizer"

	// examine DeletionTimestamp to determine if object is under deletion
	if namespaceLabel.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(&namespaceLabel, namespaceLabelFinalizerName) {
			controllerutil.AddFinalizer(&namespaceLabel, namespaceLabelFinalizerName)
			if err := r.Update(ctx, &namespaceLabel); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&namespaceLabel, namespaceLabelFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteExternalResources(ctx, &namespaceLabel, &namespace); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&namespaceLabel, namespaceLabelFinalizerName)

			if err := r.Update(ctx, &namespaceLabel); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}
	logger.Info("namespace.ObjectMeta.Labels", "req", namespace.ObjectMeta.Labels)

	// update the labels
	if namespace.ObjectMeta.Labels == nil {
		namespace.ObjectMeta.Labels = make(map[string]string)
	}

	logger.Info("namespacelabel status", "namespacelabel", namespaceLabel.Status)
	for key, value := range namespaceLabel.Spec.Labels {
		if _, exists := namespace.ObjectMeta.Labels[key]; !exists {
			logger.Info("check5", "key", exists)
			namespace.ObjectMeta.Labels[key] = value
		}
	}

	for key := range namespaceLabel.Status.LastAppliedLabels {
		if _, exists := namespaceLabel.Spec.Labels[key]; !exists {
			delete(namespace.ObjectMeta.Labels, key)
		}
	}

	// update the namespace with the new labels
	if err := r.Update(ctx, &namespace); err != nil {
		return ctrl.Result{}, err
	}

	// update the NamespaceLabel status with the total count of labels and last applied labels
	namespaceLabel.Status.LastAppliedLabels = namespaceLabel.Spec.Labels
	if err := r.Status().Update(ctx, &namespaceLabel); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&danaiodanaiov1alpha1.NamespaceLabel{}).
		Watches(&corev1.Namespace{}, handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, o client.Object) []reconcile.Request {
			namespace := o.(*corev1.Namespace)
			var requests []reconcile.Request
			var namespaceLabelList danaiodanaiov1alpha1.NamespaceLabelList
			if err := r.List(ctx, &namespaceLabelList, client.InNamespace(namespace.Name)); err != nil {
				return []reconcile.Request{}
			}
			for _, namespaceLabel := range namespaceLabelList.Items {
				requests = append(requests, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      namespaceLabel.Name,
						Namespace: namespaceLabel.Namespace,
					},
				})
			}

			return requests

		})).
		Complete(r)

}

package utils

import (
	corev1 "k8s.io/api/core/v1"
)

// Utility function to update labels on a namespace
func UpdateNamespaceLabels(namespace *corev1.Namespace, labelsToAdd map[string]string, labelsToRemove map[string]struct{}) {
	if namespace.ObjectMeta.Labels == nil {
		namespace.ObjectMeta.Labels = make(map[string]string)
	}

	for key, value := range labelsToAdd {
		namespace.ObjectMeta.Labels[key] = value
	}

	for key := range labelsToRemove {
		delete(namespace.ObjectMeta.Labels, key)
	}
}

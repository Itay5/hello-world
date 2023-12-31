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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceLabelSpec defines the desired state of NamespaceLabel
type NamespaceLabelSpec struct {

	// Lables consists of a collection of items known as labels, where each label is represented by a key-value pair.
	//
	Labels map[string]string `json:"labels,omitempty"`
}

// NamespaceLabelStatus defines the observed state of NamespaceLabel
type NamespaceLabelStatus struct {

	// LastAppliedLabels represents the last applied lables, it consists of the
	// last state of the spec Labels field before the last change.
	LastAppliedLabels map[string]string `json:"lastAppliedLabels,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Labels",type="string",JSONPath=".spec.labels",description="The labels of the namespace"

// NamespaceLabel is the Schema for the namespacelabels API
type NamespaceLabel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespaceLabelSpec   `json:"spec,omitempty"`
	Status NamespaceLabelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NamespaceLabelList contains a list of NamespaceLabel
type NamespaceLabelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceLabel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NamespaceLabel{}, &NamespaceLabelList{})
}

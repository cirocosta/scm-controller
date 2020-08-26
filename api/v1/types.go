// +kubebuilder:object:generate=true
// +groupName=experimental.kontinue.io
// +versionName=v1
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CommitStatusSpec struct {
	// +kubebuilder:validation:MinLength=1
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// +kubebuilder:validation:MinLength=1
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:MinLength=1
	Revision string `json:"revision,omitempty"`

	// +kubebuilder:validation:MinLength=1
	State string `json:"state,omitempty"`

	// +kubebuilder:validation:MinLength=1
	Label string `json:"label,omitempty"`

	// +kubebuilder:validation:MinLength=1
	Description string `json:"description,omitempty"`

	// +kubebuilder:validation:MinLength=1
	Target string `json:"target,omitempty"`
}

type CommitStatusStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource
type CommitStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CommitStatusSpec   `json:"spec,omitempty"`
	Status            CommitStatusStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type CommitStatusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CommitStatus `json:"items"`
}

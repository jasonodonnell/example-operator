package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PGClusterSpec defines the desired state of PGCluster
// +k8s:openapi-gen=true
type PGClusterSpec struct {
	Size int32 `json:"size"`
}

// PGClusterStatus defines the observed state of PGCluster
// +k8s:openapi-gen=true
type PGClusterStatus struct {
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PGCluster is the Schema for the pgclusters API
// +k8s:openapi-gen=true
type PGCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PGClusterSpec   `json:"spec,omitempty"`
	Status PGClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PGClusterList contains a list of PGCluster
type PGClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PGCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PGCluster{}, &PGClusterList{})
}

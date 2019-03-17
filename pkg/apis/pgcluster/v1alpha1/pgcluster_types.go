package v1alpha1

import (
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PGClusterSpec defines the desired state of PGCluster
// +k8s:openapi-gen=true
type PGClusterSpec struct {
	Database            string `json:"database"`
	ImageTag            string `json:"imageTag"`
	Mode                string `json:"mode"`
	PrimaryHost         string `json:"primaryHost"`
	PrimaryPort         string `json:"primaryPort"`
	ReplicationUser     string `json:"replicationUser"`
	ReplicationPassword string `json:"replicationPassword"`
	Username            string `json:"username"`
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

// Validate checks if PGCluster request is valid.
func (pg *PGCluster) Validate() error {
	if pg.Spec.ImageTag == "" {
		return errors.New("image tag not specified in spec")
	}

	if pg.Spec.Database == "" {
		return errors.New("database not specified in spec")
	}

	if pg.Spec.Mode != "primary" && pg.Spec.Mode != "replica" {
		return errors.New("postgres mode (primary or replica) must be specified in spec")
	}

	if pg.Spec.Mode == "replica" {
		if pg.Spec.PrimaryHost == "" {
			return errors.New("primary hostname must be specified to create a replica")
		}

		if pg.Spec.PrimaryPort == "" {
			return errors.New("primary port must be specified to create a replica")
		}

		if pg.Spec.ReplicationUser == "" {
			return errors.New("replication user must be specified to create a replica")
		}

		if pg.Spec.ReplicationPassword == "" {
			return errors.New("replication password must be specified to create a replica")
		}
	}

	return nil
}

/*
Copyright 2023 The Kubernetes-CSI-Addons Authors.

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

// VolumeGroupDeletionPolicy describes a policy for end-of-life maintenance of
// volume group contents
type VolumeGroupReplicationDeletionPolicy string

const (
	// VolumeGroupReplicationContentDelete means the group replication will be deleted from the
	// underlying storage system on release from its volume group replication.
	VolumeGroupReplicationContentDelete VolumeGroupReplicationDeletionPolicy = "Delete"

	// VolumeGroupReplicationContentRetain means the group replication will be left in its current
	// state on release from its volume group replication.
	VolumeGroupReplicationContentRetain VolumeGroupReplicationDeletionPolicy = "Retain"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=vgrclass
// +kubebuilder:printcolumn:name="Driver",type=string,JSONPath=`.driver`
// +kubebuilder:printcolumn:name="DeletionPolicy",type=string,JSONPath=`.volumeGroupDeletionPolicy`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// VolumeGroupReplicationClass is the Schema for the volumereplicationclasses API.
type VolumeGroupReplicationClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Driver is the driver expected to handle this VolumeGroupReplicationClass.
	Driver string `json:"driver"`

	// Parameters hold parameters for the driver.
	// These values are opaque to the system and are passed directly
	// to the driver.
	// +optional

	Parameters map[string]string `json:"parameters,omitempty"`

	// +optional
	VolumeGroupReplicationDeletionPolicy *VolumeGroupReplicationDeletionPolicy `json:"volumeGroupReplicationDeletionPolicy,omitempty"`
}

// +kubebuilder:object:root=true

// VolumeGroupReplicationClassList contains a list of VolumeGroupReplicationClass.
type VolumeGroupReplicationClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VolumeGroupReplicationClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VolumeGroupReplicationClass{}, &VolumeGroupReplicationClassList{})
}

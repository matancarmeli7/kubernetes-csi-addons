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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	VolumeGroupReplicationNameAnnotation = "replication.storage.openshift.io/volume-group-replication-name"
)

// GroupReplicationState represents the replication operations to be performed on the volume group.
// +kubebuilder:validation:Enum=primary;secondary
type GroupReplicationState string

const (
	// Primary GroupReplicationState enables mirroring and promotes the volumes under the volume group to primary.
	GroupPrimary GroupReplicationState = "primary"

	// Secondary GroupReplicationState demotes the volume to secondary and resyncs the volumes under the volume group if out of sync.
	GroupSecondary GroupReplicationState = "secondary"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".metadata.creationTimestamp",name=Age,type=date
// +kubebuilder:printcolumn:JSONPath=".spec.volumeGroupReplicationClass",name=volumeGroupReplicationClass,type=string
// +kubebuilder:printcolumn:JSONPath=".spec.replicationState",name=desiredState,type=string
// +kubebuilder:printcolumn:JSONPath=".status.state",name=currentState,type=string
// +kubebuilder:resource:shortName=vgr

// VolumeGroupReplication is the Schema for the VolumeGroupReplications API.
type VolumeGroupReplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec VolumeGroupReplicationSpec `json:"spec"`

	Status VolumeGroupReplicationStatus `json:"status,omitempty"`
}

// VolumeGroupReplicationSpec defines the desired state of VolumeGroupReplication.
type VolumeGroupReplicationSpec struct {
	// VolumeGroupReplicationClass is the VolumeGroupReplicationClass name for this VolumeGroupReplication resource
	// +kubebuilder:validation:Required
	VolumeGroupReplicationClass string `json:"volumeGroupReplicationClass"`

	// VolumeGroupReplicationClass is the VolumeGroupReplicationClass name for this VolumeGroupReplication resource
	// +kubebuilder:validation:Required
	VolumeGroupReplicationContentName *string `json:"volumeGroupReplicationContentName,omitempty"`

	// ReplicationState represents the replication operation to be performed on the volume.
	// Supported operations are "primary", "secondary" and "resync"
	// +kubebuilder:validation:Required
	ReplicationState GroupReplicationState `json:"replicationState"`

	// Label selector to identify all the PVCs that are in this group
	// that needs to be replicated to the peer cluster.
	PVCSelector *metav1.LabelSelector `json:"pvcSelector"`

	// Label selector to identify the VolumeGroupReplicationClass resources
	// that are scanned to select an appropriate VolumeGroupReplicationClass
	// for the VolumeReplication resource.
	//+optional
	ReplicationClassSelector metav1.LabelSelector `json:"replicationClassSelector,omitempty"`

	// VolumeGroupReplicationSecretRef is a reference to the secret object containing
	// sensitive information to pass to the CSI driver to complete the CSI
	// calls for VolumeGroupReplications.
	// +optional
	VolumeGroupReplicationSecretRef *corev1.SecretReference `json:"volumeGroupReplicationSecretRef,omitempty"`
}

// VolumeGroupReplicationStatus defines the observed state of VolumeGroupReplication.
type VolumeGroupReplicationStatus struct {
	State State `json:"state,omitempty"`

	// +optional
	BoundVolumeGroupReplicationContentName *string `json:"boundVolumeGroupReplicationContentName,omitempty"`

	// Conditions are the list of conditions and their status.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	Ready *bool `json:"ready,omitempty"`

	// A list of persistent volume claims
	// +optional
	ReplicatedPVCs []ReplicatedPVC `json:"replicatedPVCs,omitempty"`

	// observedGeneration is the last generation change the operator has dealt with
	// +optional
	ObservedGeneration int64        `json:"observedGeneration,omitempty"`
	LastCompletionTime *metav1.Time `json:"lastCompletionTime,omitempty"`
	// lastGroupSyncTime is the time of the most recent successful synchronization of all PVCs
	//+optional
	LastGroupSyncTime *metav1.Time `json:"lastGroupSyncTime,omitempty"`
}

type ReplicatedPVC struct {
	// Name of the VolRep/PVC resource
	//+optional
	Name string `json:"name,omitempty"`

	// Namespace of the VolRep/PVC resource
	//+optional
	Namespace string `json:"namespace,omitempty"`

	// PV name of the VolRep/PVC resource
	//+optional
	VolumeName string `json:"volumeName,omitempty"`

	// VolSyncPVC can be used to denote whether this PVC is protected by VolSync. Defaults to "false".
	//+optional
	ProtectedByVolSync bool `json:"protectedByVolSync,omitempty"`

	// Name of the StorageClass required by the claim.
	//+optional
	StorageClassName *string `json:"storageClassName,omitempty"`

	// Labels for the PVC
	//+optional
	Labels map[string]string `json:"labels,omitempty"`

	// AccessModes set in the claim to be replicated
	//+optional
	AccessModes []corev1.PersistentVolumeAccessMode `json:"accessModes,omitempty"`

	// Resources set in the claim to be replicated
	//+optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Conditions for this protected pvc
	//+optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Time of the most recent successful synchronization for the PVC, if
	// protected in the async or volsync mode
	//+optional
	LastSyncTime *metav1.Time `json:"lastSyncTime,omitempty"`
}

// +kubebuilder:object:root=true

// VolumeGroupReplicationList contains a list of VolumeGroupReplication.
type VolumeGroupReplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VolumeGroupReplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VolumeGroupReplication{}, &VolumeGroupReplicationList{})
}

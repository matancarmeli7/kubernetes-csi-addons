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

// VolumeGroupReplication is the Schema for the VolumeGroupReplications API.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".metadata.creationTimestamp",name=Age,type=date
// +kubebuilder:printcolumn:JSONPath=".spec.volumeGroupReplicationClass",name=volumeGroupReplicationClass,type=string
// +kubebuilder:printcolumn:JSONPath=".spec.replicationState",name=desiredState,type=string
// +kubebuilder:printcolumn:JSONPath=".status.state",name=currentState,type=string
// +kubebuilder:resource:shortName=vgrc
type VolumeGroupReplicationContent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired characteristics of a group replication content
	Spec VolumeGroupReplicationContentSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// Status represents the latest observed state of the group replication content
	// +optional
	Status *VolumeGroupReplicationContentStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// VolumeGroupReplicationContentSpec describes the common attributes of a group replication content
type VolumeGroupReplicationContentSpec struct {
	// VolumeGroupReplicationRef specifies the VolumeGroupReplication object
	// to which this VolumeGroupReplicationContent object is bound.
	// +optional
	VolumeGroupReplicationRef corev1.ObjectReference `json:"volumeGroupReplicationRef"`
	// +optional
	VolumeGroupReplicationDeletionPolicy VolumeGroupReplicationDeletionPolicy `json:"volumeGroupReplicationDeletionPolicy"`
	// +optional
	Driver string `json:"driver"`
	// This field may be unset for pre-provisioned replications.
	// For dynamic provisioning, this field must be set.
	// +optional
	VolumeGroupReplicationClassName *string `json:"volumeGroupReplicationClassName"`
	// VolumeGroupReplicationHandle is a unique id returned by the CSI driver
	// to identify the VolumeGroupReplication on the storage system.
	// +optional
	VolumeGroupReplicationHandle *string `json:"volumeGroupReplicationHandle"`
}

type VolumeGroupReplicationContentStatus struct {
	State State `json:"state,omitempty"`

	// Conditions are the list of conditions and their status.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	Ready *bool `json:"ready,omitempty"`

	// A list of persistent volumes
	// +optional
	ReplicatedPVs []ReplicatedPV `json:"replicatedPVs,omitempty"`

	// observedGeneration is the last generation change the operator has dealt with
	// +optional
	ObservedGeneration int64        `json:"observedGeneration,omitempty"`
	LastCompletionTime *metav1.Time `json:"lastCompletionTime,omitempty"`
	// lastGroupSyncTime is the time of the most recent successful synchronization of all PVCs
	//+optional
	LastGroupSyncTime *metav1.Time `json:"lastGroupSyncTime,omitempty"`
}

type ReplicatedPV struct {
	// Name of the VolRep/PV resource
	//+optional
	Name string `json:"name,omitempty"`

	// VolSyncPV can be used to denote whether this PV is protected by VolSync. Defaults to "false".
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
